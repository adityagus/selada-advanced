package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"go-api/config"
	"go-api/models"

	"github.com/gin-gonic/gin"
)

// GetJenisBarangHandler fetches category list from PostgreSQL
func GetJenisBarangHandler(c *gin.Context) {
	if config.DBPostgres == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database PostgreSQL tidak terhubung"})
		return
	}

	var rawList []models.JenisBarang
	err := config.DBPostgres.Select("kd_jenis_barang, nm_jenis_barang").Find(&rawList).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil jenis barang"})
		return
	}

	var filteredList []models.JenisBarang
	for _, item := range rawList {
		if len(item.KdJenisBarang) == 2 {
			filteredList = append(filteredList, item)
		}
	}

	c.JSON(http.StatusOK, filteredList)
}

// GetBarangUmumHandler queries PostgreSQL for general goods detail
func GetBarangUmumHandler(c *gin.Context) {
	if config.DBPostgres == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database PostgreSQL tidak terhubung"})
		return
	}

	jenis := c.Query("jenisBarang")
	grade := c.Query("gradeBarang")
	nama := c.Query("namaBarang")

	var results []map[string]interface{}
	query := config.DBPostgres.Table("tblbarang b").
		Select("b.nm_barang, b.kd_barang, b.status_antam, bd.*").
		Joins("LEFT JOIN tblbarang_detail bd ON bd.fk_barang = b.kd_barang").
		Where("b.barang_active = ? AND bd.fk_wilayah = ?", "true", "2")
	
	if jenis != "" {
		query = query.Where("b.fk_jenis_barang = ?", jenis)
	}
	if grade != "" {
		query = query.Where("b.fk_grade = ?", grade)
	}
	if nama != "" {
		query = query.Where("b.nm_barang ILIKE ?", "%"+nama+"%")
	}

	err := query.Order("b.kd_barang ASC").Limit(5).Scan(&results).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data barang umum"})
		return
	}

	c.JSON(http.StatusOK, results)
}

// GetBarangEmasHandler returns gold inventory information from PostgreSQL
func GetBarangEmasHandler(c *gin.Context) {
	if config.DBPostgres == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database PostgreSQL tidak terhubung"})
		return
	}

	var results []map[string]interface{}
	err := config.DBPostgres.Table("tblbarang b").
		Select("b.kd_barang, b.nm_barang, b.status_antam, t.kd_cabang, t.harga_stle_antam, t.harga_stle_non_antam, t.harga_stle_perhiasan").
		Joins("CROSS JOIN tblcabang t").
		Where("b.barang_active = ? AND b.fk_jenis_barang = ? AND t.kd_cabang = ?", "true", "0", "0229").
		Scan(&results).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data barang emas"})
		return
	}

	c.JSON(http.StatusOK, results)
}

// GetTaksiranProductHandler fetches products from PostgreSQL sam_live
func GetTaksiranProductHandler(c *gin.Context) {
	if config.DBPostgres == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database PostgreSQL tidak terhubung"})
		return
	}

	typeParam := c.Query("type")

	var products []models.Produk
	query := config.DBPostgres.Where("active = ?", true)

	if strings.ToLower(typeParam) == "emas" {
		query = query.Where("UPPER(nm_produk) LIKE ?", "%EMAS%")
	} else {
		query = query.Where("kd_produk IN ?", []string{"21", "23"})
	}

	err := query.Order("kd_produk ASC").Find(&products).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data produk"})
		return
	}

	c.JSON(http.StatusOK, products)
}

// GetBiayaAdminHandler handles administrative fee calculation requests
func GetBiayaAdminHandler(c *gin.Context) {
	if config.DBPostgres == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database PostgreSQL tidak terhubung"})
		return
	}

	totalPinjamanStr := c.Query("total_pinjaman")
	kdProduk := c.Query("kd_produk")

	if totalPinjamanStr == "" || kdProduk == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameter total_pinjaman dan kd_produk wajib diisi"})
		return
	}

	totalPinjaman, err := strconv.ParseFloat(totalPinjamanStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameter total_pinjaman tidak valid"})
		return
	}

	var adminDetail models.ProdukDetailBiayaAdmin
	err = config.DBPostgres.
		Where("fk_produk = ? AND dari <= ? AND ke >= ?", kdProduk, totalPinjaman, totalPinjaman).
		First(&adminDetail).Error

	biayaAdmin := 0.0
	if err == nil {
		biayaAdmin = adminDetail.Nilai
	}

	c.JSON(http.StatusOK, gin.H{"biaya_admin": biayaAdmin})
}

// SimpanTransaksiHandler stores the transaction details
func SimpanTransaksiHandler(c *gin.Context) {
	// Parse request details and mock simulation output
	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payload tidak valid"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Transaksi berhasil disimpan"})
}
