package database

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type SessionRepo struct {
	db *gorm.DB
}

type Session struct {
	UserID    string `gorm:"primaryKey"`
	Session   string `gorm:"primaryKey"`
	CreatedAt time.Time
}

func NewSessionRepo(db *gorm.DB) *SessionRepo {
	return &SessionRepo{db: db}
}

func (r *SessionRepo) Create(ctx context.Context, session Session) error {
	return r.db.WithContext(ctx).Create(&session).Error
}

func (r *SessionRepo) Get(ctx context.Context, userID string, session string) (*Session, error) {
	var output Session
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Where("session = ?", session).First(&output).Error
	if err != nil {
		return nil, err
	}

	return &output, nil
}
