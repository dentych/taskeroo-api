package database

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type UserRepo struct {
	db *gorm.DB
}

type User struct {
	UserID         string `gorm:"primaryKey;"`
	Email          string `gorm:"uniqueIndex;"`
	HashedPassword string `gorm:"not null;"`
	TeamID         *string
	CreatedAt      time.Time
	LastLogin      time.Time
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(ctx context.Context, user User) error {
	return r.db.WithContext(ctx).Create(&user).Error
}

func (r *UserRepo) Get(ctx context.Context, userID string) (*User, error) {
	var user User
	err := r.db.WithContext(ctx).First(&user, "user_id = ?", userID).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := r.db.WithContext(ctx).First(&user, "email = ?", email).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepo) SetTeam(ctx context.Context, userID string, teamID string) error {
	return r.db.WithContext(ctx).Model(&User{}).Where("user_id = ?", userID).Update("team_id", teamID).Error
}
