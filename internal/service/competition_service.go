package service

import (
	"context"
	"errors"
	"swim_service/internal/domain"
	"swim_service/internal/dto"
	"time"
)

type CompetitionRepository interface {
	Create(comp *domain.Competition) error
	GetByID(id int) (*domain.Competition, error)
	GetAll() ([]*domain.Competition, error)
	AddResult(compID int, result domain.Result) error
	GetAthleteResultsByDistance(athleteID int, distance string) ([]domain.Result, error)
	GetAllResultsByDistance(distance string) ([]domain.Result, error)
	GetAthleteResultsByPeriod(ctx context.Context, athleteID int, startDate, endDate time.Time) ([]domain.Result, error)
}

type CompetitionService struct {
	compRepo CompetitionRepository
}

func NewCompetitionService(compRepo CompetitionRepository) *CompetitionService {
	return &CompetitionService{
		compRepo: compRepo,
	}
}

func (s *CompetitionService) CreateCompetition(compDTO dto.CreateCompetitionRequest) (*dto.CompetitionResponse, error) {
	if compDTO.Name == "" {
		return nil, errors.New("название не может быть пустым")
	}
	if compDTO.OrganizerID <= 0 {
		return nil, errors.New("неверный ID организатора")
	}

	comp := &domain.Competition{
		Name:      compDTO.Name,
		Date:      compDTO.Date,
		Organizer: compDTO.OrganizerID,
		Results:   []domain.Result{},
	}

	err := s.compRepo.Create(comp)
	if err != nil {
		return nil, err
	}

	return &dto.CompetitionResponse{
		ID:        comp.ID,
		Name:      comp.Name,
		Date:      comp.Date,
		Organizer: comp.Organizer,
		Results:   []dto.ResultResponse{},
	}, nil
}

func (s *CompetitionService) GetCompetition(id int) (*dto.CompetitionResponse, error) {
	if id <= 0 {
		return nil, errors.New("неверный ID")
	}

	comp, err := s.compRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if comp == nil {
		return nil, nil
	}

	resultsDTO := make([]dto.ResultResponse, len(comp.Results))
	for i, result := range comp.Results {
		resultsDTO[i] = dto.ResultResponse{
			AthleteID: result.AthleteID,
			Distance:  result.Distance,
			TimeSec:   result.TimeSec,
			Place:     result.Place,
		}
	}

	return &dto.CompetitionResponse{
		ID:        comp.ID,
		Name:      comp.Name,
		Date:      comp.Date,
		Organizer: comp.Organizer,
		Results:   resultsDTO,
	}, nil
}

func (s *CompetitionService) GetAllCompetitions() ([]dto.CompetitionResponse, error) {
	comps, err := s.compRepo.GetAll()
	if err != nil {
		return nil, err
	}

	compsDTO := make([]dto.CompetitionResponse, len(comps))
	for i, comp := range comps {
		resultsDTO := make([]dto.ResultResponse, len(comp.Results))
		for j, result := range comp.Results {
			resultsDTO[j] = dto.ResultResponse{
				AthleteID: result.AthleteID,
				Distance:  result.Distance,
				TimeSec:   result.TimeSec,
				Place:     result.Place,
			}
		}

		compsDTO[i] = dto.CompetitionResponse{
			ID:        comp.ID,
			Name:      comp.Name,
			Date:      comp.Date,
			Organizer: comp.Organizer,
			Results:   resultsDTO,
		}
	}

	return compsDTO, nil
}

func (s *CompetitionService) AddResult(compID int, resultDTO dto.AddResultRequest) (*dto.ResultResponse, error) {
	if compID <= 0 {
		return nil, errors.New("неверный ID соревнования")
	}
	if resultDTO.AthleteID <= 0 {
		return nil, errors.New("неверный ID спортсмена")
	}
	if resultDTO.Distance == "" {
		return nil, errors.New("дистанция не может быть пустой")
	}
	if resultDTO.TimeSec <= 0 {
		return nil, errors.New("время должно быть больше 0")
	}

	result := domain.Result{
		AthleteID: resultDTO.AthleteID,
		Distance:  resultDTO.Distance,
		TimeSec:   resultDTO.TimeSec,
		Place:     resultDTO.Place,
	}

	err := s.compRepo.AddResult(compID, result)
	if err != nil {
		return nil, err
	}

	return &dto.ResultResponse{
		AthleteID: result.AthleteID,
		Distance:  result.Distance,
		TimeSec:   result.TimeSec,
		Place:     result.Place,
	}, nil
}

func (s *CompetitionService) GetAthleteProgress(athleteID int, distance string) ([]dto.ResultResponse, error) {
	if athleteID <= 0 {
		return nil, errors.New("неверный ID спортсмена")
	}
	if distance == "" {
		return nil, errors.New("дистанция не может быть пустой")
	}

	results, err := s.compRepo.GetAthleteResultsByDistance(athleteID, distance)
	if err != nil {
		return nil, err
	}

	resultsDTO := make([]dto.ResultResponse, len(results))
	for i, result := range results {
		resultsDTO[i] = dto.ResultResponse{
			AthleteID: result.AthleteID,
			Distance:  result.Distance,
			TimeSec:   result.TimeSec,
			Place:     result.Place,
		}
	}

	return resultsDTO, nil
}

func (s *CompetitionService) GetAllResultsByDistance(distance string) ([]dto.ResultResponse, error) {
	if distance == "" {
		return nil, errors.New("дистанция не может быть пустой")
	}

	results, err := s.compRepo.GetAllResultsByDistance(distance)
	if err != nil {
		return nil, err
	}

	resultsDTO := make([]dto.ResultResponse, len(results))
	for i, result := range results {
		resultsDTO[i] = dto.ResultResponse{
			AthleteID: result.AthleteID,
			Distance:  result.Distance,
			TimeSec:   result.TimeSec,
			Place:     result.Place,
		}
	}

	return resultsDTO, nil
}
