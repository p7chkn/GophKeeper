// Package loops пакет для хранения циклов взаимодействия
package loops

import (
	"context"
	"fmt"
	grpc_client "new_diplom_client/grpc-client"
	"new_diplom_client/handlers"
	"new_diplom_client/models"
)

var (
	accessToken  = ""
	refreshToken = ""
)

// NewUserLoop функция создания структуры цикла взаимодейтсвия с пользователями
func NewUserLoop(address string, userHandler *handlers.UserHandler) *UserLoop {
	return &UserLoop{
		address:     address,
		userHandler: userHandler,
	}
}

// UserLoop структра цикла взаимодействия с пользователем
type UserLoop struct {
	address       string
	userHandler   *handlers.UserHandler
	secretHandler *handlers.SecretHandler
}

// MainLoop функция запуска основного цикла
func (ul *UserLoop) MainLoop(ctx context.Context) {
	for {
		fmt.Println("To login input l, to register input r, to quit input q:")
		var input string
		_, err := fmt.Scan(&input)
		if err != nil {
			break
		}
		switch {
		case input == "r":
			accessToken, refreshToken, err = ul.userHandler.RegisterUser(ctx)
			if err != nil {
				fmt.Printf("Error with registration: %v \n", err)
				continue
			}
			if accessToken != "" && refreshToken != "" {
				ul.clientLoop(ctx, accessToken, refreshToken)
			} else {
				fmt.Println("Problem with register")
			}
		case input == "l":
			accessToken, refreshToken, err = ul.userHandler.AuthUser(ctx)
			if err != nil {
				fmt.Println("Wrong login or password")
				continue
			}
			if accessToken != "" && refreshToken != "" {
				ul.clientLoop(ctx, accessToken, refreshToken)
			} else {
				fmt.Println("Problem with login")
			}
		case input == "q":
			return
		default:
			fmt.Printf("Wrong input: %v \n", input)
		}
	}
}

func (ul *UserLoop) clientLoop(ctx context.Context, accessToken string, refreshToken string) {
	secretClient := grpc_client.NewSecretClient(ul.address, accessToken, refreshToken,
		ul.userHandler.UserClient)
	secretHandler := handlers.NewSecretHandler(secretClient)
	for {
		fmt.Println("To get all secrets input g, to create secret input c, to quit input q:")
		var input string
		_, err := fmt.Scan(&input)
		if err != nil {
			break
		}
		switch {
		case input == "c":
			ul.creatingLoop(ctx, secretHandler)
		case input == "g":
			result, err := secretHandler.GetSecret(ctx)
			if err != nil {
				fmt.Println("Error while get a secret ", err)
				return
			}
			if len(result) == 0 {
				fmt.Println("You have 0 secrets yet")
			} else {
				for _, r := range result {
					fmt.Println("--------------------------")
					fmt.Printf("Secret ID: %v \nSecret type: %v \nSecret Meta: %v \n", r.ID, r.Type, r.MetaData)
					fmt.Println("Data:")
					for k, v := range r.UsefulData {
						fmt.Println(" ", k, ": ", v)
					}
				}
			}
			fmt.Println()
		case input == "q":
			return
		default:
			fmt.Printf("Wrong input: %v \n", input)
		}

	}
}

func (ul *UserLoop) creatingLoop(ctx context.Context, secretHandler *handlers.SecretHandler) {
	for {
		fmt.Println("Select type of secret:")
		fmt.Println("Enter lp for login password, c for credit card, b for bytes, s for string")
		fmt.Println("Press q for go back")
		var input string
		_, err := fmt.Scan(&input)
		if err != nil {
			break
		}
		switch {
		case input == "lp":
			ul.loginPasswordLoop(ctx, secretHandler)
		case input == "s":
			ul.stringLoop(ctx, secretHandler)
		case input == "b":
			ul.binaryLoop(ctx, secretHandler)
		case input == "c":
			ul.creditCardLoop(ctx, secretHandler)
		case input == "q":
			return
		default:
			fmt.Printf("Wrong input")
		}
	}
}

func (ul *UserLoop) enterDataLoop(data []string) (map[string]string, error) {
	result := map[string]string{}
	for _, key := range data {
		var value string
		fmt.Printf("Enter %v\n", key)
		_, err := fmt.Scanln(&value)
		if err != nil {
			return nil, err
		}
		result[key] = value
	}
	return result, nil
}

func (ul *UserLoop) loginPasswordLoop(ctx context.Context, secretHandler *handlers.SecretHandler) {
	valuesMap, err := ul.enterDataLoop([]string{"metadata", "login", "password"})
	if err != nil {
		fmt.Printf("Wrong input")
		return
	}
	err = secretHandler.CreateSecret(ctx, models.Secret{
		Type:     "login_password",
		MetaData: valuesMap["metadata"],
		UsefulData: map[string]string{
			"login":    valuesMap["login"],
			"password": valuesMap["password"],
		},
	})
	if err != nil {
		fmt.Println("Something went wrong ", err)
		return
	}
	fmt.Println("Successfully created")
}

func (ul *UserLoop) stringLoop(ctx context.Context, secretHandler *handlers.SecretHandler) {
	valuesMap, err := ul.enterDataLoop([]string{"metadata", "string"})
	if err != nil {
		fmt.Printf("Wrong input")
		return
	}
	err = secretHandler.CreateSecret(ctx, models.Secret{
		Type:     "string",
		MetaData: valuesMap["metadata"],
		UsefulData: map[string]string{
			"string": valuesMap["string"],
		},
	})
	if err != nil {
		fmt.Println("Something went wrong ", err)
		return
	}
	fmt.Println("Successfully created")
}

func (ul *UserLoop) binaryLoop(ctx context.Context, secretHandler *handlers.SecretHandler) {
	valuesMap, err := ul.enterDataLoop([]string{"metadata", "binary"})
	if err != nil {
		fmt.Printf("Wrong input")
		return
	}
	err = secretHandler.CreateSecret(ctx, models.Secret{
		Type:     "binary",
		MetaData: valuesMap["metadata"],
		UsefulData: map[string]string{
			"binary": valuesMap["binary"],
		},
	})
	if err != nil {
		fmt.Println("Something went wrong ", err)
		return
	}
	fmt.Println("Successfully created")
}

func (ul *UserLoop) creditCardLoop(ctx context.Context, secretHandler *handlers.SecretHandler) {
	valuesMap, err := ul.enterDataLoop([]string{"metadata", "card_number", "expired_date", "owner", "CVV"})
	if err != nil {
		fmt.Printf("Wrong input")
		return
	}
	err = secretHandler.CreateSecret(ctx, models.Secret{
		Type:     "credit_card",
		MetaData: valuesMap["metadata"],
		UsefulData: map[string]string{
			"card_number":  valuesMap["card_number"],
			"expired_date": valuesMap["expired_date"],
			"owner":        valuesMap["owner"],
			"CVV":          valuesMap["CVV"],
		},
	})
	if err != nil {
		fmt.Println("Something went wrong ", err)
		return
	}
	fmt.Println("Successfully created")
}
