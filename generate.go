// +build ignore

package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	out, err := os.Create(filepath.Join(wd, "roler", "template.go"))
	if err != nil {
		panic(err)
	}

	roleTemplate := filepath.Join(wd, "roler", "role-cf.json.tmpl")

	out.Write([]byte("// This file is autogenerated; DO NOT EDIT\n// See generate.go for more info.\npackage roler\n\n"))

	out.Write([]byte("const RoleTemplate = `"))
	t, err := ioutil.ReadFile(roleTemplate)
	if err != nil {
		panic(err)
	}

	out.Write([]byte(t))
	out.Write([]byte("`\n"))
}
