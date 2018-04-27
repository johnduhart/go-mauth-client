// The following directive is necessary to make the package coherent:

// +build ignore

// This program generates version.go. It can be invoked by running
// go generate

package main

import (
	"io/ioutil"
	"log"
	"os"
	"text/template"
	"time"
)

var VERSION_FILE_NAME = "VERSION"

func main() {
	// Open the file
	versionString, err := ioutil.ReadFile(VERSION_FILE_NAME)
	die(err)

	f, err := os.Create("version.go")
	die(err)
	defer f.Close()

	versionTemplate.Execute(f, struct {
		Timestamp time.Time
		Version   string
	}{Timestamp: time.Now(),
		Version: string(versionString),
	})
}

func die(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var versionTemplate = template.Must(template.New("").Parse(`// Code generated by go generate; DO NOT EDIT.
// This file was generated by robots at
// {{ .Timestamp }}
package go_mauth_client

var VersionString = "{{ .Version }}"
`))