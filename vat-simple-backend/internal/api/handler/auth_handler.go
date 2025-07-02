package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hoanggtan02/my-app/vat-simple-backend/internal/models"  // Import models package
	"github.com/hoanggtan02/my-app/vat-simple-backend/internal/service" // Import service package
	"github.com/hoanggtan02/my-app/vat-simple-backend/internal/utils"   // Import utils package
	"github.com/hoanggtan02/my-app/vat-simple-backend/pkg/errors"       // Import errors package
)

// AuthHandler handles authentication related HTTP requests
type AuthHandler struct {
	AuthService service.AuthService
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		AuthService: authService,
	}
}

// Register handles user registration
// @Summary Register a new user
// @Description Register a new user with email, password, and company name
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body models.RegisterRequest true "User registration details"
// @Success 201 {object} map[string]string "message: User registered successfully"
// @Failure 400 {object} errors.AppError "Invalid input"
// @Failure 409 {object} errors.AppError "Email already exists"
// @Failure 500 {object} errors.AppError "Internal server error"
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := utils.ValidateStruct(req) // Use our custom validator
		if appErr != nil {
			c.JSON(appErr.StatusCode, appErr)
			return
		}
		// Fallback for other binding errors (e.g., malformed JSON)
		c.JSON(http.StatusBadRequest, errors.NewAppError(http.StatusBadRequest, "Invalid request body", errors.WithCause(err)))
		return
	}

	userID, appErr := h.AuthService.Register(&req)
	if appErr != nil {
		c.JSON(appErr.StatusCode, appErr)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully", "user_id": userID})
}

// Login handles user login and returns a JWT token
// @Summary Log in a user
// @Description Authenticate user with email and password, return JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param credentials body models.LoginRequest true "User login credentials"
// @Success 200 {object} map[string]string "access_token: JWT token"
// @Failure 400 {object} errors.AppError "Invalid input"
// @Failure 401 {object} errors.AppError "Invalid credentials"
// @Failure 500 {object} errors.AppError "Internal server error"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := utils.ValidateStruct(req) // Use our custom validator
		if appErr != nil {
			c.JSON(appErr.StatusCode, appErr)
			return
		}
		c.JSON(http.StatusBadRequest, errors.NewAppError(http.StatusBadRequest, "Invalid request body", errors.WithCause(err)))
		return
	}

	token, appErr := h.AuthService.Login(&req)
	if appErr != nil {
		c.JSON(appErr.StatusCode, appErr)
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": token, "token_type": "bearer"})
}

// GetUserProfile handles fetching the current authenticated user's profile
// @Summary Get current user profile
// @Description Retrieve the profile information of the authenticated user
// @Tags Users
// @Security BearerAuth
// @Produce json
// @Success 200 {object} models.User "User profile data"
// @Failure 401 {object} errors.AppError "Unauthorized"
// @Failure 500 {object} errors.AppError "Internal server error"
// @Router /users/me [get]
func (h *AuthHandler) GetUserProfile(c *gin.Context) {
	// UserID is set by AuthMiddleware
	userID, exists := c.Get("userID")
	if !exists {
		log.Println("UserID not found in context, AuthMiddleware might be missing or failed.")
		c.JSON(http.StatusInternalServerError, errors.ErrInternalServerError)
		return
	}

	user, appErr := h.AuthService.GetUserByID(userID.(string))
	if appErr != nil {
		c.JSON(appErr.StatusCode, appErr)
		return
	}

	// Do not return password hash
	user.PasswordHash = ""
	c.JSON(http.StatusOK, user)
}

// UpdateUserProfile handles updating the current authenticated user's profile
// @Summary Update current user profile
// @Description Update the company name for the authenticated user
// @Tags Users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param user body models.UpdateUserProfileRequest true "User profile update details"
// @Success 200 {object} map[string]string "message: User information updated successfully"
// @Failure 400 {object} errors.AppError "Invalid input"
// @Failure 401 {object} errors.AppError "Unauthorized"
// @Failure 500 {object} errors.AppError "Internal server error"
// @Router /users/me [put]
func (h *AuthHandler) UpdateUserProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		log.Println("UserID not found in context, AuthMiddleware might be missing or failed.")
		c.JSON(http.StatusInternalServerError, errors.ErrInternalServerError)
		return
	}

	var req models.UpdateUserProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := utils.ValidateStruct(req)
		if appErr != nil {
			c.JSON(appErr.StatusCode, appErr)
			return
		}
		c.JSON(http.StatusBadRequest, errors.NewAppError(http.StatusBadRequest, "Invalid request body", errors.WithCause(err)))
		return
	}

	appErr := h.AuthService.UpdateUser(userID.(string), &req)
	if appErr != nil {
		c.JSON(appErr.StatusCode, appErr)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User information updated successfully"})
}

// TODO: Implement ForgotPassword and ResetPassword handlers later
