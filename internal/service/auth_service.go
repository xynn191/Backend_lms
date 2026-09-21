package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/xynn191/Backend_lms/internal/config"
	"github.com/xynn191/Backend_lms/internal/models"
	"github.com/xynn191/Backend_lms/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type TokenResponse struct {
	Token     string      `json:"token"`
	ExpiresAt time.Time   `json:"expires_at"`
	User      UserProfile `json:"user"`
}

type UserProfile struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type JWTClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type AuthService interface {
	Login(req *LoginRequest) (*TokenResponse, error)
	GetProfile(userID string) (*models.User, error)
}

type authService struct {
	repo repository.AuthRepository
	cfg  *config.Config
}

func NewAuthService(repo repository.AuthRepository, cfg *config.Config) AuthService {
	return &authService{repo, cfg}
}

func (s *authService) Login(req *LoginRequest) (*TokenResponse, error) {
	user, err := s.repo.FindUserByEmail(req.Email)
	if err != nil {
		return nil, errors.New("email atau password salah")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("email atau password salah")
	}

	expiresAt := time.Now().Add(s.cfg.JWTExpiresIn)
	claims := &JWTClaims{
		UserID: user.ID,
		Role:   user.Role.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return nil, errors.New("gagal membuat token")
	}

	return &TokenResponse{
		Token:     tokenString,
		ExpiresAt: expiresAt,
		User: UserProfile{
			ID:    user.ID,
			Email: user.Email,
			Role:  user.Role.Name,
		},
	}, nil
}

func (s *authService) GetProfile(userID string) (*models.User, error) {
	return s.repo.FindUserByID(userID)
}
