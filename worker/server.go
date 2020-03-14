package main

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"

	"google.golang.org/grpc"

	"github.com/sergivb01/acmecopy/api"
)

type compileServer struct {
	Server grpc.Server
}

func (c *compileServer) CompileFiles(_ context.Context, req *api.CompileRequest) (*api.CompileResponse, error) {
	pwd, _ := os.Getwd()
	tempDir, err := ioutil.TempDir(pwd, "*")
	if err != nil {
		return &api.CompileResponse{}, err
	}
	defer os.RemoveAll(tempDir)

	for _, file := range req.Files {
		if err := ioutil.WriteFile(filepath.Join(tempDir, file.FileName), file.Content, 744); err != nil {
			return &api.CompileResponse{}, err
		}
	}

	var res api.CompileResponse
	buildRes, err := compileFiles(tempDir, req.Files)
	if err != nil {
		return &res, err
	}
	res.Build = &buildRes.apiResponse

	execRes, err := runTarget(tempDir, req.Input, req.ExpectedOutput)
	res.Execute = &execRes.apiResponse

	return &res, err
}
