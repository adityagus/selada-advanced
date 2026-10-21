package location

import (
	"context"
)

type Service interface {
	SearchKelurahan(ctx context.Context, query string) ([]KelurahanResponse, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) SearchKelurahan(ctx context.Context, query string) ([]KelurahanResponse, error) {
	items, err := s.repo.SearchVillage(ctx, query)
	if err != nil {
		return []KelurahanResponse{}, err
	}

	var res []KelurahanResponse
	for _, it := range items {
		res = append(res, KelurahanResponse{
			KdKelurahan: it.KdKelurahan,
			NmKelurahan: it.NmKelurahan,
			KdKecamatan: it.KdKecamatan,
			NmKecamatan: it.NmKecamatan,
			KdKota:      it.KdKota,
			NmKota:      it.NmKota,
		})
	}
	return res, nil
}
