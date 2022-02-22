package database

import (
	"context"
	"gorm.io/gorm"
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
	return r.db.WithContext(ctx).Create(&Telegram{
		ID:             telegram.ID,
		TelegramUserID: telegram.TelegramUserID,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}).Error
}

func (r *TelegramRepo) SetUserID(ctx context.Context, ID string, userID string) error {
	return r.db.WithContext(ctx).Model(&Telegram{ID: ID}).Update("user_id", userID).Error
}

func (r *TelegramRepo) Get(ctx context.Context, connectID string) (*Telegram, error) {
	var output Telegram
	err := r.db.WithContext(ctx).First(&output, "id = ?", connectID).Error
	if err != nil {
		return nil, err
	}

	return &output, nil
}

func (r *TelegramRepo) GetByTelegramUserID(ctx context.Context, telegramUserID int) (*Telegram, error) {
	var output Telegram
	err := r.db.WithContext(ctx).First(&output, "telegram_user_id = ?", telegramUserID).Error
	if err != nil {
		return nil, err
	}

	return &output, nil
}

func (r *TelegramRepo) DeleteAllByUserID(ctx context.Context, userID string) error {
	return r.db.WithContext(ctx).Delete(&Telegram{}, "user_id = ?", userID).Error
}

func (r *TelegramRepo) DeleteAllByTelegramUserID(ctx context.Context, telegramUserID int) error {
	return r.db.WithContext(ctx).Delete(&Telegram{}, "telegram_user_id = ?", telegramUserID).Error
}

func (r *TelegramRepo) GetByUserID(ctx context.Context, userID string) (*Telegram, error) {
	var output Telegram
	err := r.db.WithContext(ctx).First(&output, "user_id = ?", userID).Error
	if err != nil {
		return nil, err
	}

	return &output, nil
}
