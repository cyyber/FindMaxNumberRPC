package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/cyyber/FindMaxNumberRPC/config"
	"github.com/cyyber/FindMaxNumberRPC/pkg/generated"
	"github.com/cyyber/FindMaxNumberRPC/pkg/helper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io"
	"time"
)

type Client struct {
	Conn *grpc.ClientConn
}

// Stream number to server
func (c *Client) Send(stream generated.FindMaxNumber_GetMaxNumberClient, number int64) error {
	n := &generated.Number{Number: number}
	err := stream.Send(n)
	if err != nil {
		fmt.Println("Error while sending Stream: ", err)
		return err
	}
	return nil
}

// Receive Stream from server
func (c *Client) Recv(stream generated.FindMaxNumber_GetMaxNumberClient) (int64, error) {
	in, err := stream.Recv()
	if err == io.EOF {
		fmt.Print("Connection Closed")
		return 0, err
	}
	if err != nil {
		fmt.Print("Error while receiving Stream: ", err)
		return 0, err
	}
	return in.Number, err
}

// Run RPC gRPC bi-directional streaming client, connecting to server.
func (c *Client) RunClient(conf *config.Config) error {
	defer c.Conn.Close()
	client := generated.NewFindMaxNumberClient(c.Conn)
	ctx, cancel := context.WithTimeout(context.Background(), 100 * time.Second)
	defer cancel()
	stream, err := client.GetMaxNumber(ctx)
	if err != nil {
		fmt.Print("Error while creating stream GetMaxNumber: ", err)
		return err
	}

	waitc := make(chan struct{})
	go func() {
		for {
			number, err := c.Recv(stream)
			if err != nil {
				close(waitc)
			}
			fmt.Print(number, " ")
		}
	}()

	for {
		var number int64
		_, err := fmt.Scan(&number)
		err = c.Send(stream, number)
		if err != nil {
			fmt.Println("Error while sending Stream: ", err)
			break
		}
	}
	err = stream.CloseSend()
	if err != nil {
		fmt.Println("Error while closing Stream: ", err)
	}
	<-waitc
	return nil
}

// Creates NewClient struct with connection established with server
func NewClient(conf *config.Config) (*Client, error) {
	c := &Client{}
	certificate, certPool, err := helper.GetCertificateAndCertPool(
		conf.Prod.ClientCrt,
		conf.Prod.ClientKey,
		conf.Prod.CertAuth)
	if err != nil {
		return nil, err
	}

	creds := credentials.NewTLS(&tls.Config{
		ServerName:   conf.Prod.Host,
		Certificates: []tls.Certificate{*certificate},
		RootCAs:      certPool,
	})

	serverAddr := fmt.Sprintf("%s:%d", conf.Prod.Host, conf.Prod.Port)
	c.Conn, err = grpc.Dial(serverAddr, grpc.WithTransportCredentials(creds))
	if err != nil {
		fmt.Println("Error while dialing server: ", err)
		return nil, err
	}
	return c, nil
}
