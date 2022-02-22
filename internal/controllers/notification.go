package controllers

import (
	"github.com/dentych/taskeroo/internal/app"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type NotificationController struct {
	notificationLogic *app.NotificationLogic
}

func NewNotificationController(protectedRouter gin.IRouter, notificationLogic *app.NotificationLogic) *NotificationController {
	handler := &NotificationController{notificationLogic: notificationLogic}

	protectedRouter.GET("/notifications", handler.GetNotifications())

	return handler
}

func (c *NotificationController) GetNotifications() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := ctx.GetString(KeyUserID)
		notificationInfo, err := c.notificationLogic.GetNotificationInfo(ctx, userID)
		if err != nil {
			log.Printf("Failed to get group discord settings, for user=%s: %s\n", userID, err)
			HTML(ctx, http.StatusInternalServerError, "pages/notifications", gin.H{
				"title": "Notifikationsindstillinger",
				"error": "Der skete en fejl. Prøv igen om lidt.",
			})
			return
		}

		HTML(ctx, http.StatusOK, "pages/notifications", gin.H{
			"title":          "Notifikationsindstillinger",
			"groupOwner":     notificationInfo.GroupOwner,
			"telegramActive": notificationInfo.TelegramActive,
		})
	}
}
