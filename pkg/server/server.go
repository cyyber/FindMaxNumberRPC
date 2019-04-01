package server

import (
	"crypto/tls"
	"fmt"
	"github.com/cyyber/FindMaxNumberRPC/config"
	"github.com/cyyber/FindMaxNumberRPC/pkg/generated"
	"github.com/cyyber/FindMaxNumberRPC/pkg/helper"
	"github.com/cyyber/FindMaxNumberRPC/pkg/maxnumber"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io"
	"net"
	"sync"
)

type FindMaxNumberServer struct {
	lock sync.Mutex

	isRunning bool
	GrpcServer *grpc.Server
}


func (f *FindMaxNumberServer) GetMaxNumber(stream generated.FindMaxNumber_GetMaxNumberServer) error {
	m := maxnumber.NewMaxNumber()

	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		value, isNewMax := m.FindMaxNumber(in.Number)
		if isNewMax {
			response := &generated.Number{Number: value}
			if err := stream.Send(response); err != nil {
				return err
			}
		}
	}
}

func (f *FindMaxNumberServer) SetIsRunning(value bool) {
	f.lock.Lock()
	defer f.lock.Unlock()

	f.isRunning = value
}

func (f *FindMaxNumberServer) GetIsRunning() bool {
	f.lock.Lock()
	defer f.lock.Unlock()

	return f.isRunning
}

func StartMaxNumberServer(conf *config.Config, findMaxNumberServer *FindMaxNumberServer) error {
	certificate, certPool, err := helper.GetCertificateAndCertPool(
		conf.Prod.ServerCrt,
		conf.Prod.ServerKey,
		conf.Prod.CertAuth)
	if err != nil {
		return err
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", conf.Prod.Host, conf.Prod.Port))
	if err != nil {
		fmt.Println("Error while Listening: ", err)
		return err
	}

	creds := credentials.NewTLS(&tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{*certificate},
		ClientCAs:    certPool,
	})

	findMaxNumberServer.GrpcServer = grpc.NewServer(grpc.Creds(creds))
	generated.RegisterFindMaxNumberServer(findMaxNumberServer.GrpcServer, findMaxNumberServer)
	findMaxNumberServer.SetIsRunning(true)
	err = findMaxNumberServer.GrpcServer.Serve(lis)
	if err != nil {
		fmt.Println("Error while Starting Server: ", err)
		findMaxNumberServer.SetIsRunning(false)
		return err
	}
	findMaxNumberServer.SetIsRunning(false)
	return nil
}
