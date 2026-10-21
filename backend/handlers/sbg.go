package handlers

import (
	"fmt"
	"math"
	"net/http"
	"strings"
	"time"

	"go-api/config"

	"github.com/gin-gonic/gin"
)

// SopigaTransactionItem represents transaction data retrieved from Sopiga Pro (PostgreSQL)
type SopigaTransactionItem struct {
	TransactionID   int     `json:"transaction_id"`
	DocNo           string  `json:"no_sbg"`
	CustomerName    string  `json:"nama_customer"`
	CustomerPhone   string  `json:"hp"`
	CustomerNIK     string  `json:"nik"`
	CustomerAlamat  string  `json:"alamat"`
	ProductName     string  `json:"nama_produk"`
	ProductID       int     `json:"product_id"`
	JenisBarang     string  `json:"jenis_barang"`
	DeskripsiBarang string  `json:"deskripsi_barang"`
	NominalUP       float64 `json:"nominal_up"`
	BiayaAdmin      float64 `json:"biaya_admin"`
	TglPencairan    string  `json:"tgl_pencairan"`
	TglJatuhTempo   string  `json:"tgl_jatuh_tempo"`
	SelisihHari     int     `json:"selisih_hari"`
	LoanStatus      string  `json:"loan_status"`
	AppStatus       string  `json:"app_status"`
}

