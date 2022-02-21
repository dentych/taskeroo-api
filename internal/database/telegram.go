package database

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type TelegramRepo struct {
	db *gorm.DB
}

type Telegram struct {
	ID             string    `gorm:"primaryKey;"`
	TelegramUserID int       `gorm:"not null;uniqueIndex;"`
	UserID         *string   `gorm:"index;uniqueIndex"`
	CreatedAt      time.Time `gorm:"not null;"`
	UpdatedAt      time.Time `gorm:"not null;"`
}

type NewTelegram struct {
	ID             string
	TelegramUserID int
}

func NewTelegramRepo(db *gorm.DB) *TelegramRepo {
	return &TelegramRepo{db: db}
}

func (r *TelegramRepo) Create(ctx context.Context, telegram NewTelegram) error {
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "telegram_user_id"}},
		UpdateAll: true,
	}).Create(&Telegram{
		ID:             telegram.ID,
		TelegramUserID: telegram.TelegramUserID,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}).Error
}

func (r *TelegramRepo) SetUserID(ctx context.Context, ID string, userID string) error {
	return r.db.WithContext(ctx).Model(&Telegram{ID: ID}).Update("user_id", userID).Error
}
