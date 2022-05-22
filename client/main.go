package main

import (
	"context"
	"fmt"
	grpc_client "new_diplom_client/grpc-client"
	"new_diplom_client/handlers"
	"new_diplom_client/loops"
)

const address = "localhost:50051"

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	fmt.Println("InitApp")
	userClient := grpc_client.NewUserClient(address)
	userHandler := handlers.NewUserHandler(userClient)

	userLoop := loops.NewUserLoop(address, userHandler)
	userLoop.MainLoop(ctx)
	cancel()
}
