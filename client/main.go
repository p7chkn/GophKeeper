package main

import (
	"context"
	"flag"
	"fmt"
	grpc_client "new_diplom_client/grpc-client"
	"new_diplom_client/handlers"
	"new_diplom_client/loops"
	"os"
)

var (
	BuildTime  string
	AppVersion string
	address    string
)

const defaultAddress = "localhost:50051"

func init() {
	envAddress := os.Getenv("ADDRESS")
	if envAddress != "" {
		address = envAddress
	}
	address = *flag.String("a", defaultAddress, "address of gGRPC server")
}

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
