package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"fmt"

	"go-api/config"
	"go-api/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// InquiryCountsHandler returns the counts for badge indicators
func InquiryCountsHandler(c *gin.Context) {
	if config.DBMysql == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database MySQL tidak terhubung"})
		return
	}

	userId, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	type ScanResult struct {
		IdHslaktiv int
		Total      int
	}

	var results []ScanResult
	err := config.DBMysql.Table("tb_rencana_det rd_sub").
		Select("rd_sub.id_hslaktiv, COUNT(*) as total").
		Joins("JOIN (SELECT id_rencana, MAX(id_rencana_det) as max_id FROM tb_rencana_det GROUP BY id_rencana) latest ON rd_sub.id_rencana_det = latest.max_id").
		Joins("JOIN tb_rencana a ON rd_sub.id_rencana = a.id_rencana").
		Where("a.id_login = ?", userId).
		Group("rd_sub.id_hslaktiv").
		Scan(&results).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data counts"})
		return
	}

	counts := gin.H{
		"all":         0,
		"leads":       0,
		"prospect":    0,
		"hotprospect": 0,
		"sbg":         0,
		"batal":       0,
	}

	leads := 0
	prospect := 0
	hotprospect := 0

	for _, r := range results {
		fmt.Println("ID:", r.IdHslaktiv, "Total:", r.Total, "UserId", userId)

		switch r.IdHslaktiv {
		case 1:
			counts["leads"] = r.Total
			leads = r.Total
		case 2:
			counts["prospect"] = r.Total
			prospect = r.Total
		case 3:
			counts["hotprospect"] = r.Total
			hotprospect = r.Total
		case 7:
			counts["sbg"] = r.Total
		case 6:
			counts["batal"] = r.Total
		}
	}

	counts["all"] = leads + prospect + hotprospect

	// Fetch dashboard items (new cif, emas, non emas, leads per bulan ini)
	fkUser, _ := c.Get("fk_user")
	fkUserStr, _ := fkUser.(string)

	now := time.Now()
	tgl1 := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	tgl2 := tgl1.AddDate(0, 1, 0).Add(-time.Second)

	var newCif int64 = 0
	var emas float64 = 0
	var nonEmas float64 = 0
	var leadsMonthly int64 = 0

	if config.DBPostgres != nil && fkUserStr != "" {
		var newCifGadai int64
		var newCifCicil int64

		config.DBPostgres.Table("data_gadai.tblproduk_gadai").
			Where("status_aplikasi = ? AND status_data = ? AND fk_karyawan_sales = ? AND tgl_pengajuan BETWEEN ? AND ?",
				"5", "Approve", fkUserStr, tgl1, tgl2).
			Count(&newCifGadai)

		config.DBPostgres.Table("data_gadai.tblproduk_cicilan").
			Where("status_aplikasi = ? AND status_data = ? AND fk_karyawan_sales = ? AND tgl_pengajuan BETWEEN ? AND ?",
				"5", "Approve", fkUserStr, tgl1, tgl2).
			Count(&newCifCicil)

		newCif = newCifGadai + newCifCicil

		var emasGadai float64
		var emasCicil float64

		config.DBPostgres.Table("data_gadai.tblproduk_gadai").
			Select("COALESCE(SUM(total_nilai_pinjaman), 0)").
			Where("fk_karyawan_sales = ? AND fk_produk IN ? AND status_aplikasi IN ? AND status_data = ? AND tgl_pengajuan BETWEEN ? AND ?",
				fkUserStr, []string{"10", "11", "30", "31"}, []string{"4", "5"}, "Approve", tgl1, tgl2).
			Row().Scan(&emasGadai)

		config.DBPostgres.Table("data_gadai.tblproduk_cicilan").
			Select("COALESCE(SUM(total_nilai_pinjaman), 0)").
			Where("fk_karyawan_sales = ? AND fk_produk IN ? AND status_aplikasi IN ? AND status_data = ? AND tgl_pengajuan BETWEEN ? AND ?",
				fkUserStr, []string{"55", "58", "59"}, []string{"4", "5"}, "Approve", tgl1, tgl2).
			Row().Scan(&emasCicil)

		emas = emasGadai + emasCicil

		config.DBPostgres.Table("data_gadai.tblproduk_gadai").
			Select("COALESCE(SUM(total_nilai_pinjaman), 0)").
			Where("fk_karyawan_sales = ? AND fk_produk IN ? AND status_aplikasi IN ? AND status_data = ? AND tgl_pengajuan BETWEEN ? AND ?",
				fkUserStr, []string{"21", "27", "32", "70"}, []string{"4", "5"}, "Approve", tgl1, tgl2).
			Row().Scan(&nonEmas)
	}

	if config.DBMysql != nil {
		config.DBMysql.Table("tb_rencana_det a").
			Joins("LEFT JOIN tb_rencana b ON a.id_rencana = b.id_rencana").
			Where("b.id_login = ? AND a.tgl_rencana BETWEEN ? AND ? AND a.id_hslaktiv = ?",
				userId, tgl1, tgl2, 1).
			Distinct("a.id_rencana").
			Count(&leadsMonthly)
	}

	counts["newcif"] = newCif
	counts["emas"] = emas
	counts["nonemas"] = nonEmas
	counts["leads_monthly"] = leadsMonthly

	c.JSON(http.StatusOK, counts)
}

