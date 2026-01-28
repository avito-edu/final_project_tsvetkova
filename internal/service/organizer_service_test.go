package service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"swim_service/internal/domain"
	"swim_service/internal/dto"
	"swim_service/mocks"
)

func TestOrganizerService_CreateOrganizer(t *testing.T) {
	t.Run("Успешное создание организатора", func(t *testing.T) {
		mockRepo := new(mocks.OrganizerRepository)
		service := NewOrganizerService(mockRepo)

		req := dto.CreateOrganizerRequest{
			Name:  "Организатор 1",
			Email: "org1@example.com",
		}

		mockRepo.On("Create", mock.AnythingOfType("*domain.Organizer")).
			Run(func(args mock.Arguments) {
				org := args.Get(0).(*domain.Organizer)
				org.ID = 1
			}).
			Return(nil)

		resp, err := service.CreateOrganizer(req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, 1, resp.ID)
		assert.Equal(t, "Организатор 1", resp.Name)
		assert.Equal(t, "org1@example.com", resp.Email)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Пустой email", func(t *testing.T) {
		mockRepo := new(mocks.OrganizerRepository)
		service := NewOrganizerService(mockRepo)

		req := dto.CreateOrganizerRequest{
			Name:  "Организатор 1",
			Email: "",
		}

		resp, err := service.CreateOrganizer(req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, "email не может быть пустым", err.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("Ошибка репозитория", func(t *testing.T) {
		mockRepo := new(mocks.OrganizerRepository)
		service := NewOrganizerService(mockRepo)

		req := dto.CreateOrganizerRequest{
			Name:  "Организатор 1",
			Email: "org1@example.com",
		}

		mockRepo.On("Create", mock.AnythingOfType("*domain.Organizer")).
			Return(errors.New("unique constraint violation"))

		resp, err := service.CreateOrganizer(req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, "unique constraint violation", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestOrganizerService_GetOrganizer(t *testing.T) {
	t.Run("Успешное получение организатора", func(t *testing.T) {
		mockRepo := new(mocks.OrganizerRepository)
		service := NewOrganizerService(mockRepo)

		expectedOrg := &domain.Organizer{
			ID:    1,
			Name:  "Организатор 1",
			Email: "org1@example.com",
		}

		mockRepo.On("GetByID", 1).Return(expectedOrg, nil)

		resp, err := service.GetOrganizer(1)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, 1, resp.ID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Организатор не найден", func(t *testing.T) {
		mockRepo := new(mocks.OrganizerRepository)
		service := NewOrganizerService(mockRepo)

		mockRepo.On("GetByID", 999).Return(nil, nil)

		resp, err := service.GetOrganizer(999)

		assert.NoError(t, err)
		assert.Nil(t, resp)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Неверный ID", func(t *testing.T) {
		mockRepo := new(mocks.OrganizerRepository)
		service := NewOrganizerService(mockRepo)

		resp, err := service.GetOrganizer(0)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, "неверный ID", err.Error())
		mockRepo.AssertNotCalled(t, "GetByID")
	})

	t.Run("Ошибка репозитория", func(t *testing.T) {
		mockRepo := new(mocks.OrganizerRepository)
		service := NewOrganizerService(mockRepo)

		mockRepo.On("GetByID", 1).Return(nil, errors.New("database error"))

		resp, err := service.GetOrganizer(1)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, "database error", err.Error())
		mockRepo.AssertExpectations(t)
	})
}
