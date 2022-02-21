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
	ID             string `gorm:"primaryKey;"`
	Email          string `gorm:"uniqueIndex;"`
	Name           string
	HashedPassword string `gorm:"not null;"`
	GroupID        *string
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
	err := r.db.WithContext(ctx).First(&user, "id = ?", userID).Error
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

func (r *UserRepo) SetGroup(ctx context.Context, userID string, groupID string) error {
	return r.db.WithContext(ctx).Model(&User{}).Where("id = ?", userID).Update("group_id", groupID).Error
}
