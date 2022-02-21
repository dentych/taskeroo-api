package app

import (
	"github.com/dentych/taskeroo/internal/database"
	"github.com/gin-gonic/gin"
)

type TelegramLogic struct {
	telegramRepo *database.TelegramRepo
}

func NewTelegramLogic(telegramRepo *database.TelegramRepo) *TelegramLogic {
	return &TelegramLogic{telegramRepo: telegramRepo}
}

func (l *TelegramLogic) Connect(ctx *gin.Context, userID string, connectID string) error {
	return l.telegramRepo.SetUserID(ctx, connectID, userID)
}
