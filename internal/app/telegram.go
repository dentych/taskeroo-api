package app

import (
	"context"
	"fmt"
	"github.com/dentych/taskeroo/internal/database"
	"github.com/dentych/taskeroo/internal/telegram"
	"time"
)

var ErrNotValid = fmt.Errorf("connect link no longer valid")

type TelegramLogic struct {
	telegramRepo   *database.TelegramRepo
	telegramClient *telegram.Telegram
}

func NewTelegramLogic(telegramRepo *database.TelegramRepo) *TelegramLogic {
	return &TelegramLogic{telegramRepo: telegramRepo}
}

func (l *TelegramLogic) Connect(ctx context.Context, userID string, connectID string) error {
	tele, err := l.telegramRepo.Get(ctx, connectID)
	if err != nil {
		return err
	}

	if tele.UserID != nil {
		return ErrNotValid
	}

	if time.Now().Sub(tele.CreatedAt) > 24*time.Hour {
		return ErrNotValid
	}

	err = l.telegramRepo.DeleteAllByUserID(ctx, userID)
	if err != nil {
		return err
	}

	return l.telegramRepo.SetUserID(ctx, connectID, userID)
}

func (l *TelegramLogic) SendMessage(ctx context.Context, userID string, msg string) error {
	dbTelegram, err := l.telegramRepo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	return l.telegramClient.SendMessage(ctx, dbTelegram.TelegramUserID, msg)
}