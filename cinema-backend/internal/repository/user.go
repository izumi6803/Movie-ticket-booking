package repository

import (
	"cinema-backend/internal/models"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindAll(page, limit int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	offset := (page - 1) * limit

	err := r.db.Model(&models.User{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Offset(offset).Limit(limit).Find(&users).Error
	return users, total, err
}

func (r *UserRepository) FindByRole(role models.UserRole, page, limit int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	offset := (page - 1) * limit

	err := r.db.Model(&models.User{}).Where("role = ?", role).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Where("role = ?", role).Offset(offset).Limit(limit).Find(&users).Error
	return users, total, err
}

func (r *UserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) DeletePasswordResetTokens(userID uuid.UUID) error {
	return r.db.Where("user_id = ? OR expires_at <= ? OR used_at IS NOT NULL", userID, time.Now()).
		Delete(&models.PasswordResetToken{}).Error
}

func (r *UserRepository) CreatePasswordResetToken(token *models.PasswordResetToken) error {
	return r.db.Create(token).Error
}

func (r *UserRepository) ConsumePasswordResetToken(tokenHash string) (*models.PasswordResetToken, error) {
	var token models.PasswordResetToken
	if err := r.db.Where("token_hash = ? AND used_at IS NULL AND expires_at > ?", tokenHash, time.Now()).
		First(&token).Error; err != nil {
		return nil, err
	}

	now := time.Now()
	result := r.db.Model(&models.PasswordResetToken{}).
		Where("id = ? AND used_at IS NULL", token.ID).
		Update("used_at", now)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected != 1 {
		return nil, errors.New("reset token has already been used")
	}

	token.UsedAt = &now
	return &token, nil
}

func (r *UserRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.User{}, "id = ?", id).Error
}