// SearchSopigaTransactionsHandler handles GET /api/sopiga/transactions
// Searches and gets SBG transaction records from Sopiga Pro (PostgreSQL db_gadaimulia)
func SearchSopigaTransactionsHandler(c *gin.Context) {
	if config.DBPostgres == nil {
		c.JSON(http.StatusOK, gin.H{
			"data":    []SopigaTransactionItem{},
			"message": "Database Sopiga Pro (Postgres) tidak terhubung",
		})
		return
	}

	searchQuery := strings.TrimSpace(c.Query("q"))
	noSbgQuery := strings.TrimSpace(c.Query("no_sbg"))
	if noSbgQuery != "" && searchQuery == "" {
		searchQuery = noSbgQuery
	}

	// 1. Inspect column names of gadai.transactions safely
	var txCols []string
	_ = config.DBPostgres.Raw(`
		SELECT column_name 
		FROM information_schema.columns 
		WHERE table_schema = 'gadai' AND table_name = 'transactions'
	`).Scan(&txCols).Error

	colMap := make(map[string]bool)
	for _, col := range txCols {
		colMap[strings.ToLower(col)] = true
	}

	// Determine available date column names
	pencairanExpr := "COALESCE(t.created_at, NOW())"
	if colMap["submission_date"] {
		pencairanExpr = "COALESCE(t.submission_date, t.created_at, NOW())"
	} else if colMap["disbursement_date"] {
		pencairanExpr = "COALESCE(t.disbursement_date, t.created_at, NOW())"
	}

	dueExpr := fmt.Sprintf("(%s + INTERVAL '15 days')", pencairanExpr)
	if colMap["estimated_maturity_date"] {
		dueExpr = fmt.Sprintf("COALESCE(t.estimated_maturity_date, %s + INTERVAL '15 days')", pencairanExpr)
	} else if colMap["maturity_date"] {
		dueExpr = fmt.Sprintf("COALESCE(t.maturity_date, %s + INTERVAL '15 days')", pencairanExpr)
	} else if colMap["due_date"] {
		dueExpr = fmt.Sprintf("COALESCE(t.due_date, %s + INTERVAL '15 days')", pencairanExpr)
	}

	loanExpr := "0"
	if colMap["requested_loan"] {
		loanExpr = "COALESCE(t.requested_loan, 0)"
	} else if colMap["loan_amount"] {
		loanExpr = "COALESCE(t.loan_amount, 0)"
	}

	adminExpr := "0"
	if colMap["admin_cost_amount"] {
		adminExpr = "COALESCE(t.admin_cost_amount, 0)"
	} else if colMap["admin_fee"] {
		adminExpr = "COALESCE(t.admin_fee, 0)"
	}

	selectSQL := fmt.Sprintf(`
		t.transaction_id,
		COALESCE(t.doc_no, '') AS no_sbg,
		COALESCE(c.full_name, '') AS nama_customer,
		COALESCE(c.mobile_phone, c.phone, '') AS hp,
		COALESCE(c.identity_number, '') AS nik,
		COALESCE(c.address, '') AS alamat,
		COALESCE(o.offering_name, 'Gadai Reguler') AS nama_produk,
		COALESCE(o.product_id, 1) AS product_id,
		%s AS nominal_up,
		%s AS biaya_admin,
		%s AS raw_pencairan,
		%s AS raw_due,
		COALESCE(t.loan_status, 'ACTIVE') AS loan_status,
		COALESCE(t.application_status, 'NEW') AS app_status
	`, loanExpr, adminExpr, pencairanExpr, dueExpr)

	query := config.DBPostgres.Table("gadai.transactions t").
		Select(selectSQL).
		Joins("LEFT JOIN product.product_offerings o ON t.offering_id = o.offering_id").
		Joins("LEFT JOIN customer.customers c ON t.customer_id = c.customer_id")

	if searchQuery != "" {
		likePattern := "%" + searchQuery + "%"
		query = query.Where(`(
			t.doc_no ILIKE ? OR 
			c.full_name ILIKE ? OR 
			c.mobile_phone ILIKE ? OR 
			c.phone ILIKE ? OR
			c.identity_number ILIKE ?
		)`, likePattern, likePattern, likePattern, likePattern, likePattern)
	}

	type RawSopigaRow struct {
		TransactionID int       `gorm:"column:transaction_id"`
		NoSbg         string    `gorm:"column:no_sbg"`
		NamaCustomer  string    `gorm:"column:nama_customer"`
		Hp            string    `gorm:"column:hp"`
		Nik           string    `gorm:"column:nik"`
		Alamat        string    `gorm:"column:alamat"`
		NamaProduk    string    `gorm:"column:nama_produk"`
		ProductID     int       `gorm:"column:product_id"`
		NominalUP     float64   `gorm:"column:nominal_up"`
		BiayaAdmin    float64   `gorm:"column:biaya_admin"`
		RawPencairan  time.Time `gorm:"column:raw_pencairan"`
		RawDue        time.Time `gorm:"column:raw_due"`
		LoanStatus    string    `gorm:"column:loan_status"`
		AppStatus     string    `gorm:"column:app_status"`
	}

	var rows []RawSopigaRow
	err := query.Order("COALESCE(t.submission_date, t.created_at) DESC").Limit(25).Scan(&rows).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"data":  []SopigaTransactionItem{},
			"error": fmt.Sprintf("Gagal query data Sopiga Pro: %v", err),
		})
		return
	}

	results := make([]SopigaTransactionItem, 0, len(rows))
	for _, r := range rows {
		jenis := "Emas"
		if r.ProductID == 2 || r.ProductID == 3 || strings.Contains(strings.ToLower(r.NamaProduk), "elektronik") {
			jenis = "Elektronik"
		}

		tglPencairan := r.RawPencairan.Format("2006-01-02")
		tglDue := r.RawDue.Format("2006-01-02")

		// Calculate exact selisih hari
		diffDays := int(math.Round(r.RawDue.Sub(r.RawPencairan).Hours() / 24))
		if diffDays <= 0 {
			diffDays = 15
		}

		results = append(results, SopigaTransactionItem{
			TransactionID:   r.TransactionID,
			DocNo:           r.NoSbg,
			CustomerName:    r.NamaCustomer,
			CustomerPhone:   r.Hp,
			CustomerNIK:     r.Nik,
			CustomerAlamat:  r.Alamat,
			ProductName:     r.NamaProduk,
			ProductID:       r.ProductID,
			JenisBarang:     jenis,
			DeskripsiBarang: fmt.Sprintf("%s - %s", r.NamaProduk, jenis),
			NominalUP:       r.NominalUP,
			BiayaAdmin:      r.BiayaAdmin,
			TglPencairan:    tglPencairan,
			TglJatuhTempo:   tglDue,
			SelisihHari:     diffDays,
			LoanStatus:      r.LoanStatus,
			AppStatus:       r.AppStatus,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  results,
		"count": len(results),
	})
}

// SbgSelectItem matches the JSON item returned by C_sistem/getsbglist
type SbgSelectItem struct {
	ID         string `json:"id"`
	Text       string `json:"text"`
	NmCustomer string `json:"nm_customer"`
	TglCair    string `json:"tgl_cair"`
	SchemaType string `json:"schema_type"`
	JnstrxVal  string `json:"jnstrx_val"`
	NmProduk   string `json:"nm_produk"`
	KdProduk   string `json:"kd_produk"`
}

