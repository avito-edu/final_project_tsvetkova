package service

import (
	"context"
	"errors"
	"fmt"
	"swim_service/internal/domain"
	"swim_service/internal/dto"
	"swim_service/pkg/jwtutil"

	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	ExistsByLogin(ctx context.Context, login string) (bool, error)
	FindByLogin(ctx context.Context, login string) (*domain.User, error)
}

type UserService struct {
	userRepo     UserRepository
	tokenService *jwtutil.TokenService
}

func NewUserService(repo UserRepository, ts *jwtutil.TokenService) *UserService {
	return &UserService{
		userRepo:     repo,
		tokenService: ts,
	}
}

func (s *UserService) RegisterUser(ctx context.Context, req dto.RegisterUserRequest) (*domain.User, error) {
	exists, err := s.userRepo.ExistsByLogin(ctx, req.Login)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("user with this login already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &domain.User{
		Login:    req.Login,
		Password: string(hashedPassword),
		Role:     domain.RoleAthlet,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) RegisterSpecial(ctx context.Context, req dto.RegisterSpecialRequest) (*domain.User, error) {
	exists, err := s.userRepo.ExistsByLogin(ctx, req.Login)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("user with this login already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	role := domain.Role(req.Role)

	user := &domain.User{
		Login:    req.Login,
		Password: string(hashedPassword),
		Role:     role,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Login(ctx context.Context, req dto.LoginRequest) (string, error) {
	user, err := s.userRepo.FindByLogin(ctx, req.Login)
	if err != nil {
		return "", errors.New("invalid login or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return "", errors.New("invalid login or password")
	}

	token, err := s.tokenService.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, nil
}
