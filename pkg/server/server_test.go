package server

import (
	"context"
	"fmt"
	"github.com/cyyber/FindMaxNumberRPC/config"
	"github.com/cyyber/FindMaxNumberRPC/pkg/client"
	"github.com/cyyber/FindMaxNumberRPC/pkg/generated"
	"github.com/phayes/freeport"
	"testing"
	"time"
)

var (
	inputs = []int64{1, 5, 3, 6, 2, 20}  // Input test data, that will be sent by client to server
	outputs = []int64{1, 5, 6, 20}  // Expected output from server, that will be received by the client
)

func TestFindMaxNumberServer_GetMaxNumber(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test")
	}
	freePort, err := freeport.GetFreePort() // Get random unused port
	if err != nil {
		t.Error("Error GetFreePort: ", err)
	}
	conf, err := config.GetConfig()
	if err != nil {
		t.Error("Error GetConfig: ", err)
	}

	// Configuration for Integration Test
	conf.Prod.Port = uint16(freePort)
	conf.Prod.ServerCrt = "test_data/server.crt"
	conf.Prod.ServerKey = "test_data/server.key"
	conf.Prod.ClientCrt = "test_data/client.crt"
	conf.Prod.ClientKey = "test_data/client.key"
	conf.Prod.CertAuth = "test_data/My_Root_CA.crt"

	findMaxNumberServer := &FindMaxNumberServer{}
	go func() {
		err := StartMaxNumberServer(conf, findMaxNumberServer)
		if err != nil {
			fmt.Println("Error starting server: ", err)
		}
	}()
	time.Sleep(2 * time.Second)  // Wait for 2 seconds, before testing if server started
	if !findMaxNumberServer.GetIsRunning() {
		t.Error("Failed to Start Server")
	}
	defer findMaxNumberServer.GrpcServer.Stop()

	c, err := client.NewClient(conf)
	if err != nil {
		t.Error("Failed to Start Client: ", err)
	}

	defer c.Conn.Close()
	client := generated.NewFindMaxNumberClient(c.Conn)
	ctx, cancel := context.WithTimeout(context.Background(), 100 * time.Second)
	defer cancel()
	stream, err := client.GetMaxNumber(ctx)
	if err != nil {
		t.Error("Error while creating stream GetMaxNumber: ", err)
	}

	// Send all inputs to server
	for _, input := range inputs {
		err = c.Send(stream, input)
		if err != nil {
			t.Error(fmt.Sprintf("Failed to Send Input %d: %s", input, err))
			break
		}
	}

	// Receive from server and verify the output
	for _, output := range outputs {
		resultChan := make(chan int64, 1)
		errChan := make(chan error, 1)
		go func() {
			value, err := c.Recv(stream)
			if err != nil {
				errChan <- err
				return
			}
			resultChan <- value
		}()

		select {
		case value := <-resultChan:
			if output != value {
				t.Error(fmt.Sprintf("Unexpected Recv Value\nExpected %d, Found %d", output, value))
				break
			}
		case err := <-errChan:
			if err != nil {
				t.Error("Failed to Recv: ", err)
				break
			}
		case <-time.After(10 * time.Second):  // Timeout in case no data is received in 10 seconds
			t.Error("Receive Timeout")
		}
	}
}
