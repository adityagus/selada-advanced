package handlers

import (
	"crypto/md5"
	"errors"
	"fmt"
	"net/http"
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database MySQL POS tidak terhubung"})
		return
	}

	var user struct {
		models.UserPos
		NmJabatan string `gorm:"column:nm_jabatan"`
	}

	// Developer / Fallback override password check
	isDevPass := req.Pass == "selada" || req.Pass == "admin" || req.Pass == "stagingdev@6177"

	// Legacy M_login: void = 0 indicates active user
	err := config.DBMysqlPos.Table("tbl_user u").
		Select("u.*, j.nm_jabatan").
		Joins("LEFT JOIN mst_jabatan j ON j.id_jabatan = u.id_jabatan").
		Where("u.username = ? AND u.void = ?", req.User, 0).
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

			if md5Hash == user.Password || req.Pass == user.Password {
				passwordIsValid = true
				// Auto upgrade legacy MD5 password to BCrypt
				if hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Pass), bcrypt.DefaultCost); err == nil {
					config.DBMysqlPos.Table("tbl_user").Where("username = ?", user.Username).Update("password", string(hashedPassword))
				}
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

	// 3. Fetch employee reference from Postgres (sam_live auth.users left join master.employee)
	fkUser := ""
	employeeID := ""
	employeeCode := ""

	if config.DBPostgres != nil {
		var pgUser models.UserPg
		err := config.DBPostgres.Table("auth.users u").
			Select("e.employee_code AS fk_karyawan, e.employee_id AS employee_id, b.branch_name, e.is_active").
			Joins("LEFT JOIN master.employee e ON e.employee_id = u.employee_id").
			Joins("LEFT JOIN master.branch b ON e.branch_id = b.id_branch").
			Where("u.username = ?", user.Username).
			Scan(&pgUser).Error

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

			if employeeID != "" {
				fkUser = employeeID
			} else {
				fkUser = employeeCode
			}
		}
	}

	// 4. Generate custom signed session token
	claims := middleware.TokenClaims{
		Username:   user.Username,
		Nama:       user.Nama,
		Cabang:     user.IDCabang,
		BranchName: user.BranchName,
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
		"branch_name":   user.BranchName,
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
