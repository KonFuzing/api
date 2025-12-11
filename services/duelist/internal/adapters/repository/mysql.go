package repository

import (
	"api/services/duelist/internal/core/domain"
	"api/services/duelist/internal/core/ports"

	"gorm.io/gorm"
)

// DB Entity (Infrastructure Layer)
type cowboyModel struct {
	ID       string `gorm:"primaryKey"`
	Name     string
	Health   int
	Damage   int
	Speed    int
	Accuracy float64
}

func (cowboyModel) TableName() string {
	return "cowboys"
}

// แปลงจาก Model -> Domain
func (m *cowboyModel) toDomain() *domain.Cowboy {
	return &domain.Cowboy{
		ID:       m.ID,
		Name:     m.Name,
		Health:   m.Health,
		Damage:   m.Damage,
		Speed:    m.Speed,
		Accuracy: m.Accuracy,
	}
}

// แปลงจาก Domain -> Model
func fromDomain(d *domain.Cowboy) *cowboyModel {
	return &cowboyModel{
		ID:       d.ID,
		Name:     d.Name,
		Health:   d.Health,
		Damage:   d.Damage,
		Speed:    d.Speed,
		Accuracy: d.Accuracy,
	}
}

type mysqlRepo struct {
	db *gorm.DB
}

func NewMySQLRepository(db *gorm.DB) ports.CowboyRepository {
	db.AutoMigrate(&cowboyModel{})
	return &mysqlRepo{db: db}
}

func (r *mysqlRepo) Save(cowboy *domain.Cowboy) error {
	model := fromDomain(cowboy)
	return r.db.Create(model).Error
}

func (r *mysqlRepo) FindByID(id string) (*domain.Cowboy, error) {
	var model cowboyModel
	if err := r.db.First(&model, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return model.toDomain(), nil
}
