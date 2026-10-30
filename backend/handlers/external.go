package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"go-api/config"
	"go-api/models"

	"github.com/gin-gonic/gin"
)

// ValidateExternalAPIKey checks for valid authorization token in header, query, or payload
func ValidateExternalAPIKey(c *gin.Context, payloadKey string) bool {
	// If invoked via webhook endpoint, allow webhook access
	if strings.HasPrefix(c.Request.URL.Path, "/api/webhook") {
		return true
	}

	expectedKey := os.Getenv("EXTERNAL_API_KEY")
	if expectedKey == "" {
		expectedKey = "selada_external_secret_2026" // Default fallback secret
	}

	headerKey := strings.TrimSpace(c.GetHeader("X-Api-Key"))
	if headerKey == "" {
		headerKey = strings.TrimSpace(c.GetHeader("Authorization"))
		headerKey = strings.TrimPrefix(headerKey, "Bearer ")
		headerKey = strings.TrimSpace(headerKey)
	}

	providedKey := headerKey
	if providedKey == "" {
		providedKey = strings.TrimSpace(payloadKey)
	}
	if providedKey == "" {
		providedKey = strings.TrimSpace(c.Query("api_key"))
		if providedKey == "" {
			providedKey = strings.TrimSpace(c.Query("token"))
		}
	}

	return providedKey == expectedKey
}

// BookingWebhookVerifyHandler handles GET /api/webhook/booking
func BookingWebhookVerifyHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":      "ok",
		"service":     "Selada Booking Webhook Listener",
		"timestamp":   time.Now().Format(time.RFC3339),
		"description": "Ready to receive external booking events and update customer status to BOOKING",
	})
}

// UpdateCustomerPhoneRequest represents the third party payload for phone updates
type UpdateCustomerPhoneRequest struct {
	ApiKey     string `json:"api_key"`
	IDCustomer int    `json:"id_customer"`
	IDRencana  int    `json:"id_rencana"`
	OldPhone   string `json:"old_phone"`
	NewPhone   string `json:"new_phone" binding:"required"`
	Reason     string `json:"reason"`
	UpdatedBy  string `json:"updated_by"`
}

