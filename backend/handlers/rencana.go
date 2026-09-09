package handlers

import (
	"errors"
	"fmt"
	"log"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

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

	targetUser := userId.(string)
	if salesParam := c.Query("sales_user"); salesParam != "" {
		targetUser = salesParam
	} else if userParam := c.Query("user"); userParam != "" {
		targetUser = userParam
	}

	type ScanResult struct {
		IdHslaktiv int
		Total      int
	}

	var results []ScanResult
	query := config.DBMysql.Table("tb_rencana_det rd_sub").
		Select("rd_sub.id_hslaktiv, COUNT(*) as total").
		Joins("JOIN (SELECT id_rencana, MAX(id_rencana_det) as max_id FROM tb_rencana_det GROUP BY id_rencana) latest ON rd_sub.id_rencana_det = latest.max_id").
		Joins("JOIN tb_rencana a ON rd_sub.id_rencana = a.id_rencana")

	// log.Printf("isi query", query)
	if targetUser != "all" && targetUser != "*" {
		query = query.Where("a.id_login = ?", targetUser)
	}

	err := query.Group("rd_sub.id_hslaktiv").
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
		// case 7:
		// 	counts["sbg"] = r.Total
		case 6:
			counts["batal"] = r.Total
		}
	}

	counts["all"] = leads + prospect + hotprospect
	c.JSON(http.StatusOK, counts)
}

// InquiryItem represents inquiry record in Datatables response
type InquiryItem struct {
	IDRencana    int    `json:"id_rencana" gorm:"column:id_rencana"`
	IDCustomer   int    `json:"id_customer" gorm:"column:id_customer"`
	Nama         string `json:"nama" gorm:"column:nama"`
	Hp           string `json:"hp" gorm:"column:hp"`
	TglRencana   string `json:"tgl_rencana" gorm:"column:tgl_rencana"`
	NamaHslaktiv string `json:"nama_hslaktiv" gorm:"column:nama_hslaktiv"`
	IDHslaktiv   int    `json:"id_hslaktiv" gorm:"column:id_hslaktiv"`
	KetRencana   string `json:"ket_rencana" gorm:"column:ket_rencana"`
	IDRencanaDet int    `json:"id_rencana_det" gorm:"column:id_rencana_det"`
	KetAktivitas string `json:"ket_aktivitas" gorm:"column:ket_aktivitas"`
	IDSumbercust string `json:"id_sumbercust" gorm:"column:id_sumbercust"`
}

// InquiryListRequest represents Datatables query parameters
type InquiryListRequest struct {
	StatusFilter string `form:"status_filter"`
	Status       string `form:"status"`
	Search       string `form:"search[value]"`
	SearchDirect string `form:"search"`
	Start        int    `form:"start"`
	Length       int    `form:"length"`
}

// InquiryListHandler handles GET /api/inquiry/filter
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameter tidak valid"})
		return
	}

	if req.Length <= 0 {
		req.Length = 10
	}

	var items []InquiryItem

	targetUser := userId.(string)
	if salesParam := c.Query("sales_user"); salesParam != "" {
		targetUser = salesParam
	} else if userParam := c.Query("user"); userParam != "" {
		targetUser = userParam
	}

	baseQuery := config.DBMysql.Table("tb_rencana a").
		Joins("JOIN tb_customer b ON b.id_customer = a.id_customer").
		Joins("JOIN tb_rencana_det rd ON rd.id_rencana = a.id_rencana").
		Joins("JOIN (SELECT id_rencana, MAX(id_rencana_det) as max_id FROM tb_rencana_det GROUP BY id_rencana) latest ON rd.id_rencana_det = latest.max_id").
		Joins("LEFT JOIN tb_hslaktivitas ha ON ha.id_hslaktiv = rd.id_hslaktiv").
		Where("rd.id_hslaktiv NOT IN ?", []int{6, 7})

	if targetUser != "all" && targetUser != "*" {
		baseQuery = baseQuery.Where("a.id_login = ?", targetUser)
	}

	filterStatus := req.StatusFilter
	if filterStatus == "" {
		filterStatus = req.Status
	}

	switch strings.ToLower(filterStatus) {
	case "leads":
		baseQuery = baseQuery.Where("(rd.id_hslaktiv = 1 OR ha.nama_hslaktiv LIKE ?)", "%lead%")
	case "prospect":
		baseQuery = baseQuery.Where("(rd.id_hslaktiv = 2 OR ha.nama_hslaktiv = ?)", "Prospect")
	case "hot prospect", "hotprospect", "hot_prospect":
		baseQuery = baseQuery.Where("(rd.id_hslaktiv = 3 OR ha.nama_hslaktiv LIKE ?)", "%hot%")
	// case "sbg", "deal", "do":
	// 	baseQuery = baseQuery.Where("(rd.id_hslaktiv = 7 OR ha.nama_hslaktiv LIKE ? OR ha.nama_hslaktiv LIKE ?)", "%sbg%", "%deal%")
	case "batal":
		baseQuery = baseQuery.Where("(rd.id_hslaktiv = 6 OR ha.nama_hslaktiv LIKE ?)", "%batal%")
	case "all", "":
		baseQuery = baseQuery.Where("rd.id_hslaktiv != ?", 6)
	default:
		baseQuery = baseQuery.Where("(ha.nama_hslaktiv LIKE ? OR CAST(rd.id_hslaktiv AS CHAR) = ?)", "%"+filterStatus+"%", filterStatus)
	}

	searchVal := req.Search
	if searchVal == "" {
		searchVal = req.SearchDirect
	}

	if searchVal != "" {
		likeSearch := "%" + searchVal + "%"
		baseQuery = baseQuery.Where("(b.nama LIKE ? OR b.hp LIKE ? OR ha.nama_hslaktiv LIKE ?)", likeSearch, likeSearch, likeSearch)
	}

	var totalRecords int64
	countQuery := baseQuery
	if err := countQuery.Count(&totalRecords).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghitung total data"})
		return
	}

	err := baseQuery.
		Select("a.id_rencana, a.id_customer, b.nama, b.hp, rd.tgl_rencana, ha.nama_hslaktiv, rd.id_hslaktiv, rd.ket_rencana, rd.id_rencana_det, rd.ket_aktivitas, a.id_sumbercust").
		Order("rd.id_rencana_det DESC").
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

	log.Println("Inquiry list fetched successfully", items)
}

