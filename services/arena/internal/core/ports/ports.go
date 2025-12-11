package ports

import (
	"api/services/arena/internal/core/domain"
	"api/services/arena/internal/core/domain/entity"
)

// Secondary Port (Outbound) - สำหรับดึงข้อมูล Cowboy (เช่นจาก gRPC)
type CowboyProvider interface {
	GetCowboy(id string) (*entity.Cowboy, error)
}

// Secondary Port (Outbound) - สำหรับเก็บผล (Database)
type ArenaService interface {
	Duel(fighter1ID, fighter2ID string) (*domain.BattleResult, error)
	GetHistory(limit int, fighterID string) ([]domain.BattleResult, error)
}

type BattleRepository interface {
	Save(result *domain.BattleResult, f1, f2 string) error
	GetHistory(limit int, fighterID string) ([]domain.BattleResult, error)
}