// GetSbgListHandler handles GET /api/sopiga/getsbglist
// Exactly matches C_sistem/getsbglist in legacy programming/selada
func GetSbgListHandler(c *gin.Context) {
	if config.DBPostgres == nil {
		c.JSON(http.StatusOK, []SbgSelectItem{})
		return
	}

	searchTerm := strings.TrimSpace(strings.ToLower(c.Query("term")))
	jnstrx := strings.TrimSpace(c.Query("jnstrx"))

	// 1. Get already used SBG numbers (no_do in tb_rencana_det)
	var usedSbg []string
	if config.DBSelada != nil {
		_ = config.DBSelada.Table("tb_rencana_det").
			Where("no_do IS NOT NULL AND no_do != ''").
			Pluck("DISTINCT no_do", &usedSbg).Error
	}
	usedMap := make(map[string]bool)
	for _, u := range usedSbg {
		usedMap[strings.TrimSpace(strings.ToUpper(u))] = true
	}

	// 2. Query transactions from Sopiga Pro (gadai.transactions)
	type RawSbgRow struct {
		NoSbg      string    `gorm:"column:no_sbg"`
		NmCustomer string    `gorm:"column:nm_customer"`
		TglCair    time.Time `gorm:"column:tgl_cair"`
		SchemaType string    `gorm:"column:schema_type"`
		NmProduk   string    `gorm:"column:nm_produk"`
		KdProduk   string    `gorm:"column:kd_produk"`
	}

	query := config.DBPostgres.Table("gadai.transactions t").
		Select(`
			COALESCE(t.doc_no, '') AS no_sbg,
			COALESCE(c.full_name, '') AS nm_customer,
			COALESCE(t.submission_date, t.created_at, NOW()) AS tgl_cair,
			COALESCE(t.schema_type, 'GADAI') AS schema_type,
			COALESCE(o.offering_name, '') AS nm_produk,
			COALESCE(o.offering_code, '') AS kd_produk
		`).
		Joins("LEFT JOIN customer.customers c ON c.customer_id = t.customer_id").
		Joins("LEFT JOIN product.product_offerings o ON t.offering_id = o.offering_id")

	// Filter by submission date
	query = query.Where("t.doc_no IS NOT NULL AND t.doc_no != ''")

	if jnstrx == "2" || strings.ToUpper(jnstrx) == "CICILAN" {
		query = query.Where("UPPER(t.schema_type) = 'CICILAN'")
	} else if jnstrx == "1" || strings.ToUpper(jnstrx) == "GADAI" {
		query = query.Where("UPPER(t.schema_type) = 'GADAI'")
	}

	if searchTerm != "" {
		likeTerm := "%" + searchTerm + "%"
		query = query.Where("(LOWER(t.doc_no) LIKE ? OR LOWER(c.full_name) LIKE ? OR LOWER(c.mobile_phone) LIKE ? OR LOWER(c.phone) LIKE ?)",
			likeTerm, likeTerm, likeTerm, likeTerm)
	}

	var rows []RawSbgRow
	_ = query.Order("COALESCE(t.submission_date, t.created_at) DESC").Limit(30).Scan(&rows).Error

	results := make([]SbgSelectItem, 0)
	for _, r := range rows {
		trimmedNo := strings.TrimSpace(r.NoSbg)
		if trimmedNo == "" || usedMap[strings.ToUpper(trimmedNo)] {
			continue
		}

		schemaType := strings.ToUpper(strings.TrimSpace(r.SchemaType))
		jnstrxVal := "1"
		if schemaType == "CICILAN" {
			jnstrxVal = "2"
		}

		displayTxt := trimmedNo
		if r.NmCustomer != "" {
			displayTxt = fmt.Sprintf("%s - %s", trimmedNo, r.NmCustomer)
		}

		tglCairStr := r.TglCair.Format("2006-01-02")

		results = append(results, SbgSelectItem{
			ID:         trimmedNo,
			Text:       displayTxt,
			NmCustomer: r.NmCustomer,
			TglCair:    tglCairStr,
			SchemaType: schemaType,
			JnstrxVal:  jnstrxVal,
			NmProduk:   r.NmProduk,
			KdProduk:   r.KdProduk,
		})
	}

	c.JSON(http.StatusOK, results)
}

