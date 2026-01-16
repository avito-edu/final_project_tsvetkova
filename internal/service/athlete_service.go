package service

import (
	"errors"
	"swim_service/internal/domain"
	"swim_service/internal/dto"
)

type AthleteRepository interface {
	CreateAthlete(athlete *domain.Athlete) error
	GetByID(id int) (*domain.Athlete, error)
	UpdateStatus(id int, newStatus string) error
}

type AthleteService struct {
	athleteRepo AthleteRepository
}

func NewAthleteService(athleteRepo AthleteRepository) *AthleteService {
	return &AthleteService{
		athleteRepo: athleteRepo,
	}
}

func (s *AthleteService) CreateAthlete(athleteDTO dto.CreateAthleteRequest) (*dto.AthleteResponse, error) {
	athlete := &domain.Athlete{
		Name:      athleteDTO.Name,
		BirthDate: athleteDTO.BirthDate,
		Status:    athleteDTO.Status,
		Gender:    athleteDTO.Gender,
	}

	err := s.athleteRepo.CreateAthlete(athlete)
	if err != nil {
		return nil, err
	}

	return &dto.AthleteResponse{
		ID:        athlete.ID,
		Name:      athlete.Name,
		BirthDate: athlete.BirthDate,
		Status:    athlete.Status,
		Gender:    athlete.Gender,
	}, nil
}

func (s *AthleteService) GetAthlete(id int) (*dto.AthleteResponse, error) {
	athlete, err := s.athleteRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if athlete == nil {
		return nil, nil
	}

	return &dto.AthleteResponse{
		ID:        athlete.ID,
		Name:      athlete.Name,
		BirthDate: athlete.BirthDate,
		Status:    athlete.Status,
		Gender:    athlete.Gender,
	}, nil
}

func (s *AthleteService) UpdateAthleteStatus(id int, status string) (*dto.AthleteResponse, error) {
	athlete, err := s.athleteRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if athlete == nil {
		return nil, errors.New("спортсмен не найден")
	}

	err = s.athleteRepo.UpdateStatus(id, status)
	if err != nil {
		return nil, err
	}

	updatedAthlete, err := s.athleteRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return &dto.AthleteResponse{
		ID:        updatedAthlete.ID,
		Name:      updatedAthlete.Name,
		BirthDate: updatedAthlete.BirthDate,
		Status:    updatedAthlete.Status,
		Gender:    updatedAthlete.Gender,
	}, nil
}
