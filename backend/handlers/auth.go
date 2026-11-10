package handlers

import (
	"crypto/md5"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"go-api/config"
	"go-api/middleware"
	"go-api/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// LoginRequest represents the user login payload
type LoginRequest struct {
	User string `json:"user" form:"user" binding:"required"`
	Pass string `json:"pass" form:"pass" binding:"required"`
}

// LoginHandler handles POST /api/login matching legacy M_login and C_login logic
func LoginHandler(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// 1. Fetch user from MySQL POS database (sam_pos.tbl_user joined with mst_jabatan)
	if config.DBMysqlPos == nil {
		config.ConnectMySQLPos()
	}

	var user struct {
		models.UserPos
		NmJabatan string `gorm:"column:nm_jabatan"`
	}

	// Developer / Fallback override password check
	loginBypass := strings.EqualFold(strings.Trim(strings.TrimSpace(os.Getenv("LOGIN_BYPASS")), "\"'"), "true")
	isDevPass := loginBypass || req.Pass == "selada" || req.Pass == "admin"

	if config.DBMysqlPos == nil {
		// Fallback to PostgreSQL auth.users / dev bypass when MySQL POS is unreachable from server
		if config.DBPostgres != nil && isDevPass {
			var pgFallback struct {
				Username   string `gorm:"column:username"`
				FullName   string `gorm:"column:full_name"`
				BranchCode string `gorm:"column:branch_code"`
			}
			errPg := config.DBPostgres.Table("auth.users u").
				Select("u.username, COALESCE(e.full_name, u.username) AS full_name, COALESCE(b.branch_code, '001') AS branch_code").
				Joins("LEFT JOIN master.employee e ON e.employee_id = u.employee_id").
				Joins("LEFT JOIN master.branch b ON e.branch_id = b.branch_id").
				Where("LOWER(u.username) = LOWER(?)", req.User).
				Limit(1).
				Scan(&pgFallback).Error
			if errPg == nil && pgFallback.Username != "" {
				user.Username = pgFallback.Username
				user.Nama = pgFallback.FullName
				user.IDCabang = pgFallback.BranchCode
				user.IDJabatan = 17
				user.NmJabatan = "Sales"
			} else {
				user.Username = req.User
				user.Nama = req.User
				user.IDCabang = "001"
				user.IDJabatan = 1
				user.NmJabatan = "Administrator"
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database MySQL POS tidak terhubung"})
			return
		}
	} else {
		// Legacy M_login: void = 1 indicates active user
		err := config.DBMysqlPos.Table("tbl_user u").
			Select("u.*, j.nm_jabatan").
			Joins("LEFT JOIN mst_jabatan j ON j.id_jabatan = u.id_jabatan").
			Where("u.username = ? AND u.void = ?", req.User, 1).
			First(&user).Error

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if isDevPass {
					// Fallback user details if username not found in DB but dev password is used
					user.Username = req.User
					user.Nama = req.User
					user.IDCabang = "001"
					user.IDJabatan = 1
					user.NmJabatan = "Administrator"
				} else {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Username atau password anda salah"})
					return
				}
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Database user error: " + err.Error()})
				return
			}
		} else {
			// 2. Verify password (bcrypt with MD5 legacy fallback & auto upgrade)
			passwordIsValid := false

			if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Pass)); err == nil {
				passwordIsValid = true
			} else {
				hasher := md5.New()
				hasher.Write([]byte(req.Pass))
				md5Hash := fmt.Sprintf("%x", hasher.Sum(nil))

				if md5Hash == user.Password {
					passwordIsValid = true
				}
			}

			if isDevPass {
				passwordIsValid = true
			}

			if !passwordIsValid {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Username atau password anda salah"})
				return
			}
		}
	}

	// 3. Fetch employee reference from Postgres (sam_live auth.users left join master.employee)
	fkUser := ""
	employeeID := ""
	employeeCode := ""
	branchName := ""

	// Always resolve branch name from user.IDCabang in sam_pos (mst_cabang) first
	if user.IDCabang != "" {
		branchName = ResolveBranchName(user.IDCabang, user.Username)
	}

	if config.DBPostgres != nil {
		// var pgUser map[string]interface{}
		var pgUser models.UserPg
		err := config.DBPostgres.Table("auth.users u").
			Select("e.employee_code AS fk_karyawan, e.employee_id AS employee_id, b.branch_name AS branch_name, e.is_active").
			Joins("LEFT JOIN master.employee e ON e.employee_id = u.employee_id").
			Joins("LEFT JOIN master.branch b ON e.branch_id = b.branch_id").
			Where("LOWER(u.username) = LOWER(?)", user.Username).
			Scan(&pgUser).Error

		fmt.Printf("data %v\n", pgUser)
		fmt.Printf("data %v\n %v\n %v\n %v", pgUser.BranchName, pgUser.EmployeeCode, pgUser.EmployeeID, pgUser.IsActive)

		if err == nil {
			// Legacy C_login: check if user/employee is active
			if pgUser.IsActive != nil && !*pgUser.IsActive {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Oops... Akun Anda sudah tidak aktif."})
				return
			}
			if pgUser.EmployeeID != nil {
				employeeID = *pgUser.EmployeeID
			}
			if pgUser.EmployeeCode != nil {
				employeeCode = *pgUser.EmployeeCode
			}
			if branchName == "" && pgUser.BranchName != nil {
				branchName = *pgUser.BranchName
			}

			if employeeID != "" {
				fkUser = employeeID
			} else {
				fkUser = employeeCode
			}
		}
	}

	// Fallback lookup for branchName if still empty
	if branchName == "" {
		branchName = ResolveBranchName(user.IDCabang, user.Username)
	}

	// 4. Generate custom signed session token
	claims := middleware.TokenClaims{
		Username:   user.Username,
		Nama:       user.Nama,
		Cabang:     user.IDCabang,
		BranchName: branchName,
		FkUser:     fkUser,
		Expiry:     time.Now().Add(24 * time.Hour),
	}

	token, err := middleware.GenerateToken(claims)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal generate session token"})
		return
	}

	isReview := user.IDJabatan == 4 || user.IDJabatan == 21 || user.IDJabatan == 74

	c.JSON(http.StatusOK, gin.H{
		"token":         token,
		"username":      user.Username,
		"nama":          user.Nama,
		"cabang":        user.IDCabang,
		"branch_name":   branchName,
		"fk_user":       fkUser,
		"employee_id":   employeeID,
		"employee_code": employeeCode,
		"id_jabatan":    user.IDJabatan,
		"nm_jabatan":    user.NmJabatan,
		"is_review":     isReview,
	})
}

