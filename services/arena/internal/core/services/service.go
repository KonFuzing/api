package services

import (
	"api/services/arena/internal/core/domain"
	"api/services/arena/internal/core/ports"
	"errors"
)

type service struct {
	provider ports.CowboyProvider
	repo     ports.BattleRepository
}

func NewArenaService(p ports.CowboyProvider, r ports.BattleRepository) ports.ArenaService {
	return &service{provider: p, repo: r}
}

func (s *service) GetHistory(limit int, fighterID string) ([]domain.BattleResult, error) {
	return s.repo.GetHistory(limit, fighterID)
}

func (s *service) Duel(id1, id2 string) (*domain.BattleResult, error) {
	// 1. เรียกข้อมูลจาก Port (Adapter จะไปเรียก gRPC)
	c1, err := s.provider.GetCowboy(id1)
	if err != nil {
		return nil, err
	}

	c2, err := s.provider.GetCowboy(id2)
	if err != nil {
		return nil, err
	}

	// 2. รัน Domain Logic
	result := domain.SimulateFight(c1, c2)

	// 3. บันทึกผ่าน Port (Adapter จะไปลง DB)
	if err := s.repo.Save(&result, id1, id2); err != nil {
		return nil, errors.New("failed to save battle record")
	}

	return &result, nil
}
