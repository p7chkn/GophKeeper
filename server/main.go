package main

import (
	"context"
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/jmoiron/sqlx"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"log"
	"net"
	"new_diplom/authorization"
	"new_diplom/configuration"
	"new_diplom/database"
	"new_diplom/handlers"
	"new_diplom/pb"
	"new_diplom/services"
	"new_diplom/setup"
	"os"
)

var (
	grpcServer *grpc.Server
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := configuration.NewConfig()

	g, ctx := errgroup.WithContext(ctx)

	interrupt := make(chan os.Signal, 1)

	db, err := sqlx.Connect("postgres", cfg.PostgresString)
	if err != nil {
		log.Fatal(err)
	}
	err = setup.MustSetupDatabase(db.DB)
	if err != nil {
		log.Fatal(err)
	}
	repo := database.NewPostgresDataBase(db)
	secretService := services.NewSecretService(repo)
	userService := services.NewUserService(repo, cfg.AccessTokenLiveTimeMinutes, cfg.RefreshTokenLiveTimeDays,
		cfg.AccessTokenSecret, cfg.RefreshTokenSecret)

	grpcSecrets := handlers.NewGrpcSecrets(secretService)
	grpcUsers := handlers.NewGrpcUsers(userService)

	g.Go(func() error {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GrpcPort))
		if err != nil {
			log.Printf("gRPC server failed to listen: %v", err.Error())
			return err
		}
		grpcServer = grpc.NewServer(grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				// auth
				grpc_auth.UnaryServerInterceptor(func(ctx context.Context) (context.Context, error) {

					userID, err := authorization.TokenValid(ctx, cfg.AccessTokenSecret)
					if err != nil {
						userID = ""
					}
					newCtx := context.WithValue(ctx, "userID", userID)
					return newCtx, nil
				}),
			)))
		//grpcServer = grpc-client.NewServer()
		pb.RegisterSecretsServer(grpcServer, grpcSecrets)
		pb.RegisterUsersServer(grpcServer, grpcUsers)
		log.Printf("server listening at %v", lis.Addr())
		return grpcServer.Serve(lis)
	})

	select {
	case <-interrupt:
		break
	case <-ctx.Done():
		break
	}

	log.Println("Receive shutdown signal")

	cancel()

	if grpcServer != nil {
		grpcServer.GracefulStop()
	}

	err = g.Wait()
	if err != nil {
		log.Printf("server returning an error: %v", err)
	}
}
