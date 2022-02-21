package app

import (
	"context"
	"github.com/dentych/taskeroo/internal/database"
	internalerrors "github.com/dentych/taskeroo/internal/errors"
	"time"
)

type NotificationLogic struct {
	notificationRepo *database.NotificationRepo
	userRepo         *database.UserRepo
	groupRepo        *database.GroupRepo
}

func NewNotificationLogic(
	notificationRepo *database.NotificationRepo,
	userRepo *database.UserRepo,
	groupRepo *database.GroupRepo,
) *NotificationLogic {
	return &NotificationLogic{notificationRepo: notificationRepo, userRepo: userRepo, groupRepo: groupRepo}
}

func (n *NotificationLogic) SetupGroupDiscord(ctx context.Context, userID string, webhook string) error {
	user, err := n.userRepo.Get(ctx, userID)
	if err != nil {
		return err
	}

	if user.GroupID == nil {
		return internalerrors.ErrUserNotInGroup
	}

	group, err := n.groupRepo.Get(ctx, *user.GroupID)
	if err != nil {
		return err
	}

	if user.ID != group.OwnerUserID {
		return internalerrors.ErrUserNotOwner
	}

	return n.notificationRepo.CreateDiscord(ctx, database.GroupDiscord{
		GroupID:   group.ID,
		Webhook:   webhook,
		CreatedAt: time.Now(),
	})
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
		discord, err := n.notificationRepo.GetGroupDiscord(ctx, group.ID)
		if err != nil {
			return nil, err
		}
		if discord != nil {
			output.DiscordActive = true
		}
	}

	username, err := n.notificationRepo.GetDiscordUsername(ctx, userID)
	if err != nil {
		return nil, err
	}

	if username != nil {
		output.Username = username.DiscordUsername
	}

	return output, nil
}

func (n *NotificationLogic) SetupDiscordUsername(ctx context.Context, userID string, discordUsername string) error {
	err := n.notificationRepo.CreateDiscordUsername(ctx, userID, discordUsername)
	if err != nil {
		return err
	}

	return nil
}

type NotificationInfo struct {
	GroupOwner    bool
	DiscordActive bool
	Username      string
}
