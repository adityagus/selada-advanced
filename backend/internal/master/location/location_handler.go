package location

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

// SearchKelurahan handles GET /api/kelurahan
func (h *Handler) SearchKelurahan(c *gin.Context) {
	query := c.Query("q")
	ctx := c.Request.Context()

	res, err := h.service.SearchKelurahan(ctx, query)
	if err != nil {
		c.JSON(http.StatusOK, []KelurahanResponse{})
		return
	}

	c.JSON(http.StatusOK, res)
}

// RegisterRoutes registers endpoints for location domain
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/kelurahan", h.SearchKelurahan)
}