// InputRencanaRequest represents the payload for inputting plans
type InputRencanaRequest struct {
	IDRencana  int    `json:"id_rencana"`
	Cust       string `json:"cust"` // "Lama" or "Baru"
	CustLama   int    `json:"custlama"`
	Nama       string `json:"nm"`
	Hp         string `json:"hp"`
	SumberCust string `json:"sumbercust"`
	HslAktiv   int    `json:"hslaktiv"`
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payload tidak valid: " + err.Error()})
		return
	}

	// Default values
	if req.Cust == "" {
		if req.CustLama > 0 {
			req.Cust = "Lama"
		} else {
			req.Cust = "Baru"
		}
	}

	hslAktiv := req.HslAktiv
	if hslAktiv <= 0 {
		hslAktiv = 1
	}

	var idCustomer int
	var ro string = "N"
	var targetIDRencana int = req.IDRencana

	tx := config.DBMysql.Begin()
	tx.Exec("SET SESSION sql_mode = ''")
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if req.Cust == "Lama" {
		if targetIDRencana > 0 {
			var existingRencana models.Rencana
			if err := tx.Table("tb_rencana").Where("id_rencana = ?", targetIDRencana).First(&existingRencana).Error; err == nil {
				idCustomer = existingRencana.IDCustomer
			}
		}

		if idCustomer == 0 && req.CustLama > 0 {
			// Check if req.CustLama is an id_customer in tb_customer
			var cust models.Customer
			if err := tx.Table("tb_customer").Where("id_customer = ?", req.CustLama).First(&cust).Error; err == nil && cust.IDCustomer > 0 {
				idCustomer = cust.IDCustomer
			} else {
				// Fallback: check if req.CustLama is an id_rencana
				var renc models.Rencana
				if err := tx.Table("tb_rencana").Where("id_rencana = ?", req.CustLama).First(&renc).Error; err == nil && renc.IDCustomer > 0 {
					idCustomer = renc.IDCustomer
					if targetIDRencana == 0 {
						targetIDRencana = renc.IDRencana
					}
				}
			}
		}

		if targetIDRencana == 0 && idCustomer > 0 {
			var activePlanCount int64
			err := tx.Table("tb_rencana").Where("id_customer = ?", idCustomer).Count(&activePlanCount).Error
			if err == nil && activePlanCount >= 2 {
				tx.Rollback()
				c.JSON(http.StatusBadRequest, gin.H{"error": "Batas maksimal rencana / leads per customer hanya boleh 2. Customer ini sudah memiliki 2 rencana."})
				return
			}
		}
	} else {
		if strings.TrimSpace(req.Nama) == "" {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Nama customer baru wajib diisi dan tidak boleh kosong."})
			return
		}

		// Clean phone number: keep only digits
		telp := regexp.MustCompile(`\D`).ReplaceAllString(req.Hp, "")
		if len(telp) >= 2 && telp[0:2] == "62" {
			telp = "0" + telp[2:]
		}

		if telp == "" {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Nomor HP customer baru wajib diisi dan tidak boleh kosong."})
			return
		}

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
			c.JSON(http.StatusBadRequest, gin.H{"error": "Nomor HP sudah terdaftar di data Customer Lama (Atas nama: " + dupCust.Nama + "). Silakan pilih menu Pelanggan Lama."})
			return
		}

		newCust := models.Customer{
			Nama:    strings.TrimSpace(req.Nama),
			Hp:      telp,
			IDLogin: userId.(string),
		}
		if err := tx.Create(&newCust).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal membuat customer baru", "error": err.Error()})
			return
		}
		idCustomer = newCust.IDCustomer
	}

	if req.Cust == "Lama" && hslAktiv == 1 && idCustomer > 0 {
		var hpCount int64
		err := tx.Table("tb_rencana r").
			Joins("JOIN tb_rencana_det rd ON r.id_rencana = rd.id_rencana").
			Where("r.id_customer = ? AND rd.id_hslaktiv = 3", idCustomer).
			Count(&hpCount).Error

		if err == nil && hpCount > 0 {
			hslAktiv = 2
		}
	}

	if targetIDRencana == 0 {
		newRencana := models.Rencana{
			IDCustomer:   idCustomer,
			IDSumbercust: req.SumberCust,
			IDLogin:      userId.(string),
		}
		if err := tx.Create(&newRencana).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat rencana: " + err.Error()})
			return
		}
		targetIDRencana = newRencana.IDRencana
	}

	newRencanaDet := models.RencanaDet{
		IDRencana:  targetIDRencana,
		TglRencana: time.Now(),
		IDHslaktiv: hslAktiv,
		Status:     1,
		KetRencana: req.Ket,
	}

	if req.Cust == "Lama" {
		ro = "Y"
	}
	if err := tx.Create(&newRencanaDet).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat rencana detil: " + err.Error()})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan transaksi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Rencana berhasil disimpan",
		"id_rencana":   targetIDRencana,
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

func resolveSalesNumericIDs(userStr string) []int {
	var validIDs []int
	if config.DBPostgres == nil || userStr == "" {
		return validIDs
	}

	if numID, err := strconv.Atoi(userStr); err == nil && numID > 0 {
		validIDs = append(validIDs, numID)
	}

	type EmpResult struct {
		EmployeeID int `gorm:"column:employee_id"`
	}
	var empResults []EmpResult

	query := config.DBPostgres.Table("master.employee e").
		Select("e.employee_id").
		Joins("LEFT JOIN auth.users u ON u.employee_id = e.employee_id")

	if numID, err := strconv.Atoi(userStr); err == nil && numID > 0 {
		query = query.Where("e.employee_id = ?", numID)
	} else {
		query = query.Where("e.employee_code = ? OR u.username = ?", userStr, userStr)
	}

	_ = query.Scan(&empResults).Error

	for _, r := range empResults {
		if r.EmployeeID > 0 {
			validIDs = append(validIDs, r.EmployeeID)
		}
	}

	if len(validIDs) == 0 {
		var uRow struct {
			EmployeeID int `gorm:"column:employee_id"`
		}
		err := config.DBPostgres.Table("auth.users").
			Select("employee_id").
			Where("username = ?", userStr).
			First(&uRow).Error
		if err == nil && uRow.EmployeeID > 0 {
			validIDs = append(validIDs, uRow.EmployeeID)
		}
	}

	seen := make(map[int]bool)
	var result []int
	for _, id := range validIDs {
		if !seen[id] {
			seen[id] = true
			result = append(result, id)
		}
	}
	return result
}

