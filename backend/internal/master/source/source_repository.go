package source

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

type Repository interface {
	FindReferenceSources(ctx context.Context) ([]ReferenceSource, error)
	FindApplicationSources(ctx context.Context) ([]ApplicationSource, error)
}

type repository struct {
	dbPG *gorm.DB
}

func NewRepository(dbPG *gorm.DB) Repository {
	return &repository{dbPG: dbPG}
}

func (r *repository) FindReferenceSources(ctx context.Context) ([]ReferenceSource, error) {
	if r.dbPG == nil {
		return nil, fmt.Errorf("postgres database is not initialized")
	}

	var sources []ReferenceSource
	err := r.dbPG.WithContext(ctx).
		Table("master.reference_sources").
		Select("code, name, is_active").
		Where("is_active = ?", true).
		Order("name ASC").
		Find(&sources).Error

	return sources, err
}

func (r *repository) FindApplicationSources(ctx context.Context) ([]ApplicationSource, error) {
	if r.dbPG == nil {
		return nil, fmt.Errorf("postgres database is not initialized")
	}

	var sources []ApplicationSource
	err := r.dbPG.WithContext(ctx).
		Table("master.application_sources").
		Select("code, name, is_active").
		Where("is_active = ?", true).
		Order("name ASC").
		Find(&sources).Error

	return sources, err
}
