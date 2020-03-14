package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/sergivb01/acmecopy/api"
)

var cli api.CompilerClient

var fileNames = []string{
	"main.cpp",
	"PilaString.cpp",
	"PilaString.h",
}

func main() {
	creds, err := credentials.NewClientTLSFromFile("C:\\Users\\Sergi\\Desktop\\acmecopy\\certs\\certificate.pem", "")
	if err != nil {
		log.Fatalf("error creating credentials from TLS file: %v", err)
		return
	}

	conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("error dialing to server: %s", err)
		return
	}
	defer conn.Close()

	cli = api.NewCompilerClient(conn)
	requestUpload()
}

func requestUpload() {
	var files []*api.File
	for _, fileName := range fileNames {
		b, err := ioutil.ReadFile(fileName)
		if err != nil {
			log.Fatalf("error reading file: %s", err)
			return
		}

		files = append(files, &api.File{
			FileName: fileName,
			Content:  b,
		})
	}

	response, err := cli.CompileFiles(context.TODO(), &api.CompileRequest{
		Files: files,
		Input: []string{
			"hola test123",
			"#",
		},
		ExpectedOutput: []string{
			"ENTRA TEXT ACABAT EN #:",
			"TEXT REVES:",
			"test123 hola ",
		},
	})
	if err != nil {
		log.Fatalf("error uploading: %s", err)
	}
	b, err := json.Marshal(response)
	if err != nil {
		log.Fatalf("couldn't marshal to json: %s", err)
	}
	fmt.Printf("%s\n", b)
}
