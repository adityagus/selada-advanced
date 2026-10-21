package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"go-api/config"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CustomerListItem represents a full customer record with latest inquiry status
type CustomerListItem struct {
	IDCustomer     int    `json:"id_customer" gorm:"column:id_customer"`
	Cif            string `json:"cif" gorm:"column:identitas"`
	Nama           string `json:"nama" gorm:"column:nama"`
	Hp             string `json:"hp" gorm:"column:hp"`
	Alamat         string `json:"alamat" gorm:"column:alamat"`
	Kota           string `json:"kota" gorm:"column:kota"`
	Aktivitas      string `json:"aktivitas" gorm:"column:aktivitas"`
	Jenis          string `json:"jenis" gorm:"column:jenis"`
	IDLogin        string `json:"id_login" gorm:"column:id_login"`
	IDSumbercust   string `json:"id_sumbercust" gorm:"column:id_sumbercust"`
	LatestStatusID int    `json:"latest_status_id" gorm:"column:latest_status_id"`
	LatestStatus   string `json:"latest_status" gorm:"column:latest_status"`
	LatestTgl      string `json:"latest_tgl" gorm:"column:latest_tgl"`
	TotalInquiries int64  `json:"total_inquiries" gorm:"column:total_inquiries"`
}

// CustomerInquiryHistory represents inquiry records for a customer
type CustomerInquiryHistory struct {
	IDRencana    int    `json:"id_rencana" gorm:"column:id_rencana"`
	IDRencanaDet int    `json:"id_rencana_det" gorm:"column:id_rencana_det"`
	TglRencana   string `json:"tgl_rencana" gorm:"column:tgl_rencana"`
	IDHslaktiv   int    `json:"id_hslaktiv" gorm:"column:id_hslaktiv"`
	NamaHslaktiv string `json:"nama_hslaktiv" gorm:"column:nama_hslaktiv"`
	KetRencana   string `json:"ket_rencana" gorm:"column:ket_rencana"`
	KetAktivitas string `json:"ket_aktivitas" gorm:"column:ket_aktivitas"`
	IDSumbercust string `json:"id_sumbercust" gorm:"column:id_sumbercust"`
	NoSbg        string `json:"no_sbg" gorm:"column:no_sbg"`
}