func getLeaderboardPosition(validIDs []int, tgl1Str, tgl2Str string) gin.H {
	defaultRes := gin.H{
		"rank":               1,
		"total_participants": 1,
		"total_omzet":        0.0,
		"total_emas":         0.0,
		"total_non_emas":     0.0,
		"gap_to_next":        0.0,
	}

	if config.DBPostgres == nil {
		return defaultRes
	}

	type LeaderboardRow struct {
		SalesEmployeeID int     `gorm:"column:sales_employee_id"`
		TotalEmas       float64 `gorm:"column:total_emas"`
		TotalNonEmas    float64 `gorm:"column:total_non_emas"`
		TotalOmzet      float64 `gorm:"column:total_omzet"`
	}

	var rows []LeaderboardRow
	err := config.DBPostgres.Table("gadai.transactions t").
		Select(`
			t.sales_employee_id,
			SUM(CASE WHEN o.product_id = 1 THEN t.requested_loan ELSE 0 END) AS total_emas,
			SUM(CASE WHEN o.product_id IN (2, 3) THEN t.requested_loan ELSE 0 END) AS total_non_emas,
			SUM(t.requested_loan) AS total_omzet
		`).
		Joins("JOIN product.product_offerings o ON t.offering_id = o.offering_id").
		Where("t.application_status IN ('NEW', 'ADO', 'REO') AND t.loan_status IN ('ACTIVE', 'REDEEM', 'LUNAS', 'EXTEND') AND t.submission_date BETWEEN ? AND ? AND t.sales_employee_id IS NOT NULL", tgl1Str, tgl2Str).
		Group("t.sales_employee_id").
		Order("total_omzet DESC").
		Scan(&rows).Error

	if err != nil || len(rows) == 0 {
		return defaultRes
	}

	targetMap := make(map[int]bool)
	for _, id := range validIDs {
		targetMap[id] = true
	}

	rank := 0
	var userStats *LeaderboardRow
	var nextRankStats *LeaderboardRow
	totalParticipants := len(rows)

	for index, r := range rows {
		if targetMap[r.SalesEmployeeID] {
			rank = index + 1
			userStats = &rows[index]
			if index > 0 {
				nextRankStats = &rows[index-1]
			}
			break
		}
	}

	if rank == 0 {
		totalParticipants += 1
		rank = totalParticipants
		userStats = &LeaderboardRow{
			SalesEmployeeID: 0,
			TotalEmas:       0,
			TotalNonEmas:    0,
			TotalOmzet:      0,
		}
		if len(rows) > 0 {
			nextRankStats = &rows[len(rows)-1]
		}
	}

	gapToNext := 0.0
	if nextRankStats != nil {
		gapToNext = math.Max(0, nextRankStats.TotalOmzet-userStats.TotalOmzet)
	}

	return gin.H{
		"rank":               rank,
		"total_participants": math.Max(1, float64(totalParticipants)),
		"total_omzet":        userStats.TotalOmzet,
		"total_emas":         userStats.TotalEmas,
		"total_non_emas":     userStats.TotalNonEmas,
		"gap_to_next":        gapToNext,
	}
}

func getDailyPacingBooking(validIDs []int, tgl1Str, tgl2Str string, daysInMonth int) gin.H {
	pacingSopiga := make(map[string]float64)
	pacingSelada := make(map[string]float64)
	for i := 1; i <= daysInMonth; i++ {
		dayKey := fmt.Sprintf("%d", i)
		pacingSopiga[dayKey] = 0
		pacingSelada[dayKey] = 0
	}

	if config.DBPostgres != nil && len(validIDs) > 0 {
		type PacingRow struct {
			TglHari      int     `gorm:"column:tgl_hari"`
			TotalBooking float64 `gorm:"column:total_booking"`
		}

		var rows []PacingRow
		err := config.DBPostgres.Table("gadai.transactions t").
			Select("EXTRACT(DAY FROM COALESCE(t.submission_date, t.created_at)) AS tgl_hari, SUM(t.requested_loan) AS total_booking").
			Where("t.sales_employee_id IN ? AND t.application_status IN ('NEW', 'ADO', 'REO') AND t.loan_status IN ('ACTIVE', 'REDEEM', 'LUNAS', 'EXTEND') AND COALESCE(t.submission_date, t.created_at) BETWEEN ? AND ?", validIDs, tgl1Str, tgl2Str).
			Group("EXTRACT(DAY FROM COALESCE(t.submission_date, t.created_at))").
			Scan(&rows).Error

		if err == nil {
			for _, r := range rows {
				if r.TglHari >= 1 && r.TglHari <= daysInMonth {
					dayKey := fmt.Sprintf("%d", r.TglHari)
					pacingSopiga[dayKey] = r.TotalBooking
				}
			}
		}
	}

	return gin.H{
		"selada":    pacingSelada,
		"sopigapro": pacingSopiga,
		"sopiga":    pacingSopiga,
	}
}

