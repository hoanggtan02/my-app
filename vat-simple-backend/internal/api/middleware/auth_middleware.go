package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hoanggtan02/my-app/vat-simple-backend/internal/utils" // Import utilities for JWT
	"github.com/hoanggtan02/my-app/vat-simple-backend/pkg/errors"     // Import custom errors
)

// AuthMiddleware is a Gin middleware to authenticate requests using JWT
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			err := errors.ErrUnauthorized
			c.JSON(err.StatusCode, err)
			c.Abort() // Stop processing this request
			return
		}

		// Expected format: "Bearer <token>"
		tokenString := ""
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenString = authHeader[7:]
		} else {
			err := errors.NewAppError(http.StatusUnauthorized, "Invalid authorization header format", errors.WithCode("INVALID_AUTH_HEADER"))
			c.JSON(err.StatusCode, err)
			c.Abort()
			return
		}

		claims, appErr := utils.ParseToken(tokenString)
		if appErr != nil {
			log.Printf("Failed to parse or validate token: %v\n", appErr.Error())
			c.JSON(appErr.StatusCode, appErr)
			c.Abort()
			return
		}

		// Store user ID and other claims in context for handlers to access
		c.Set("userID", claims.UserID)
		c.Set("email", claims.Email)
		if claims.CompanyID != "" {
			c.Set("companyID", claims.CompanyID)
		}

		c.Next() // Continue to the next handler
	}
}
