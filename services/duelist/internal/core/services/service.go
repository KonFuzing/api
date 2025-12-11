package services

import (
	"api/services/duelist/internal/core/domain"
	"api/services/duelist/internal/core/ports"
	"errors"
)

type service struct {
	repo ports.CowboyRepository
}

func NewDuelistService(repo ports.CowboyRepository) ports.DuelistService {
	return &service{repo: repo}
}

func (s *service) Create(cowboy *domain.Cowboy) (*domain.Cowboy, error) {
	if cowboy.ID == "" {
		return nil, errors.New("ID is required")
	}
	if err := s.repo.Save(cowboy); err != nil {
		return nil, err
	}
	return cowboy, nil
}

func (s *service) Get(id string) (*domain.Cowboy, error) {
	return s.repo.FindByID(id)
}