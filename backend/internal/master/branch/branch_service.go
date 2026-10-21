package branch

import (
	"context"
)

type Service interface {
	GetBranches(ctx context.Context) ([]BranchResponse, error)
	ResolveCurrentBranch(ctx context.Context, branchCode, username, cachedBranchName string) CurrentBranchResponse
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetBranches(ctx context.Context) ([]BranchResponse, error) {
	branches, err := s.repo.FindAll(ctx)
	if err != nil {
		return []BranchResponse{}, err
	}

	var res []BranchResponse
	for _, b := range branches {
		res = append(res, BranchResponse{
			IDCabang:   b.BranchID,
			KodeCabang: b.BranchCode,
			NmCabang:   b.BranchName,
		})
	}
	return res, nil
}

func (s *service) ResolveCurrentBranch(ctx context.Context, branchCode, username, cachedBranchName string) CurrentBranchResponse {
	finalName := cachedBranchName

	if finalName == "" {
		name, err := s.repo.FindNameByCodeOrUser(ctx, branchCode, username)
		if err == nil && name != "" {
			finalName = name
		} else {
			finalName = branchCode
		}
	}

	return CurrentBranchResponse{
		Cabang:     branchCode,
		BranchName: finalName,
	}
}
