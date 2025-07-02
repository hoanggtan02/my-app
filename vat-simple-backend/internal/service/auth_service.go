package service

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/go-sql-driver/mysql"
	"github.com/hoanggtan02/my-app/vat-simple-backend/internal/models"
	"github.com/hoanggtan02/my-app/vat-simple-backend/internal/repository"
	"github.com/hoanggtan02/my-app/vat-simple-backend/internal/utils"
	"github.com/hoanggtan02/my-app/vat-simple-backend/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// AuthService defines the interface for authentication services.
// This makes it easy to mock for testing.
type AuthService interface {
	Register(req *models.RegisterRequest) (string, *errors.AppError)
	Login(req *models.LoginRequest) (string, *errors.AppError)
	GetUserByID(userID string) (*models.User, *errors.AppError)
	UpdateUser(userID string, req *models.UpdateUserProfileRequest) *errors.AppError
}

type authServiceImpl struct {
	userRepo repository.UserRepository
}

// NewAuthService creates a new instance of AuthService.
func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authServiceImpl{
		userRepo: userRepo,
	}
}

// Register handles the business logic for user registration.
func (s *authServiceImpl) Register(req *models.RegisterRequest) (string, *errors.AppError) {
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return "", errors.ErrInternalServerError
	}

	// Create user and company profile in a transaction
	userID, err := s.userRepo.CreateUserAndCompany(req.Email, string(hashedPassword), req.CompanyName)
	if err != nil {
		// Check for duplicate entry error
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			return "", errors.NewAppError(http.StatusConflict, "Email already exists", errors.WithCode("EMAIL_EXISTS"))
		}
		log.Printf("Error creating user and company: %v", err)
		return "", errors.ErrInternalServerError
	}

	return userID, nil
}

// Login handles the business logic for user login.
func (s *authServiceImpl) Login(req *models.LoginRequest) (string, *errors.AppError) {
	// Find user by email
	user, err := s.userRepo.FindUserByEmail(req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.NewAppError(http.StatusUnauthorized, "Invalid credentials", errors.WithCode("INVALID_CREDENTIALS"))
		}
		log.Printf("Error finding user by email: %v", err)
		return "", errors.ErrInternalServerError
	}

	// Compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		// If passwords don't match, return the same error as user not found for security
		return "", errors.NewAppError(http.StatusUnauthorized, "Invalid credentials", errors.WithCode("INVALID_CREDENTIALS"))
	}

	// Generate JWT token
	token, appErr := utils.GenerateToken(user.ID, user.Email, user.CompanyID)
	if appErr != nil {
		log.Printf("Error generating JWT token: %v", appErr)
		return "", appErr
	}

	return token, nil
}

// GetUserByID retrieves a user's profile by their ID.
func (s *authServiceImpl) GetUserByID(userID string) (*models.User, *errors.AppError) {
	user, err := s.userRepo.FindUserByID(userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrNotFound
		}
		log.Printf("Error finding user by ID: %v", err)
		return nil, errors.ErrInternalServerError
	}
	return user, nil
}

// UpdateUser updates a user's profile.
func (s *authServiceImpl) UpdateUser(userID string, req *models.UpdateUserProfileRequest) *errors.AppError {
	// For now, we only update the company name which is in a different table.
	// This logic might need adjustment depending on what user fields are updatable.
	err := s.userRepo.UpdateCompanyNameByUserID(userID, req.CompanyName)
	if err != nil {
		log.Printf("Error updating company name for user ID %s: %v", userID, err)
		return errors.ErrInternalServerError
	}
	return nil
}
