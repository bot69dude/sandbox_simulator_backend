package services

import (
	"context"
	"errors"
	"time"

	"github.com/bot69dude/sandbox_simulator_backend.git/internal/config"
	"github.com/bot69dude/sandbox_simulator_backend.git/internal/dbschema"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	db        *dbschema.Queries
	jwtConfig *config.JWTConfig
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	User        *dbschema.User `json:"user"`
	AccessToken string         `json:"access_token"`
}

func NewAuthService(db *dbschema.Queries, jwtConfig *config.JWTConfig) *AuthService {
	return &AuthService{
		db:        db,
		jwtConfig: jwtConfig,
	}
}

func (s *AuthService) Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error) {
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user
	user, err := s.db.CreateUser(ctx, dbschema.CreateUserParams{
		Email:          req.Email,
		Name:           req.Name,
		HashedPassword: string(hashedPassword),
	})
	if err != nil {
		return nil, err
	}

	// Generate token
	token, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		User:        &user,
		AccessToken: token,
	}, nil
}

func (s *AuthService) Login(ctx context.Context, req LoginRequest) (*AuthResponse, error) {
	// Get user by email
	user, err := s.db.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	// Compare passwords
	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Generate token
	token, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		User:        &user,
		AccessToken: token,
	}, nil
}

func (s *AuthService) generateToken(user dbschema.User) (string, error) {
	// Create claims
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(s.jwtConfig.AccessTokenDuration).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString([]byte(s.jwtConfig.SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *AuthService) ValidateToken(tokenString string) (*jwt.Token, error) {
	// Parse token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.jwtConfig.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

func (s *AuthService) GetUserFromToken(token *jwt.Token) (*dbschema.User, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	// Get user ID from claims
	userID := int32(claims["user_id"].(float64))

	// Get user from database
	user, err := s.db.GetUserByID(context.Background(), userID)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