// UpdateCustomerPhoneHandler handles POST /api/external/customer/update-phone
// Allows third-party systems to update customer phone number and resets WA validation
func UpdateCustomerPhoneHandler(c *gin.Context) {
	if config.DBSelada == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database Selada tidak terhubung"})
		return
	}

	var req UpdateCustomerPhoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Payload tidak valid. Pastikan new_phone terisi: " + err.Error(),
		})
		return
	}

	if !ValidateExternalAPIKey(c, req.ApiKey) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Akses ditolak: API Key tidak valid atau tidak disediakan",
		})
		return
	}

	cleanNewPhone := CleanPhoneNumber(req.NewPhone)
	if cleanNewPhone == "" || len(cleanNewPhone) < 9 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Nomor telepon baru tidak valid: " + req.NewPhone,
		})
		return
	}

	cleanOldPhone := CleanPhoneNumber(req.OldPhone)

	// Step 1: Find Customer
	var customer models.Customer
	var err error

	if req.IDCustomer > 0 {
		err = config.DBSelada.Where("id_customer = ?", req.IDCustomer).First(&customer).Error
	} else if req.IDRencana > 0 {
		err = config.DBSelada.Table("tb_customer c").
			Joins("JOIN tb_rencana r ON r.id_customer = c.id_customer").
			Where("r.id_rencana = ?", req.IDRencana).
			Select("c.*").First(&customer).Error
	} else if cleanOldPhone != "" {
		phoneSuffix := cleanOldPhone
		if len(phoneSuffix) > 8 {
			phoneSuffix = phoneSuffix[len(phoneSuffix)-8:]
		}
		err = config.DBSelada.Where("hp LIKE ?", "%"+phoneSuffix+"%").
			Order("id_customer DESC").First(&customer).Error
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Harus menyertakan salah satu: id_customer, id_rencana, atau old_phone untuk identifikasi nasabah",
		})
		return
	}

	if err != nil || customer.IDCustomer == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Data customer tidak ditemukan di database Selada",
		})
		return
	}

	oldPhoneRecord := customer.Hp
	now := time.Now()

	// Step 2: Update customer phone number and reset status_wa
	// When phone number changes, status_wa is reset to 'belum_validasi' for re-verification
	updateData := map[string]interface{}{
		"hp":              cleanNewPhone,
		"status_wa":       "belum_validasi",
		"wa_validated_at": nil,
	}

	if err := config.DBSelada.Table("tb_customer").
		Where("id_customer = ?", customer.IDCustomer).
		Updates(updateData).Error; err != nil {
		log.Printf("[API UPDATE PHONE ERROR] Failed to update customer #%d: %v", customer.IDCustomer, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengupdate nomor telepon customer di database",
		})
		return
	}

	// Step 3: Record change in active inquiry detail (tb_rencana_det)
	var latestRencana models.Rencana
	_ = config.DBSelada.Where("id_customer = ?", customer.IDCustomer).
		Order("id_rencana DESC").First(&latestRencana)

	if latestRencana.IDRencana > 0 {
		var lastDet models.RencanaDet
		_ = config.DBSelada.Where("id_rencana = ?", latestRencana.IDRencana).
			Order("id_rencana_det DESC").First(&lastDet)

		sourceNote := req.UpdatedBy
		if sourceNote == "" {
			sourceNote = "Pihak Ketiga"
		}
		reasonNote := req.Reason
		if reasonNote == "" {
			reasonNote = "Permintaan perubahan nomor"
		}

		changeActivity := fmt.Sprintf("[API Pihak Ketiga - %s] Nomor telepon nasabah diperbarui dari %s menjadi %s. Alasan: %s. Status WA direset.",
			sourceNote, oldPhoneRecord, cleanNewPhone, reasonNote)

		currentStatus := lastDet.IDHslaktiv
		if currentStatus == 0 {
			currentStatus = 1
		}

		newDet := models.RencanaDet{
			IDRencana:    latestRencana.IDRencana,
			TglRencana:   now,
			IDHslaktiv:   currentStatus,
			Status:       1,
			KetRencana:   lastDet.KetRencana,
			KetAktivitas: changeActivity,
		}
		_ = config.DBSelada.Create(&newDet)
	}

	log.Printf("[API UPDATE PHONE SUCCESS] Customer #%d (%s) phone changed from %s to %s",
		customer.IDCustomer, customer.Nama, oldPhoneRecord, cleanNewPhone)

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Nomor telepon customer berhasil diperbarui dan status WA direset ke belum_validasi",
		"data": gin.H{
			"id_customer": customer.IDCustomer,
			"nama":        customer.Nama,
			"old_phone":   oldPhoneRecord,
			"new_phone":   cleanNewPhone,
			"status_wa":   "belum_validasi",
			"updated_at":  now.Format("2006-01-02 15:04:05"),
		},
	})
}

// ExternalBookingDetails represents booking transaction info from third-party or web form
type ExternalBookingDetails struct {
	TanggalBooking   string  `json:"tanggal_booking"`
	JenisBarang      string  `json:"jenis_barang"`
	EstimasiPinjaman float64 `json:"estimasi_pinjaman"`
	CabangTujuan     string  `json:"cabang_tujuan"`
	Catatan          string  `json:"catatan"`
}

// ExternalBookingRequest represents payload from API / Webhook when a booking succeeds
type ExternalBookingRequest struct {
	ApiKey         string                 `json:"api_key"`
	Token          string                 `json:"token"`
	Event          string                 `json:"event"` // e.g. "booking.success", "booking_success"
	ReferenceID    string                 `json:"reference_id"`
	BookingID      string                 `json:"booking_id"`
	BookingCode    string                 `json:"booking_code"`
	IDRencana      int                    `json:"id_rencana"`
	IDCustomer     int                    `json:"id_customer"`
	Phone          string                 `json:"phone"`
	HP             string                 `json:"hp"`
	NoHP           string                 `json:"no_hp"`
	Nama           string                 `json:"nama"`
	NamaCustomer   string                 `json:"nama_customer"`
	CustomerName   string                 `json:"customer_name"`
	Alamat         string                 `json:"alamat"`
	Kota           string                 `json:"kota"`
	Source         string                 `json:"source"`
	BookingDetails ExternalBookingDetails `json:"booking_details"`

	// Direct top-level fields fallback if not nested under booking_details
	TanggalBooking   string  `json:"tanggal_booking"`
	Jadwal           string  `json:"jadwal"`
	JenisBarang      string  `json:"jenis_barang"`
	Barang           string  `json:"barang"`
	EstimasiPinjaman float64 `json:"estimasi_pinjaman"`
	Nominal          float64 `json:"nominal"`
	CabangTujuan     string  `json:"cabang_tujuan"`
	Cabang           string  `json:"cabang"`
	Catatan          string  `json:"catatan"`
	Notes            string  `json:"notes"`
}

