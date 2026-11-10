package source

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

type Repository interface {
	FindReferenceSources(ctx context.Context) ([]ReferenceSource, error)
	FindApplicationSources(ctx context.Context) ([]ApplicationSource, error)
	FindApplicationTypes(ctx context.Context) ([]ApplicationType, error)
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
		Select("id, code, name, is_active").
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
		Select("id, code, name, is_active").
		Where("is_active = ?", true).
		Order("name ASC").
		Find(&sources).Error

	return sources, err
}

func (r *repository) FindApplicationTypes(ctx context.Context) ([]ApplicationType, error) {
	if r.dbPG == nil {
		return nil, fmt.Errorf("postgres database is not initialized")
	}

	var types []ApplicationType
	err := r.dbPG.WithContext(ctx).
		Table("master.application_types").
		Select("id, code, name, is_active").
		Where("is_active = ?", true).
		Order("name ASC").
		Find(&types).Error

	return types, err
}
