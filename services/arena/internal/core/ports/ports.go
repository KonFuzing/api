package ports

import "api/services/arena/internal/core/domain"

// Primary Port (Inbound)
type ArenaService interface {
	Duel(fighter1ID, fighter2ID string) (*domain.BattleResult, error)
	GetHistory() ([]domain.BattleResult, error)
}

// Secondary Port (Outbound) - สำหรับดึงข้อมูล Cowboy (เช่นจาก gRPC)
type CowboyProvider interface {
	GetCowboy(id string) (*domain.Cowboy, error)
}

// Secondary Port (Outbound) - สำหรับเก็บผล (Database)
type BattleRepository interface {
	Save(result *domain.BattleResult, f1, f2 string) error
	GetAll() ([]domain.BattleResult, error)
}
