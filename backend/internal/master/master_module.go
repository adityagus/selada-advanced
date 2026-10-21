package master

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"go-api/internal/master/branch"
	"go-api/internal/master/item"
	"go-api/internal/master/location"
	"go-api/internal/master/sales"
	"go-api/internal/master/source"
)

// MasterModule aggregates all master data sub-domains
type MasterModule struct {
	BranchHandler   *branch.Handler
	SourceHandler   *source.Handler
	LocationHandler *location.Handler
	SalesHandler    *sales.Handler
	ItemHandler     *item.Handler
}

// InitMasterModule constructs all master dependencies via Dependency Injection
func InitMasterModule(dbPG *gorm.DB, dbSelada *gorm.DB, dbMysqlPos *gorm.DB) *MasterModule {
	// Branch Domain Wiring
	branchRepo := branch.NewRepository(dbPG)
	branchService := branch.NewService(branchRepo)
	branchHandler := branch.NewHandler(branchService)

	// Source Domain Wiring
	sourceRepo := source.NewRepository(dbPG)
	sourceService := source.NewService(sourceRepo)
	sourceHandler := source.NewHandler(sourceService)

	// Location Domain Wiring
	locRepo := location.NewRepository(dbPG)
	locService := location.NewService(locRepo)
	locHandler := location.NewHandler(locService)

	// Sales Domain Wiring
	salesRepo := sales.NewRepository(dbPG, dbMysqlPos)
	salesService := sales.NewService(salesRepo)
	salesHandler := sales.NewHandler(salesService)

	// Item / Katalog Domain Wiring
	itemRepo := item.NewRepository(dbMysqlPos, dbSelada)
	itemService := item.NewService(itemRepo)
	itemHandler := item.NewHandler(itemService)

	return &MasterModule{
		BranchHandler:   branchHandler,
		SourceHandler:   sourceHandler,
		LocationHandler: locHandler,
		SalesHandler:    salesHandler,
		ItemHandler:     itemHandler,
	}
}

// RegisterRoutes registers all master endpoints to the provided protected router group
func (m *MasterModule) RegisterRoutes(protected *gin.RouterGroup) {
	m.BranchHandler.RegisterRoutes(protected)
	m.SourceHandler.RegisterRoutes(protected)
	m.LocationHandler.RegisterRoutes(protected)
	m.SalesHandler.RegisterRoutes(protected)
	m.ItemHandler.RegisterRoutes(protected)
}
