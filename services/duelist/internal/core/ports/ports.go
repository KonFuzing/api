package ports

import "api/services/duelist/internal/core/domain"

// Primary Port (Inbound): สิ่งที่ Service นี้ทำได้
type DuelistService interface {
	Create(cowboy *domain.Cowboy) (*domain.Cowboy, error)
	Get(id string) (*domain.Cowboy, error)
}

// Secondary Port (Outbound): สิ่งที่ Service นี้ต้องการจากภายนอก (DB)
type CowboyRepository interface {
	Save(cowboy *domain.Cowboy) error
	FindByID(id string) (*domain.Cowboy, error)
}