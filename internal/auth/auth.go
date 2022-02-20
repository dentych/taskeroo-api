package auth

import (
	"github.com/dentych/taskeroo/internal/database"
	"github.com/dentych/taskeroo/internal/errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Auth struct {
	sessionRepo *database.SessionRepo
	userRepo    *database.UserRepo
}

func New(sessionRepo *database.SessionRepo, userRepo *database.UserRepo) *Auth {
	return &Auth{
		sessionRepo: sessionRepo,
		userRepo:    userRepo,
	}
}

func (a *Auth) IsAuthenticated(userID string, session string) (bool, error) {
	_, err := a.sessionRepo.Get(userID, session)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (a *Auth) Login(userID string, password string) (string, error) {
	user, err := a.userRepo.Get(userID)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	if err != nil {
		return "", errors.ErrEmailOrPasswordIncorrect
	}

	session := uuid.NewString()
	err = a.sessionRepo.Create(database.Session{
		UserID:    userID,
		Session:   session,
		CreatedAt: time.Now(),
	})
	if err != nil {
		return "", err
	}

	return session, nil
}
