package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"

	"google.golang.org/grpc"

	"github.com/sergivb01/acmecopy/api"
)

var cli api.CompilerClient

var fileNames = []string{
	"main.cpp",
	"PilaString.cpp",
	"PilaString.h",
}

func main() {
	conn, err := grpc.Dial(":8080", grpc.WithInsecure(), grpc.WithBlock())
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

	response, err := cli.CompileFiles(context.TODO(), &api.CompileRequest{Files: files})
	if err != nil {
		log.Fatalf("error uploading: %s", err)
	}
	fmt.Println(response.String())
}
