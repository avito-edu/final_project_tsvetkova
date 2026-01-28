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

func TestCompetitionService_CreateCompetition(t *testing.T) {
	t.Run("Успешное создание соревнования", func(t *testing.T) {
		mockRepo := new(mocks.CompetitionRepository)
		service := NewCompetitionService(mockRepo)

		date := time.Now()
		req := dto.CreateCompetitionRequest{
			Name:        "Чемпионат",
			Date:        date,
			OrganizerID: 1,
		}

		mockRepo.On("Create", mock.AnythingOfType("*domain.Competition")).
			Run(func(args mock.Arguments) {
				comp := args.Get(0).(*domain.Competition)
				comp.ID = 1
			}).
			Return(nil)

		resp, err := service.CreateCompetition(req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, 1, resp.ID)
		assert.Equal(t, "Чемпионат", resp.Name)
		assert.Equal(t, date, resp.Date)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Пустое название", func(t *testing.T) {
		mockRepo := new(mocks.CompetitionRepository)
		service := NewCompetitionService(mockRepo)

		req := dto.CreateCompetitionRequest{
			Name:        "",
			Date:        time.Now(),
			OrganizerID: 1,
		}

		resp, err := service.CreateCompetition(req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, "название не может быть пустым", err.Error())
		mockRepo.AssertNotCalled(t, "Create")
	})

	t.Run("Неверный ID организатора", func(t *testing.T) {
		mockRepo := new(mocks.CompetitionRepository)
		service := NewCompetitionService(mockRepo)

		req := dto.CreateCompetitionRequest{
			Name:        "Чемпионат",
			Date:        time.Now(),
			OrganizerID: 0,
		}

		resp, err := service.CreateCompetition(req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, "неверный ID организатора", err.Error())
		mockRepo.AssertNotCalled(t, "Create")
	})

	t.Run("Ошибка репозитория", func(t *testing.T) {
		mockRepo := new(mocks.CompetitionRepository)
		service := NewCompetitionService(mockRepo)

		req := dto.CreateCompetitionRequest{
			Name:        "Чемпионат",
			Date:        time.Now(),
			OrganizerID: 1,
		}

		mockRepo.On("Create", mock.AnythingOfType("*domain.Competition")).
			Return(errors.New("database error"))

		resp, err := service.CreateCompetition(req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, "database error", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestCompetitionService_GetCompetition(t *testing.T) {
	t.Run("Успешное получение соревнования", func(t *testing.T) {
		mockRepo := new(mocks.CompetitionRepository)
		service := NewCompetitionService(mockRepo)

		date := time.Now()
		expectedComp := &domain.Competition{
			ID:        1,
			Name:      "Чемпионат",
			Date:      date,
			Organizer: 1,
			Results: []domain.Result{
				{AthleteID: 1, Distance: "100m", TimeSec: 12.5, Place: 1},
			},
		}

		mockRepo.On("GetByID", 1).Return(expectedComp, nil)

		resp, err := service.GetCompetition(1)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, 1, resp.ID)
		assert.Equal(t, 1, len(resp.Results))
		assert.Equal(t, "100m", resp.Results[0].Distance)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Неверный ID", func(t *testing.T) {
		mockRepo := new(mocks.CompetitionRepository)
		service := NewCompetitionService(mockRepo)

		resp, err := service.GetCompetition(0)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, "неверный ID", err.Error())
		mockRepo.AssertNotCalled(t, "GetByID")
	})

	t.Run("Соревнование не найдено", func(t *testing.T) {
		mockRepo := new(mocks.CompetitionRepository)
		service := NewCompetitionService(mockRepo)

		mockRepo.On("GetByID", 999).Return(nil, nil)

		resp, err := service.GetCompetition(999)

		assert.NoError(t, err)
		assert.Nil(t, resp)
		mockRepo.AssertExpectations(t)
	})
}

func TestCompetitionService_AddResult(t *testing.T) {
	t.Run("Успешное добавление результата", func(t *testing.T) {
		mockRepo := new(mocks.CompetitionRepository)
		service := NewCompetitionService(mockRepo)

		resultReq := dto.AddResultRequest{
			AthleteID: 1,
			Distance:  "100m",
			TimeSec:   12.5,
			Place:     1,
		}

		mockRepo.On("AddResult", 1, mock.AnythingOfType("domain.Result")).
			Return(nil)

		resp, err := service.AddResult(1, resultReq)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, 1, resp.AthleteID)
		assert.Equal(t, "100m", resp.Distance)
		assert.Equal(t, 12.5, resp.TimeSec)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Неверный ID соревнования", func(t *testing.T) {
		mockRepo := new(mocks.CompetitionRepository)
		service := NewCompetitionService(mockRepo)

		resultReq := dto.AddResultRequest{
			AthleteID: 1,
			Distance:  "100m",
			TimeSec:   12.5,
			Place:     1,
		}

		resp, err := service.AddResult(0, resultReq)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, "неверный ID соревнования", err.Error())
		mockRepo.AssertNotCalled(t, "AddResult")
	})

	t.Run("Неверный ID спортсмена", func(t *testing.T) {
		mockRepo := new(mocks.CompetitionRepository)
		service := NewCompetitionService(mockRepo)

		resultReq := dto.AddResultRequest{
			AthleteID: 0,
			Distance:  "100m",
			TimeSec:   12.5,
			Place:     1,
		}

		resp, err := service.AddResult(1, resultReq)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, "неверный ID спортсмена", err.Error())
		mockRepo.AssertNotCalled(t, "AddResult")
	})

	t.Run("Пустая дистанция", func(t *testing.T) {
		mockRepo := new(mocks.CompetitionRepository)
		service := NewCompetitionService(mockRepo)

		resultReq := dto.AddResultRequest{
			AthleteID: 1,
			Distance:  "",
			TimeSec:   12.5,
			Place:     1,
		}

		resp, err := service.AddResult(1, resultReq)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, "дистанция не может быть пустой", err.Error())
		mockRepo.AssertNotCalled(t, "AddResult")
	})

	t.Run("Неверное время", func(t *testing.T) {
		mockRepo := new(mocks.CompetitionRepository)
		service := NewCompetitionService(mockRepo)

		resultReq := dto.AddResultRequest{
			AthleteID: 1,
			Distance:  "100m",
			TimeSec:   0,
			Place:     1,
		}

		resp, err := service.AddResult(1, resultReq)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, "время должно быть больше 0", err.Error())
		mockRepo.AssertNotCalled(t, "AddResult")
	})
}

