package database

import (
	"gorm.io/gorm"
	"time"
)

type UserRepo struct {
	db *gorm.DB
}

type User struct {
	Email          string `gorm:"primaryKey"`
	HashedPassword string
	CreatedAt      time.Time
	LastLogin      time.Time
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(user User) error {
	return r.db.Create(&user).Error
}

func (r *UserRepo) Get(userID string) (*User, error) {
	var user User
	err := r.db.First(&user, "id = ?", userID).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}
