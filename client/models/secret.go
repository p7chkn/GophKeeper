// Package models пакет для хранения моделей
package models

import "new_diplom_client/pb"

// Secret структура секрета
type Secret struct {
	ID         string            `json:"id"`
	Type       string            `json:"type"`
	MetaData   string            `json:"meta_data"`
	UsefulData map[string]string `json:"useful_data"`
}

// TransferUsefulData функция преобразования полезных данных в формат для передачи по gRPC
func (s *Secret) TransferUsefulData() []*pb.Data {
	var result []*pb.Data
	for k, v := range s.UsefulData {
		result = append(result, &pb.Data{
			Title: k,
			Value: v,
		})
	}
	return result
}
