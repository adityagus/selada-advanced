package sales

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

// GetSalesList handles GET /api/sales/list
func (h *Handler) GetSalesList(c *gin.Context) {
	ctx := c.Request.Context()
	res, err := h.service.GetSalesList(ctx)
	if err != nil {
		c.JSON(http.StatusOK, []SalesResponse{})
		return
	}
	c.JSON(http.StatusOK, res)
}

// RegisterRoutes registers endpoints for sales domain
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/sales/list", h.GetSalesList)
}
