package repository

import (
	"database/sql"
	"log"

	"github.com/google/uuid"
	"github.com/hoanggtan02/my-app/vat-simple-backend/internal/models"
)

// UserRepository defines the interface for user data operations.
type UserRepository interface {
	CreateUserAndCompany(email, passwordHash, companyName string) (string, error)
	FindUserByEmail(email string) (*models.User, error)
	FindUserByID(userID string) (*models.User, error)
	UpdateCompanyNameByUserID(userID, companyName string) error
}

type userRepositoryImpl struct {
	db *sql.DB
}

// NewUserRepository creates a new instance of UserRepository.
func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepositoryImpl{db: db}
}

// CreateUserAndCompany creates a new user and a new company profile within a single transaction.
func (r *userRepositoryImpl) CreateUserAndCompany(email, passwordHash, companyName string) (string, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return "", err
	}
	defer tx.Rollback() // Rollback on error

	// 1. Create Company Profile
	companyID := uuid.New().String()
	companyQuery := "INSERT INTO company_profiles (id, name) VALUES (?, ?)"
	_, err = tx.Exec(companyQuery, companyID, companyName)
	if err != nil {
		return "", err
	}

	// 2. Create User
	userID := uuid.New().String()
	userQuery := "INSERT INTO users (id, email, password_hash, company_id) VALUES (?, ?, ?, ?)"
	_, err = tx.Exec(userQuery, userID, email, passwordHash, companyID)
	if err != nil {
		return "", err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return "", err
	}

	return userID, nil
}

// FindUserByEmail finds a user by their email address.
func (r *userRepositoryImpl) FindUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := "SELECT id, email, password_hash, company_id, created_at, updated_at FROM users WHERE email = ? LIMIT 1"
	err := r.db.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CompanyID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err // Returns sql.ErrNoRows if not found
	}
	return user, nil
}

// FindUserByID finds a user by their ID.
func (r *userRepositoryImpl) FindUserByID(userID string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT u.id, u.email, u.company_id, cp.name as company_name, u.created_at, u.updated_at
		FROM users u
		JOIN company_profiles cp ON u.company_id = cp.id
		WHERE u.id = ?
		LIMIT 1
	`
	err := r.db.QueryRow(query, userID).Scan(&user.ID, &user.Email, &user.CompanyID, &user.CompanyName, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// UpdateCompanyNameByUserID updates a company's name based on the user ID.
func (r *userRepositoryImpl) UpdateCompanyNameByUserID(userID, companyName string) error {
	query := `
		UPDATE company_profiles cp
		JOIN users u ON cp.id = u.company_id
		SET cp.name = ?
		WHERE u.id = ?
	`
	result, err := r.db.Exec(query, companyName, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		log.Printf("No company profile found to update for user ID: %s", userID)
		// Not returning an error, as the operation itself didn't fail.
	}

	return nil
}
