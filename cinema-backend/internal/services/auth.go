package services

import (
	"cinema-backend/internal/models"
	"cinema-backend/internal/repository"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo     *repository.UserRepository
	jwtSecret    string
	emailService *EmailService
	frontendURL  string
}

func NewAuthService(userRepo *repository.UserRepository, jwtSecret string, emailService *EmailService, frontendURL string) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		jwtSecret:    jwtSecret,
		emailService: emailService,
		frontendURL:  frontendURL,
	}
}

const passwordResetTokenTTL = 30 * time.Minute

var ErrUserAlreadyExists = errors.New("an account with this email already exists. Sign in or reset your password")

type Claims struct {
	UserID uuid.UUID `json:"userId"`
	Email  string    `json:"email"`
	Role   string    `json:"role"`
	jwt.RegisteredClaims
}

func (s *AuthService) Register(name, email, password, phone string) (*models.User, string, error) {
	user, err := s.CreateUser(name, email, password, phone, models.RoleCustomer)
	if err != nil {
		return nil, "", err
	}

	// Generate token
	token, err := s.generateToken(user)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *AuthService) CreateUser(name, email, password, phone string, role models.UserRole) (*models.User, error) {
	normalizedEmail := strings.ToLower(strings.TrimSpace(email))
	if existingUser, _ := s.userRepo.FindByEmailIncludingDeleted(normalizedEmail); existingUser != nil {
		return nil, ErrUserAlreadyExists
	}

	if role != models.RoleAdmin && role != models.RoleCustomer {
		return nil, errors.New("invalid user role")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Name:     strings.TrimSpace(name),
		Email:    normalizedEmail,
		Password: string(hashedPassword),
		Role:     role,
	}
	if phone != "" {
		user.Phone = &phone
	}

	if err := s.userRepo.Create(user); err != nil {
		if existingUser, _ := s.userRepo.FindByEmailIncludingDeleted(normalizedEmail); existingUser != nil {
			return nil, ErrUserAlreadyExists
		}
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(email, password string) (*models.User, string, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	// Generate token
	token, err := s.generateToken(user)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *AuthService) RequestPasswordReset(email string) (string, error) {
	normalizedEmail := strings.ToLower(strings.TrimSpace(email))
	user, err := s.userRepo.FindByEmail(normalizedEmail)
	if err != nil {
		// Do not reveal whether an account exists.
		return "", nil
	}

	rawToken := make([]byte, 32)
	if _, err := rand.Read(rawToken); err != nil {
		return "", err
	}
	resetToken := url.QueryEscape(encodeResetToken(rawToken))
	if err := s.userRepo.DeletePasswordResetTokens(user.ID); err != nil {
		return "", err
	}

	if err := s.userRepo.CreatePasswordResetToken(&models.PasswordResetToken{
		UserID:    user.ID,
		TokenHash: hashResetToken(resetToken),
		ExpiresAt: time.Now().Add(passwordResetTokenTTL),
	}); err != nil {
		return "", err
	}

	resetURL := fmt.Sprintf("%s/auth/reset-password?token=%s", strings.TrimRight(s.frontendURL, "/"), resetToken)
	if s.emailService == nil {
		log.Printf("Password reset email service is unavailable for %s", normalizedEmail)
		return resetURL, nil
	}
	if err := s.emailService.SendPasswordReset(user.Email, user.Name, resetURL); err != nil {
		log.Printf("Failed to send password reset email: %v", err)
	}

	if !s.emailService.IsConfigured() && getEnv("APP_ENV", "development") != "production" {
		return resetURL, nil
	}

	return "", nil
}

func (s *AuthService) ResetPassword(token, newPassword string) error {
	token = strings.TrimSpace(token)
	if token == "" {
		return errors.New("invalid or expired reset token")
	}

	resetToken, err := s.userRepo.ConsumePasswordResetToken(hashResetToken(token))
	if err != nil {
		return errors.New("invalid or expired reset token")
	}

	user, err := s.userRepo.FindByID(resetToken.UserID)
	if err != nil {
		return errors.New("invalid or expired reset token")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	return s.userRepo.Update(user)
}

func encodeResetToken(token []byte) string {
	return hex.EncodeToString(token)
}

func hashResetToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func (s *AuthService) generateToken(user *models.User) (string, error) {
	claims := Claims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   string(user.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func (s *AuthService) GetUserByID(userID uuid.UUID) (*models.User, error) {
	return s.userRepo.FindByID(userID)
}

func (s *AuthService) GetAllUsers(page, limit int) ([]models.User, int64, error) {
	return s.userRepo.FindAll(page, limit)
}

func (s *AuthService) GetCustomers(page, limit int) ([]models.User, int64, error) {
	return s.userRepo.FindByRole(models.RoleCustomer, page, limit)
}

func (s *AuthService) DeleteUser(userID uuid.UUID) error {
	return s.userRepo.Delete(userID)
}

func (s *AuthService) UpdateUser(userID uuid.UUID, name, email, phone string, role models.UserRole, password string) (*models.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	normalizedEmail := strings.ToLower(strings.TrimSpace(email))
	if normalizedEmail != user.Email {
		otherUser, lookupErr := s.userRepo.FindByEmailIncludingDeleted(normalizedEmail)
		if lookupErr == nil && otherUser.ID != userID {
			return nil, errors.New("email is already in use")
		}
	}
	if role != models.RoleAdmin && role != models.RoleCustomer {
		return nil, errors.New("invalid user role")
	}

	user.Name = strings.TrimSpace(name)
	user.Email = normalizedEmail
	user.Role = role
	if phone == "" {
		user.Phone = nil
	} else {
		user.Phone = &phone
	}
	if password != "" {
		hashedPassword, hashErr := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if hashErr != nil {
			return nil, hashErr
		}
		user.Password = string(hashedPassword)
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *AuthService) UpdateProfile(userID uuid.UUID, name, phone string) (*models.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	user.Name = name
	if phone != "" {
		user.Phone = &phone
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) ChangePassword(userID uuid.UUID, currentPassword, newPassword string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(currentPassword)); err != nil {
		return errors.New("current password is incorrect")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	return s.userRepo.Update(user)
}
