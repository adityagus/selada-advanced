package sales

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	GetActiveSales(ctx context.Context) ([]SalesUser, error)
}

type repository struct {
	dbPG  *gorm.DB
	dbPOS *gorm.DB
}

func NewRepository(dbPG *gorm.DB, dbPOS *gorm.DB) Repository {
	return &repository{
		dbPG:  dbPG,
		dbPOS: dbPOS,
	}
}

func (r *repository) fetchMonthlyOmzetMap(ctx context.Context) map[string]float64 {
	userOmzetMap := make(map[string]float64)
	if r.dbPG == nil {
		return userOmzetMap
	}

	loc := time.FixedZone("WIB", 7*3600)
	now := time.Now().In(loc)
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, loc)
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Nanosecond)

	type OmzetRow struct {
		SalesEmployeeID int     `gorm:"column:sales_employee_id"`
		TotalOmzet      float64 `gorm:"column:total_omzet"`
	}

	var rows []OmzetRow
	_ = r.dbPG.WithContext(ctx).Table("transaction.pawn_transaction").
		Select("sales_employee_id, SUM(up_amount) as total_omzet").
		Where("status = ?", "ACTIVE").
		Where("pawn_date >= ? AND pawn_date <= ?", startOfMonth, endOfMonth).
		Group("sales_employee_id").
		Scan(&rows).Error

	empToUsers := make(map[int][]string)
	type EmpUser struct {
		EmployeeID int    `gorm:"column:employee_id"`
		Username   string `gorm:"column:username"`
	}
	var empUsers []EmpUser
	_ = r.dbPG.WithContext(ctx).Table("auth.users").
		Select("employee_id, username").
		Where("employee_id IS NOT NULL").
		Scan(&empUsers).Error

	for _, eu := range empUsers {
		u := strings.ToLower(strings.TrimSpace(eu.Username))
		if u != "" {
			empToUsers[eu.EmployeeID] = append(empToUsers[eu.EmployeeID], u)
		}
	}

	for _, row := range rows {
		if row.SalesEmployeeID <= 0 || row.TotalOmzet <= 0 {
			continue
		}
		idents := empToUsers[row.SalesEmployeeID]
		seen := make(map[string]bool)
		for _, ident := range idents {
			if ident != "" && !seen[ident] {
				seen[ident] = true
				userOmzetMap[ident] += row.TotalOmzet
			}
		}
	}

	return userOmzetMap
}

func (r *repository) GetActiveSales(ctx context.Context) ([]SalesUser, error) {
	if r.dbPOS == nil {
		return []SalesUser{}, fmt.Errorf("POS database is not initialized")
	}

	// 1. Fetch active usernames in Sopiga Pro
	activeSopigaUserMap := make(map[string]bool)
	if r.dbPG != nil {
		type PgUser struct {
			Username string `gorm:"column:username"`
		}
		var pgUsers []PgUser
		_ = r.dbPG.WithContext(ctx).Table("auth.users u").
			Select("u.username").
			Joins("JOIN master.employee e ON e.employee_id = u.employee_id").
			Where("COALESCE(e.is_active, true) = ?", true).
			Scan(&pgUsers).Error

		for _, p := range pgUsers {
			if p.Username != "" {
				activeSopigaUserMap[strings.ToLower(strings.TrimSpace(p.Username))] = true
			}
		}
	}

	// 2. Fetch omzet map for the current month
	omzetMap := r.fetchMonthlyOmzetMap(ctx)

	// 3. Fetch active users in POS
	var posRows []SalesUser
	query := r.dbPOS.WithContext(ctx).Table("tbl_user u").
		Select("u.username, u.nama, u.id_jabatan, COALESCE(j.nm_jabatan, 'Sales') AS nm_jabatan").
		Joins("LEFT JOIN mst_jabatan j ON j.id_jabatan = u.id_jabatan").
		Where("u.void = ?", 1).
		Where("u.id_jabatan IN (6, 17, 29, 74) OR LOWER(j.nm_jabatan) LIKE '%sales%' OR LOWER(j.nm_jabatan) LIKE '%trainee%' OR LOWER(j.nm_jabatan) LIKE '%junior%' OR LOWER(j.nm_jabatan) LIKE '%senior%'")

	_ = query.Scan(&posRows).Error

	var salesList []SalesUser
	for _, s := range posRows {
		uLower := strings.ToLower(strings.TrimSpace(s.Username))
		nLower := strings.ToLower(strings.TrimSpace(s.Nama))

		if len(activeSopigaUserMap) > 0 && !activeSopigaUserMap[uLower] {
			continue
		}
		s.FKUser = s.Username

		if val, exists := omzetMap[uLower]; exists {
			s.TotalOmzet = val
		} else if val, exists := omzetMap[nLower]; exists {
			s.TotalOmzet = val
		}

		salesList = append(salesList, s)
	}

	// 4. Sort sales list by highest omzet first (descending), then alphabetically by nama
	sort.SliceStable(salesList, func(i, j int) bool {
		if salesList[i].TotalOmzet != salesList[j].TotalOmzet {
			return salesList[i].TotalOmzet > salesList[j].TotalOmzet
		}
		return strings.ToLower(salesList[i].Nama) < strings.ToLower(salesList[j].Nama)
	})

	return salesList, nil
}
