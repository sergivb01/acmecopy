package main

import (
	"context"
	"io/ioutil"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/sergivb01/acmecopy/api"
)

var lines = []string{"hola test123", "#"}
var expected = []string{
	"ENTRA TEXT ACABAT EN #:",
	"TEXT REVES:",
	"test123 hola ",
}

func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer lis.Close()

	srv := grpc.NewServer()
	api.RegisterCompilerServer(srv, &compileServer{})

	if err := srv.Serve(lis); err != nil {
		log.Fatalf("error serving grpc server: %s", err)
	}
}

type compileServer struct{}

func (c *compileServer) CompileFiles(ctx context.Context, req *api.CompileRequest) (*api.CompileResponse, error) {
	for _, compileFile := range req.Files {
		if err := ioutil.WriteFile(compileFile.FileName, compileFile.Content, 744); err != nil {
			return &api.CompileResponse{}, err
		}
	}

	var res api.CompileResponse

	buildRes, err := compileFiles(req.Files)
	if err != nil {
		return nil, err
	}
	res.Build = &buildRes.apiResponse

	execRes, err := runTarget()
	if err != nil {
		return nil, err
	}
	res.Execute = &execRes.apiResponse

	return &res, nil
}

//
// func mainshit() {
// 	start := time.Now()
// 	if err := compileProject(); err != nil {
// 		log.Fatalf("error running command: %s", err)
// 		return
// 	}
// 	log.Printf("compiled in %s", time.Since(start))
//
// 	start = time.Now()
// 	buff, err := runProject()
// 	if err != nil {
// 		log.Fatalf("error running project: %s", err)
// 		return
// 	}
//
// 	log.Printf("ran project in %s", time.Since(start))
//
// 	scan := bufio.NewScanner(buff)
// 	i := 0
// 	for scan.Scan() {
// 		str := scan.Text()
// 		if str != expected[i] {
// 			log.Printf("output mismatch, expected %q and received %q", expected[i], str)
// 		}
// 		i++
// 		fmt.Println(str)
// 	}
// }
