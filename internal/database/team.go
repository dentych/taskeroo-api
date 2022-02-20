package database

import "gorm.io/gorm"

type TeamRepo struct {
	db *gorm.DB
}

type Team struct {
	ID          string `gorm:"primaryKey;"`
	TeamName    string
	OwnerUserID string
}

func NewTeamRepo(db *gorm.DB) *TeamRepo {
	return &TeamRepo{db: db}
}

func (r *TeamRepo) Create(team Team) error {
	return r.db.Create(&team).Error
}
