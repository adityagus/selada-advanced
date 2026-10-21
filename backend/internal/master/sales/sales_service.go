package sales

import (
	"context"
)

type Service interface {
	GetSalesList(ctx context.Context) ([]SalesResponse, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetSalesList(ctx context.Context) ([]SalesResponse, error) {
	users, err := s.repo.GetActiveSales(ctx)
	if err != nil {
		return []SalesResponse{}, err
	}

	var res []SalesResponse
	for _, u := range users {
		res = append(res, SalesResponse{
			Username:   u.Username,
			Nama:       u.Nama,
			IDJabatan:  u.IDJabatan,
			NmJabatan:  u.NmJabatan,
			FKUser:     u.FKUser,
			TotalOmzet: u.TotalOmzet,
		})
	}
	return res, nil
}