func TestCompetitionService_GetAthleteProgress(t *testing.T) {
	t.Run("Успешное получение прогресса", func(t *testing.T) {
		mockRepo := new(mocks.CompetitionRepository)
		service := NewCompetitionService(mockRepo)

		expectedResults := []domain.Result{
			{AthleteID: 1, Distance: "100m", TimeSec: 12.5, Place: 1},
			{AthleteID: 1, Distance: "100m", TimeSec: 12.3, Place: 2},
		}

		mockRepo.On("GetAthleteResultsByDistance", 1, "100m").
			Return(expectedResults, nil)

		resp, err := service.GetAthleteProgress(1, "100m")

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, 2, len(resp))
		assert.Equal(t, 12.5, resp[0].TimeSec)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Неверный ID спортсмена", func(t *testing.T) {
		mockRepo := new(mocks.CompetitionRepository)
		service := NewCompetitionService(mockRepo)

		resp, err := service.GetAthleteProgress(0, "100m")

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, "неверный ID спортсмена", err.Error())
		mockRepo.AssertNotCalled(t, "GetAthleteResultsByDistance")
	})

	t.Run("Пустая дистанция", func(t *testing.T) {
		mockRepo := new(mocks.CompetitionRepository)
		service := NewCompetitionService(mockRepo)

		resp, err := service.GetAthleteProgress(1, "")

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, "дистанция не может быть пустой", err.Error())
		mockRepo.AssertNotCalled(t, "GetAthleteResultsByDistance")
	})
}

func TestCompetitionService_GetAllResultsByDistance(t *testing.T) {
	t.Run("Успешное получение всех результатов по дистанции", func(t *testing.T) {
		mockRepo := new(mocks.CompetitionRepository)
		service := NewCompetitionService(mockRepo)

		expectedResults := []domain.Result{
			{AthleteID: 1, Distance: "100m", TimeSec: 12.5, Place: 1},
			{AthleteID: 2, Distance: "100m", TimeSec: 12.8, Place: 2},
		}

		mockRepo.On("GetAllResultsByDistance", "100m").
			Return(expectedResults, nil)

		resp, err := service.GetAllResultsByDistance("100m")

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, 2, len(resp))
		mockRepo.AssertExpectations(t)
	})

	t.Run("Пустая дистанция", func(t *testing.T) {
		mockRepo := new(mocks.CompetitionRepository)
		service := NewCompetitionService(mockRepo)

		resp, err := service.GetAllResultsByDistance("")

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, "дистанция не может быть пустой", err.Error())
		mockRepo.AssertNotCalled(t, "GetAllResultsByDistance")
	})
}