// LogoutHandler handles POST /api/logout
func LogoutHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Logout berhasil"})
}

// ResolveBranchName resolves the textual branch name from MySQL POS (mst_cabang) or Postgres master.branch
func ResolveBranchName(cabangID string, username string) string {
	branchName := ""

	// 1. Prioritize MySQL POS (sam_pos: tbl_user -> mst_cabang)
	if config.DBMysqlPos != nil {
		type CabangPos struct {
			NmCabang string `gorm:"column:nm_cabang"`
		}
		var cPos CabangPos
		if username != "" {
			if err := config.DBMysqlPos.Table("tbl_user u").
				Select("c.nm_cabang").
				Joins("JOIN mst_cabang c ON c.id_cabang = u.id_cabang OR c.nick_cabang = u.id_cabang").
				Where("LOWER(u.username) = LOWER(?)", username).
				Limit(1).Scan(&cPos).Error; err == nil && cPos.NmCabang != "" {
				branchName = cPos.NmCabang
			}
		}
		if branchName == "" && cabangID != "" {
			if err := config.DBMysqlPos.Table("mst_cabang").
				Select("nm_cabang").
				Where("id_cabang = ? OR nick_cabang = ?", cabangID, cabangID).
				Limit(1).Scan(&cPos).Error; err == nil && cPos.NmCabang != "" {
				branchName = cPos.NmCabang
			}
		}
	}

	// 2. If still empty, check employee in Postgres auth.users
	if branchName == "" && config.DBPostgres != nil {
		if username != "" {
			_ = config.DBPostgres.Table("auth.users u").
				Select("b.branch_name").
				Joins("LEFT JOIN master.employee e ON e.employee_id = u.employee_id").
				Joins("LEFT JOIN master.branch b ON e.branch_id = b.branch_id").
				Where("LOWER(u.username) = LOWER(?)", username).
				Limit(1).
				Pluck("b.branch_name", &branchName).Error
		}

		// 3. Fallback to Postgres master.branch
		if branchName == "" && cabangID != "" {
			idInt, errConv := strconv.Atoi(cabangID)
			query := config.DBPostgres.Table("master.branch").Select("branch_name")
			if errConv == nil {
				query = query.Where("branch_id = ? OR branch_code = ?", idInt, cabangID)
			} else {
				query = query.Where("branch_code = ? OR LOWER(branch_name) = LOWER(?)", cabangID, cabangID)
			}
			_ = query.Limit(1).Scan(&branchName).Error
		}
	}

	return branchName
}

// GetCurrentBranchHandler handles GET /api/branch/current
func GetCurrentBranchHandler(c *gin.Context) {
	cabang := c.GetString("cabang")
	username := c.GetString("username")
	branchName := c.GetString("branch_name")

	if branchName == "" {
		branchName = ResolveBranchName(cabang, username)
	}

	c.JSON(http.StatusOK, gin.H{
		"cabang":      cabang,
		"branch_name": branchName,
	})
}

