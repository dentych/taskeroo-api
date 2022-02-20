package auth

import (
	"context"
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

func (a *Auth) IsAuthenticated(ctx context.Context, userID string, session string) (bool, error) {
	_, err := a.sessionRepo.Get(ctx, userID, session)
	if err != nil {
		return false, err
	}

	return true, nil
}

type UserSession struct {
	UserID  string
	Session string
}

func (a *Auth) Login(ctx context.Context, email string, password string) (UserSession, error) {
	user, err := a.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return UserSession{}, internalerrors.InvalidEmailOrPassword
		}
		return UserSession{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	if err != nil {
		return UserSession{}, internalerrors.InvalidEmailOrPassword
	}

	session := uuid.NewString()
	err = a.sessionRepo.Create(ctx, database.Session{
		UserID:    user.UserID,
		Session:   session,
		CreatedAt: time.Now(),
	})
	if err != nil {
		return UserSession{}, err
	}

	return UserSession{UserID: user.UserID, Session: session}, nil
}

func (a *Auth) Register(ctx context.Context, email string, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 0)
	if err != nil {
		return err
	}

	err = a.userRepo.Create(ctx, database.User{
		UserID:         uuid.NewString(),
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
