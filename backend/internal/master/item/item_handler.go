package item

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{service: s}
}

// GetKatalogList handles GET /api/katalog/list/:category
func (h *Handler) GetKatalogList(c *gin.Context) {
	category := c.Param("category")
	ctx := c.Request.Context()

	items, err := h.service.GetKatalogList(ctx, category)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, items)
}

// GetKatalogSpec handles GET /api/katalog/spec/:id
func (h *Handler) GetKatalogSpec(c *gin.Context) {
	id := c.Param("id")
	ctx := c.Request.Context()

	spec, err := h.service.GetKatalogSpec(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Motor tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, spec)
}

// GetKatalogPoint handles GET /api/katalog/point
func (h *Handler) GetKatalogPoint(c *gin.Context) {
	userID, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	ctx := c.Request.Context()
	point, err := h.service.GetSalesPoint(ctx, userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, point)
}

// RegisterRoutes registers endpoints for catalog item domain
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/katalog/point", h.GetKatalogPoint)
	rg.GET("/katalog/list/:category", h.GetKatalogList)
	rg.GET("/katalog/spec/:id", h.GetKatalogSpec)
}
