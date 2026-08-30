package handlers

import (
	"math"
	"net/http"
	"time"

	"go-api/config"
	"go-api/models"

	"github.com/gin-gonic/gin"
)

// GetKatalogPointHandler calculates sales activity score/points
func GetKatalogPointHandler(c *gin.Context) {
	if config.DBMysql == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database MySQL tidak terhubung"})
		return
	}

	userId, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	tgl := time.Now()
	// Set threshold date as in legacy PHP code
	limitTgl, _ := time.Parse("2006-01-02", "2021-12-31")
	if tgl.After(limitTgl) {
		tgl = limitTgl
	}

	n := int(tgl.Month())
	thn := tgl.Year()

	var j int = 0
	var tpoin float64 = 0.0

	// Loop months from March (3) to current month (n)
	for i := 3; i <= n; i++ {
		var q []models.CronDetPointFn
		err := config.DBMysql.Where("bln = ? AND thn = ? AND id_sales = ?", i, thn, userId).Find(&q).Error
		if err == nil {
			for _, r := range q {
				mpoin := float64(r.Pleads + r.Pfu + r.Pdo + r.Pakurdat + r.Leadstodo)
				if r.Leadspal > 0 {
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

	c.JSON(http.StatusOK, gin.H{"point": poin})
}

// GetKatalogListHandler fetches list of motorcycles based on category slug
func GetKatalogListHandler(c *gin.Context) {
	if config.DBMysqlPos == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database MySQL POS tidak terhubung"})
		return
	}

	category := c.Param("category")
	var categoryId int

	switch category {
	case "sport":
		categoryId = 10201
	case "trail":
		categoryId = 10202
	case "moped":
		categoryId = 10203
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kategori tidak valid"})
		return
	}

	type KatalogItem struct {
		IDMotor     string `json:"id_motor"`
		NamaMotor   string `json:"nama_motor"`
		GambarColor string `json:"gambar_color"`
	}

	var items []KatalogItem
	err := config.DBMysqlPos.Table("motorcycle m").
		Select("m.id_motor, m.nama_motor, mc.gambar_color").
		Joins("LEFT JOIN motcol mc ON mc.id_motor = m.id_motor AND mc.default_color = 'Y'").
		Where("m.id_category = ?", categoryId).
		Scan(&items).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil katalog motor"})
		return
	}

	c.JSON(http.StatusOK, items)
}

// SpecResponse details spec data for a motorcycle
type SpecResponse struct {
	DataMotor   models.Motorcycle     `json:"datamotor"`
	GbrMotor    []GbrMotorDetail      `json:"gbrmotor"`
	WarnaMotor  []models.MasterColor  `json:"warnamotor"`
	HargaMotor  models.MotFitur       `json:"hargamotor"`
	FiturMotor  []models.MotFiturList `json:"fiturmotor"`
	BrosurMotor []models.FiturMotor   `json:"brosurmotor"`
}

type GbrMotorDetail struct {
	IDMotor     string `json:"id_motor"`
	NamaMotor   string `json:"nama_motor"`
	IDIcon      string `json:"id_icon"`
	GambarColor string `json:"gambar_color"`
	ColorName   string `json:"color_name"`
}

// GetKatalogSpecHandler returns specifications, pricing, brochures, and colors for a motor ID
func GetKatalogSpecHandler(c *gin.Context) {
	if config.DBMysqlPos == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database MySQL POS tidak terhubung"})
		return
	}

	idMotor := c.Param("id")

	// 1. Data Motor
	var dataMotor models.Motorcycle
	if err := config.DBMysqlPos.Where("id_motor = ?", idMotor).First(&dataMotor).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Motor tidak ditemukan"})
		return
	}

	// 2. Gambar Motor Details
	var gbrMotor []GbrMotorDetail
	config.DBMysqlPos.Table("motorcycle m").
		Select("m.id_motor, m.nama_motor, mc.id_icon, mc.gambar_color, c.color_name").
		Joins("JOIN motcol mc ON mc.id_motor = m.id_motor").
		Joins("JOIN master_color c ON c.id_master = mc.id_icon").
		Where("m.id_motor = ?", idMotor).
		Scan(&gbrMotor)

	// 3. Warna Motor lists
	var warnaMotor []models.MasterColor
	config.DBMysqlPos.Table("motcol mc").
		Select("c.id_master, c.color_name").
		Joins("JOIN master_color c ON c.id_master = mc.id_icon").
		Where("mc.id_motor = ?", idMotor).
		Scan(&warnaMotor)

	// 4. Harga Motor
	var hargaMotor models.MotFitur
	config.DBMysqlPos.Where("id_motor = ?", idMotor).First(&hargaMotor)

	// 5. Fitur Motor
	var fiturMotor []models.MotFiturList
	config.DBMysqlPos.Where("id_motor = ?", idMotor).Find(&fiturMotor)

	// 6. Brosur Motor
	var brosurMotor []models.FiturMotor
	config.DBMysqlPos.Where("id_motor = ?", idMotor).Find(&brosurMotor)

	c.JSON(http.StatusOK, SpecResponse{
		DataMotor:   dataMotor,
		GbrMotor:    gbrMotor,
		WarnaMotor:  warnaMotor,
		HargaMotor:  hargaMotor,
		FiturMotor:  fiturMotor,
		BrosurMotor: brosurMotor,
	})
}
