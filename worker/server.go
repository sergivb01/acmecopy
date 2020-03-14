package main

import (
	"context"
	"io/ioutil"
	"log"
	"os"

	"google.golang.org/grpc"

	"github.com/sergivb01/acmecopy/api"
)

type compileServer struct {
	Server grpc.Server
}

func (c *compileServer) CompileFiles(_ context.Context, req *api.CompileRequest) (*api.CompileResponse, error) {
	for _, compileFile := range req.Files {
		if err := ioutil.WriteFile(compileFile.FileName, compileFile.Content, 744); err != nil {
			return &api.CompileResponse{}, err
		}
	}
	defer func() {
		pwd, _ := os.Getwd()
		if err := os.RemoveAll(pwd); err != nil {
			log.Printf("error cleaning up %s files: %v\n", pwd, err)
		}
	}()

	var res api.CompileResponse
	buildRes, err := compileFiles(req.Files)
	if err != nil {
		return &res, err
	}
	res.Build = &buildRes.apiResponse

	execRes, err := runTarget(req.Input, req.ExpectedOutput)
	res.Execute = &execRes.apiResponse

	return &res, err
}
