package main

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/sergivb01/acmecopy/api"
)

func (s *Server) createGRPCClient() error {
	creds, err := credentials.NewClientTLSFromFile("C:\\Users\\Sergi\\Desktop\\acmecopy\\certs\\certificate.pem", "")
	if err != nil {
		return fmt.Errorf("cannot read credentials from TLS file: %w", err)
	}

	s.clientConn, err = grpc.Dial("localhost:3000", grpc.WithTransportCredentials(creds))
	if err != nil {
		return fmt.Errorf("cannot dial to GRPC server: %w", err)
	}

	s.cli = api.NewCompilerClient(s.clientConn)

	return nil
}
