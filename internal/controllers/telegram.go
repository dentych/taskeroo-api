package controllers

import (
	"errors"
	"github.com/dentych/taskeroo/internal/app"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type TelegramController struct {
	telegramLogic *app.TelegramLogic
}

func NewTelegramController(protectedRouter gin.IRouter, telegramLogic *app.TelegramLogic) *TelegramController {
	handler := &TelegramController{telegramLogic: telegramLogic}

	protectedRouter.GET("/telegram/connect/:connectID", handler.GetTelegramConnect())

	return handler
}

func (c *TelegramController) GetTelegramConnect() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := ctx.GetString(KeyUserID)
		connectID := ctx.Param("connectID")

		err := c.telegramLogic.Connect(ctx, userID, connectID)
		if err != nil {
			if errors.Is(err, app.ErrNotValid) {
				HTML(ctx, http.StatusBadRequest, "pages/telegram-connect", gin.H{
					"title": "Telegram forbindelse",
					"error": "Connect linket er ikke længere gyldigt. Opret et nyt link gennem Telegram botten.",
				})
				return
			}
			log.Printf("Failed to connect telegram for user=%s, connectID=%s: %s\n", userID, connectID, err)
			HTML(ctx, http.StatusInternalServerError, "pages/telegram-connect", gin.H{
				"title": "Telegram forbindelse",
				"error": "Der var en fejl med at forbinde din Telegram konto til din Taskeroo konto. Prøv igen om lidt.",
			})
			return
		}

		HTML(ctx, http.StatusOK, "pages/telegram-connect", gin.H{
			"title":   "Telegram forbindelse",
			"success": "Telegram er nu forbundet til din Taskeroo konto! 👏",
		})
	}
}
