package handlers

import (
	"encoding/json"
	"go-api/models"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go-api/config"

	"github.com/gin-gonic/gin"
)

// GetJenisBarangHandler fetches electronics category list from product.collateral_items (product_id = 2)
func GetJenisBarangHandler(c *gin.Context) {
	if config.DBPostgres == nil {
		c.JSON(http.StatusOK, []gin.H{})
		return
	}

	var results []struct {
		KdJenisBarang string `json:"kd_jenis_barang" gorm:"column:kd_jenis_barang"`
		NmJenisBarang string `json:"nm_jenis_barang" gorm:"column:nm_jenis_barang"`
	}

	err := config.DBPostgres.Table("product.collateral_items").
		Select("CAST(collateral_items_id AS TEXT) AS kd_jenis_barang, name AS nm_jenis_barang").
		Where("product_id = ? AND is_active = ?", 2, true).
		Order("name ASC").
		Scan(&results).Error

	if err != nil || len(results) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		log.Println("Terjadi eror saat mengambil data:", err)
		return
		// 	// Fallback without is_active filter
		// 	_ = config.DBPostgres.Table("product.collateral_items").
		// 		Select("CAST(collateral_items_id AS TEXT) AS kd_jenis_barang, name AS nm_jenis_barang").
		// 		Where("product_id = ? AND is_active = ?", 2, true).
		// 		Order("name ASC").
		// 		Scan(&results).Error
	}

	c.JSON(http.StatusOK, results)
}

// GetGradeHandler fetches grade list from master.grades or returns fallback A, B, C, D
func GetGradeHandler(c *gin.Context) {
	if config.DBPostgres == nil {
		c.JSON(http.StatusOK, defaultGrades())
		return
	}

	var results []struct {
		KdGrade string `json:"kd_grade" gorm:"column:kd_grade"`
		NmGrade string `json:"nm_grade" gorm:"column:nm_grade"`
	}

	err := config.DBPostgres.Table("master.grades").
		Select("grade_code AS kd_grade, grade_name AS nm_grade").
		Where("is_active = ?", true).
		Order("grade_code ASC").
		Scan(&results).Error

	if err != nil || len(results) == 0 {
		c.JSON(http.StatusOK, defaultGrades())
		return
	}

	c.JSON(http.StatusOK, results)
}

func defaultGrades() []gin.H {
	return []gin.H{
		{"kd_grade": "A", "nm_grade": "A"},
		{"kd_grade": "B", "nm_grade": "B"},
		{"kd_grade": "C", "nm_grade": "C"},
		{"kd_grade": "D", "nm_grade": "D"},
	}
}