// DashboardMetricsHandler returns real-time metrics matching legacy C_sistem and IncentiveCalculator
func DashboardMetricsHandler(c *gin.Context) {
	userId, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	loggedInUser := fmt.Sprintf("%v", userId)
	userStr := loggedInUser

	// Check review mode & sales user selection
	paramSales := c.Query("sales_user")
	isReview := false
	var salesList []gin.H

	// 1. Fetch user grade & info from MySQL POS
	var idJabatan int = 17
	var nmJabatan string = "JUNIOR"

	if config.DBMysqlPos != nil {
		var userPos struct {
			IDJabatan int    `gorm:"column:id_jabatan"`
			NmJabatan string `gorm:"column:nm_jabatan"`
		}
		err := config.DBMysqlPos.Table("tbl_user u").
			Select("u.id_jabatan, j.nm_jabatan").
			Joins("LEFT JOIN mst_jabatan j ON j.id_jabatan = u.id_jabatan").
			Where("u.username = ?", loggedInUser).
			First(&userPos).Error

		if err == nil && userPos.IDJabatan > 0 {
			if userPos.IDJabatan == 4 || userPos.IDJabatan == 21 || userPos.IDJabatan == 74 {
				isReview = true
			}
		}

		if isReview {
			type SalesRow struct {
				Username  string `gorm:"column:username"`
				Nama      string `gorm:"column:nama"`
				IDJabatan int    `gorm:"column:id_jabatan"`
				NmJabatan string `gorm:"column:nm_jabatan"`
			}
			var rows []SalesRow
			_ = config.DBMysqlPos.Table("tbl_user u").
				Select("u.username, u.nama, u.id_jabatan, j.nm_jabatan").
				Joins("LEFT JOIN mst_jabatan j ON j.id_jabatan = u.id_jabatan").
				Where("u.void = ? AND u.id_jabatan IN ?", 1, []int{6, 17, 29, 74}).
				Order("u.nama ASC").
				Scan(&rows).Error

			for _, r := range rows {
				salesList = append(salesList, gin.H{
					"username":   r.Username,
					"nama":       r.Nama,
					"id_jabatan": r.IDJabatan,
					"nm_jabatan": r.NmJabatan,
					"fkuser":     r.Username,
				})
			}

			if paramSales != "" {
				userStr = paramSales
			} else if len(salesList) > 0 {
				userStr = salesList[0]["username"].(string)
			}
		}
	}

	// Fetch selected target sales info if different
	if config.DBMysqlPos != nil {
		var targetPos struct {
			IDJabatan int    `gorm:"column:id_jabatan"`
			NmJabatan string `gorm:"column:nm_jabatan"`
		}
		err := config.DBMysqlPos.Table("tbl_user u").
			Select("u.id_jabatan, j.nm_jabatan").
			Joins("LEFT JOIN mst_jabatan j ON j.id_jabatan = u.id_jabatan").
			Where("u.username = ?", userStr).
			First(&targetPos).Error

		if err == nil && targetPos.IDJabatan > 0 {
			idJabatan = targetPos.IDJabatan
			if targetPos.NmJabatan != "" {
				nmJabatan = targetPos.NmJabatan
			}
		}
	}

	// Date calculations (Current Month & Last Month - matching PHP get_dashboard_date_range)
	now := time.Now()
	tgl1 := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	tgl2 := time.Date(now.Year(), now.Month()+1, 0, 23, 59, 59, 0, now.Location())
	daysInMonth := tgl2.Day()

	lastMonth := now.AddDate(0, -1, 0)
	lastTgl1 := time.Date(lastMonth.Year(), lastMonth.Month(), 1, 0, 0, 0, 0, now.Location())
	lastTgl2 := time.Date(lastMonth.Year(), lastMonth.Month()+1, 0, 23, 59, 59, 0, now.Location())

	tgl1Str := tgl1.Format("2006-01-02 15:04:05")
	tgl2Str := tgl2.Format("2006-01-02 15:04:05")
	lastTgl1Str := lastTgl1.Format("2006-01-02 15:04:05")
	lastTgl2Str := lastTgl2.Format("2006-01-02 15:04:05")

	// 2. Query Inquiry Badge Counts & Sales Funnel from MySQL (tb_rencana & tb_rencana_det)
	overallLeads, overallProspect, overallHotProspect, overallSBG, overallBatal := 0, 0, 0, 0, 0
	leadsMonthly, prospectMonthly, hotprospectMonthly, sbgMonthly := 0, 0, 0, 0
	lastMonthLeads := 0

	if config.DBMysql != nil {
		// --- Overall Accumulation (All-Time History) ---
		type ActiveResult struct {
			IdHslaktiv int
			Total      int
		}
		var activeResults []ActiveResult
		_ = config.DBMysql.Table("tb_rencana_det rd_sub").
			Select("rd_sub.id_hslaktiv, COUNT(*) as total").
			Joins("JOIN (SELECT id_rencana, MAX(id_rencana_det) as max_id FROM tb_rencana_det GROUP BY id_rencana) latest ON rd_sub.id_rencana_det = latest.max_id").
			Joins("JOIN tb_rencana a ON rd_sub.id_rencana = a.id_rencana").
			Where("a.id_login = ?", userStr).
			Group("rd_sub.id_hslaktiv").
			Scan(&activeResults).Error

		cLeads, cProspect, cHotProspect, cSBG, cBatal := 0, 0, 0, 0, 0
		for _, r := range activeResults {
			switch r.IdHslaktiv {
			case 1:
				cLeads = r.Total
			case 2:
				cProspect = r.Total
			case 3:
				cHotProspect = r.Total
			case 7:
				cSBG = r.Total
			case 6:
				cBatal = r.Total
			}
		}

		// Overall Funnel (Cumulative all-time)
		overallSBG = cSBG
		overallHotProspect = cHotProspect + cSBG
		overallProspect = cProspect + cHotProspect + cSBG
		overallLeads = cLeads + cProspect + cHotProspect + cSBG
		overallBatal = cBatal

		// --- Current Month Data (New data for current month only, e.g. 2026-09) ---
		// 1. Monthly Leads created in current month (id_hslaktiv = 1)
		var currLeadsRes struct{ Total int }
		_ = config.DBMysql.Table("tb_rencana_det a").
			Select("COUNT(DISTINCT a.id_rencana) as total").
			Joins("JOIN tb_rencana b ON a.id_rencana = b.id_rencana").
			Where("b.id_login = ? AND a.tgl_rencana BETWEEN ? AND ? AND a.id_hslaktiv = 1", userStr, tgl1Str, tgl2Str).
			Scan(&currLeadsRes).Error
		leadsMonthly = currLeadsRes.Total

		// 2. Monthly Prospect reached in current month (id_hslaktiv = 2)
		var currPRes struct{ Total int }
		_ = config.DBMysql.Table("tb_rencana_det a").
			Select("COUNT(DISTINCT a.id_rencana) as total").
			Joins("JOIN tb_rencana b ON a.id_rencana = b.id_rencana").
			Where("b.id_login = ? AND a.tgl_rencana BETWEEN ? AND ? AND a.id_hslaktiv = 2", userStr, tgl1Str, tgl2Str).
			Scan(&currPRes).Error
		prospectMonthly = currPRes.Total

		// 3. Monthly Hot Prospect reached in current month (id_hslaktiv = 3)
		var currHPRes struct{ Total int }
		_ = config.DBMysql.Table("tb_rencana_det a").
			Select("COUNT(DISTINCT a.id_rencana) as total").
			Joins("JOIN tb_rencana b ON a.id_rencana = b.id_rencana").
			Where("b.id_login = ? AND a.tgl_rencana BETWEEN ? AND ? AND a.id_hslaktiv = 3", userStr, tgl1Str, tgl2Str).
			Scan(&currHPRes).Error
		hotprospectMonthly = currHPRes.Total

		// 4. Monthly SBG reached in current month (id_hslaktiv = 7)
		var currSBGRes struct{ Total int }
		_ = config.DBMysql.Table("tb_rencana_det a").
			Select("COUNT(DISTINCT a.id_rencana) as total").
			Joins("JOIN tb_rencana b ON a.id_rencana = b.id_rencana").
			Where("b.id_login = ? AND a.tgl_rencana BETWEEN ? AND ? AND a.id_hslaktiv = 7", userStr, tgl1Str, tgl2Str).
			Scan(&currSBGRes).Error
		sbgMonthly = currSBGRes.Total

		// Monthly leads created for last month
		var lastLeadsRes struct{ Total int }
		_ = config.DBMysql.Table("tb_rencana_det a").
			Select("COUNT(DISTINCT a.id_rencana) as total").
			Joins("JOIN tb_rencana b ON a.id_rencana = b.id_rencana").
			Where("b.id_login = ? AND a.tgl_rencana BETWEEN ? AND ? AND a.id_hslaktiv = 1", userStr, lastTgl1Str, lastTgl2Str).
			Scan(&lastLeadsRes).Error
		lastMonthLeads = lastLeadsRes.Total
	}

	// 3. Resolve all numeric employee_ids from DB2 (PostgreSQL sam_live) matching M_incentive.php
	validNumericIDs := resolveSalesNumericIDs(userStr)

	// 4. Query Transactions & Omzet from PostgreSQL sam_live
	emas := 0.0
	nonemas := 0.0
	biayaAdminEmas := 0.0
	biayaAdminNonEmas := 0.0
	totalTrxCount := 0
	gramasiEmas := 0.0
	newcif := 0
	lastMonthCif := 0

	if config.DBPostgres != nil && len(validNumericIDs) > 0 {
		var txSummary struct {
			NewCIF            int     `gorm:"column:newcif"`
			Emas              float64 `gorm:"column:emas"`
			BiayaAdminEmas    float64 `gorm:"column:biaya_admin_emas"`
			NonEmas           float64 `gorm:"column:nonemas"`
			BiayaAdminNonEmas float64 `gorm:"column:biaya_admin_nonemas"`
			TotalTrx          int     `gorm:"column:total_trx"`
		}

		_ = config.DBPostgres.Table("gadai.transactions t").
			Select(`
				COUNT(CASE WHEN t.loan_status = 'ACTIVE' AND t.application_status IN ('NEW') THEN t.doc_no END) AS newcif,
				SUM(CASE WHEN o.product_id = 1 AND t.application_status IN ('NEW', 'ADO', 'REO') AND t.loan_status IN ('ACTIVE', 'REDEEM', 'LUNAS', 'EXTEND') THEN t.requested_loan ELSE 0 END) AS emas,
				SUM(CASE WHEN o.product_id = 1 AND t.application_status IN ('NEW', 'ADO', 'REO') AND t.loan_status IN ('ACTIVE', 'REDEEM', 'LUNAS', 'EXTEND') THEN t.admin_cost_amount ELSE 0 END) AS biaya_admin_emas,
				SUM(CASE WHEN o.product_id IN (2, 3) AND t.application_status IN ('NEW', 'ADO', 'REO') AND t.loan_status IN ('ACTIVE', 'REDEEM', 'LUNAS', 'EXTEND') THEN t.requested_loan ELSE 0 END) AS nonemas,
				SUM(CASE WHEN o.product_id IN (2, 3) AND t.application_status IN ('NEW', 'ADO', 'REO') AND t.loan_status IN ('ACTIVE', 'REDEEM', 'LUNAS', 'EXTEND') THEN t.admin_cost_amount ELSE 0 END) AS biaya_admin_nonemas,
				COUNT(CASE WHEN t.application_status IN ('NEW', 'ADO', 'REO') AND t.loan_status IN ('ACTIVE', 'REDEEM', 'LUNAS', 'EXTEND') THEN t.doc_no END) AS total_trx
			`).
			Joins("LEFT JOIN product.product_offerings o ON t.offering_id = o.offering_id").
			Where("t.sales_employee_id IN ? AND t.submission_date BETWEEN ? AND ?", validNumericIDs, tgl1Str, tgl2Str).
			Scan(&txSummary).Error

		newcif = txSummary.NewCIF
		emas = txSummary.Emas
		biayaAdminEmas = txSummary.BiayaAdminEmas
		nonemas = txSummary.NonEmas
		biayaAdminNonEmas = txSummary.BiayaAdminNonEmas
		totalTrxCount = txSummary.TotalTrx

		var gramRes struct{ TotalGramasi float64 }
		_ = config.DBPostgres.Table("gadai.transactions t").
			Select("SUM(rg.net_weight) AS total_gramasi").
			Joins("JOIN product.product_offerings o ON t.offering_id = o.offering_id").
			Joins("JOIN gadai.transaction_items rg ON rg.transaction_id = t.transaction_id").
			Where("o.product_id = 1 AND t.loan_status IN ('ACTIVE', 'REDEEM', 'LUNAS', 'EXTEND') AND t.sales_employee_id IN ? AND t.submission_date BETWEEN ? AND ?", validNumericIDs, tgl1Str, tgl2Str).
			Scan(&gramRes).Error
		gramasiEmas = gramRes.TotalGramasi

		var lmCifRes struct{ TotalBaris int }
		_ = config.DBPostgres.Table("gadai.transactions t").
			Select("COUNT(t.doc_no) AS total_baris").
			Where("t.loan_status = 'ACTIVE' AND t.application_status IN ('NEW') AND t.sales_employee_id IN ? AND t.submission_date BETWEEN ? AND ?", validNumericIDs, lastTgl1Str, lastTgl2Str).
			Scan(&lmCifRes).Error
		lastMonthCif = lmCifRes.TotalBaris
	}

	counts := gin.H{
		"all":                 overallLeads,
		"leads":               overallLeads,
		"prospect":            overallProspect,
		"hotprospect":         overallHotProspect,
		"sbg":                 overallSBG,
		"batal":               overallBatal,
		"leads_monthly":       leadsMonthly,
		"prospect_monthly":    prospectMonthly,
		"hotprospect_monthly": hotprospectMonthly,
		"sbg_monthly":         sbgMonthly,
	}

	// 5. Incentive Engine & Target Progress Engine Calculations (100% matching IncentiveCalculator.php)
	totalEmasPoints := emas / 15000000.0
	totalNonEmasPoints := nonemas / 2000000.0
	poinTotal := totalEmasPoints + totalNonEmasPoints
	poinTotalFormatted := math.Round(poinTotal*10) / 10

	totalAdminAlokasi := (biayaAdminEmas + biayaAdminNonEmas) * 0.45

	appliedRate := 0.0
	isQualified := false
	projectedRate := 0.25

	if idJabatan == 6 { // Trainee
		projectedRate = 0.20
		if poinTotal >= 20.0 {
			appliedRate = 0.20
			isQualified = true
		}
	} else if idJabatan == 29 { // Senior
		projectedRate = 0.35
		if poinTotal >= 46.0 {
			appliedRate = 0.45
			isQualified = true
		} else if poinTotal >= 35.0 {
			appliedRate = 0.30
			isQualified = true
		}
	} else { // Junior (ID 17 or default)
		projectedRate = 0.25
		if poinTotal >= 41.0 {
			appliedRate = 0.40
			isQualified = true
		} else if poinTotal >= 25.0 {
			appliedRate = 0.25
			isQualified = true
		}
	}

	bonusGramasi := 0.0
	if gramasiEmas >= 1000.0 {
		bonusGramasi = 3500000.0
	} else if gramasiEmas >= 500.0 {
		bonusGramasi = 2000000.0
	} else if gramasiEmas >= 300.0 {
		bonusGramasi = 1000000.0
	}

	bonusElektronik := 0.0
	if nonemas >= 50000000.0 {
		bonusElektronik = 500000.0
	}

	rateForIncentive := projectedRate
	if isQualified {
		rateForIncentive = appliedRate
	}
	mainIncentiveAmount := (projectedRate * totalAdminAlokasi) + bonusGramasi + bonusElektronik
	totalTakeHome := 0.0
	if isQualified {
		totalTakeHome = (appliedRate * totalAdminAlokasi) + bonusGramasi + bonusElektronik
	}
	estIncentive := mainIncentiveAmount
	if isQualified {
		estIncentive = totalTakeHome
	}

	targetGoalPoints := 25.0
	dailyTargetProgress := 0.0
	targetDescription := ""
	motivationText := ""

	if idJabatan == 6 { // Trainee
		targetGoalPoints = 20.0
		dailyTargetProgress = (poinTotal / 20.0) * 100.0
		if poinTotal >= 20.0 {
			targetDescription = "Target Qualified Tercapai"
			motivationText = "Target Qualified telah dicapai."
		} else {
			targetDescription = "Mengejar Qualified (Target: 20 Poin)"
			motivationText = fmt.Sprintf("Kurang %.1f Poin lagi untuk mencapai Qualified.", 20.0-poinTotal)
		}
	} else if idJabatan == 29 { // Senior
		if poinTotal < 35.0 {
			targetGoalPoints = 35.0
			dailyTargetProgress = (poinTotal / 35.0) * 100.0
			targetDescription = "Mengejar Qualified 1 (Target: 35 Poin)"
			motivationText = fmt.Sprintf("Kurang %.1f Poin lagi untuk Qualified 1.", 35.0-poinTotal)
		} else {
			targetGoalPoints = 46.0
			dailyTargetProgress = ((poinTotal - 35.0) / 11.0) * 100.0
			if poinTotal >= 46.0 {
				targetDescription = "Target Qualified 2 Tercapai"
				motivationText = "Target Qualified 2 telah dicapai."
			} else {
				targetDescription = "Q1 Selesai - Mengejar Qualified 2 (Target: 46 Poin)"
				motivationText = fmt.Sprintf("Kurang %.1f Poin lagi untuk Bonus 45%%.", 46.0-poinTotal)
			}
		}
	} else { // Junior (ID 17 or default)
		if poinTotal < 25.0 {
			targetGoalPoints = 25.0
			dailyTargetProgress = (poinTotal / 25.0) * 100.0
			targetDescription = "Mengejar Qualified 1 (Target: 25 Poin)"
			motivationText = fmt.Sprintf("Kurang %.1f Poin lagi untuk Qualified 1.", 25.0-poinTotal)
		} else {
			targetGoalPoints = 41.0
			dailyTargetProgress = ((poinTotal - 25.0) / 16.0) * 100.0
			if poinTotal >= 41.0 {
				targetDescription = "Target Qualified 2 Tercapai"
				motivationText = "Target Qualified 2 telah dicapai."
			} else {
				targetDescription = "Q1 Selesai - Mengejar Qualified 2 (Target: 41 Poin)"
				motivationText = fmt.Sprintf("Kurang %.1f Poin lagi untuk Bonus 40%%.", 41.0-poinTotal)
			}
		}
	}

	dailyTargetProgress = math.Min(100.0, math.Max(0.0, dailyTargetProgress))

	targetEmasG := 300.1
	pctEmas := 0.0
	bonusEmasTitle := "Bonus Rp 1.000.000 (Target 300.1g)"
	subEmasText := fmt.Sprintf("Butuh %.1fg lagi untuk klaim bonus Rp 1 Jt", math.Max(0, 300.1-gramasiEmas))
	isEmasBonusDone := false

	if gramasiEmas >= 1000.1 {
		targetEmasG = 1000.1
		pctEmas = 100.0
		bonusEmasTitle = "Bonus Rp 3.500.000 Tercapai"
		subEmasText = "Bonus Gramasi Emas Maksimal (1000g+) telah dicapai"
		isEmasBonusDone = true
	} else if gramasiEmas >= 500.0 {
		targetEmasG = 1000.1
		pctEmas = math.Min(100.0, (gramasiEmas/1000.1)*100.0)
		bonusEmasTitle = "Bonus Rp 2.000.000 Tercapai"
		subEmasText = fmt.Sprintf("Lanjut kejar 1000g untuk Bonus Rp 3.5 Jt (Kurang %.1fg)", 1000.1-gramasiEmas)
		isEmasBonusDone = true
	} else if gramasiEmas >= 300.0 {
		targetEmasG = 500.1
		pctEmas = math.Min(100.0, (gramasiEmas/500.1)*100.0)
		bonusEmasTitle = "Bonus Rp 1.000.000 Tercapai"
		subEmasText = fmt.Sprintf("Lanjut kejar 500g untuk Bonus Rp 2 Jt (Kurang %.1fg)", 500.1-gramasiEmas)
		isEmasBonusDone = true
	} else {
		targetEmasG = 300.1
		pctEmas = math.Min(100.0, (gramasiEmas/300.1)*100.0)
	}
	pctEmas = math.Round(pctEmas*10) / 10

	pctNonEmas := math.Min(100.0, (nonemas/50000000.0)*100.0)
	pctNonEmas = math.Round(pctNonEmas*10) / 10
	isNonEmasBonusDone := nonemas >= 50000000.0
	gapNonEmas := math.Max(0, 50000000.0-nonemas)
	nonEmasBonusTitle := "Bonus Rp 500.000 (Target 50 Jt)"
	nonEmasSubText := fmt.Sprintf("Butuh Rp %.0f lagi untuk klaim bonus", gapNonEmas)
	if isNonEmasBonusDone {
		nonEmasBonusTitle = "Bonus Rp 500.000 Tercapai"
		nonEmasSubText = "Bonus UP Elektronik 50 Jt telah dicapai"
	}

	leaderboardRes := getLeaderboardPosition(validNumericIDs, tgl1Str, tgl2Str)
	pacingDailyRes := getDailyPacingBooking(validNumericIDs, tgl1Str, tgl2Str, daysInMonth)

	// 6. Query Live Incentive Feed real-time from DB2 Postgres
	var liveFeeds []gin.H
	if config.DBPostgres != nil && len(validNumericIDs) > 0 {
		type TrxItem struct {
			DocNo           string    `gorm:"column:doc_no"`
			SubmissionDate  time.Time `gorm:"column:submission_date"`
			RequestedLoan   float64   `gorm:"column:requested_loan"`
			AdminCostAmount float64   `gorm:"column:admin_cost_amount"`
			OfferingName    string    `gorm:"column:offering_name"`
			ProductID       int       `gorm:"column:product_id"`
			CustomerName    string    `gorm:"column:full_name"`
		}
		var items []TrxItem
		_ = config.DBPostgres.Table("gadai.transactions t").
			Select("t.doc_no, t.submission_date, t.requested_loan, t.admin_cost_amount, o.offering_name, o.product_id, c.full_name").
			Joins("LEFT JOIN product.product_offerings o ON t.offering_id = o.offering_id").
			Joins("LEFT JOIN customer.customers c ON t.customer_id = c.customer_id").
			Where("t.sales_employee_id IN ? AND t.application_status IN ('NEW', 'ADO', 'REO') AND t.loan_status IN ('ACTIVE', 'REDEEM', 'LUNAS', 'EXTEND') AND t.submission_date BETWEEN ? AND ?", validNumericIDs, tgl1Str, tgl2Str).
			Order("t.submission_date DESC").
			Limit(10).
			Scan(&items).Error

		for _, item := range items {
			poin := item.RequestedLoan / 15000000.0
			if item.ProductID != 1 {
				poin = item.RequestedLoan / 2000000.0
			}
			adminAlokasi := item.AdminCostAmount * 0.45
			potensi := adminAlokasi * rateForIncentive

			iconText := "🥇"
			if item.ProductID != 1 {
				iconText = "💻"
			}

			liveFeeds = append(liveFeeds, gin.H{
				"id":          item.DocNo,
				"code":        item.DocNo,
				"date":        item.SubmissionDate.Format("02 Jan 2006, 15:04"),
				"type":        item.OfferingName,
				"customer":    item.CustomerName,
				"amount":      math.Round(potensi),
				"poin":        math.Round(poin*100) / 100,
				"sub":         fmt.Sprintf("%.1f Poin • Pinjaman Rp %.0f", poin, item.RequestedLoan),
				"loan_amount": item.RequestedLoan,
				"iconText":    iconText,
				"is_bonus":    false,
			})
		}
	}

	if len(liveFeeds) == 0 {
		liveFeeds = []gin.H{
			{"type": "Gadai Emas", "code": "TRX-8412", "customer": "Budi Santoso", "amount": 125000, "sub": "2 Poin • Pinjaman Rp 5.000.000", "iconText": "🥇", "isBonus": false},
			{"type": "Nasabah Baru", "code": "TRX-7741", "customer": "Siti Rahma", "amount": 25000, "sub": "New CIF Bonus", "iconText": "👤", "isBonus": false},
			{"type": "Gadai Elektronik", "code": "TRX-6109", "customer": "Ahmad Fauzi", "amount": 50000, "sub": "1 Poin • Laptop Asus", "iconText": "💻", "isBonus": false},
			{"type": "Bonus Gramasi", "code": "BONUS-101", "customer": "Achievement 300g", "amount": 1000000, "sub": "Target Emas Level 1", "iconText": "🏆", "isBonus": true},
		}
	}

	// Motivational Carousel Labels
	ptsNeeded := math.Max(0, targetGoalPoints-poinTotal)
	ptsNeededFmt := fmt.Sprintf("%.1f", ptsNeeded)
	labelPoint := gin.H{
		"icon": "🎯",
		"text": fmt.Sprintf("Kejar %s Poin lagi sebelum akhir periode agar insentif cair!", ptsNeededFmt),
	}
	if isQualified {
		labelPoint = gin.H{
			"icon": "🎉",
			"text": "Target Qualified Tercapai! Tambah transaksi untuk bonus maksimal.",
		}
	}

	userRank := 1
	if rVal, ok := leaderboardRes["rank"].(int); ok {
		userRank = rVal
	}
	labelLeaderboardText := "Peringkat #1 Leaderboard! Pertahankan posisimu di puncak."
	if userRank > 1 {
		labelLeaderboardText = fmt.Sprintf("Peringkat #%d! Tambah transaksi untuk naik posisi di Leaderboard.", userRank)
	}
	labelLeaderboard := gin.H{
		"icon": "🏆",
		"text": labelLeaderboardText,
	}

	c.JSON(http.StatusOK, gin.H{
		"counts":                counts,
		"newcif":                newcif,
		"last_month_cif":        lastMonthCif,
		"leads":                 counts["leads"],
		"leads_monthly":         leadsMonthly,
		"prospect_monthly":      prospectMonthly,
		"hotprospect_monthly":   hotprospectMonthly,
		"sbg_monthly":           sbgMonthly,
		"last_month_leads":      lastMonthLeads,
		"emas":                  emas,
		"nonemas":               nonemas,
		"gramasi_emas":          gramasiEmas,
		"total_trx_count":       totalTrxCount,
		"est_incentive":         estIncentive,
		"main_incentive_amount": mainIncentiveAmount,
		"total_take_home":       totalTakeHome,
		"is_qualified":          isQualified,
		"user_grade":            nmJabatan,
		"applied_rate":          int(appliedRate * 100),
		"poin_total":            poinTotalFormatted,
		"target_goal_points":    int(targetGoalPoints),
		"daily_target_progress": math.Round(dailyTargetProgress),
		"target_description":    targetDescription,
		"motivation_text":       motivationText,
		"is_review":             isReview,
		"sales_list":            salesList,
		"selected_sales_user":   userStr,
		"label_point":           labelPoint,
		"label_leaderboard":     labelLeaderboard,
		"emas_card": gin.H{
			"target_emas_g":      targetEmasG,
			"pct_emas":           pctEmas,
			"bonus_emas_title":   bonusEmasTitle,
			"sub_emas_text":      subEmasText,
			"is_emas_bonus_done": isEmasBonusDone,
		},
		"nonemas_card": gin.H{
			"pct_nonemas":           pctNonEmas,
			"gap_nonemas":           gapNonEmas,
			"bonus_title":           nonEmasBonusTitle,
			"sub_text":              nonEmasSubText,
			"is_nonemas_bonus_done": isNonEmasBonusDone,
		},
		"leaderboard":  leaderboardRes,
		"pacing_daily": pacingDailyRes,
		"live_feeds":   liveFeeds,
		"username":     userStr,
	})
}

