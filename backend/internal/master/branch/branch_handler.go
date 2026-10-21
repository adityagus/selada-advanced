package branch

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

// GetCurrentBranch handles GET /api/branch/current
func (h *Handler) GetCurrentBranch(c *gin.Context) {
	cabang := c.GetString("cabang")
	username := c.GetString("username")
	branchName := c.GetString("branch_name")

	ctx := c.Request.Context()
	res := h.service.ResolveCurrentBranch(ctx, cabang, username, branchName)

	c.JSON(http.StatusOK, res)
}

// GetBranchList handles GET /api/master/cabang
func (h *Handler) GetBranchList(c *gin.Context) {
	ctx := c.Request.Context()
	branches, err := h.service.GetBranches(ctx)
	if err != nil {
		c.JSON(http.StatusOK, []BranchResponse{})
		return
	}

	c.JSON(http.StatusOK, branches)
}

// RegisterRoutes registers endpoints for branch domain
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/branch/current", h.GetCurrentBranch)
	rg.GET("/master/cabang", h.GetBranchList)
}
