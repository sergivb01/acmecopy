package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/sergivb01/acmecopy/api"
)

func main() {
	creds, err := credentials.NewServerTLSFromFile("C:\\Users\\Sergi\\Desktop\\acmecopy\\certs\\certificate.pem", "C:\\Users\\Sergi\\Desktop\\acmecopy\\certs\\key.pem")
	if err != nil {
		log.Fatalf("error creating credentials from TLS file: %v", err)
		return
	}

	srv := grpc.NewServer(grpc.Creds(creds))
	api.RegisterCompilerServer(srv, &compileServer{})

	lis, err := net.Listen("tcp", "localhost:3000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer lis.Close()

	if err := srv.Serve(lis); err != nil {
		log.Fatalf("error serving grpc server: %s", err)
	}
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
