package service

import (
	"context"
	"swim_service/internal/domain"
	"swim_service/internal/dto"
	"time"
)

type AnalyticsService struct {
	competitionRepo CompetitionRepository
}

func NewAnalyticsService(competitionRepo CompetitionRepository) *AnalyticsService {
	return &AnalyticsService{
		competitionRepo: competitionRepo,
	}
}

func (s *AnalyticsService) GetAthleteAnalytics(ctx context.Context, req dto.AthleteAnalyticsRequest) (dto.AthleteAnalyticsResponse, error) {
	results, err := s.competitionRepo.GetAthleteResultsByPeriod(
		ctx,
		req.AthleteID,
		req.PeriodStart,
		req.PeriodEnd,
	)
	if err != nil {
		return dto.AthleteAnalyticsResponse{}, err
	}

	stats := calculateAnalytics(results, req.PeriodStart, req.PeriodEnd)

	response := dto.AthleteAnalyticsResponse{
		AthleteID:         req.AthleteID,
		PeriodStart:       req.PeriodStart,
		PeriodEnd:         req.PeriodEnd,
		TotalDistance:     stats.TotalDistance,
		TotalCompetitions: len(stats.CompetitionIDs),
		AverageTime:       stats.AverageTime,
		BestTime:          stats.BestTime,
		WorstTime:         stats.WorstTime,
		Results:           convertResultsToDTO(results),
		DistanceBreakdown: calculateDistanceBreakdown(results),
	}

	return response, nil
}

type athleteStats struct {
	TotalDistance      string
	CompetitionIDs     map[int]bool
	TotalTime          float64
	AverageTime        float64
	Count              int
	BestTime           float64
	WorstTime          float64
	DistanceTimeSums   map[string]float64
	DistanceCounts     map[string]int
	DistanceBestTimes  map[string]float64
	DistanceWorstTimes map[string]float64
}

func calculateAnalytics(results []domain.Result, start, end time.Time) athleteStats {
	if len(results) == 0 {
		return athleteStats{
			BestTime:           0,
			WorstTime:          0,
			DistanceTimeSums:   make(map[string]float64),
			DistanceCounts:     make(map[string]int),
			DistanceBestTimes:  make(map[string]float64),
			DistanceWorstTimes: make(map[string]float64),
		}
	}

	stats := athleteStats{
		CompetitionIDs:     make(map[int]bool),
		BestTime:           results[0].TimeSec,
		WorstTime:          results[0].TimeSec,
		DistanceTimeSums:   make(map[string]float64),
		DistanceCounts:     make(map[string]int),
		DistanceBestTimes:  make(map[string]float64),
		DistanceWorstTimes: make(map[string]float64),
	}

	for _, result := range results {
		stats.TotalTime += result.TimeSec
		stats.Count++

		if result.TimeSec < stats.BestTime {
			stats.BestTime = result.TimeSec
		}
		if result.TimeSec > stats.WorstTime {
			stats.WorstTime = result.TimeSec
		}

		if _, exists := stats.DistanceTimeSums[result.Distance]; !exists {
			stats.DistanceTimeSums[result.Distance] = 0
			stats.DistanceCounts[result.Distance] = 0
			stats.DistanceBestTimes[result.Distance] = result.TimeSec
			stats.DistanceWorstTimes[result.Distance] = result.TimeSec
		}

		stats.DistanceTimeSums[result.Distance] += result.TimeSec
		stats.DistanceCounts[result.Distance]++

		if result.TimeSec < stats.DistanceBestTimes[result.Distance] {
			stats.DistanceBestTimes[result.Distance] = result.TimeSec
		}
		if result.TimeSec > stats.DistanceWorstTimes[result.Distance] {
			stats.DistanceWorstTimes[result.Distance] = result.TimeSec
		}
	}

	if stats.Count > 0 {
		stats.AverageTime = stats.TotalTime / float64(stats.Count)
	}

	return stats
}

func calculateDistanceBreakdown(results []domain.Result) []dto.DistanceStats {
	distanceMap := make(map[string]*dto.DistanceStats)

	for _, result := range results {
		if _, exists := distanceMap[result.Distance]; !exists {
			distanceMap[result.Distance] = &dto.DistanceStats{
				Distance:  result.Distance,
				Count:     0,
				BestTime:  result.TimeSec,
				WorstTime: result.TimeSec,
			}
		}

		stats := distanceMap[result.Distance]
		stats.Count++
		stats.AverageTime += result.TimeSec

		if result.TimeSec < stats.BestTime {
			stats.BestTime = result.TimeSec
		}
		if result.TimeSec > stats.WorstTime {
			stats.WorstTime = result.TimeSec
		}
	}

	var breakdown []dto.DistanceStats
	for _, stats := range distanceMap {
		if stats.Count > 0 {
			stats.AverageTime = stats.AverageTime / float64(stats.Count)
		}
		breakdown = append(breakdown, *stats)
	}

	return breakdown
}

func convertResultsToDTO(results []domain.Result) []dto.ResultResponse {
	var dtoResults []dto.ResultResponse
	for _, result := range results {
		dtoResults = append(dtoResults, dto.ResultResponse{
			AthleteID: result.AthleteID,
			Distance:  result.Distance,
			TimeSec:   result.TimeSec,
			Place:     result.Place,
		})
	}
	return dtoResults
}
