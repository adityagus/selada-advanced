package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// TokenClaims represents the payload of our custom secure token
type TokenClaims struct {
	Username string    `json:"username"`
	Nama     string    `json:"nama"`
	Cabang   string    `json:"cabang"`
	FkUser   string    `json:"fk_user"`
	Expiry   time.Time `json:"expiry"`
}

func getSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "super-secret-key-selada-v2-123456789"
	}
	return []byte(secret)
}

// GenerateToken creates a signed token string for a user session
func GenerateToken(claims TokenClaims) (string, error) {
	claimsBytes, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	// Base64 encode the claims
	payloadEncoded := base64.RawURLEncoding.EncodeToString(claimsBytes)

	// Sign the encoded claims with HMAC-SHA256
	h := hmac.New(sha256.New, getSecret())
	h.Write([]byte(payloadEncoded))
	signature := h.Sum(nil)
	signatureEncoded := base64.RawURLEncoding.EncodeToString(signature)

	// Combine as token: payload.signature
	token := payloadEncoded + "." + signatureEncoded
	return token, nil
}

// ParseAndVerifyToken verifies the token signature and returns claims if valid
func ParseAndVerifyToken(tokenStr string) (*TokenClaims, error) {
	parts := strings.Split(tokenStr, ".")
	if len(parts) != 2 {
		return nil, errors.New("invalid token format")
	}

	payloadEncoded := parts[0]
	signatureEncoded := parts[1]

	// Recalculate signature
	h := hmac.New(sha256.New, getSecret())
	h.Write([]byte(payloadEncoded))
	expectedSignature := h.Sum(nil)
	expectedSignatureEncoded := base64.RawURLEncoding.EncodeToString(expectedSignature)

	// Constant-time comparison to prevent timing attacks (SonarQube compliance)
	if !hmac.Equal([]byte(signatureEncoded), []byte(expectedSignatureEncoded)) {
		return nil, errors.New("invalid token signature")
	}

	// Decode claims
	claimsBytes, err := base64.RawURLEncoding.DecodeString(payloadEncoded)
	if err != nil {
		return nil, err
	}

	var claims TokenClaims
	if err := json.Unmarshal(claimsBytes, &claims); err != nil {
		return nil, err
	}

	// Verify expiration
	if time.Now().After(claims.Expiry) {
		return nil, errors.New("token expired")
	}

	return &claims, nil
}

// AuthMiddleware intercepts requests and checks for a valid Bearer token
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header must be Bearer <token>"})
			c.Abort()
			return
		}

		tokenStr := parts[1]
		claims, err := ParseAndVerifyToken(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: " + err.Error()})
			c.Abort()
			return
		}

		// Store user data in Gin context for controllers
		c.Set("username", claims.Username)
		c.Set("nama", claims.Nama)
		c.Set("cabang", claims.Cabang)
		c.Set("fk_user", claims.FkUser)

		c.Next()
	}
}
