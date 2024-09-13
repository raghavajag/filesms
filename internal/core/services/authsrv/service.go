package authsrv

import (
	"context"
	"errors"
	"filesms/internal/core/domain"
	"filesms/internal/core/ports"
	"fmt"
	"time"

	"filesms/pkg/jwt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo ports.UserRepository
	jwtMaker jwt.Maker
}

func NewAuthService(userRepo ports.UserRepository, jwtMaker jwt.Maker) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		jwtMaker: jwtMaker,
	}
}
func (s *AuthService) Register(ctx context.Context, email, password string) (*domain.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Email:     email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		ID:        uuid.New(),
	}
	fmt.Printf("%+v", user)
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}
func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	token, err := s.jwtMaker.CreateToken(user.ID, time.Hour*24)
	if err != nil {
		return "", err
	}

	return token, nil
}
