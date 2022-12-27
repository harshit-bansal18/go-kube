package main

import (
	"fmt"
	"time"

	"github.com/harshit-bansal18/go-kube/config"
)

func multiply() {
	fmt.Println("Enter two numbers (separated by space) in one line\n Program outputs their product.")
	for {
		var a, b int
		fmt.Scanln(&a, &b)
		fmt.Println("Output: ", a*b)
	}
}

func main() {
	fmt.Println("Starting Go module....")

	configclient := config.GetClient("default")
	config := configclient.ReadConfig("sample-config")
	fmt.Println(config)
	go configclient.WatchConfig("sample-config")
	for {
		// do nothing
		time.Sleep(10 * time.Second)
	}
}
