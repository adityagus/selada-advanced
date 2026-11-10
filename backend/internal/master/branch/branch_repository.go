package branch

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

type Repository interface {
	FindAll(ctx context.Context) ([]Branch, error)
	FindNameByCodeOrUser(ctx context.Context, branchCode, username string) (string, error)
}

type repository struct {
	dbPG  *gorm.DB
	dbPOS *gorm.DB
}

func NewRepository(dbPG *gorm.DB, dbPOS *gorm.DB) Repository {
	return &repository{dbPG: dbPG, dbPOS: dbPOS}
}

func (r *repository) FindAll(ctx context.Context) ([]Branch, error) {
	if r.dbPG == nil {
		return nil, fmt.Errorf("postgres database is not initialized")
	}

	var branches []Branch
	err := r.dbPG.WithContext(ctx).
		Table("master.branch").
		Select("branch_id, branch_code, branch_name, is_active").
		Where("is_active = ?", true).
		Order("branch_name ASC").
		Find(&branches).Error

	return branches, err
}

func (r *repository) FindNameByCodeOrUser(ctx context.Context, branchCode, username string) (string, error) {
	var branchName string

	// 1. Coba lookup via MySQL POS (mst_cabang) jika branchCode dari POS
	if branchCode != "" && r.dbPOS != nil {
		_ = r.dbPOS.WithContext(ctx).Table("mst_cabang").
			Select("nm_cabang").
			Where("id_cabang = ? OR nick_cabang = ?", branchCode, branchCode).
			Limit(1).
			Pluck("nm_cabang", &branchName).Error

		if branchName != "" {
			return branchName, nil
		}
	}

	if r.dbPG == nil {
		return "", fmt.Errorf("postgres database is not initialized")
	}

	// 2. Coba lookup via auth.users -> employee -> branch di Postgres
	if username != "" {
		err := r.dbPG.WithContext(ctx).
			Table("auth.users u").
			Select("b.branch_name").
			Joins("JOIN master.employee e ON e.employee_id = u.employee_id").
			Joins("JOIN master.branch b ON b.branch_id = e.branch_id").
			Where("LOWER(u.username) = LOWER(?)", username).
			Limit(1).
			Pluck("b.branch_name", &branchName).Error

		if err == nil && branchName != "" {
			return branchName, nil
		}
	}

	// 3. Coba lookup via branch_code atau branch_id di Postgres master.branch
	if branchCode != "" {
		err := r.dbPG.WithContext(ctx).
			Table("master.branch").
			Select("branch_name").
			Where("branch_code = ? OR CAST(branch_id AS text) = ?", branchCode, branchCode).
			Limit(1).
			Pluck("branch_name", &branchName).Error

		if err == nil && branchName != "" {
			return branchName, nil
		}
	}

	return "", fmt.Errorf("branch not found")
}
