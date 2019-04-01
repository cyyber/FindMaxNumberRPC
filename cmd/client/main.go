package main

import (
	"fmt"
	"github.com/cyyber/FindMaxNumberRPC/config"
	"github.com/cyyber/FindMaxNumberRPC/pkg/client"
)

func main() {
	conf, err := config.GetConfig()
	if err != nil {
		fmt.Println("Error GetConfig: ", err)
		return
	}
	c, err := client.NewClient(conf)
	if err != nil {
		fmt.Println("Error NewClient: ", err)
	}
	c.RunClient(conf)
}
