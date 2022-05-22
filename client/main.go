package main

import (
	"context"
	"fmt"
	grpc_client "new_diplom_client/grpc-client"
	"new_diplom_client/handlers"
	"new_diplom_client/loops"
)

var (
	BuildTime  string
	AppVersion string
)

const address = "localhost:50051"

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	fmt.Println("InitApp")
	fmt.Printf("App version: %v, Date compile: %v\n", AppVersion, BuildTime)
	userClient := grpc_client.NewUserClient(address)
	userHandler := handlers.NewUserHandler(userClient)

	userLoop := loops.NewUserLoop(address, userHandler)
	userLoop.MainLoop(ctx)
	cancel()
}
