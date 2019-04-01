package main

import (
	"fmt"
	"github.com/cyyber/FindMaxNumberRPC/config"
	"github.com/cyyber/FindMaxNumberRPC/pkg/server"
)

func main() {
	c, err := config.GetConfig()
	if err != nil {
		fmt.Print("Error getting config: ", err)
		return
	}
	findMaxNumberServer := &server.FindMaxNumberServer{}
	server.StartMaxNumberServer(c, findMaxNumberServer)
}
