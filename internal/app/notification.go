package app

import (
	"context"
	"errors"
	"github.com/dentych/taskeroo/internal/database"
	internalerrors "github.com/dentych/taskeroo/internal/errors"
	"gorm.io/gorm"
	"log"
)

type NotificationLogic struct {
	notificationRepo *database.NotificationRepo
	userRepo         *database.UserRepo
	groupRepo        *database.GroupRepo
	telegramRepo     *database.TelegramRepo
	telegramLogic    *TelegramLogic
}

func NewNotificationLogic(
	notificationRepo *database.NotificationRepo,
	userRepo *database.UserRepo,
	groupRepo *database.GroupRepo,
	telegramRepo *database.TelegramRepo,
	telegramLogic *TelegramLogic,
) *NotificationLogic {
	return &NotificationLogic{
		notificationRepo: notificationRepo,
		userRepo:         userRepo,
		groupRepo:        groupRepo,
		telegramRepo:     telegramRepo,
		telegramLogic:    telegramLogic,
	}
}

func (n *NotificationLogic) GetNotificationInfo(ctx context.Context, userID string) (*NotificationInfo, error) {
	user, err := n.userRepo.Get(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user.GroupID == nil {
		return nil, internalerrors.ErrUserNotInGroup
	}

	group, err := n.groupRepo.Get(ctx, *user.GroupID)
	if err != nil {
		return nil, err
	}

	output := &NotificationInfo{}

	if user.ID == group.OwnerUserID {
		output.GroupOwner = true
	}

	dbTelegram, err := n.telegramRepo.GetByUserID(ctx, userID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if dbTelegram != nil {
		output.TelegramActive = true
	}

	return output, nil
}

func (n *NotificationLogic) SendNotification(ctx context.Context, userID string, msg string) error {
	log.Printf("Sending notification to user=%s, msg:\n%s\n", userID, msg)
	return n.telegramLogic.SendMessage(ctx, userID, msg)
}

func (n *NotificationLogic) NotifyAllInGroup(ctx context.Context, groupID string, msg string) error {
	users, err := n.userRepo.GetByGroup(ctx, groupID)
	if err != nil {
		return err
	}

	for _, user := range users {
		err = n.SendNotification(ctx, user.ID, msg)
		if err != nil {
			log.Printf("Failed to send message to a member of group=%s, user=%s: %s", groupID, user.ID, err)
		}
	}

	return nil
}

type NotificationInfo struct {
	GroupOwner     bool
	TelegramActive bool
}
