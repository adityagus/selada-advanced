package handlers

import (
	"crypto/md5"
	"errors"
	"fmt"
	"net/http"
	"strconv"
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
			fmt.Printf("md5Hash %v\n", md5Hash)
			fmt.Printf("user.Password %v\n", user.Password)
			fmt.Printf("req.Pass %v\n", req.Pass)
			fmt.Printf("hasher %v\n", hasher)

			if md5Hash == user.Password {
				passwordIsValid = true
				// Auto upgrade legacy MD5 password to BCrypt | jangan diupgrade
				// if hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Pass), bcrypt.DefaultCost); err == nil {
				// 	config.DBMysqlPos.Table("tbl_user").Where("username = ?", user.Username).Update("password", string(hashedPassword))
				// }
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
	branchName := ""

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
			if pgUser.BranchName != nil {
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

// ResolveBranchName resolves the textual branch name from Postgres master.branch or MySQL POS
func ResolveBranchName(cabangID string, username string) string {
	branchName := ""

	if config.DBPostgres != nil {
		// 1. Try resolving via employee associated with the username in auth.users
		if username != "" {
			_ = config.DBPostgres.Table("auth.users u").
				Select("b.branch_name").
				Joins("LEFT JOIN master.employee e ON e.employee_id = u.employee_id").
				Joins("LEFT JOIN master.branch b ON e.branch_id = b.branch_id").
				Where("LOWER(u.username) = LOWER(?)", username).
				Limit(1).
				Scan(&branchName).Error
		}

		// 2. Fallback to master.branch lookup using cabangID (branch_id or branch_code)
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

	// 3. Fallback to MySQL POS (sam_pos) tbl_cabang or mst_cabang
	if branchName == "" && cabangID != "" && config.DBMysqlPos != nil {
		type CabangPos struct {
			NmCabang string `gorm:"column:nm_cabang"`
		}
		var cPos CabangPos
		if err := config.DBMysqlPos.Table("tbl_cabang").
			Select("nm_cabang").
			Where("id_cabang = ? OR kd_cabang = ?", cabangID, cabangID).
			Limit(1).Scan(&cPos).Error; err == nil && cPos.NmCabang != "" {
			branchName = cPos.NmCabang
		} else if err := config.DBMysqlPos.Table("mst_cabang").
			Select("nm_cabang").
			Where("id_cabang = ? OR kd_cabang = ?", cabangID, cabangID).
			Limit(1).Scan(&cPos).Error; err == nil && cPos.NmCabang != "" {
			branchName = cPos.NmCabang
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

