package routes

import (
	"go-api/handlers"
	"go-api/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRouter configures endpoints and applies global security and auth middlewares
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Apply security and CORS middlewares globally
	r.Use(middleware.SecurityHeadersMiddleware())
	r.Use(middleware.CORSMiddleware())

	// Public routes
	r.POST("/api/login", handlers.LoginHandler)
	r.POST("/api/logout", handlers.LogoutHandler)

	// Protected routes group
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		// Dashboard Metrics
		protected.GET("/dashboard/metrics", handlers.DashboardMetricsHandler)

		// Inquiry / Rencana
		protected.GET("/inquiry/counts", handlers.InquiryCountsHandler)
		protected.GET("/inquiry/filter", handlers.InquiryListHandler)
		protected.POST("/rencana/input", handlers.InputRencanaHandler)
		protected.DELETE("/rencana/:id", handlers.DeleteRencanaHandler)

		// Catalog / Katalog
		protected.GET("/katalog/point", handlers.GetKatalogPointHandler)
		protected.GET("/katalog/list/:category", handlers.GetKatalogListHandler)
		protected.GET("/katalog/spec/:id", handlers.GetKatalogSpecHandler)

		// Pawn Estimation / Taksiran
		protected.GET("/taksiran/jenis-barang", handlers.GetJenisBarangHandler)
		protected.GET("/taksiran/grade", handlers.GetGradeHandler)
		protected.GET("/taksiran/barang-umum", handlers.GetBarangUmumHandler)
		protected.GET("/taksiran/barang-emas", handlers.GetBarangEmasHandler)
		protected.GET("/taksiran/products", handlers.GetTaksiranProductHandler)
		protected.GET("/taksiran/biaya-admin", handlers.GetBiayaAdminHandler)
		protected.GET("/taksiran/biaya-admin-emas", handlers.GetBiayaAdminEmasHandler)
		protected.GET("/taksiran/wilayah", handlers.GetWilayahHandler)
		protected.GET("/taksiran/nilai-taksir", handlers.GetNilaiTaksiranHandler)
		protected.GET("/taksiran/gold-stle", handlers.GetGoldStleHandler)
		protected.GET("/taksiran/barang-by-kode", handlers.GetBarangByKodeHandler)
		protected.GET("/taksiran/electronic-scoring", handlers.GetElectronicScoringHandler)
		protected.POST("/taksiran/simpan", handlers.SimpanTransaksiHandler)
		protected.GET("/taksiran/history", handlers.GetTaksiranHistoryHandler)

		// Reference & Lookup Endpoints (PHP Parity)
		protected.GET("/customer/select", handlers.GetCustomerSelectHandler)
		protected.GET("/sumbercust", handlers.GetSumberCustHandler)
		protected.GET("/sales/list", handlers.GetSalesListHandler)
		protected.GET("/kelurahan", handlers.GetKelurahanHandler)

		// master sumber source
		protected.GET("/master/sumber_source", handlers.SumberSourceList)
	}

	return r
}
