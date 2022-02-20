package database

import (
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

func (r *SessionRepo) Create(session Session) error {
	return r.db.Create(&session).Error
}

func (r *SessionRepo) Get(userID string, session string) (*Session, error) {
	var output Session
	err := r.db.Where("user_id = ?").Where("session = ?").First(&output).Error
	if err != nil {
		return nil, err
	}

	return &output, nil
}
