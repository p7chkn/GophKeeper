package models

import "new_diplom_client/pb"

type Secret struct {
	ID         string            `json:"id"`
	Type       string            `json:"type"`
	MetaData   string            `json:"meta_data"`
	UsefulData map[string]string `json:"useful_data"`
}

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
