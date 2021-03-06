package app

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

type AuthLogic struct {
	sessionRepo *database.SessionRepo
	userRepo    *database.UserRepo
	groupRepo   *database.GroupRepo
}

func NewAuthLogic(sessionRepo *database.SessionRepo, userRepo *database.UserRepo, groupRepo *database.GroupRepo) *AuthLogic {
	return &AuthLogic{
		sessionRepo: sessionRepo,
		userRepo:    userRepo,
		groupRepo:   groupRepo,
	}
}

func (a *AuthLogic) IsAuthenticated(ctx context.Context, userID string, session string) (bool, error) {
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

func (a *AuthLogic) Login(ctx context.Context, email string, password string) (UserSession, error) {
	user, err := a.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return UserSession{}, internalerrors.ErrInvalidEmailOrPassword
		}
		return UserSession{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	if err != nil {
		return UserSession{}, internalerrors.ErrInvalidEmailOrPassword
	}

	session := uuid.NewString()
	err = a.sessionRepo.Create(ctx, database.Session{
		UserID:    user.ID,
		Session:   session,
		CreatedAt: time.Now(),
	})
	if err != nil {
		return UserSession{}, err
	}

	return UserSession{UserID: user.ID, Session: session}, nil
}

func (a *AuthLogic) Register(ctx context.Context, email, name, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 0)
	if err != nil {
		return err
	}

	err = a.userRepo.Create(ctx, database.User{
		ID:             uuid.NewString(),
		Email:          email,
		Name:           name,
		HashedPassword: string(hashedPassword),
		CreatedAt:      time.Now(),
		LastLogin:      time.Now(),
	})
	if err != nil {
		return err
	}

	return nil
}

type Profile struct {
	Email      string
	Name       string
	GroupID    *string
	GroupName  string
	GroupOwner bool
	Members    []string
}

func (a *AuthLogic) GetProfile(ctx context.Context, userID string) (Profile, error) {
	user, err := a.userRepo.Get(ctx, userID)
	if err != nil {
		return Profile{}, err
	}

	groupName := ""
	groupOwner := false
	var members []string
	if user.GroupID != nil {
		group, err := a.groupRepo.Get(ctx, *user.GroupID)
		if err != nil {
			return Profile{}, err
		}
		groupName = group.Name
		groupOwner = group.OwnerUserID == user.ID
		users, err := a.userRepo.GetByGroup(ctx, *user.GroupID)
		if err != nil {
			return Profile{}, err
		}

		for _, user := range users {
			members = append(members, user.Name)
		}
	}

	return Profile{
		Email:      user.Email,
		Name:       user.Name,
		GroupID:    user.GroupID,
		GroupName:  groupName,
		GroupOwner: groupOwner,
		Members:    members,
	}, nil
}

func (a *AuthLogic) LeaveGroup(ctx context.Context, userID string) error {
	return a.userRepo.SetGroup(ctx, userID, nil)
}
