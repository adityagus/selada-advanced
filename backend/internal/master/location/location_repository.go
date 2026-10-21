package location

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

type Repository interface {
	SearchVillage(ctx context.Context, query string) ([]VillageResult, error)
}

type repository struct {
	dbPG *gorm.DB
}

func NewRepository(dbPG *gorm.DB) Repository {
	return &repository{dbPG: dbPG}
}

func (r *repository) SearchVillage(ctx context.Context, query string) ([]VillageResult, error) {
	if r.dbPG == nil {
		return nil, fmt.Errorf("postgres database is not initialized")
	}

	var results []VillageResult
	dbQuery := r.dbPG.WithContext(ctx).
		Table("master.village v").
		Select("v.village_id AS kd_kelurahan, v.village_name AS nm_kelurahan, d.district_id AS kd_kecamatan, d.district_name AS nm_kecamatan, c.city_id AS kd_kota, c.city_name AS nm_kota").
		Joins("LEFT JOIN master.district d ON d.district_id = v.district_id").
		Joins("LEFT JOIN master.city c ON c.city_id = d.city_id")

	if query != "" {
		searchTerm := "%" + query + "%"
		dbQuery = dbQuery.Where("v.village_name ILIKE ? OR d.district_name ILIKE ? OR c.city_name ILIKE ?", searchTerm, searchTerm, searchTerm)
	}

	err := dbQuery.Order("v.village_name ASC").Limit(20).Scan(&results).Error
	return results, err
}