func SumberSourceList(c *gin.Context) {
	fmt.Println("testingg")
	if config.DBPostgres == nil {
		c.JSON(http.StatusOK, []gin.H{})
		return
	}

	var results []models.ApplicationSource

	err := config.DBPostgres.Table("master.application_sources").
		Select("code, name").
		Where("is_active = ?", true).
		Order("name ASC").
		Scan(&results).Error
	if err != nil {
		c.JSON(http.StatusOK, []gin.H{})
		return
	}
	c.JSON(http.StatusOK, results)

}

// GetCustomerSelectHandler returns list of customers for select dropdown matching M_rencana getSelectCustomer
func GetCustomerSelectHandler(c *gin.Context) {
	userId, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	search := c.Query("q")
	type CustomerSelect struct {
		ID   int    `json:"id"`
		Text string `json:"text"`
	}

	if config.DBMysql == nil {
		c.JSON(http.StatusOK, []CustomerSelect{})
		return
	}

	var customers []models.Customer
	query := config.DBMysql.Where("id_login = ?", userId)
	if search != "" {
		likeTerm := "%" + search + "%"
		query = query.Where("(nama LIKE ? OR hp LIKE ? OR identitas LIKE ?)", likeTerm, likeTerm, likeTerm)
	}

	err := query.Order("nama ASC").Limit(10).Find(&customers).Error
	if err != nil {
		c.JSON(http.StatusOK, []CustomerSelect{})
		return
	}

	var results []CustomerSelect
	for _, cust := range customers {
		results = append(results, CustomerSelect{
			ID:   cust.IDCustomer,
			Text: fmt.Sprintf("%s | %s", cust.Nama, cust.Hp),
		})
	}

	c.JSON(http.StatusOK, results)
}