// GetCustomerListHandler fetches Customer & CIF records from database
func GetCustomerListHandler(c *gin.Context) {
	userId, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if config.DBSelada == nil {
		c.JSON(http.StatusOK, gin.H{
			"data":            []CustomerListItem{},
			"recordsTotal":    0,
			"recordsFiltered": 0,
		})
		return
	}

	search := strings.TrimSpace(c.Query("search"))
	startStr := c.Query("start")
	lengthStr := c.Query("length")

	start := 0
	length := 100

	if s, err := strconv.Atoi(startStr); err == nil && s >= 0 {
		start = s
	}
	if l, err := strconv.Atoi(lengthStr); err == nil && l > 0 {
		length = l
	}

	targetUser := userId.(string)
	if salesParam := c.Query("sales_user"); salesParam != "" {
		targetUser = salesParam
	} else if userParam := c.Query("user"); userParam != "" {
		targetUser = userParam
	}

	baseQuery := config.DBSelada.Table("tb_customer c").
		Select(`
			c.id_customer,
			COALESCE(c.identitas, '') as identitas,
			c.nama,
			COALESCE(c.hp, '') as hp,
			COALESCE(c.alamat, '') as alamat,
			COALESCE(c.kota, '') as kota,
			COALESCE(c.aktivitas, '') as aktivitas,
			COALESCE(c.jenis, '') as jenis,
			COALESCE(c.id_login, '') as id_login,
			COALESCE(r.id_sumbercust, '') as id_sumbercust,
			COALESCE(rd.id_hslaktiv, 1) as latest_status_id,
			COALESCE(ha.nama_hslaktiv, '') as latest_status,
			COALESCE(CAST(rd.tgl_rencana AS TEXT), '') as latest_tgl,
			COALESCE(inq_cnt.cnt, 0) as total_inquiries
		`).
		Joins(`LEFT JOIN (
			SELECT r1.id_customer, MAX(r1.id_rencana) as max_id_rencana
			FROM tb_rencana r1
			GROUP BY r1.id_customer
		) latest_r ON latest_r.id_customer = c.id_customer`).
		Joins("LEFT JOIN tb_rencana r ON r.id_rencana = latest_r.max_id_rencana").
		Joins(`LEFT JOIN (
			SELECT rd1.id_rencana, MAX(rd1.id_rencana_det) as max_id_det
			FROM tb_rencana_det rd1
			GROUP BY rd1.id_rencana
		) latest_rd ON latest_rd.id_rencana = r.id_rencana`).
		Joins("LEFT JOIN tb_rencana_det rd ON rd.id_rencana_det = latest_rd.max_id_det").
		Joins("LEFT JOIN tb_hslaktivitas ha ON ha.id_hslaktiv = rd.id_hslaktiv").
		Joins(`LEFT JOIN (
			SELECT id_customer, COUNT(id_rencana) as cnt
			FROM tb_rencana
			GROUP BY id_customer
		) inq_cnt ON inq_cnt.id_customer = c.id_customer`)

	uLower := strings.ToLower(targetUser)
	if uLower != "all" && uLower != "*" && uLower != "admin" && uLower != "superadmin" && uLower != "administrator" {
		baseQuery = baseQuery.Where("c.id_login = ? OR c.id_login IS NULL OR c.id_login = ''", targetUser)
	}

	if search != "" {
		likeTerm := "%" + search + "%"
		baseQuery = baseQuery.Where("(c.nama ILIKE ? OR c.hp ILIKE ? OR c.identitas ILIKE ? OR c.alamat ILIKE ? OR c.kota ILIKE ?)", likeTerm, likeTerm, likeTerm, likeTerm, likeTerm)
	}

	var totalRecords int64
	if err := baseQuery.Session(&gorm.Session{}).Count(&totalRecords).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghitung total customer"})
		return
	}

	var items []CustomerListItem
	err := baseQuery.Order("c.id_customer DESC").Offset(start).Limit(length).Scan(&items).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil list customer"})
		return
	}

	// Status fallback check
	for i := range items {
		if items[i].LatestStatus == "" {
			switch items[i].LatestStatusID {
			case 1:
				items[i].LatestStatus = "Leads"
			case 2:
				items[i].LatestStatus = "Prospect"
			case 3:
				items[i].LatestStatus = "Hot Prospect"
			case 7:
				items[i].LatestStatus = "SBG"
			case 6:
				items[i].LatestStatus = "Batal"
			default:
				items[i].LatestStatus = "Leads"
			}
		}
		if items[i].Cif == "" {
			items[i].Cif = fmt.Sprintf("CIF-%06d", items[i].IDCustomer)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data":            items,
		"recordsTotal":    totalRecords,
		"recordsFiltered": totalRecords,
	})
}

// GetCustomerInquiriesHandler returns all inquiries from DB for a given customer
func GetCustomerInquiriesHandler(c *gin.Context) {
	_, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	idParam := c.Param("id")
	idCust, err := strconv.Atoi(idParam)
	if err != nil || idCust <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer ID"})
		return
	}

	if config.DBSelada == nil {
		c.JSON(http.StatusOK, []CustomerInquiryHistory{})
		return
	}

	var history []CustomerInquiryHistory
	err = config.DBSelada.Table("tb_rencana r").
		Select(`
			r.id_rencana,
			rd.id_rencana_det,
			COALESCE(CAST(rd.tgl_rencana AS TEXT), '') as tgl_rencana,
			COALESCE(rd.id_hslaktiv, 1) as id_hslaktiv,
			COALESCE(ha.nama_hslaktiv, '') as nama_hslaktiv,
			COALESCE(rd.ket_rencana, '') as ket_rencana,
			COALESCE(rd.ket_aktivitas, '') as ket_aktivitas,
			COALESCE(r.id_sumbercust, '') as id_sumbercust,
			COALESCE(rd.no_do, '') as no_sbg
		`).
		Joins("JOIN tb_rencana_det rd ON rd.id_rencana = r.id_rencana").
		Joins("LEFT JOIN tb_hslaktivitas ha ON ha.id_hslaktiv = rd.id_hslaktiv").
		Where("r.id_customer = ?", idCust).
		Order("rd.id_rencana_det DESC").
		Scan(&history).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil riwayat inquiry"})
		return
	}

	for i := range history {
		if history[i].NamaHslaktiv == "" {
			switch history[i].IDHslaktiv {
			case 1:
				history[i].NamaHslaktiv = "Leads"
			case 2:
				history[i].NamaHslaktiv = "Prospect"
			case 3:
				history[i].NamaHslaktiv = "Hot Prospect"
			case 7:
				history[i].NamaHslaktiv = "SBG"
			case 6:
				history[i].NamaHslaktiv = "Batal"
			default:
				history[i].NamaHslaktiv = "Leads"
			}
		}
	}

	c.JSON(http.StatusOK, history)
}
