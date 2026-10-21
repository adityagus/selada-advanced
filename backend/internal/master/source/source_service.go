package source

import (
	"context"
)

type Service interface {
	GetReferenceSources(ctx context.Context) ([]SumberCustResponse, error)
	GetApplicationSources(ctx context.Context) ([]ApplicationSourceResponse, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetReferenceSources(ctx context.Context) ([]SumberCustResponse, error) {
	items, err := s.repo.FindReferenceSources(ctx)
	if err != nil {
		return []SumberCustResponse{}, err
	}

	var res []SumberCustResponse
	for _, it := range items {
		res = append(res, SumberCustResponse{
			KdRefCust: it.Code,
			NmRefCust: it.Name,
		})
	}
	return res, nil
}

func (s *service) GetApplicationSources(ctx context.Context) ([]ApplicationSourceResponse, error) {
	items, err := s.repo.FindApplicationSources(ctx)
	if err != nil {
		return []ApplicationSourceResponse{}, err
	}

	var res []ApplicationSourceResponse
	for _, it := range items {
		res = append(res, ApplicationSourceResponse{
			Code: it.Code,
			Name: it.Name,
		})
	}
	return res, nil
}