// GetBarangUmumHandler searches electronics (product_id = 2) from product.collateral_item_details
func GetBarangUmumHandler(c *gin.Context) {
	if config.DBPostgres == nil {
		c.JSON(http.StatusOK, []gin.H{})
		return
	}

	jenis := c.Query("jenisBarang")
	nama := c.Query("namaBarang")

	type ItemResult struct {
		IDBarangDetail   int     `json:"id_barang_detail" gorm:"column:id_detail"`
		CollateralItemID int     `json:"collateral_item_id" gorm:"column:collateral_item_id"`
		KdBarang         string  `json:"kd_barang" gorm:"column:kd_barang"`
		NmBarang         string  `json:"nm_barang" gorm:"column:nm_barang"`
		Harga            float64 `json:"harga" gorm:"column:harga"`
		HargaTaksir      float64 `json:"harga_taksir" gorm:"column:harga_taksir"`
	}

	var results []ItemResult
	query := config.DBPostgres.Table("product.collateral_item_details cid").
		Select("cid.collateral_item_detail_id as id_detail, cid.collateral_item_id, cid.code AS kd_barang, cid.name AS nm_barang, MAX(COALESCE(cip.price, 0)) AS harga, MAX(COALESCE(cip.price, 0)) AS harga_taksir").
		Joins("JOIN product.collateral_items ci ON ci.collateral_items_id = cid.collateral_item_id").
		Joins("LEFT JOIN product.collateral_item_prices cip ON cip.collateral_item_detail_id = cid.collateral_item_detail_id AND cip.is_active = true").
		Where("ci.product_id = ? AND cid.is_active = ?", 2, true).
		Group("cid.collateral_item_detail_id, cid.collateral_item_id, cid.code, cid.name")

	if jenis != "" && jenis != "- Pilih -" && jenis != "- Pilih Jenis -" {
		if jInt, err := strconv.Atoi(jenis); err == nil {
			query = query.Where("cid.collateral_item_id = ?", jInt)
		} else {
			query = query.Where("ci.code = ?", jenis)
		}
	}

	if nama != "" {
		upperNama := "%" + strings.ToUpper(nama) + "%"
		query = query.Where("(UPPER(cid.name) LIKE ? OR UPPER(cid.code) LIKE ?)", upperNama, upperNama)
	}

	err := query.Limit(20).Scan(&results).Error
	if err != nil || len(results) == 0 {
		// Fallback without price join
		fbQuery := config.DBPostgres.Table("product.collateral_item_details cid").
			Select("cid.collateral_item_detail_id as id_detail, cid.collateral_item_id, cid.code AS kd_barang, cid.name AS nm_barang, 0 AS harga, 0 AS harga_taksir").
			Joins("JOIN product.collateral_items ci ON ci.collateral_items_id = cid.collateral_item_id").
			Where("ci.product_id = ?", 2)

		if jenis != "" && jenis != "- Pilih -" && jenis != "- Pilih Jenis -" {
			if jInt, err := strconv.Atoi(jenis); err == nil {
				fbQuery = fbQuery.Where("cid.collateral_item_id = ?", jInt)
			} else {
				fbQuery = fbQuery.Where("ci.code = ?", jenis)
			}
		}

		if nama != "" {
			upperNama := "%" + strings.ToUpper(nama) + "%"
			fbQuery = fbQuery.Where("(UPPER(cid.name) LIKE ? OR UPPER(cid.code) LIKE ?)", upperNama, upperNama)
		}

		_ = fbQuery.Limit(20).Scan(&results).Error
	}

	c.JSON(http.StatusOK, results)
}

// GetBarangEmasHandler returns gold items (product_id = 1) from product.collateral_item_details
func GetBarangEmasHandler(c *gin.Context) {
	if config.DBPostgres == nil {
		c.JSON(http.StatusOK, []gin.H{})
		return
	}

	nama := c.Query("namaBarang")

	type GoldRow struct {
		IDDetail int    `gorm:"column:collateral_item_detail_id"`
		Code     string `gorm:"column:code"`
		Name     string `gorm:"column:name"`
		IsCert   bool   `gorm:"column:is_cert"`
	}

	var rows []GoldRow
	query := config.DBPostgres.Table("product.collateral_item_details cid").
		Select("cid.collateral_item_detail_id, cid.code, cid.name, ci.is_cert").
		Joins("JOIN product.collateral_items ci ON ci.collateral_items_id = cid.collateral_item_id").
		Where("ci.product_id = ? AND cid.is_active = ?", 1, true)

	if nama != "" {
		upperNama := "%" + strings.ToUpper(nama) + "%"
		query = query.Where("(UPPER(cid.name) LIKE ? OR UPPER(cid.code) LIKE ?)", upperNama, upperNama)
	}

	err := query.Order("cid.name ASC").Limit(20).Scan(&rows).Error
	if err != nil {
		c.JSON(http.StatusOK, []gin.H{})
		return
	}

	var output []gin.H
	for _, r := range rows {
		nmLower := strings.ToLower(r.Name)
		statusAntam := "Non Antam"
		if r.IsCert {
			statusAntam = "Antam"
		} else if strings.Contains(nmLower, "perhiasan") || strings.Contains(nmLower, "jewelry") ||
			strings.Contains(nmLower, "kalung") || strings.Contains(nmLower, "cincin") || strings.Contains(nmLower, "gelang") {
			statusAntam = "Perhiasan"
		}

		output = append(output, gin.H{
			"id_barang_detail": r.IDDetail,
			"kode_barang":      r.Code,
			"nama_barang":      strings.TrimSpace(r.Name),
			"status_antam":     statusAntam,
			"stle":             0,
		})
	}

	c.JSON(http.StatusOK, output)
}

