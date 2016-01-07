package roler

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/opsee/basic/com"
	"github.com/opsee/spanx/store"
	log "github.com/sirupsen/logrus"
	"os"
	"regexp"
	"strconv"
)

var (
	awsSession *session.Session
	arnRegexp  = regexp.MustCompile(`^arn:aws:iam::(\d{12}):user.+$`)

	AccountNotFound         = errors.New("AWS account for that customer not found.")
	InsufficientPermissions = errors.New("IAM role or user provided has insufficient permissions to provision a role.")
)

func init() {
	var (
		ec2meta *ec2metadata.EC2Metadata
		region  string
		creds   *credentials.Credentials
	)

	ec2meta = ec2metadata.New(session.New())
	if ec2meta.Available() {
		// ignoring error here since we'll try to get region from env later
		region, _ = ec2meta.Region()
	}

	if region == "" {
		region = getEnvRegion()
	}

	creds = credentials.NewChainCredentials(
		[]credentials.Provider{
			&ec2rolecreds.EC2RoleProvider{
				Client: ec2meta,
			},
			&credentials.EnvProvider{},
		},
	)

	awsSession = session.New(&aws.Config{
		Credentials: creds,
		MaxRetries:  aws.Int(11),
		Region:      aws.String(region),
	})
}

func ResolveCredentials(db store.Store, customerID, accessKey, secretKey string) (credentials.Value, error) {
	var (
		account *com.Account
		creds   credentials.Value
		err     error
		logger  = log.WithFields(log.Fields{"customer-id": customerID})
	)

	account, err = db.GetAccount(&store.GetAccountRequest{CustomerID: customerID, Active: true})
	if err != nil {
		return creds, err
	}

	if account != nil {
		// test if the stored account access is still valid,
		// if so, just return those creds. if not, recover
		// from the error by attempting to provision another role
		creds, err = getAccountCredentials(db, account)

		if err != nil {
			logger.WithError(err).Error("attempt to access stored account role failed, going to provision another")

			// yeah, so just delete this sucker and we'll put a new one
			err = deleteAccountCredentials(db, account)
			if err != nil {
				return creds, err
			}

			// nil out our value here
			creds = credentials.Value{}
		} else {
			return creds, nil
		}
	}

	// grab the account id and persist it in an account object
	iamClient := iam.New(session.New(&aws.Config{
		Credentials: credentials.NewStaticCredentials(
			accessKey,
			secretKey,
			"",
		),
		MaxRetries: aws.Int(11),
		Region:     aws.String(getEnvRegion()),
	}))

	user, err := iamClient.GetUser(nil)
	if err != nil {
		return creds, handleAWSError("GetUser", err)
	}

	arn := aws.StringValue(user.User.Arn)
	if arn == "" {
		return creds, fmt.Errorf("No user found when fetching the current user from provided credentials")
	}

	accountID, err := parseARNAccount(arn)
	if err != nil {
		return creds, err
	}

	account = &com.Account{
		ID:         accountID,
		CustomerID: customerID,
		Active:     true,
	}

	err = db.PutAccount(account)
	if err != nil {
		return creds, err
	}

	// time 2 provision a policy / role for us in their aws account
	_, err = iamClient.CreatePolicy(&iam.CreatePolicyInput{
		PolicyDocument: aws.String(Policy),
		PolicyName:     aws.String(PolicyName),
		Description:    aws.String("A policy for Opsee monitoring"),
	})

	if err = handleAWSError("CreatePolicy", err); err != nil {
		return creds, err
	}

	_, err = iamClient.CreateRole(&iam.CreateRoleInput{
		AssumeRolePolicyDocument: aws.String(fmt.Sprintf(AssumeRolePolicy, customerID)),
		RoleName:                 aws.String(RoleName),
	})

	if err = handleAWSError("CreateRole", err); err != nil {
		return creds, err
	}

	_, err = iamClient.AttachRolePolicy(&iam.AttachRolePolicyInput{
		PolicyArn: aws.String(policyARN(account)),
		RoleName:  aws.String(RoleName),
	})

	if err = handleAWSError("AttachRolePolicy", err); err != nil {
		return creds, err
	}

	// go ahead and return some creds
	return getAccountCredentials(db, account)
}

func GetCredentials(db store.Store, customerID string) (credentials.Value, error) {
	account, err := db.GetAccount(&store.GetAccountRequest{CustomerID: customerID, Active: true})
	if err != nil {
		return credentials.Value{}, err
	}

	return getAccountCredentials(db, account)
}

func DeleteCredentials(db store.Store, customerID string) error {
	account, err := db.GetAccount(&store.GetAccountRequest{CustomerID: customerID, Active: true})
	if err != nil {
		return err
	}

	return deleteAccountCredentials(db, account)
}

func deleteAccountCredentials(db store.Store, account *com.Account) error {
	return db.DeleteAccount(account)
}

func getAccountCredentials(db store.Store, account *com.Account) (credentials.Value, error) {
	if account != nil {
		return stscreds.NewCredentials(awsSession, roleARN(account), func(arp *stscreds.AssumeRoleProvider) {
			arp.ExternalID = aws.String(account.CustomerID)
		}).Get()
	} else {
		return credentials.Value{}, AccountNotFound
	}
}

func handleAWSError(meth string, err error) error {
	if awsErr, ok := err.(awserr.Error); ok {
		log.WithError(err).Error("IAM error")

		switch awsErr.Code() {
		case "AccessDenied":
			return InsufficientPermissions

		case "EntityAlreadyExists": // ignoring this
			return nil
		}
	}

	return err
}

func parseARNAccount(arn string) (int, error) {
	matches := arnRegexp.FindStringSubmatch(arn)
	if matches[1] == "" {
		return 0, fmt.Errorf("No account ID match in ARN")
	}

	id, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, err
	}

	return id, nil
}

func roleARN(a *com.Account) string {
	return fmt.Sprintf("arn:aws:iam::%d:role/%s", a.ID, RoleName)
}

func policyARN(a *com.Account) string {
	return fmt.Sprintf("arn:aws:iam::%d:policy/%s", a.ID, PolicyName)
}

func getEnvRegion() string {
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = os.Getenv("AWS_DEFAULT_REGION")
	}

	if region == "" {
		// not sure it matters really
		region = "us-west-1"
	}

	return region
}
