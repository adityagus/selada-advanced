package item

import (
	"context"
)

type Service interface {
	GetKatalogList(ctx context.Context, categorySlug string) ([]KatalogItemResponse, error)
	GetKatalogSpec(ctx context.Context, motorID string) (*MotorcycleSpecResponse, error)
	GetSalesPoint(ctx context.Context, userID string) (PointResponse, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetKatalogList(ctx context.Context, categorySlug string) ([]KatalogItemResponse, error) {
	return s.repo.FindKatalogList(ctx, categorySlug)
}

func (s *service) GetKatalogSpec(ctx context.Context, motorID string) (*MotorcycleSpecResponse, error) {
	return s.repo.FindKatalogSpec(ctx, motorID)
}

func (s *service) GetSalesPoint(ctx context.Context, userID string) (PointResponse, error) {
	p, err := s.repo.CalculateSalesPoint(ctx, userID)
	if err != nil {
		return PointResponse{Point: 0}, err
	}
	return PointResponse{Point: p}, nil
}
