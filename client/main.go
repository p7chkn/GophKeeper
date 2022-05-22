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
	//log := logger.InitLogger()
	fmt.Println("InitApp")
	userClient := grpc_client.NewUserClient(address)
	userHandler := handlers.NewUserHandler(userClient)

	userLoop := loops.NewUserLoop(address, userHandler)
	userLoop.MainLoop(ctx)
	//for {
	//	fmt.Println("To login input l, to register input r, to quit input q:")
	//	var input string
	//	_, err := fmt.Scan(&input)
	//	if err != nil {
	//		break
	//	}
	//	switch {
	//	case input == "r":
	//		AccessToken, RefreshToken, err = userHandler.RegisterUser(ctx)
	//		if err != nil {
	//			fmt.Printf("Error with registration: %v \n", err)
	//			continue
	//		}
	//		if AccessToken != "" && RefreshToken != "" {
	//			clientLoop(ctx, address, AccessToken, RefreshToken)
	//		} else {
	//			fmt.Println("Problem with register")
	//		}
	//	case input == "l":
	//		AccessToken, RefreshToken, err = userHandler.AuthUser(ctx)
	//		if err != nil {
	//			fmt.Println("Wrong login or password")
	//			continue
	//		}
	//		if AccessToken != "" && RefreshToken != "" {
	//			clientLoop(ctx, address, AccessToken, RefreshToken)
	//		} else {
	//			fmt.Println("Problem with login")
	//		}
	//	case input == "q":
	//		return
	//	default:
	//		fmt.Printf("Wrong input: %v \n", input)
	//	}
	//}
	cancel()
}

//func clientLoop(ctx context.Context, address string, accessToken string, refreshToken string) {
//	secretClient := grpc_client.NewSecretClient(address, accessToken, refreshToken)
//	secretHandler := handlers.NewSecretHandler(secretClient)
//	for {
//		fmt.Println("To get all secrets input g, to create secret input c, to quit input q:")
//		var input string
//		_, err := fmt.Scan(&input)
//		if err != nil {
//			break
//		}
//		switch {
//		case input == "c":
//			creatingLoop(ctx, secretHandler)
//		case input == "g":
//			result, err := secretHandler.GetSecret(ctx)
//			if err != nil {
//				fmt.Println("Error while get a secret ", err)
//				return
//			}
//			if len(result) == 0 {
//				fmt.Println("You have 0 secrets yet")
//			} else {
//				for _, r := range result {
//					fmt.Println("--------------------------")
//					fmt.Printf("Secret ID: %v \nSecret type: %v \nSecret Meta: %v \n", r.ID, r.Type, r.MetaData)
//					fmt.Println("Data:")
//					for k, v := range r.UsefulData {
//						fmt.Println(" ", k, ": ", v)
//					}
//				}
//			}
//			fmt.Println()
//		case input == "q":
//			return
//		default:
//			fmt.Printf("Wrong input: %v \n", input)
//		}
//
//	}
//}
//
//func creatingLoop(ctx context.Context, secretHandler *handlers.SecretHandler) {
//	for {
//		fmt.Println("Select type of secret:")
//		fmt.Println("Enter lp for login password, c for credit card, b for bytes, s for string")
//		fmt.Println("Press q for go back")
//		var input string
//		_, err := fmt.Scan(&input)
//		if err != nil {
//			break
//		}
//		switch {
//		case input == "lp":
//			loginPasswordLoop(ctx, secretHandler)
//		case input == "s":
//			stringLoop(ctx, secretHandler)
//		case input == "b":
//			binaryLoop(ctx, secretHandler)
//		case input == "c":
//			creditCardLoop(ctx, secretHandler)
//		case input == "q":
//			return
//		default:
//			fmt.Printf("Wrong input")
//		}
//	}
//}
//
//func enterDataLoop(data []string) (map[string]string, error) {
//	result := map[string]string{}
//	for _, key := range data {
//		var value string
//		fmt.Printf("Enter %v\n", key)
//		_, err := fmt.Scanln(&value)
//		if err != nil {
//			return nil, err
//		}
//		result[key] = value
//	}
//	return result, nil
//}
//
//func loginPasswordLoop(ctx context.Context, secretHandler *handlers.SecretHandler) {
//	valuesMap, err := enterDataLoop([]string{"metadata", "login", "password"})
//	if err != nil {
//		fmt.Printf("Wrong input")
//		return
//	}
//	err = secretHandler.CreateSecret(ctx, models.Secret{
//		Type:     "login_password",
//		MetaData: valuesMap["metadata"],
//		UsefulData: map[string]string{
//			"login":    valuesMap["login"],
//			"password": valuesMap["password"],
//		},
//	})
//	if err != nil {
//		fmt.Println("Something went wrong ", err)
//		return
//	}
//	fmt.Println("Successfully created")
//}
//
//func stringLoop(ctx context.Context, secretHandler *handlers.SecretHandler) {
//	valuesMap, err := enterDataLoop([]string{"metadata", "string"})
//	if err != nil {
//		fmt.Printf("Wrong input")
//		return
//	}
//	err = secretHandler.CreateSecret(ctx, models.Secret{
//		Type:     "string",
//		MetaData: valuesMap["metadata"],
//		UsefulData: map[string]string{
//			"string": valuesMap["string"],
//		},
//	})
//	if err != nil {
//		fmt.Println("Something went wrong ", err)
//		return
//	}
//	fmt.Println("Successfully created")
//}
//
//func binaryLoop(ctx context.Context, secretHandler *handlers.SecretHandler) {
//	valuesMap, err := enterDataLoop([]string{"metadata", "binary"})
//	if err != nil {
//		fmt.Printf("Wrong input")
//		return
//	}
//	err = secretHandler.CreateSecret(ctx, models.Secret{
//		Type:     "binary",
//		MetaData: valuesMap["metadata"],
//		UsefulData: map[string]string{
//			"binary": valuesMap["binary"],
//		},
//	})
//	if err != nil {
//		fmt.Println("Something went wrong ", err)
//		return
//	}
//	fmt.Println("Successfully created")
//}
//
//func creditCardLoop(ctx context.Context, secretHandler *handlers.SecretHandler) {
//	valuesMap, err := enterDataLoop([]string{"metadata", "card_number", "expired_date", "owner", "CVV"})
//	if err != nil {
//		fmt.Printf("Wrong input")
//		return
//	}
//	err = secretHandler.CreateSecret(ctx, models.Secret{
//		Type:     "credit_card",
//		MetaData: valuesMap["metadata"],
//		UsefulData: map[string]string{
//			"card_number":  valuesMap["card_number"],
//			"expired_date": valuesMap["expired_date"],
//			"owner":        valuesMap["owner"],
//			"CVV":          valuesMap["CVV"],
//		},
//	})
//	if err != nil {
//		fmt.Println("Something went wrong ", err)
//		return
//	}
//	fmt.Println("Successfully created")
//}
