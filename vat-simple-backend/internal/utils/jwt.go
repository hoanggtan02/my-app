package utils

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hoanggtan02/my-app/vat-simple-backend/internal/config"
	apperrors "github.com/hoanggtan02/my-app/vat-simple-backend/pkg/errors" // TYPO FIXED HERE
)

// Claims represents the JWT claims
type Claims struct {
	UserID    string `json:"user_id"`
	CompanyID string `json:"company_id,omitempty"`
	Email     string `json:"email"`
	jwt.RegisteredClaims
}

// JWTSecretKey stores the JWT secret key
var JWTSecretKey []byte

// InitJWT initializes the JWT secret key from config
func InitJWT(cfg *config.Config) {
	JWTSecretKey = []byte(cfg.JWTSecret)
	if len(JWTSecretKey) == 0 {
		log.Fatalf("JWT_SECRET is empty. Please set it in .env file.")
	}
	log.Println("JWT secret key initialized.")
}

// GenerateToken generates a new JWT token
func GenerateToken(userID, email, companyID string) (string, *apperrors.AppError) {
	if len(JWTSecretKey) == 0 {
		return "", apperrors.NewAppError(http.StatusInternalServerError, "JWT secret key not initialized", apperrors.WithCode("JWT_ERROR"))
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID:    userID,
		CompanyID: companyID,
		Email:     email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "vat-simple-backend",
			Subject:   userID,
			ID:        fmt.Sprintf("%d", time.Now().UnixNano()),
			Audience:  []string{"users"},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JWTSecretKey)
	if err != nil {
		return "", apperrors.NewAppError(http.StatusInternalServerError, "Failed to sign token", apperrors.WithCode("JWT_SIGN_ERROR"), apperrors.WithCause(err))
	}

	return tokenString, nil
}

// ParseToken parses and validates a JWT token
func ParseToken(tokenString string) (*Claims, *apperrors.AppError) {
	if len(JWTSecretKey) == 0 {
		return nil, apperrors.NewAppError(http.StatusInternalServerError, "JWT secret key not initialized", apperrors.WithCode("JWT_ERROR"))
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return JWTSecretKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, apperrors.NewAppError(http.StatusUnauthorized, "Invalid token format", apperrors.WithCode("TOKEN_MALFORMED"))
		} else if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
			return nil, apperrors.NewAppError(http.StatusUnauthorized, "Token is expired or not valid yet", apperrors.WithCode("TOKEN_EXPIRED"))
		}
		// For any other parsing error
		return nil, apperrors.NewAppError(http.StatusUnauthorized, "Invalid token", apperrors.WithCode("TOKEN_INVALID"), apperrors.WithCause(err))
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, apperrors.NewAppError(http.StatusUnauthorized, "Invalid token claims", apperrors.WithCode("INVALID_CLAIMS"))
}
