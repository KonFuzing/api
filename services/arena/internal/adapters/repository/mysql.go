package repository

import (
	"api/services/arena/internal/core/domain"
	"api/services/arena/internal/core/ports"
	"strings"
	"time"
	"gorm.io/gorm"
)

type battleModel struct {
	ID         uint `gorm:"primaryKey"`
	Fighter1ID string
	Fighter2ID string
	Winner     string
	Logs       string `gorm:"type:text"`
	CreatedAt  time.Time
}

type mysqlRepo struct {
	db *gorm.DB
}

func NewMySQLRepository(db *gorm.DB) ports.BattleRepository {
	db.AutoMigrate(&battleModel{})
	return &mysqlRepo{db: db}
}

func (r *mysqlRepo) Save(res *domain.BattleResult, f1, f2 string) error {
	m := battleModel{
		Fighter1ID: f1,
		Fighter2ID: f2,
		Winner:     res.Winner,
		Logs:       strings.Join(res.Logs, "\n"),
	}
	return r.db.Create(&m).Error
}

func (r *mysqlRepo) GetAll() ([]domain.BattleResult, error) {
	var models []battleModel
	if err := r.db.Order("created_at desc").Find(&models).Error; err != nil {
		return nil, err
	}
	
	var results []domain.BattleResult
	for _, m := range models {
		results = append(results, domain.BattleResult{
			Winner: m.Winner,
			Logs:   strings.Split(m.Logs, "\n"),
		})
	}
	return results, nil
}

func (r *mysqlRepo) GetHistory(limit int, fighterID string) ([]domain.BattleResult, error) {
	var models []battleModel
	
	// เริ่มต้น Query
	query := r.db.Order("created_at desc")

	// 1. ถ้ามี limit ให้ใส่ limit (ถ้าเป็น 0 ให้ default สัก 50 กันบึ้ม)
	if limit > 0 {
		query = query.Limit(limit)
	} else {
		query = query.Limit(50) 
	}

	// 2. ถ้าระบุ fighterID ให้หาทั้งช่อง fighter1 หรือ fighter2
	if fighterID != "" {
		query = query.Where("fighter1_id = ? OR fighter2_id = ?", fighterID, fighterID)
	}

	// รัน Query
	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}
	
	// แปลงเป็น Domain Object (เหมือนเดิม)
	var results []domain.BattleResult
	for _, m := range models {
		results = append(results, domain.BattleResult{
			Winner: m.Winner,
			Logs:   strings.Split(m.Logs, "\n"),
		})
	}
	return results, nil
}