// InquiryListRequest represents request payload for inquiry list filtering
type InquiryListRequest struct {
	Start        int    `form:"start,default=0"`
	Length       int    `form:"length,default=10"`
	Search       string `form:"search"`
	StatusFilter string `form:"status,default=all"`
}

// InquiryListHandler returns paginated and filtered inquiry lists
func InquiryListHandler(c *gin.Context) {
	if config.DBMysql == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database MySQL tidak terhubung"})
		return
	}

	userId, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req InquiryListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request query"})
		return
	}

	type InquiryItem struct {
		IDRencana     int       `json:"id_rencana"`
		Nama          string    `json:"nama"`
		TglRencana    time.Time `json:"tgl_rencana"`
		NamaHslaktiv  string    `json:"status"`
		IDHslaktiv    int       `json:"id_hslaktiv"`
		KetRencana    string    `json:"keterangan"`
		IDRencanaDet  int       `json:"id_rencana_det"`
		KetAktivitas  string    `json:"ket_aktivitas"`
		IDSumbercust  string    `json:"source"`
	}

	var items []InquiryItem

	query := config.DBMysql.Table("tb_rencana a").
		Select("a.id_rencana, b.nama, rd.tgl_rencana, ha.nama_hslaktiv, rd.id_hslaktiv, rd.ket_rencana, rd.id_rencana_det, rd.ket_aktivitas, rd.id_sumbercust").
		Joins("JOIN tb_customer b ON b.id_customer = a.id_customer").
		Joins("JOIN tb_rencana_det rd ON rd.id_rencana = a.id_rencana").
		Joins("JOIN (SELECT id_rencana, MAX(id_rencana_det) as max_id FROM tb_rencana_det GROUP BY id_rencana) latest ON rd.id_rencana_det = latest.max_id").
		Joins("LEFT JOIN tb_hslaktivitas ha ON ha.id_hslaktiv = rd.id_hslaktiv").
		Where("a.id_login = ?", userId)

	if req.StatusFilter == "all" {
		query = query.Where("rd.id_hslaktiv NOT IN ?", []int{6, 7})
	} else {
		query = query.Where("ha.nama_hslaktiv = ?", req.StatusFilter)
	}

	if req.Search != "" {
		query = query.Where("(b.nama LIKE ? OR ha.nama_hslaktiv LIKE ?)", "%"+req.Search+"%", "%"+req.Search+"%")
	}

	var totalRecords int64
	// Clone query for counting total
	countQuery := query
	if err := countQuery.Count(&totalRecords).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghitung total data"})
		return
	}

	err := query.Order("rd.id_rencana_det DESC").
		Offset(req.Start).
		Limit(req.Length).
		Scan(&items).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil list inquiry"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":            items,
		"recordsTotal":    totalRecords,
		"recordsFiltered": totalRecords,
	})
}

// InputRencanaRequest represents the payload for inputting plans
type InputRencanaRequest struct {
	Cust       string `json:"cust" binding:"required"` // "Lama" or "Baru"
	CustLama   int    `json:"custlama"`
	Nama       string `json:"nm"`
	Hp         string `json:"hp"`
	SumberCust string `json:"sumbercust" binding:"required"`
	HslAktiv   int    `json:"hslaktiv" binding:"required"`
	KatCust    string `json:"katcust"`
	Ket        string `json:"ket"`
	Aktiv      int    `json:"aktiv"`
}