// ExternalBookingHandler handles POST /api/external/booking and POST /api/webhook/booking
// Automatically updates inquiry status to BOOKING (IDHslaktiv = 7) when a booking succeeds.
// If the customer or inquiry is new, it creates them automatically so no booking is lost.
func ExternalBookingHandler(c *gin.Context) {
	if config.DBSelada == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database Selada tidak terhubung"})
		return
	}

	var req ExternalBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Payload booking tidak valid: " + err.Error(),
		})
		return
	}

	payloadAuthKey := req.ApiKey
	if payloadAuthKey == "" {
		payloadAuthKey = req.Token
	}

	if !ValidateExternalAPIKey(c, payloadAuthKey) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Akses ditolak: API Key tidak valid atau tidak disediakan",
		})
		return
	}

	// Resolve phone
	rawPhone := strings.TrimSpace(req.Phone)
	if rawPhone == "" {
		rawPhone = strings.TrimSpace(req.HP)
	}
	if rawPhone == "" {
		rawPhone = strings.TrimSpace(req.NoHP)
	}
	cleanPhone := CleanPhoneNumber(rawPhone)

	if cleanPhone == "" && req.IDRencana == 0 && req.IDCustomer == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Nomor handphone nasabah atau ID inquiry wajib disertakan",
		})
		return
	}

	// Resolve customer name
	custName := strings.TrimSpace(req.Nama)
	if custName == "" {
		custName = strings.TrimSpace(req.NamaCustomer)
	}
	if custName == "" {
		custName = strings.TrimSpace(req.CustomerName)
	}
	if custName == "" {
		custName = "Nasabah Booking"
	}

	// Resolve reference ID
	refStr := strings.TrimSpace(req.ReferenceID)
	if refStr == "" {
		refStr = strings.TrimSpace(req.BookingID)
	}
	if refStr == "" {
		refStr = strings.TrimSpace(req.BookingCode)
	}
	if refStr == "" {
		refStr = fmt.Sprintf("BK-%s", time.Now().Format("0601021504"))
	}

	// Resolve branch
	branchStr := strings.TrimSpace(req.BookingDetails.CabangTujuan)
	if branchStr == "" {
		branchStr = strings.TrimSpace(req.CabangTujuan)
	}
	if branchStr == "" {
		branchStr = strings.TrimSpace(req.Cabang)
	}
	if branchStr == "" {
		branchStr = "Kantor / Cabang Terdekat"
	}

	// Resolve date
	tglBookingStr := strings.TrimSpace(req.BookingDetails.TanggalBooking)
	if tglBookingStr == "" {
		tglBookingStr = strings.TrimSpace(req.TanggalBooking)
	}
	if tglBookingStr == "" {
		tglBookingStr = strings.TrimSpace(req.Jadwal)
	}
	if tglBookingStr == "" {
		tglBookingStr = time.Now().Format("2006-01-02 15:04:05")
	}

	// Resolve item
	itemStr := strings.TrimSpace(req.BookingDetails.JenisBarang)
	if itemStr == "" {
		itemStr = strings.TrimSpace(req.JenisBarang)
	}
	if itemStr == "" {
		itemStr = strings.TrimSpace(req.Barang)
	}
	if itemStr == "" {
		itemStr = "Barang Gadai"
	}

	// Resolve estimate
	estAmount := req.BookingDetails.EstimasiPinjaman
	if estAmount == 0 {
		estAmount = req.EstimasiPinjaman
	}
	if estAmount == 0 {
		estAmount = req.Nominal
	}

	// Resolve note
	noteStr := strings.TrimSpace(req.BookingDetails.Catatan)
	if noteStr == "" {
		noteStr = strings.TrimSpace(req.Catatan)
	}
	if noteStr == "" {
		noteStr = strings.TrimSpace(req.Notes)
	}

	now := time.Now()

	// Step 1: Find or Create Customer
	var customer models.Customer
	if req.IDCustomer > 0 {
		_ = config.DBSelada.Where("id_customer = ?", req.IDCustomer).First(&customer).Error
	} else if cleanPhone != "" {
		phoneSuffix := cleanPhone
		if len(phoneSuffix) > 8 {
			phoneSuffix = phoneSuffix[len(phoneSuffix)-8:]
		}
		_ = config.DBSelada.Where("hp LIKE ?", "%"+phoneSuffix+"%").
			Order("id_customer DESC").First(&customer).Error
	}

	if customer.IDCustomer == 0 {
		// Auto-create new customer from booking
		customer = models.Customer{
			Nama:          custName,
			Hp:            cleanPhone,
			StatusWa:      "tervalidasi",
			WaValidatedAt: &now,
			Alamat:        req.Alamat,
			Kota:          req.Kota,
			IDLogin:       "system_booking",
			Aktivitas:     "1",
		}
		if err := config.DBSelada.Create(&customer).Error; err != nil {
			log.Printf("[API BOOKING] Gagal membuat customer baru: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Gagal menyimpan customer baru: " + err.Error(),
			})
			return
		}
		log.Printf("[API BOOKING] Customer baru berhasil dibuat: #%d (%s - %s)", customer.IDCustomer, customer.Nama, customer.Hp)
	} else {
		// Validasi nomor WhatsApp customer yang berhasil booking
		updates := map[string]interface{}{
			"status_wa":       "tervalidasi",
			"wa_validated_at": now,
		}
		if (customer.Nama == "" || customer.Nama == "Customer") && custName != "Nasabah Booking" {
			updates["nama"] = custName
			customer.Nama = custName
		}
		_ = config.DBSelada.Model(&customer).Updates(updates)
	}

	// Step 2: Find or Create Inquiry (Rencana)
	var rencana models.Rencana
	if req.IDRencana > 0 {
		_ = config.DBSelada.Where("id_rencana = ?", req.IDRencana).First(&rencana).Error
	}
	if rencana.IDRencana == 0 {
		_ = config.DBSelada.Where("id_customer = ?", customer.IDCustomer).
			Order("id_rencana DESC").First(&rencana).Error
	}

	previousStatus := 1
	sourceName := strings.TrimSpace(req.Source)
	if sourceName == "" {
		sourceName = "Web Booking"
	}

	if rencana.IDRencana == 0 {
		rencana = models.Rencana{
			IDCustomer:   customer.IDCustomer,
			IDLogin:      "system_booking",
			IDSumbercust: sourceName,
		}
		if err := config.DBSelada.Create(&rencana).Error; err != nil {
			log.Printf("[API BOOKING] Gagal membuat rencana: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Gagal membuat inquiry baru: " + err.Error(),
			})
			return
		}
		log.Printf("[API BOOKING] Rencana baru dibuat: #%d untuk customer #%d", rencana.IDRencana, customer.IDCustomer)
	} else {
		// Cek status terakhir
		var lastDet models.RencanaDet
		err := config.DBSelada.Where("id_rencana = ?", rencana.IDRencana).
			Order("id_rencana_det DESC").First(&lastDet).Error
		if err == nil && lastDet.IDHslaktiv > 0 {
			previousStatus = lastDet.IDHslaktiv
		}
	}

	// Step 3: Insert detail dengan status BOOKING (7)
	actionNote := fmt.Sprintf("[API Hit Booking #%s] Booking berhasil via Webhook/API. Cabang: %s, Jadwal: %s. Barang: %s, Est Pinjaman: Rp %.0f. Catatan: %s",
		refStr, branchStr, tglBookingStr, itemStr, estAmount, noteStr)

	newDet := models.RencanaDet{
		IDRencana:    rencana.IDRencana,
		TglRencana:   now,
		IDHslaktiv:   7, // Status Booking / SBG Deal
		Status:       1,
		KetRencana:   noteStr,
		KetAktivitas: actionNote,
	}

	if err := config.DBSelada.Create(&newDet).Error; err != nil {
		log.Printf("[API BOOKING ERROR] Failed to create detail: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal menyimpan status booking ke database",
		})
		return
	}

	log.Printf("[API BOOKING SUCCESS] Rencana #%d (Cust: %s - %s) status berhasil diubah ke BOOKING (7) dari status %d",
		rencana.IDRencana, customer.Nama, customer.Hp, previousStatus)

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Status inquiry berhasil diubah ke BOOKING",
		"data": gin.H{
			"id_rencana":      rencana.IDRencana,
			"id_customer":     customer.IDCustomer,
			"nama_customer":   customer.Nama,
			"hp":              customer.Hp,
			"reference_id":    refStr,
			"previous_status": previousStatus,
			"current_status":  7,
			"status_label":    "Booking / SBG",
			"status_wa":       "tervalidasi",
			"booking_details": gin.H{
				"cabang_tujuan":     branchStr,
				"tanggal_booking":   tglBookingStr,
				"jenis_barang":      itemStr,
				"estimasi_pinjaman": estAmount,
				"catatan":           noteStr,
			},
			"note":       actionNote,
			"updated_at": now.Format("2006-01-02 15:04:05"),
		},
	})
}
