package source

import (
	"context"
)

type Service interface {
	GetReferenceSources(ctx context.Context) ([]SumberCustResponse, error)
	GetApplicationSources(ctx context.Context) ([]ApplicationSourceResponse, error)
	GetApplicationTypes(ctx context.Context) ([]ApplicationTypeResponse, error)
	GetSoaOptions(ctx context.Context) (SoaOptionsResponse, error)
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
			ID:        it.ID,
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
			ID:   it.ID,
			Code: it.Code,
			Name: it.Name,
		})
	}
	return res, nil
}

func (s *service) GetApplicationTypes(ctx context.Context) ([]ApplicationTypeResponse, error) {
	items, err := s.repo.FindApplicationTypes(ctx)
	if err != nil {
		return []ApplicationTypeResponse{}, err
	}

	var res []ApplicationTypeResponse
	for _, it := range items {
		res = append(res, ApplicationTypeResponse{
			ID:   it.ID,
			Code: it.Code,
			Name: it.Name,
		})
	}
	return res, nil
}

func (s *service) GetSoaOptions(ctx context.Context) (SoaOptionsResponse, error) {
	appSources, _ := s.GetApplicationSources(ctx)
	refSources, _ := s.GetReferenceSources(ctx)
	appTypes, _ := s.GetApplicationTypes(ctx)

	return SoaOptionsResponse{
		SOA:           appSources,
		Aktifitas:     refSources,
		JenisAplikasi: appTypes,
	}, nil
}
