package service

import (
	"errors"
	"swim_service/internal/domain"
	"swim_service/internal/dto"
)

type OrganizerRepository interface {
	Create(organizer *domain.Organizer) error
	GetByID(id int) (*domain.Organizer, error)
}

type OrganizerService struct {
	organizerRepo OrganizerRepository
}

func NewOrganizerService(organizerRepo OrganizerRepository) *OrganizerService {
	return &OrganizerService{
		organizerRepo: organizerRepo,
	}
}

func (s *OrganizerService) CreateOrganizer(organizerDTO dto.CreateOrganizerRequest) (*dto.OrganizerResponse, error) {
	if organizerDTO.Email == "" {
		return nil, errors.New("email не может быть пустым")
	}
	organizer := &domain.Organizer{
		Name:  organizerDTO.Name,
		Email: organizerDTO.Email,
	}

	err := s.organizerRepo.Create(organizer)
	if err != nil {
		return nil, err
	}

	return &dto.OrganizerResponse{
		ID:    organizer.ID,
		Name:  organizer.Name,
		Email: organizer.Email,
	}, nil
}

func (s *OrganizerService) GetOrganizer(id int) (*dto.OrganizerResponse, error) {
	if id <= 0 {
		return nil, errors.New("неверный ID")
	}

	organizer, err := s.organizerRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if organizer == nil {
		return nil, nil
	}

	return &dto.OrganizerResponse{
		ID:    organizer.ID,
		Name:  organizer.Name,
		Email: organizer.Email,
	}, nil
}