// SbgProdukItem represents offering products
type SbgProdukItem struct {
	KdProduk   string `json:"kd_produk"`
	NmProduk   string `json:"nm_produk"`
	SchemaType string `json:"schema_type"`
}

// GetSbgProdukHandler handles GET/POST /api/sopiga/getproduk
// Exactly matches C_sistem/getproduk in legacy programming/selada
func GetSbgProdukHandler(c *gin.Context) {
	idVal := strings.TrimSpace(c.Query("id"))
	if idVal == "" {
		idVal = strings.TrimSpace(c.PostForm("id"))
	}

	if config.DBPostgres == nil {
		// Fallback static list
		fallback := []SbgProdukItem{
			{KdProduk: "1", NmProduk: "Gadai Emas", SchemaType: "GADAI"},
			{KdProduk: "2", NmProduk: "Gadai Non Emas", SchemaType: "GADAI"},
			{KdProduk: "3", NmProduk: "Cicilan Emas", SchemaType: "CICILAN"},
		}
		c.JSON(http.StatusOK, fallback)
		return
	}

	type RawProduk struct {
		KdProduk   string `gorm:"column:kd_produk"`
		NmProduk   string `gorm:"column:nm_produk"`
		SchemaType string `gorm:"column:schema_type"`
	}

	query := config.DBPostgres.Table("product.product_offerings").
		Select("offering_code AS kd_produk, offering_name AS nm_produk, schema_type").
		Where("is_active = true")

	if idVal == "1" || strings.ToUpper(idVal) == "GADAI" {
		query = query.Where("UPPER(schema_type) = 'GADAI'")
	} else if idVal == "2" || strings.ToUpper(idVal) == "CICILAN" {
		query = query.Where("UPPER(schema_type) = 'CICILAN'")
	}

	var rows []RawProduk
	err := query.Order("offering_name ASC").Scan(&rows).Error
	if err != nil || len(rows) == 0 {
		// Fallback
		var list []SbgProdukItem
		if idVal == "2" || strings.ToUpper(idVal) == "CICILAN" {
			list = []SbgProdukItem{
				{KdProduk: "CICILAN-EMAS", NmProduk: "Cicilan Logam Mulia", SchemaType: "CICILAN"},
			}
		} else {
			list = []SbgProdukItem{
				{KdProduk: "1", NmProduk: "Gadai Emas", SchemaType: "GADAI"},
				{KdProduk: "2", NmProduk: "Gadai Non Emas", SchemaType: "GADAI"},
			}
		}
		c.JSON(http.StatusOK, list)
		return
	}

	items := make([]SbgProdukItem, len(rows))
	for i, r := range rows {
		items[i] = SbgProdukItem{
			KdProduk:   r.KdProduk,
			NmProduk:   r.NmProduk,
			SchemaType: r.SchemaType,
		}
	}
	c.JSON(http.StatusOK, items)
}

// InputDORequest represents input DO payload
type InputDORequest struct {
	IDRendet     string `json:"idrendet" form:"idrendet"`
	NoSbg        string `json:"nosbg" form:"nosbg"`
	JnsTrx       string `json:"jnstrx" form:"jnstrx"`
	Produk       string `json:"produk" form:"produk"`
	NmCustomer   string `json:"nm_customer" form:"nm_customer"`
	SelisihHari  int    `json:"selisih_hari" form:"selisih_hari"`
	KetRencana   string `json:"ket_rencana" form:"ket_rencana"`
	KetAktivitas string `json:"ket_aktivitas" form:"ket_aktivitas"`
}