// GetTaksiranProductHandler fetches products from product.product_offerings
func GetTaksiranProductHandler(c *gin.Context) {
	if config.DBPostgres == nil {
		c.JSON(http.StatusOK, []gin.H{})
		return
	}

	typeParam := c.Query("type")
	// log.Println("typeParam", typeParam)

	var results []struct {
		KdProduk string `json:"kd_produk" gorm:"column:kd_produk"`
		NmProduk string `json:"nm_produk" gorm:"column:nm_produk"`
	}

	query := config.DBPostgres.Table("product.product_offerings po").
		Select("po.offering_code AS kd_produk, po.offering_name AS nm_produk").
		Where("po.is_active = ?", true)

	if strings.ToLower(typeParam) == "emas" {
		query = query.Where("po.product_id = ?", 1)
	} else {
		query = query.Where("po.product_id = ?", 2)
	}

	err := query.Order("po.offering_name ASC").Scan(&results).Error
	if err != nil || len(results) == 0 {
		// Fallback by name
		q2 := config.DBPostgres.Table("product.product_offerings").
			Select("offering_code AS kd_produk, offering_name AS nm_produk").
			Where("is_active = ?", true)

		if strings.ToLower(typeParam) == "emas" {
			q2 = q2.Where("UPPER(offering_name) LIKE ?", "%EMAS%")
		} else {
			q2 = q2.Where("UPPER(offering_name) NOT LIKE ?", "%EMAS%")
		}
		_ = q2.Order("offering_name ASC").Scan(&results).Error
	}

	c.JSON(http.StatusOK, results)
}

// GetBiayaAdminHandler handles administrative fee calculation
func GetBiayaAdminHandler(c *gin.Context) {
	totalPinjamanStr := c.Query("total_pinjaman")
	kdProduk := c.Query("kd_produk")

	totalPinjaman, _ := strconv.ParseFloat(totalPinjamanStr, 64)
	biayaAdmin := calculateBiayaAdmin(totalPinjaman, kdProduk)

	c.JSON(http.StatusOK, gin.H{"biaya_admin": biayaAdmin})
}

// GetBiayaAdminEmasHandler handles admin fee calculation specifically for gold
func GetBiayaAdminEmasHandler(c *gin.Context) {
	totalPinjamanStr := c.Query("total_pinjaman")
	totalPinjaman, _ := strconv.ParseFloat(totalPinjamanStr, 64)
	biayaAdmin := calculateBiayaAdmin(totalPinjaman, "")

	c.JSON(http.StatusOK, gin.H{"biaya_admin": biayaAdmin})
}

func calculateBiayaAdmin(totalPinjaman float64, kdProduk string) float64 {
	if config.DBPostgres != nil && kdProduk != "" {
		var offering struct {
			OfferingID int `gorm:"column:offering_id"`
		}
		err := config.DBPostgres.Table("product.product_offerings").
			Select("offering_id").
			Where("offering_code = ?", kdProduk).
			First(&offering).Error

	if err == nil && offering.OfferingID > 0 {
			var adminTier struct {
				AdminFee float64 `gorm:"column:admin_fee"`
			}
			errTier := config.DBPostgres.Table("product.offering_admin_tiers").
				Select("admin_fee").
				Where("offering_id = ? AND min_amount <= ? AND max_amount >= ?", offering.OfferingID, totalPinjaman, totalPinjaman).
				First(&adminTier).Error

			if errTier == nil && adminTier.AdminFee > 0 {
				return adminTier.AdminFee
			}
		}
	}

	// Default fallback limits
	if totalPinjaman <= 5000000 {
		return 10000
	} else if totalPinjaman <= 10000000 {
		return 20000
	}
	return 50000
}

