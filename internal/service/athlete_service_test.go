package service

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"swim_service/internal/domain"
	"swim_service/internal/dto"
	"swim_service/mocks"
)

func TestAthleteService_CreateAthlete(t *testing.T) {
	t.Run("Успешное создание спортсмена", func(t *testing.T) {
		mockRepo := new(mocks.AthleteRepository)
		service := NewAthleteService(mockRepo)

		birthDate := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
		req := dto.CreateAthleteRequest{
			Name:      "Иван Иванов",
			BirthDate: &birthDate,
			Status:    "active",
			Gender:    "male",
		}

		mockRepo.On("CreateAthlete", mock.AnythingOfType("*domain.Athlete")).
			Run(func(args mock.Arguments) {
				athlete := args.Get(0).(*domain.Athlete)
				athlete.ID = "123"
			}).
			Return(nil)

		resp, err := service.CreateAthlete(req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "123", resp.ID)
		assert.Equal(t, "Иван Иванов", resp.Name)
		assert.Equal(t, "active", resp.Status)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Ошибка при создании в репозитории", func(t *testing.T) {
		mockRepo := new(mocks.AthleteRepository)
		service := NewAthleteService(mockRepo)

		req := dto.CreateAthleteRequest{
			Name:   "Иван Иванов",
			Status: "active",
			Gender: "male",
		}

		mockRepo.On("CreateAthlete", mock.AnythingOfType("*domain.Athlete")).
			Return(errors.New("database error"))

		resp, err := service.CreateAthlete(req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, "database error", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestAthleteService_GetAthlete(t *testing.T) {
	t.Run("Успешное получение спортсмена", func(t *testing.T) {
		mockRepo := new(mocks.AthleteRepository)
		service := NewAthleteService(mockRepo)

		birthDate := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
		expectedAthlete := &domain.Athlete{
			ID:        "123",
			Name:      "Иван Иванов",
			BirthDate: &birthDate,
			Status:    "active",
			Gender:    "male",
		}

		mockRepo.On("GetByID", 123).Return(expectedAthlete, nil)

		resp, err := service.GetAthlete(123)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "123", resp.ID)
		assert.Equal(t, "Иван Иванов", resp.Name)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Спортсмен не найден", func(t *testing.T) {
		mockRepo := new(mocks.AthleteRepository)
		service := NewAthleteService(mockRepo)

		mockRepo.On("GetByID", 999).Return(nil, nil)

		resp, err := service.GetAthlete(999)

		assert.NoError(t, err)
		assert.Nil(t, resp)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Ошибка репозитория", func(t *testing.T) {
		mockRepo := new(mocks.AthleteRepository)
		service := NewAthleteService(mockRepo)

		mockRepo.On("GetByID", 123).Return(nil, errors.New("database error"))

		resp, err := service.GetAthlete(123)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, "database error", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestAthleteService_UpdateAthleteStatus(t *testing.T) {
	t.Run("Успешное обновление статуса", func(t *testing.T) {
		mockRepo := new(mocks.AthleteRepository)
		service := NewAthleteService(mockRepo)

		birthDate := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
		existingAthlete := &domain.Athlete{
			ID:        "123",
			Name:      "Иван Иванов",
			BirthDate: &birthDate,
			Status:    "active",
			Gender:    "male",
		}

		updatedAthlete := &domain.Athlete{
			ID:        "123",
			Name:      "Иван Иванов",
			BirthDate: &birthDate,
			Status:    "inactive",
			Gender:    "male",
		}

		mockRepo.On("GetByID", 123).Return(existingAthlete, nil)
		mockRepo.On("UpdateStatus", 123, "inactive").Return(nil)
		mockRepo.On("GetByID", 123).Return(updatedAthlete, nil)

		resp, err := service.UpdateAthleteStatus(123, "inactive")

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "inactive", "inactive")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Спортсмен не найден при обновлении", func(t *testing.T) {
		mockRepo := new(mocks.AthleteRepository)
		service := NewAthleteService(mockRepo)

		mockRepo.On("GetByID", 999).Return(nil, nil)

		resp, err := service.UpdateAthleteStatus(999, "inactive")

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, "спортсмен не найден", err.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("Ошибка при обновлении статуса в БД", func(t *testing.T) {
		mockRepo := new(mocks.AthleteRepository)
		service := NewAthleteService(mockRepo)

		existingAthlete := &domain.Athlete{
			ID:     "123",
			Name:   "Иван Иванов",
			Status: "active",
		}

		mockRepo.On("GetByID", 123).Return(existingAthlete, nil)
		mockRepo.On("UpdateStatus", 123, "inactive").Return(errors.New("update error"))

		resp, err := service.UpdateAthleteStatus(123, "inactive")

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, "update error", err.Error())
		mockRepo.AssertExpectations(t)
	})
}
