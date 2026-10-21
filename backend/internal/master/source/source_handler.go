package source

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

// GetSumberCust handles GET /api/sumbercust
func (h *Handler) GetSumberCust(c *gin.Context) {
	ctx := c.Request.Context()
	res, err := h.service.GetReferenceSources(ctx)
	if err != nil {
		c.JSON(http.StatusOK, []SumberCustResponse{})
		return
	}
	c.JSON(http.StatusOK, res)
}

// GetSumberSource handles GET /api/master/sumber_source
func (h *Handler) GetSumberSource(c *gin.Context) {
	ctx := c.Request.Context()
	res, err := h.service.GetApplicationSources(ctx)
	if err != nil {
		c.JSON(http.StatusOK, []ApplicationSourceResponse{})
		return
	}
	c.JSON(http.StatusOK, res)
}

// RegisterRoutes registers endpoints for source domain
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/sumbercust", h.GetSumberCust)
	rg.GET("/master/sumber_source", h.GetSumberSource)
}