// GetWilayahHandler returns region list from master.region matching M_taksiran.php
func GetWilayahHandler(c *gin.Context) {
	if config.DBPostgres == nil {
		c.JSON(http.StatusOK, []gin.H{})
		return
	}

	id := c.Query("wilayah")
	itemType := c.Query("type")

	type RegionResult struct {
		RegionID   int    `json:"region_id" gorm:"column:region_id"`
		RegionName string `json:"region_name" gorm:"column:region_name"`
	}

	var results []RegionResult

	if id != "" && id != "undefined" {
		parentId, _ := strconv.Atoi(id)
		if itemType == "emas" {
			var detail struct {
				CollateralItemID int `gorm:"column:collateral_item_id"`
			}
			err := config.DBPostgres.Table("product.collateral_item_details").
				Select("collateral_item_id").
				Where("collateral_item_detail_id = ?", parentId).
				First(&detail).Error

			if err == nil && detail.CollateralItemID > 0 {
				parentId = detail.CollateralItemID
			}

			_ = config.DBPostgres.Table("master.region rg").
				Select("rg.region_id, rg.region_name").
				Joins("JOIN master.branch b ON b.region_id = rg.region_id").
				Joins("JOIN product.gold_stle_rules gsr ON gsr.branch_id = b.branch_id").
				Where("gsr.collateral_item_id = ? AND gsr.is_active = ?", parentId, true).
				Group("rg.region_id, rg.region_name").
				Order("rg.region_name ASC").
				Scan(&results).Error
		} else {
			_ = config.DBPostgres.Table("master.region rg").
				Select("rg.region_id, rg.region_name").
				Joins("JOIN product.collateral_item_prices cip ON cip.region_id = rg.region_id").
				Where("cip.collateral_item_detail_id = ? AND cip.is_active = ?", parentId, true).
				Group("rg.region_id, rg.region_name").
				Order("rg.region_name ASC").
				Scan(&results).Error
		}
	}

	// Fallback to all active regions if results empty or id empty
	if len(results) == 0 {
		_ = config.DBPostgres.Table("master.region").
			Select("region_id, region_name").
			Order("region_name ASC").
			Scan(&results).Error
	}

	c.JSON(http.StatusOK, results)
}

// GetNilaiTaksiranHandler gets estimated price for electronics matching M_taksiran.php
func GetNilaiTaksiranHandler(c *gin.Context) {
	if config.DBPostgres == nil {
		c.JSON(http.StatusOK, gin.H{"price": 0})
		return
	}

	idDetail := c.Query("id_barang_detail")
	kode := c.Query("kode_barang")
	regionID := c.Query("wilayah")

	if idDetail == "" && kode == "" {
		c.JSON(http.StatusOK, gin.H{"price": 0})
		return
	}

	var price float64
	query := config.DBPostgres.Table("product.collateral_item_prices cip").
		Select("COALESCE(cip.price, 0) AS price").
		Where("cip.is_active = ?", true)

	if idDetail != "" {
		detInt, _ := strconv.Atoi(idDetail)
		query = query.Where("cip.collateral_item_detail_id = ?", detInt)
	} else {
		query = query.Joins("JOIN product.collateral_item_details cid ON cid.collateral_item_detail_id = cip.collateral_item_detail_id").
			Where("cid.code = ?", kode)
	}

	if regionID != "" && regionID != "undefined" {
		regInt, _ := strconv.Atoi(regionID)
		query = query.Where("cip.region_id = ?", regInt)
	}

	err := query.Scan(&price).Error
	if err != nil || price == 0 {
		// Fallback without region_id filter
		fbQuery := config.DBPostgres.Table("product.collateral_item_prices cip").
			Select("COALESCE(cip.price, 0) AS price").
			Where("cip.is_active = ?", true)

		if idDetail != "" {
			detInt, _ := strconv.Atoi(idDetail)
			fbQuery = fbQuery.Where("cip.collateral_item_detail_id = ?", detInt)
		} else {
			fbQuery = fbQuery.Joins("JOIN product.collateral_item_details cid ON cid.collateral_item_detail_id = cip.collateral_item_detail_id").
				Where("cid.code = ?", kode)
		}
		_ = fbQuery.Scan(&price).Error
	}

	c.JSON(http.StatusOK, gin.H{"price": price})
}