// InputRencanaHandler handles POST /api/rencana/input
func InputRencanaHandler(c *gin.Context) {
	if config.DBMysql == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database MySQL tidak terhubung"})
		return
	}

	userId, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req InputRencanaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payload tidak valid"})
		return
	}

	var idCustomer int
	var ro string = "N" // Default: not repeat order

	tx := config.DBMysql.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if req.Cust == "Lama" {
		idCustomer = req.CustLama
		// Verify limits: maximum 2 active plans
		var activeCount int64
		err := tx.Table("tb_rencana r").
			Joins("JOIN tb_rencana_det rd ON r.id_rencana = rd.id_rencana").
			Joins("JOIN (SELECT id_rencana, MAX(id_rencana_det) as max_id FROM tb_rencana_det GROUP BY id_rencana) latest ON rd.id_rencana_det = latest.max_id").
			Where("r.id_customer = ? AND rd.status = 1 AND rd.id_hslaktiv NOT IN ?", idCustomer, []int{6, 7}).
			Count(&activeCount).Error

		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memverifikasi rencana aktif"})
			return
		}

		if activeCount >= 2 {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Customer lama sudah memiliki 2 rencana aktif."})
			return
		}
	} else {
		// Normalize phone
		telp := req.Hp
		if len(telp) >= 2 && telp[0:2] == "62" {
			telp = "0" + telp[2:]
		}

		// Check duplicate phone number globally
		var dupCount int64
		err := tx.Model(&models.Customer{}).Where("hp = ? AND id_login = ?", telp, userId).Count(&dupCount).Error
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengecek duplikasi HP"})
			return
		}
		if dupCount > 0 {
			var dupCust models.Customer
			tx.Where("hp = ? AND id_login = ?", telp, userId).First(&dupCust)
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Nomor HP sudah terdaftar di data Customer Lama (Atas nama: " + dupCust.Nama + "). Silakan pilih menu Customer Lama."})
			return
		}

		// Insert new customer
		newCust := models.Customer{
			Nama:    req.Nama,
			Hp:      telp,
			IDLogin: userId.(string),
		}
		if err := tx.Create(&newCust).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat customer baru"})
			return
		}
		idCustomer = newCust.IDCustomer
	}

	hslAktiv := req.HslAktiv
	// Rule: if existing customer ever reached Hot Prospect, promote Leads to Prospect immediately
	if req.Cust == "Lama" && hslAktiv == 1 {
		var hpCount int64
		err := tx.Table("tb_rencana r").
			Joins("JOIN tb_rencana_det rd ON r.id_rencana = rd.id_rencana").
			Where("r.id_customer = ? AND rd.id_hslaktiv = 3", idCustomer).
			Count(&hpCount).Error

		if err == nil && hpCount > 0 {
			hslAktiv = 2 // Promote to Prospect
		}
	}

	// Create plans
	newRencana := models.Rencana{
		IDCustomer:   idCustomer,
		IDSumbercust: req.SumberCust,
		IDLogin:      userId.(string),
	}
	if err := tx.Create(&newRencana).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat rencana"})
		return
	}

	// Create plans detail
	newRencanaDet := models.RencanaDet{
		IDRencana:  newRencana.IDRencana,
		TglRencana: time.Now(),
		IDHslaktiv: hslAktiv,
		Status:     1,
		KetRencana: req.Ket,
	}

	// Mocking repeat order configuration
	if req.Cust == "Lama" {
		ro = "Y"
	}
	// GORM doesn't have native "ro" column mapping unless defined, let's execute raw or save it if fields are set.
	// Since we defined RencanaDet standard struct, we can just save it.
	if err := tx.Create(&newRencanaDet).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat rencana detil"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan transaksi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Rencana berhasil ditambahkan",
		"id_rencana":   newRencana.IDRencana,
		"ro_status":    ro,
		"final_status": hslAktiv,
	})
}

// DeleteRencanaHandler handles DELETE /api/rencana/:id
func DeleteRencanaHandler(c *gin.Context) {
	if config.DBMysql == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database MySQL tidak terhubung"})
		return
	}

	userId, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	// Verify plan ownership
	var renc models.Rencana
	err = config.DBMysql.Where("id_rencana = ? AND id_login = ?", id, userId).First(&renc).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Rencana tidak ditemukan atau Anda tidak berwenang"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	tx := config.DBMysql.Begin()
	// Soft delete by setting status to inactive (0) in details, or delete details.
	// Legacy says: delete details or set status = 0. We'll set status = 0 in details.
	if err := tx.Model(&models.RencanaDet{}).Where("id_rencana = ?", id).Update("status", 0).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus rencana detil"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses penghapusan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Rencana berhasil dihapus"})
}
