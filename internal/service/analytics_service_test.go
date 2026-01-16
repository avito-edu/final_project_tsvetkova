package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"swim_service/internal/domain"
	"swim_service/internal/dto"
	"swim_service/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	startTime = time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	endTime   = time.Date(2025, 12, 31, 23, 59, 59, 0, time.UTC)
)

func TestAnalyticsService_GetAthleteAnalytics_EmptyResults(t *testing.T) {
	mockRepo := new(mocks.CompetitionRepository)
	service := NewAnalyticsService(mockRepo)

	athleteID := 123
	req := dto.AthleteAnalyticsRequest{
		AthleteID:   athleteID,
		PeriodStart: startTime,
		PeriodEnd:   endTime,
	}

	mockRepo.On("GetAthleteResultsByPeriod", mock.Anything, athleteID, startTime, endTime).
		Return([]domain.Result{}, nil)

	resp, err := service.GetAthleteAnalytics(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, athleteID, resp.AthleteID)
	assert.Equal(t, startTime, resp.PeriodStart)
	assert.Equal(t, endTime, resp.PeriodEnd)
	assert.Equal(t, "", resp.TotalDistance)
	assert.Equal(t, 0, resp.TotalCompetitions)
	assert.Equal(t, 0.0, resp.AverageTime)
	assert.Equal(t, 0.0, resp.BestTime)
	assert.Equal(t, 0.0, resp.WorstTime)
	assert.Empty(t, resp.Results)
	assert.Empty(t, resp.DistanceBreakdown)

	mockRepo.AssertExpectations(t)
}

func TestAnalyticsService_GetAthleteAnalytics_SingleResult(t *testing.T) {
	mockRepo := new(mocks.CompetitionRepository)
	service := NewAnalyticsService(mockRepo)

	athleteID := 456
	result := domain.Result{
		AthleteID: athleteID,
		Distance:  "100m",
		TimeSec:   55.5,
		Place:     1,
	}

	req := dto.AthleteAnalyticsRequest{
		AthleteID:   athleteID,
		PeriodStart: startTime,
		PeriodEnd:   endTime,
	}

	mockRepo.On("GetAthleteResultsByPeriod", mock.Anything, athleteID, startTime, endTime).
		Return([]domain.Result{result}, nil)

	resp, err := service.GetAthleteAnalytics(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, athleteID, resp.AthleteID)
	assert.Equal(t, 0, resp.TotalCompetitions)
	assert.Equal(t, 55.5, resp.AverageTime)
	assert.Equal(t, 55.5, resp.BestTime)
	assert.Equal(t, 55.5, resp.WorstTime)
	assert.Len(t, resp.Results, 1)
	assert.Equal(t, result.TimeSec, resp.Results[0].TimeSec)

	assert.Len(t, resp.DistanceBreakdown, 1)
	db := resp.DistanceBreakdown[0]
	assert.Equal(t, "100m", db.Distance)
	assert.Equal(t, 1, db.Count)
	assert.Equal(t, 55.5, db.AverageTime)
	assert.Equal(t, 55.5, db.BestTime)
	assert.Equal(t, 55.5, db.WorstTime)

	mockRepo.AssertExpectations(t)
}

func TestAnalyticsService_GetAthleteAnalytics_MultipleResults(t *testing.T) {
	mockRepo := new(mocks.CompetitionRepository)
	service := NewAnalyticsService(mockRepo)

	athleteID := 789
	results := []domain.Result{
		{AthleteID: athleteID, Distance: "100m", TimeSec: 52.0, Place: 1},
		{AthleteID: athleteID, Distance: "100m", TimeSec: 54.0, Place: 2},
		{AthleteID: athleteID, Distance: "200m", TimeSec: 110.0, Place: 1},
	}

	req := dto.AthleteAnalyticsRequest{
		AthleteID:   athleteID,
		PeriodStart: startTime,
		PeriodEnd:   endTime,
	}

	mockRepo.On("GetAthleteResultsByPeriod", mock.Anything, athleteID, startTime, endTime).
		Return(results, nil)

	resp, err := service.GetAthleteAnalytics(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, 0, resp.TotalCompetitions)
	assert.InDelta(t, (52.0+54.0+110.0)/3, resp.AverageTime, 0.01)
	assert.Equal(t, 52.0, resp.BestTime)
	assert.Equal(t, 110.0, resp.WorstTime)

	assert.Len(t, resp.DistanceBreakdown, 2)

	var stats100m, stats200m *dto.DistanceStats
	for i := range resp.DistanceBreakdown {
		if resp.DistanceBreakdown[i].Distance == "100m" {
			stats100m = &resp.DistanceBreakdown[i]
		} else if resp.DistanceBreakdown[i].Distance == "200m" {
			stats200m = &resp.DistanceBreakdown[i]
		}
	}

	assert.NotNil(t, stats100m)
	assert.Equal(t, 2, stats100m.Count)
	assert.InDelta(t, 53.0, stats100m.AverageTime, 0.01)
	assert.Equal(t, 52.0, stats100m.BestTime)
	assert.Equal(t, 54.0, stats100m.WorstTime)

	assert.NotNil(t, stats200m)
	assert.Equal(t, 1, stats200m.Count)
	assert.Equal(t, 110.0, stats200m.AverageTime)
	assert.Equal(t, 110.0, stats200m.BestTime)
	assert.Equal(t, 110.0, stats200m.WorstTime)

	mockRepo.AssertExpectations(t)
}

func TestAnalyticsService_GetAthleteAnalytics_RepoError(t *testing.T) {
	mockRepo := new(mocks.CompetitionRepository)
	service := NewAnalyticsService(mockRepo)

	expectedErr := errors.New("database timeout")
	athleteID := 101

	req := dto.AthleteAnalyticsRequest{
		AthleteID:   athleteID,
		PeriodStart: startTime,
		PeriodEnd:   endTime,
	}

	mockRepo.On("GetAthleteResultsByPeriod", mock.Anything, athleteID, startTime, endTime).
		Return([]domain.Result{}, expectedErr)

	resp, err := service.GetAthleteAnalytics(context.Background(), req)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Empty(t, resp)

	mockRepo.AssertExpectations(t)
}
