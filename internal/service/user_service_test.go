package service

import (
	"context"
	"errors"
	"testing"

	"swim_service/internal/domain"
	"swim_service/internal/dto"
	"swim_service/mocks"
	"swim_service/pkg/jwtutil"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestUserService_RegisterUser_Success(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	tokenService := jwtutil.NewTokenService([]byte("test-secret"))

	service := NewUserService(mockRepo, tokenService)

	login := "testuser"
	password := "password123"

	mockRepo.On("ExistsByLogin", mock.Anything, login).Return(false, nil)
	mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(u *domain.User) bool {
		return u.Login == login &&
			bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) == nil &&
			u.Role == domain.RoleAthlet
	})).Return(nil)

	req := dto.RegisterUserRequest{
		Login:    login,
		Password: password,
	}

	user, err := service.RegisterUser(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, login, user.Login)
	assert.Equal(t, domain.RoleAthlet, user.Role)
	mockRepo.AssertExpectations(t)
}

func TestUserService_RegisterUser_UserExists(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	tokenService := jwtutil.NewTokenService([]byte("test-secret"))
	service := NewUserService(mockRepo, tokenService)

	login := "existinguser"
	mockRepo.On("ExistsByLogin", mock.Anything, login).Return(true, nil)

	req := dto.RegisterUserRequest{
		Login:    login,
		Password: "password123",
	}

	_, err := service.RegisterUser(context.Background(), req)

	assert.Error(t, err)
	assert.Equal(t, "user with this login already exists", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestUserService_RegisterUser_ExistsCheckError(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	tokenService := jwtutil.NewTokenService([]byte("test-secret"))
	service := NewUserService(mockRepo, tokenService)

	expectedErr := errors.New("database error")
	mockRepo.On("ExistsByLogin", mock.Anything, "testuser").Return(false, expectedErr)

	req := dto.RegisterUserRequest{
		Login:    "testuser",
		Password: "password123",
	}

	_, err := service.RegisterUser(context.Background(), req)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	mockRepo.AssertExpectations(t)
}

func TestUserService_RegisterUser_CreateError(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	tokenService := jwtutil.NewTokenService([]byte("test-secret"))
	service := NewUserService(mockRepo, tokenService)

	login := "newuser"
	mockRepo.On("ExistsByLogin", mock.Anything, login).Return(false, nil)
	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.User")).Return(errors.New("db create failed"))

	req := dto.RegisterUserRequest{
		Login:    login,
		Password: "password123",
	}

	_, err := service.RegisterUser(context.Background(), req)

	assert.Error(t, err)
	assert.Equal(t, "db create failed", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestUserService_RegisterSpecial_Success(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	tokenService := jwtutil.NewTokenService([]byte("test-secret"))
	service := NewUserService(mockRepo, tokenService)

	login := "coach123"
	password := "pass"
	role := "coach"

	mockRepo.On("ExistsByLogin", mock.Anything, login).Return(false, nil)
	mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(u *domain.User) bool {
		return u.Login == login &&
			bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) == nil &&
			u.Role == domain.Role(role)
	})).Return(nil)

	req := dto.RegisterSpecialRequest{
		Login:    login,
		Password: password,
		Role:     role,
	}

	user, err := service.RegisterSpecial(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, login, user.Login)
	assert.Equal(t, domain.Role(role), user.Role)
	mockRepo.AssertExpectations(t)
}

func TestUserService_Login_Success(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	tokenService := jwtutil.NewTokenService([]byte("test-secret"))
	service := NewUserService(mockRepo, tokenService)

	login := "athlete1"
	password := "mypass"
	hashedPass, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user := &domain.User{
		ID:       123,
		Login:    login,
		Password: string(hashedPass),
		Role:     domain.RoleAthlet,
	}

	mockRepo.On("FindByLogin", mock.Anything, login).Return(user, nil)

	req := dto.LoginRequest{
		Login:    login,
		Password: password,
	}

	token, err := service.Login(context.Background(), req)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	parsedToken, _ := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte("test-secret"), nil
	})
	assert.True(t, parsedToken.Valid)

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	assert.True(t, ok)
	assert.Equal(t, float64(123), claims["user_id"])
	assert.Equal(t, "athlete", claims["role"])
	mockRepo.AssertExpectations(t)
}

func TestUserService_Login_InvalidPassword(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	tokenService := jwtutil.NewTokenService([]byte("test-secret"))
	service := NewUserService(mockRepo, tokenService)

	login := "user1"
	wrongPass := "wrongpass"
	correctHash, _ := bcrypt.GenerateFromPassword([]byte("correctpass"), bcrypt.DefaultCost)

	user := &domain.User{
		Login:    login,
		Password: string(correctHash),
		Role:     domain.RoleAthlet,
	}

	mockRepo.On("FindByLogin", mock.Anything, login).Return(user, nil)

	req := dto.LoginRequest{
		Login:    login,
		Password: wrongPass,
	}

	_, err := service.Login(context.Background(), req)

	assert.Error(t, err)
	assert.Equal(t, "invalid login or password", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestUserService_Login_UserNotFound(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	tokenService := jwtutil.NewTokenService([]byte("test-secret"))
	service := NewUserService(mockRepo, tokenService)

	login := "nonexistent"
	mockRepo.On("FindByLogin", mock.Anything, login).Return(nil, errors.New("not found"))

	req := dto.LoginRequest{
		Login:    login,
		Password: "anything",
	}

	_, err := service.Login(context.Background(), req)

	assert.Error(t, err)
	assert.Equal(t, "invalid login or password", err.Error())
	mockRepo.AssertExpectations(t)
}