// GetSumberCustHandler fetches reference sources from master.reference_sources (DB2 Postgres)
func GetSumberCustHandler(c *gin.Context) {
	if config.DBPostgres == nil {
		c.JSON(http.StatusOK, []gin.H{})
		return
	}

	var results []struct {
		KdRefCust string `json:"kd_ref_cust" gorm:"column:kd_ref_cust"`
		NmRefCust string `json:"nm_ref_cust" gorm:"column:nm_ref_cust"`
	}

	err := config.DBPostgres.Table("master.reference_sources").
		Select("code AS kd_ref_cust, name AS nm_ref_cust").
		Where("is_active = ?", true).
		Order("name ASC").
		Scan(&results).Error

	if err != nil {
		c.JSON(http.StatusOK, []gin.H{})
		return
	}

	c.JSON(http.StatusOK, results)
}

// GetSalesListHandler returns active sales users for review mode
func GetSalesListHandler(c *gin.Context) {
	if config.DBMysqlPos == nil {
		c.JSON(http.StatusOK, []gin.H{})
		return
	}

	type SalesItem struct {
		Username  string `json:"username" gorm:"column:username"`
		Nama      string `json:"nama" gorm:"column:nama"`
		IDJabatan int    `json:"id_jabatan" gorm:"column:id_jabatan"`
		NmJabatan string `json:"nm_jabatan" gorm:"column:nm_jabatan"`
		FkUser    string `json:"fkuser" gorm:"column:username"`
	}

	var sales []SalesItem
	err := config.DBMysqlPos.Table("tbl_user u").
		Select("u.username, u.nama, u.id_jabatan, j.nm_jabatan").
		Joins("LEFT JOIN mst_jabatan j ON j.id_jabatan = u.id_jabatan").
		Where("u.void = ? AND u.id_jabatan IN ?", 1, []int{6, 17, 29, 74}).
		Order("u.nama ASC").
		Scan(&sales).Error

	if err != nil {
		c.JSON(http.StatusOK, []gin.H{})
		return
	}

	c.JSON(http.StatusOK, sales)
}