// InputDOHandler handles POST /api/sopiga/inputdo
// Exactly matches C_sistem/inputdo in legacy programming/selada
func InputDOHandler(c *gin.Context) {
	if config.DBSelada == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Database Selada tidak terhubung"})
		return
	}

	var req InputDORequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Format payload tidak valid: " + err.Error()})
		return
	}

	req.NoSbg = strings.TrimSpace(req.NoSbg)
	req.NmCustomer = strings.TrimSpace(req.NmCustomer)
	req.IDRendet = strings.TrimSpace(req.IDRendet)

	if req.NoSbg == "" || req.JnsTrx == "" || req.JnsTrx == "0" || req.Produk == "" || req.Produk == "0" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Lengkapi semua data terlebih dahulu."})
		return
	}

	// 1. Fetch customer from DB for this inquiry/rendet
	type CustCheck struct {
		IDRencana    int    `gorm:"column:id_rencana"`
		IDRencanaDet int    `gorm:"column:id_rencana_det"`
		Nama         string `gorm:"column:nama"`
		RO           string `gorm:"column:ro"`
	}

	var found CustCheck
	err := config.DBSelada.Table("tb_rencana_det rd").
		Select("rd.id_rencana, rd.id_rencana_det, c.nama, COALESCE(rd.ro, 0) as ro").
		Joins("JOIN tb_rencana r ON rd.id_rencana = r.id_rencana").
		Joins("JOIN tb_customer c ON r.id_customer = c.id_customer").
		Where("rd.id_rencana_det = ? OR r.id_rencana = ?", req.IDRendet, req.IDRendet).
		Order("rd.id_rencana_det DESC").
		Limit(1).
		Scan(&found).Error

	if err != nil || found.IDRencana == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": fmt.Sprintf("Data customer tidak ditemukan untuk idrendet: %s", req.IDRendet),
		})
		return
	}

	dbCustomerName := strings.TrimSpace(found.Nama)
	sbgCustomerName := strings.TrimSpace(req.NmCustomer)

	// Validate customer name matches SBG customer name (per legacy C_sistem.php line 1378)
	if sbgCustomerName != "" && !strings.EqualFold(dbCustomerName, sbgCustomerName) {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Nama Customer Tidak Sesuai dengan SBG.",
		})
		return
	}

	// 2. Insert into tb_rencana_det
	now := time.Now()
	insertMap := map[string]interface{}{
		"id_rencana":    found.IDRencana,
		"tgl_rencana":   now,
		"id_hslaktiv":   7, // Realisasi / SBG Deal
		"id_aktiv":      0,
		"id_trx":        req.JnsTrx,
		"id_produk":     req.Produk,
		"ket_aktivitas": req.KetAktivitas,
		"ket_rencana":   req.KetRencana,
		"tahun_motor":   0,
		"status":        1,
		"ro":            found.RO,
		"no_do":         req.NoSbg,
		"selisih_hari":  req.SelisihHari,
	}

	if err := config.DBSelada.Table("tb_rencana_det").Create(&insertMap).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal menyimpan data DO ke database: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "Data DO berhasil disimpan.",
	})
}

// GetInquiryContextHandler handles GET /api/sopiga/inquiry-context?id=...
// Returns current inquiry info for pre-populating pageinputdo
func GetInquiryContextHandler(c *gin.Context) {
	idVal := strings.TrimSpace(c.Query("id"))
	if idVal == "" || config.DBSelada == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID required"})
		return
	}

	type InqCtx struct {
		IDRencana    int    `json:"id_rencana" gorm:"column:id_rencana"`
		IDRencanaDet int    `json:"id_rencana_det" gorm:"column:id_rencana_det"`
		NamaCustomer string `json:"nama" gorm:"column:nama"`
		HP           string `json:"hp" gorm:"column:hp"`
		KetRencana   string `json:"ket_rencana" gorm:"column:ket_rencana"`
		KetAktivitas string `json:"ket_aktivitas" gorm:"column:ket_aktivitas"`
		Status       string `json:"status" gorm:"column:status"`
		NamaHslaktiv string `json:"nama_hslaktiv" gorm:"column:nama_hslaktiv"`
		IDHslaktiv   int    `json:"id_hslaktiv" gorm:"column:id_hslaktiv"`
	}

	var ctx InqCtx
	err := config.DBSelada.Table("tb_rencana_det rd").
		Select("rd.id_rencana, rd.id_rencana_det, c.nama, c.hp, rd.ket_rencana, rd.ket_aktivitas, rd.id_hslaktiv, ha.nama_hslaktiv").
		Joins("JOIN tb_rencana r ON rd.id_rencana = r.id_rencana").
		Joins("JOIN tb_customer c ON r.id_customer = c.id_customer").
		Joins("LEFT JOIN tb_hslaktivitas ha ON ha.id_hslaktiv = rd.id_hslaktiv").
		Where("rd.id_rencana_det = ? OR r.id_rencana = ?", idVal, idVal).
		Order("rd.id_rencana_det DESC").
		Limit(1).
		Scan(&ctx).Error

	if err != nil || ctx.IDRencana == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data inquiry tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, ctx)
}
