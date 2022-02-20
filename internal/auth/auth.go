package auth

import (
	"errors"
	"github.com/dentych/taskeroo/internal/database"
	internalerrors "github.com/dentych/taskeroo/internal/errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
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

func (a *Auth) IsAuthenticated(email string, session string) (bool, error) {
	_, err := a.sessionRepo.Get(email, session)
	if err != nil {
		return false, err
	}

	return true, nil
}

type UserSession struct {
	Email   string
	Session string
}

func (a *Auth) Login(email string, password string) (string, error) {
	user, err := a.userRepo.Get(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", internalerrors.InvalidEmailOrPassword
		}
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	if err != nil {
		return "", internalerrors.InvalidEmailOrPassword
	}

	session := uuid.NewString()
	err = a.sessionRepo.Create(database.Session{
		UserID:    email,
		Session:   session,
		CreatedAt: time.Now(),
	})
	if err != nil {
		return "", err
	}

	return session, nil
}

func (a *Auth) Register(email string, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 0)
	if err != nil {
		return err
	}

	err = a.userRepo.Create(database.User{
		Email:          email,
		HashedPassword: string(hashedPassword),
		CreatedAt:      time.Now(),
		LastLogin:      time.Now(),
	})
	if err != nil {
		return err
	}

	return nil
}
