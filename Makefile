APPENV ?= testenv
PROJECT := spanx
REV ?= latest

all: build

clean:
	rm -fr target bin pkg

fmt:
	@gofmt -w ./

deps:
	docker-compose up -d
	docker run --link $(PROJECT)_postgres_1:postgres aanand/wait

migrate:
	migrate -url $(POSTGRES_CONN) -path ./migrations up

build: deps $(APPENV)
	docker run \
		--link $(PROJECT)_postgres_1:postgres \
		--env-file ./$(APPENV) \
		-e AWS_DEFAULT_REGION \
		-e AWS_ACCESS_KEY_ID \
		-e AWS_SECRET_ACCESS_KEY \
		-e "TARGETS=linux/amd64" \
		-e PROJECT=github.com/opsee/$(PROJECT) \
		-v `pwd`:/gopath/src/github.com/opsee/$(PROJECT) \
		quay.io/opsee/build-go:16
	docker build -t quay.io/opsee/$(PROJECT):$(REV) .

run: docker
	docker run \
		--link $(PROJECT)_postgres_1:postgres \
		--env-file ./$(APPENV) \
		-e AWS_DEFAULT_REGION \
		-e AWS_ACCESS_KEY_ID \
		-e AWS_SECRET_ACCESS_KEY \
		-p 9095:9095 \
		--rm \
		quay.io/opsee/spanx:latest

deploy-role:
	docker run -it quay.io/opsee/spanx /roler | aws s3 cp --content-disposition inline --content-type application/json --region us-east-1 --acl public-read - s3://opsee-bastion-cf-us-east-1/beta/opsee-role.json
	for region in ap-northeast-1 ap-northeast-2 ap-southeast-1 ap-southeast-2 eu-central-1 eu-west-1 sa-east-1 us-west-1 us-west-2; do \
		aws s3 cp --source-region us-east-1 --content-disposition inline --content-type application/json --region $$region --acl public-read s3://opsee-bastion-cf-us-east-1/beta/opsee-role.json s3://opsee-bastion-cf-$$region/beta/ ; \
	done

.PHONY: docker run migrate clean all
