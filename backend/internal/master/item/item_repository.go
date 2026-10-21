package item

import (
	"context"
	"fmt"
	"math"
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	FindKatalogList(ctx context.Context, categorySlug string) ([]KatalogItemResponse, error)
	FindKatalogSpec(ctx context.Context, motorID string) (*MotorcycleSpecResponse, error)
	CalculateSalesPoint(ctx context.Context, userID string) (float64, error)
}

type repository struct {
	dbPOS    *gorm.DB
	dbSelada *gorm.DB
}

func NewRepository(dbPOS *gorm.DB, dbSelada *gorm.DB) Repository {
	return &repository{
		dbPOS:    dbPOS,
		dbSelada: dbSelada,
	}
}

func (r *repository) FindKatalogList(ctx context.Context, categorySlug string) ([]KatalogItemResponse, error) {
	if r.dbPOS == nil {
		return []KatalogItemResponse{}, fmt.Errorf("POS database is not initialized")
	}

	var categoryID int
	switch categorySlug {
	case "sport":
		categoryID = 10201
	case "trail":
		categoryID = 10202
	case "moped":
		categoryID = 10203
	default:
		return []KatalogItemResponse{}, fmt.Errorf("kategori tidak valid")
	}

	var items []KatalogItemResponse
	err := r.dbPOS.WithContext(ctx).Table("motorcycle m").
		Select("m.id_motor, m.nama_motor, mc.gambar_color").
		Joins("LEFT JOIN motcol mc ON mc.id_motor = m.id_motor AND mc.default_color = 'Y'").
		Where("m.id_category = ?", categoryID).
		Scan(&items).Error

	return items, err
}

func (r *repository) FindKatalogSpec(ctx context.Context, motorID string) (*MotorcycleSpecResponse, error) {
	if r.dbPOS == nil {
		return nil, fmt.Errorf("POS database is not initialized")
	}

	var motor Motorcycle
	if err := r.dbPOS.WithContext(ctx).Where("id_motor = ?", motorID).First(&motor).Error; err != nil {
		return nil, err
	}

	var fitur MotFitur
	_ = r.dbPOS.WithContext(ctx).Where("id_motor = ?", motorID).First(&fitur).Error

	var fiturBrosur FiturMotor
	_ = r.dbPOS.WithContext(ctx).Where("id_motor = ?", motorID).First(&fiturBrosur).Error

	var colors []ColorSpec
	_ = r.dbPOS.WithContext(ctx).Table("motcol mc").
		Select("m.color_name, mc.gambar_color").
		Joins("LEFT JOIN master_color m ON m.id_master = mc.id_icon").
		Where("mc.id_motor = ?", motorID).
		Scan(&colors).Error

	var fiturList []string
	_ = r.dbPOS.WithContext(ctx).Table("motfitur_list").
		Where("id_motor = ?", motorID).
		Pluck("fitur", &fiturList).Error

	return &MotorcycleSpecResponse{
		IDMotor:   motor.IDMotor,
		NamaMotor: motor.NamaMotor,
		Harga:     fitur.Harga,
		Detail:    fitur.Detail,
		Brosur:    fiturBrosur.Brosur,
		Colors:    colors,
		Fiturs:    fiturList,
	}, nil
}

func (r *repository) CalculateSalesPoint(ctx context.Context, userID string) (float64, error) {
	if r.dbSelada == nil {
		return 0, fmt.Errorf("selada database is not initialized")
	}

	tgl := time.Now()
	limitTgl, _ := time.Parse("2006-01-02", "2021-12-31")
	if tgl.After(limitTgl) {
		tgl = limitTgl
	}

	n := int(tgl.Month())
	thn := tgl.Year()

	var j int = 0
	var tpoin float64 = 0.0

	type CronDetPointFn struct {
		Bln       int     `gorm:"column:bln"`
		Thn       int     `gorm:"column:thn"`
		IDSales   string  `gorm:"column:id_sales"`
		Pleads    float64 `gorm:"column:pleads"`
		Pfu       float64 `gorm:"column:pfu"`
		Pdo       float64 `gorm:"column:pdo"`
		Pakurdat  float64 `gorm:"column:pakurdat"`
		Leadstodo float64 `gorm:"column:leadstodo"`
		Leadspal  float64 `gorm:"column:leadspal"`
	}

	for i := 3; i <= n; i++ {
		var q []CronDetPointFn
		err := r.dbSelada.WithContext(ctx).Table("cron_det_point_fn").
			Where("bln = ? AND thn = ? AND id_sales = ?", i, thn, userID).
			Find(&q).Error
		if err == nil {
			for _, row := range q {
				mpoin := row.Pleads + row.Pfu + row.Pdo + row.Pakurdat + row.Leadstodo
				if row.Leadspal > 0 {
					tpoin += mpoin / 2.0
				} else {
					tpoin += mpoin
				}
				j++
			}
		}
	}

	poin := 0.0
	if j > 0 {
		poin = math.Round(tpoin / float64(j))
	}

	return poin, nil
}
