// Package models пакет для хранения моделей
package models

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/json"
	"errors"
	"fmt"
	customerrors "new_diplom/errors"
	"strconv"
	"strings"
	"unicode"
)

var (
	key   = []byte{240, 43, 127, 3, 22, 181, 93, 105, 162, 19, 180, 125, 207, 77, 209, 70}
	nonce = []byte{161, 154, 38, 17, 9, 137, 119, 105, 204, 99, 67, 14}
)

// NewRawSecretData функия для создания нового объекта зашифрованного секрета
func NewRawSecretData(secret Secret) (*RawSecretData, error) {
	err := secret.Data.Validate()
	if err != nil {
		return nil, err
	}
	allData := map[string]string{}
	allData["meta_data"] = secret.Data.MetaData
	allData["type"] = secret.Data.Type
	useFullData := map[string]string{}
	for k, v := range secret.Data.UsefulData {
		useFullData[k] = v
	}
	marshalUseFullData, err := json.Marshal(useFullData)
	if err != nil {
		return nil, customerrors.NewCustomError(err, "can't marshal data")
	}
	allData["useful_data"] = string(marshalUseFullData)
	data, err := json.Marshal(allData)
	if err != nil {
		return nil, customerrors.NewCustomError(err, "can't marshal data")
	}
	rsd := RawSecretData{
		UserID: secret.User,
		Data:   data,
	}
	err = rsd.Encrypt()
	return &rsd, err
}

// RawSecretData структура зашифрованного секрета
type RawSecretData struct {
	ID     string `db:"id"`
	UserID string `db:"user_id"`
	Data   []byte `db:"secret_data"`
}

// DecryptToSecretData функция расшифровки секрета в обычную структуру секрета
func (rd *RawSecretData) DecryptToSecretData() (*SecretData, error) {
	aesblock, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return nil, err
	}
	decryptData, err := aesgcm.Open(nil, nonce, rd.Data, nil)
	if err != nil {
		return nil, err
	}
	rd.Data = decryptData
	return rd.turnToSecretData()
}

func (rd *RawSecretData) turnToSecretData() (*SecretData, error) {

	result := SecretData{}
	err := json.Unmarshal(rd.Data, &result)
	if err != nil {
		return nil, customerrors.NewCustomError(err, "problem with data format")
	}
	result.ID = rd.ID
	return &result, err
}

// Encrypt функция зашифровки секрета
func (rd *RawSecretData) Encrypt() error {
	aesblock, err := aes.NewCipher(key)
	if err != nil {
		return customerrors.NewCustomError(err, "error with encrypt")
	}
	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return customerrors.NewCustomError(err, "error with encrypt")
	}
	encryptData := aesgcm.Seal(nil, nonce, rd.Data, nil)
	rd.Data = encryptData
	return nil
}

// Secret структура хранения секрета
type Secret struct {
	User string     `db:"user_id"`
	Data SecretData `db:"secret_data"`
}

// SecretData структура для хранения данных секрета
type SecretData struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	MetaData   string `json:"meta_data"`
	UsefulData map[string]string
}

// UnmarshalJSON функция десереализации данных секрета
func (sd *SecretData) UnmarshalJSON(data []byte) error {
	type SecretDataAlias SecretData

	aliasValue := &struct {
		*SecretDataAlias
		TempData string `json:"useful_data"`
	}{
		SecretDataAlias: (*SecretDataAlias)(sd),
	}
	if err := json.Unmarshal(data, aliasValue); err != nil {
		return err
	}
	if err := json.Unmarshal([]byte(aliasValue.TempData), &aliasValue.UsefulData); err != nil {
		return err
	}
	return nil
}

// Validate функия валидации данных секрета
func (sd *SecretData) Validate() error {
	switch sd.Type {
	case "binary":
		return sd.validateBinary()
	case "login_password":
		return sd.validateLoginPassword()
	case "credit_card":
		return sd.validateCreditCard()
	case "string":
		return sd.validateString()
	default:
		return customerrors.NewCustomError(errors.New("wrong type of secret"), "wrong type")
	}
}

func (sd *SecretData) validateBinary() error {
	return sd.checkUsefulData([]string{"binary"})
}

func (sd *SecretData) validateLoginPassword() error {
	return sd.checkUsefulData([]string{"login", "password"})
}

func (sd *SecretData) validateCreditCard() error {
	err := sd.checkUsefulData([]string{"card_number", "expired_date",
		"owner", "CVV"})
	if err != nil {
		return err
	}
	err = sd.checkFiledIsNumber([]string{"card_number", "CVV"})
	if err != nil {
		return err
	}
	cardNumber, _ := strconv.Atoi(sd.UsefulData["card_number"])
	if !validLuhnNumber(cardNumber) {
		return customerrors.NewCustomError(
			errors.New("wrong credit card number"),
			"wrong credit card number")
	}
	return nil
}

func (sd *SecretData) validateString() error {
	return sd.checkUsefulData([]string{"string"})
}

func (sd *SecretData) checkUsefulData(fields []string) error {
	var missingFields []string
	for _, field := range fields {
		if _, ok := sd.UsefulData[field]; !ok {
			missingFields = append(missingFields, field)
		}
	}
	if len(missingFields) != 0 {
		text := fmt.Sprintf("wrong format of data, missing fields %v",
			strings.Trim(fmt.Sprint(missingFields), "[]"))
		return customerrors.NewCustomError(errors.New(text), text)
	}
	return nil
}

func (sd *SecretData) checkFiledIsNumber(fields []string) error {
	var errorFields []string
	for _, field := range fields {
		if !isInt(field) {
			errorFields = append(errorFields, field)
		}
	}
	if len(errorFields) != 0 {
		text := fmt.Sprintf("this field must consist of numbers %v",
			strings.Trim(fmt.Sprint(errorFields), "[]"))
		return customerrors.NewCustomError(errors.New(text), text)
	}
	return nil
}

func isInt(s string) bool {
	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}

func validLuhnNumber(number int) bool {
	return (number%10+checksum(number/10))%10 == 0
}

func checksum(number int) int {
	var luhn int

	for i := 0; number > 0; i++ {
		cur := number % 10

		if i%2 == 0 {
			cur = cur * 2
			if cur > 9 {
				cur = cur%10 + cur/10
			}
		}

		luhn += cur
		number = number / 10
	}
	return luhn % 10
}
