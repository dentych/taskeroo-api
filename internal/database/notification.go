package database

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type NotificationRepo struct {
	db *gorm.DB
}

type GroupDiscord struct {
	GroupID   string `gorm:"primaryKey"`
	Webhook   string `gorm:"not null;"`
	CreatedAt time.Time
}

type DiscordUsername struct {
	UserID          string `gorm:"primaryKey"`
	DiscordUsername string `gorm:"not null;"`
}

func NewNotificationRepo(db *gorm.DB) *NotificationRepo {
	return &NotificationRepo{db: db}
}

func (r *NotificationRepo) CreateDiscord(ctx context.Context, discord GroupDiscord) error {
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "group_id"}},
		UpdateAll: true,
	}).Create(&discord).Error
}

func (r *NotificationRepo) GetGroupDiscord(ctx context.Context, groupID string) (*GroupDiscord, error) {
	var discord GroupDiscord
	err := r.db.WithContext(ctx).First(&discord, "group_id = ?", groupID).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return &discord, nil
}

func (r *NotificationRepo) CreateDiscordUsername(ctx context.Context, userID string, discordUsername string) error {
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}},
		UpdateAll: true,
	}).Create(&DiscordUsername{
		UserID:          userID,
		DiscordUsername: discordUsername,
	}).Error
}

func (r *NotificationRepo) GetDiscordUsername(ctx context.Context, userID string) (*DiscordUsername, error) {
	var username *DiscordUsername
	err := r.db.WithContext(ctx).First(&username, "user_id = ?", userID).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return username, nil
}
