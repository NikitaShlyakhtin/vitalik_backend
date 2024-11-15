package auth_service

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"vitalik_backend/internal/dependencies"
	auth_types "vitalik_backend/internal/pkg/services/auth_service/types"
	"vitalik_backend/internal/pkg/services/store"
	"vitalik_backend/internal/pkg/types"
)

var (
	ErrTokenExpired  = errors.New("Token expired")
	ErrUserNotFound  = errors.New("User not found")
	ErrAlreadyExists = errors.New("User already exists")
)

type AuthService struct {
	store dependencies.IStore
}

func NewAuthService(store dependencies.IStore) (*AuthService, error) {
	if store == nil {
		return nil, errors.New("failed to initialize auth service")
	}

	return &AuthService{
		store: store,
	}, nil
}

var _ dependencies.IAuthService = (*AuthService)(nil)

func (s *AuthService) Register(ctx context.Context, args auth_types.AuthArgs) error {
	if args.UserID == "" || args.Password == "" {
		return errors.New("invalid user_id or password")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(args.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}

	user := types.User{
		UserID:         args.UserID,
		HashedPassword: string(hashedPassword),
	}

	if err = s.store.SaveUser(ctx, user); err != nil {
		if errors.Is(err, store.ErrAlreadyExists) {
			return fmt.Errorf("user already exists: %w", ErrAlreadyExists)
		}
		return fmt.Errorf("failed to save user: %w", err)
	}

	return nil
}

var jwtSecret = []byte("Innopolis")

func (s *AuthService) Login(ctx context.Context, args auth_types.AuthArgs) (*auth_types.LoginResponse, error) {
	if args.UserID == "" || args.Password == "" {
		return nil, errors.New("invalid user_id or password")
	}

	user, err := s.store.GetUser(ctx, args.UserID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return nil, fmt.Errorf("user not found: %w", ErrUserNotFound)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(args.Password)); err != nil {
		return nil, errors.New("invalid password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.UserID,
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to sign token: %w", err)
	}

	return &auth_types.LoginResponse{
		Token: tokenString,
	}, nil
}

func (s *AuthService) Authenticate(ctx context.Context, tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return "", errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", errors.New("invalid token claims")
	}

	userID, ok := claims["user_id"].(string)
	if !ok || userID == "" {
		return "", fmt.Errorf("invalid user ID in token claims: %w", ErrUserNotFound)
	}

	_, err = s.store.GetUser(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("user not found: %w", ErrUserNotFound)
	}

	return userID, nil
}
