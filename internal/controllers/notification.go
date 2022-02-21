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
	protectedRouter.POST("/notifications", handler.PostNotifications())

	protectedRouter.GET("/notifications/group-setup", handler.GetSetupGroupNotifications())
	protectedRouter.POST("/notifications/group-setup", handler.PostSetupGroupNotifications())

	protectedRouter.POST("/debug/notify", handler.PostDebugNotify())

	return handler
}

func (c *NotificationController) GetSetupGroupNotifications() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		HTML(ctx, http.StatusOK, "pages/group-notifications-setup", gin.H{
			"title": "Notifikationsindstillinger",
		})
	}
}

func (c *NotificationController) PostSetupGroupNotifications() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := ctx.GetString(KeyUserID)
		webhook := ctx.PostForm("webhook")
		if webhook == "" {
			HTML(ctx, http.StatusBadRequest, "pages/group-notifications-setup", gin.H{
				"title": "Notifikationsindstillinger",
				"error": "Webhook er et påkrævet felt og skal udfyldes.",
			})
			return
		}
		err := c.notificationLogic.SetupGroupDiscord(ctx, userID, webhook)
		if err != nil {
			HTML(ctx, http.StatusInternalServerError, "pages/group-notifications-setup", gin.H{
				"title": "Notifikationsindstillinger",
				"error": "Der skete en fejl ved oprettelse af notifikationsindstillinger. Prøv igen om lidt.",
			})
			return
		}

		ctx.Redirect(http.StatusFound, "/notifications")
	}
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
			"title":           "Notifikationsindstillinger",
			"groupOwner":      notificationInfo.GroupOwner,
			"discordActive":   notificationInfo.DiscordActive,
			"discordUsername": notificationInfo.Username,
		})
	}
}

func (c *NotificationController) PostNotifications() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := ctx.GetString(KeyUserID)
		discordUsername := ctx.PostForm("username")
		err := c.notificationLogic.SetupDiscordUsername(ctx, userID, discordUsername)
		if err != nil {
			log.Printf("Failed to setup discord username for user=%s: %s\n", userID, err)
			HTML(ctx, http.StatusInternalServerError, "pages/notifications", gin.H{
				"title": "Notifikationsindstillinger",
				"error": "Der skete en fejl. Prøv igen om lidt.",
			})
			return
		}

		ctx.Status(http.StatusOK)
	}
}

func (c *NotificationController) PostDebugNotify() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := ctx.GetString(KeyUserID)
		err := c.notificationLogic.SendNotification(ctx, userID, "Dette er en test notifikation")
		if err != nil {
			log.Printf("Error sending notification to user=%s: %s\n", userID, err)
			ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "There was an error"})
			return
		}

		ctx.Status(http.StatusOK)
	}
}