// GetGoldStleHandler calculates gold buyback value matching M_taksiran.php
func GetGoldStleHandler(c *gin.Context) {
	if config.DBPostgres == nil {
		c.JSON(http.StatusOK, gin.H{"stle": 0})
		return
	}

	collateralDetailID := c.Query("id_barang_detail")
	karatStr := c.Query("karat")
	regionID := c.Query("wilayah")
	if regionID == "" {
		regionID = c.Query("region_id")
	}

	karat, _ := strconv.ParseFloat(karatStr, 64)

	if collateralDetailID == "" || karat == 0 {
		c.JSON(http.StatusOK, gin.H{"stle": 0})
		return
	}

	parentId, _ := strconv.Atoi(collateralDetailID)
	var detail struct {
		CollateralItemID int `gorm:"column:collateral_item_id"`
	}
	err := config.DBPostgres.Table("product.collateral_item_details").
		Select("collateral_item_id").
		Where("collateral_item_detail_id = ?", parentId).
		First(&detail).Error

	if err == nil && detail.CollateralItemID > 0 {
		parentId = detail.CollateralItemID
	}

	var stle float64
	query := config.DBPostgres.Table("product.gold_stle_rules gsr").
		Select("gsr.stle_buyback_per_gram").
		Where("gsr.collateral_item_id = ? AND gsr.min_karat_buyback <= ? AND gsr.max_karat_buyback >= ? AND gsr.is_active = ?", parentId, karat, karat, true)

	if regionID != "" && regionID != "undefined" {
		regInt, _ := strconv.Atoi(regionID)
		query = query.Joins("INNER JOIN master.branch b ON b.branch_id = gsr.branch_id").
			Where("b.region_id = ?", regInt)
	}

	errStle := query.Limit(1).Scan(&stle).Error
	if errStle != nil || stle == 0 {
		// Fallback without region_id join
		_ = config.DBPostgres.Table("product.gold_stle_rules gsr").
			Select("gsr.stle_buyback_per_gram").
			Where("gsr.collateral_item_id = ? AND gsr.min_karat_buyback <= ? AND gsr.max_karat_buyback >= ? AND gsr.is_active = ?", parentId, karat, karat, true).
			Limit(1).
			Scan(&stle).Error
	}

	var minKarat float64 = 8
	var minRow struct {
		MinKarat float64 `gorm:"column:min_karat"`
	}
	_ = config.DBPostgres.Table("product.gold_stle_rules gsr").
		Select("MIN(gsr.min_karat_buyback) as min_karat").
		Where("gsr.collateral_item_id = ? AND gsr.is_active = ?", parentId, true).
		Scan(&minRow).Error
	if minRow.MinKarat > 0 {
		minKarat = minRow.MinKarat
	}

	c.JSON(http.StatusOK, gin.H{"stle": stle, "min_karat": minKarat})
}

// GetBarangByKodeHandler fetches item detail by code
func GetBarangByKodeHandler(c *gin.Context) {
	if config.DBPostgres == nil {
		c.JSON(http.StatusOK, nil)
		return
	}

	kode := c.Query("kode")
	if kode == "" {
		c.JSON(http.StatusOK, nil)
		return
	}

	type BarangResult struct {
		IDBarangDetail int     `json:"id_barang_detail" gorm:"column:id"`
		KdBarang       string  `json:"kd_barang" gorm:"column:code"`
		NmBarang       string  `json:"nm_barang" gorm:"column:name"`
		Harga          float64 `json:"harga" gorm:"column:price"`
	}

	var res BarangResult
	err := config.DBPostgres.Table("product.collateral_item_details cid").
		Select("cid.collateral_item_detail_id AS id, cid.code, cid.name, COALESCE(cip.price, 0) AS price").
		Joins("LEFT JOIN product.collateral_item_prices cip ON cip.collateral_item_detail_id = cid.collateral_item_detail_id AND cip.is_active = true").
		Where("cid.code = ?", kode).
		First(&res).Error

	if err == nil {
		c.JSON(http.StatusOK, res)
		return
	}

	c.JSON(http.StatusOK, nil)
}