// GetKelurahanHandler searches village/district/city from master.village
func GetKelurahanHandler(c *gin.Context) {
	if config.DBPostgres == nil {
		c.JSON(http.StatusOK, []gin.H{})
		return
	}

	q := c.Query("q")
	type KelurahanResult struct {
		KdKelurahan int    `json:"kd_kelurahan" gorm:"column:kd_kelurahan"`
		NmKelurahan string `json:"nm_kelurahan" gorm:"column:nm_kelurahan"`
		KdKecamatan int    `json:"kd_kecamatan" gorm:"column:kd_kecamatan"`
		NmKecamatan string `json:"nm_kecamatan" gorm:"column:nm_kecamatan"`
		KdKota      int    `json:"kd_kota" gorm:"column:kd_kota"`
		NmKota      string `json:"nm_kota" gorm:"column:nm_kota"`
	}

	var results []KelurahanResult
	query := config.DBPostgres.Table("master.village v").
		Select("v.village_id AS kd_kelurahan, v.village_name AS nm_kelurahan, d.district_id AS kd_kecamatan, d.district_name AS nm_kecamatan, c.city_id AS kd_kota, c.city_name AS nm_kota").
		Joins("LEFT JOIN master.district d ON d.district_id = v.district_id").
		Joins("LEFT JOIN master.city c ON c.city_id = d.city_id").
		Where("v.is_active = ?", true)

	if q != "" {
		upperQ := "%" + strings.ToUpper(q) + "%"
		query = query.Where("UPPER(v.village_name) LIKE ?", upperQ)
	}

	err := query.Order("v.village_name ASC").Limit(20).Scan(&results).Error
	if err != nil {
		c.JSON(http.StatusOK, []gin.H{})
		return
	}

	c.JSON(http.StatusOK, results)
}
