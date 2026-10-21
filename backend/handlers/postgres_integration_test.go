package handlers

import (
	"fmt"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"go-api/config"
	"go-api/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Opt-in: requires a migrated Selada database. All test writes are rolled back;
// PostgreSQL identity sequences may advance even after rollback.
func TestSeladaPostgresIntegration(t *testing.T) {
	if os.Getenv("SELADA_INTEGRATION") != "1" {
		t.Skip("set SELADA_INTEGRATION=1 to test migrated database")
	}
	if err := config.LoadEnv("../.env"); err != nil {
		t.Fatal(err)
	}
	dsn, err := config.PostgresDSN("DB_PG2", "")
	if err != nil {
		t.Fatal(err)
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatal(err)
	}
	defer sqlDB.Close()
	tx := db.Begin()
	if tx.Error != nil {
		t.Fatal(tx.Error)
	}
	defer tx.Rollback()
	previous := config.DBSelada
	config.DBSelada = tx
	defer func() { config.DBSelada = previous }()
	gin.SetMode(gin.TestMode)
	user := fmt.Sprintf("migration-test-%d", time.Now().UnixNano())
	customer := models.Customer{Nama: "Migration Test", Hp: "0800000000", IDLogin: user}
	if err := tx.Create(&customer).Error; err != nil {
		t.Fatal(err)
	}
	plan := models.Rencana{IDCustomer: customer.IDCustomer, IDLogin: user, IDSumbercust: "1"}
	if err := tx.Create(&plan).Error; err != nil {
		t.Fatal(err)
	}
	detail := models.RencanaDet{IDRencana: plan.IDRencana, TglRencana: time.Now(), IDHslaktiv: 1, Status: 1}
	if err := tx.Create(&detail).Error; err != nil {
		t.Fatal(err)
	}
	if customer.IDCustomer == 0 || plan.IDRencana == 0 || detail.IDRencanaDet == 0 {
		t.Fatal("identity not returned")
	}
	r := gin.New()
	r.Use(func(c *gin.Context) { c.Set("username", user); c.Next() })
	r.GET("/customers", GetCustomerListHandler)
	r.GET("/customers/:id/inquiries", GetCustomerInquiriesHandler)
	r.GET("/counts", InquiryCountsHandler)
	r.GET("/inquiries", InquiryListHandler)
	r.GET("/select", GetCustomerSelectHandler)
	r.GET("/context", GetInquiryContextHandler)
	r.POST("/do", InputDOHandler)
	for _, path := range []string{"/customers?search=migration", fmt.Sprintf("/customers/%d/inquiries", customer.IDCustomer), "/counts", "/inquiries?search=migration", "/inquiries?status_filter=12", "/select?q=migration", fmt.Sprintf("/context?id=%d", detail.IDRencanaDet)} {
		t.Run(path, func(t *testing.T) {
			w := httptest.NewRecorder()
			r.ServeHTTP(w, httptest.NewRequest("GET", path, nil))
			if w.Code != 200 {
				t.Fatalf("HTTP %d: %s", w.Code, w.Body.String())
			}
		})
	}
	t.Run("insert SBG", func(t *testing.T) {
		body := fmt.Sprintf(`{"idrendet":"%d","nosbg":"MIGRATION-TEST","jnstrx":"1","produk":"1","nm_customer":"Migration Test"}`, detail.IDRencanaDet)
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/do", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)
		if w.Code != 200 {
			t.Fatalf("HTTP %d: %s", w.Code, w.Body.String())
		}
		var count int64
		if err := tx.Table("tb_rencana_det").Where("id_rencana = ? AND no_do = ?", plan.IDRencana, "MIGRATION-TEST").Count(&count).Error; err != nil || count != 1 {
			t.Fatalf("SBG count=%d, error=%v", count, err)
		}
	})
}