// SimpanTransaksiHandler stores transaction simulation and saves history
func SimpanTransaksiHandler(c *gin.Context) {
	userId, exists := c.Get("username")
	username := "system"
	if exists {
		username = userId.(string)
	}

	var payload struct {
		JenisBarang   string      `json:"jenis_barang"`
		NamaBarang    string      `json:"nama_barang"`
		NilaiTaksir   float64     `json:"nilai_taksir"`
		NilaiPinjaman float64     `json:"nilai_pinjaman"`
		ItemsEmas     interface{} `json:"items_emas"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payload tidak valid"})
		return
	}

	detailBytes, _ := json.Marshal(payload.ItemsEmas)
	detailStr := string(detailBytes)

	history := models.TaksiranHistory{
		CreatedBy:     username,
		JenisBarang:   payload.JenisBarang,
		NamaBarang:    payload.NamaBarang,
		NilaiTaksir:   payload.NilaiTaksir,
		NilaiPinjaman: payload.NilaiPinjaman,
		DetailData:    detailStr,
		CreatedAt:     time.Now(),
	}

	if config.DBMysql != nil {
		_ = config.DBMysql.AutoMigrate(&models.TaksiranHistory{})
		if err := config.DBMysql.Create(&history).Error; err != nil {
			log.Printf("Gagal simpan taksiran history: %v\n", err)
		}
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Transaksi taksiran berhasil disimpan"})
}

// GetTaksiranHistoryHandler fetches estimation history records filtered by creator/user
func GetTaksiranHistoryHandler(c *gin.Context) {
	if config.DBMysql == nil {
		c.JSON(http.StatusOK, gin.H{"data": []gin.H{}})
		return
	}

	_ = config.DBMysql.AutoMigrate(&models.TaksiranHistory{})

	userId, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	targetUser := userId.(string)
	if salesParam := c.Query("sales_user"); salesParam != "" {
		targetUser = salesParam
	} else if userParam := c.Query("user"); userParam != "" {
		targetUser = userParam
	}

	var histories []models.TaksiranHistory
	query := config.DBMysql.Table("tb_taksiran_history").Order("id DESC")

	if targetUser != "all" && targetUser != "*" && targetUser != "" {
		query = query.Where("created_by = ?", targetUser)
	}

	search := c.Query("search")
	if search != "" {
		likeSearch := "%" + search + "%"
		query = query.Where("(nama_barang LIKE ? OR jenis_barang LIKE ? OR created_by LIKE ?)", likeSearch, likeSearch, likeSearch)
	}

	if err := query.Limit(100).Find(&histories).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"data": []gin.H{}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": histories})
}

// GetElectronicScoringHandler fetches scoring sections, options, items, and option scores from DB2 Postgres matching new-feature/simulasi-taksiran
func GetElectronicScoringHandler(c *gin.Context) {
	if config.DBPostgres == nil {
		c.JSON(http.StatusOK, gin.H{"sections": []gin.H{}})
		return
	}

	collateralIDStr := c.Query("collateral_id")
	if collateralIDStr == "" {
		c.JSON(http.StatusOK, gin.H{"sections": []gin.H{}})
		return
	}

	collID, _ := strconv.Atoi(collateralIDStr)

	// Check if input ID is collateral_item_detail_id, resolve parent collateral_item_id
	var qDetail struct {
		CollateralItemID int `gorm:"column:collateral_item_id"`
	}
	_ = config.DBPostgres.Table("product.collateral_item_details").
		Select("collateral_item_id").
		Where("collateral_item_detail_id = ?", collID).
		Scan(&qDetail).Error
	if qDetail.CollateralItemID > 0 {
		collID = qDetail.CollateralItemID
	}

	// 1. Fetch collateral_scoring_sections
	type SectionRow struct {
		SectionID    int    `json:"section_id" gorm:"column:section_id"`
		SectionName  string `json:"section_name" gorm:"column:section_name"`
		SectionType  string `json:"section_type" gorm:"column:section_type"`
		DisplayOrder int    `json:"display_order" gorm:"column:display_order"`
	}

	var sections []SectionRow
	_ = config.DBPostgres.Table("product.collateral_scoring_sections").
		Select("id AS section_id, section_name, section_type, display_order").
		Where("collateral_item_id = ? AND is_active = ?", collID, true).
		Order("display_order ASC").
		Scan(&sections).Error

	if len(sections) == 0 {
		// Fallback without collateral_item_id filter
		_ = config.DBPostgres.Table("product.collateral_scoring_sections").
			Select("id AS section_id, section_name, section_type, display_order").
			Where("is_active = ?", true).
			Order("display_order ASC").
			Scan(&sections).Error
	}

	var resultSections []gin.H

	for _, sec := range sections {
		// 2. Fetch items for section
		type ItemRow struct {
			ItemID       int    `json:"item_id" gorm:"column:item_id"`
			ItemName     string `json:"item_name" gorm:"column:item_name"`
			DisplayOrder int    `json:"display_order" gorm:"column:display_order"`
		}

		var items []ItemRow
		_ = config.DBPostgres.Table("product.collateral_scoring_items").
			Select("id AS item_id, item_name, display_order").
			Where("section_id = ? AND is_active = ?", sec.SectionID, true).
			Order("display_order ASC, id DESC").
			Scan(&items).Error

		var mappedItems []gin.H
		sectionOptionsMap := make(map[int]gin.H)

		for _, itm := range items {
			// 3. Fetch scores for item
			type ScoreRow struct {
				ScoreID      int     `gorm:"column:score_id"`
				OptionID     int     `gorm:"column:option_id"`
				ScorePercent float64 `gorm:"column:score_percent"`
				IsReject     bool    `gorm:"column:is_reject"`
				OptionName   string  `gorm:"column:option_name"`
				DisplayOrder int     `gorm:"column:display_order"`
			}

			var scores []ScoreRow
			_ = config.DBPostgres.Table("product.collateral_scoring_item_option_scores s").
				Select("s.id AS score_id, s.option_id, s.score_percent, s.is_reject, o.option_name, o.display_order").
				Joins("INNER JOIN product.collateral_scoring_options o ON o.id = s.option_id").
				Where("s.item_id = ? AND s.is_active = ? AND o.is_active = ?", itm.ItemID, true, true).
				Order("o.display_order ASC, o.id ASC").
				Scan(&scores).Error

			if len(scores) == 0 {
				continue
			}

			var mappedScores []gin.H
			for _, sc := range scores {
				optUpper := strings.ToUpper(sc.OptionName)
				if sc.OptionID == 626 || sc.OptionID == 643 || optUpper == "BARET & LECET" || optUpper == "BERJAMUR & DEN" {
					continue
				}

				if _, exists := sectionOptionsMap[sc.OptionID]; !exists {
					sectionOptionsMap[sc.OptionID] = gin.H{
						"option_id":     sc.OptionID,
						"option_name":   sc.OptionName,
						"display_order": sc.DisplayOrder,
						"selected":      nil,
					}
				}

				mappedScores = append(mappedScores, gin.H{
					"score_id":      sc.ScoreID,
					"option_id":     sc.OptionID,
					"option_name":   sc.OptionName,
					"display_order": sc.DisplayOrder,
					"score_percent": sc.ScorePercent,
					"is_reject":     sc.IsReject,
				})
			}

			mappedItems = append(mappedItems, gin.H{
				"item_id":   itm.ItemID,
				"item_name": itm.ItemName,
				"scores":    mappedScores,
				"selected":  nil,
			})
		}

		var finalOptions []gin.H
		for _, opt := range sectionOptionsMap {
			finalOptions = append(finalOptions, opt)
		}

		if len(mappedItems) == 0 {
			if strings.ToLower(sec.SectionType) == "reject" {
				type RejectOpt struct {
					OptionID     int    `gorm:"column:option_id"`
					OptionName   string `gorm:"column:option_name"`
					DisplayOrder int    `gorm:"column:display_order"`
				}
				var rejectOpts []RejectOpt
				_ = config.DBPostgres.Table("product.collateral_scoring_options").
					Select("id AS option_id, option_name, display_order").
					Where("section_id = ? AND is_active = ?", sec.SectionID, true).
					Order("display_order ASC, id ASC").
					Scan(&rejectOpts).Error

				for _, ro := range rejectOpts {
					finalOptions = append(finalOptions, gin.H{
						"option_id":     ro.OptionID,
						"option_name":   ro.OptionName,
						"display_order": ro.DisplayOrder,
						"selected":      nil,
					})
				}
			} else {
				continue
			}
		}

		resultSections = append(resultSections, gin.H{
			"section_id":   sec.SectionID,
			"section_name": sec.SectionName,
			"section_type": strings.ToLower(sec.SectionType),
			"options":      finalOptions,
			"items":        mappedItems,
		})
	}

	c.JSON(http.StatusOK, gin.H{"sections": resultSections})
}
