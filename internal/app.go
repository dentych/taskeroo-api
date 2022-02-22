package internal

import (
	"fmt"
	"github.com/dentych/taskeroo/internal/app"
	"github.com/dentych/taskeroo/internal/controllers"
	"github.com/dentych/taskeroo/internal/database"
	"github.com/dentych/taskeroo/internal/telegram"
	"github.com/foolin/goview"
	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"time"
)

const (
	CookieKeyUserID  = "auth_userid"
	CookieKeySession = "auth_session"

	KeyUserID  = "userID"
	KeySession = "session"
)

var (
	Time31Days = 31 * 24 * time.Hour
)

var secureCookies bool

func Run() {
	router := gin.Default()

	var dsn string
	if os.Getenv("ENVIRONMENT") == "prod" {
		dsn = os.Getenv("DATABASE_URL")
		if dsn == "" {
			log.Fatalf("DATABASE_URL is not set, but is required.")
		}
		secureCookies = true
	} else {
		dsn = "postgres://postgres:postgres@localhost/postgres"
	}
	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		log.Fatalf("Failed to establish database connection: %s\n", err)
	}

	err = db.AutoMigrate(
		&database.User{},
		&database.Session{},
		&database.Group{},
		&database.Task{},
		&database.GroupDiscord{},
		&database.DiscordUsername{},
		&database.Telegram{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database models: %s\n", err)
	}

	userRepo := database.NewUserRepo(db)
	sessionRepo := database.NewSessionRepo(db)
	groupRepo := database.NewGroupRepo(db)
	taskRepo := database.NewTaskRepo(db)
	notificationRepo := database.NewNotificationRepo(db)
	telegramRepo := database.NewTelegramRepo(db)
	telegramClient := telegram.NewTelegram(telegramRepo, os.Getenv("TELEGRAM_TOKEN"))

	authService := app.NewAuthLogic(sessionRepo, userRepo, groupRepo)
	taskLogic := app.NewTaskLogic(taskRepo, userRepo)
	notificationLogic := app.NewNotificationLogic(notificationRepo, userRepo, groupRepo, telegramRepo)
	telegramLogic := app.NewTelegramLogic(telegramRepo, telegramClient)

	goviewConfig := goview.DefaultConfig
	if os.Getenv("ENVIRONMENT") != "prod" {
		goviewConfig.DisableCache = true
	}
	router.HTMLRender = ginview.New(goviewConfig)

	protectedRouter := router.Group("")
	protectedRouter.Use(controllers.AuthMiddleware(authService))

	controllers.NewAuthController(router, protectedRouter, authService, secureCookies)
	controllers.NewGroupController(protectedRouter, groupRepo, userRepo)
	controllers.NewTaskController(protectedRouter, userRepo, taskLogic)
	controllers.NewNotificationController(protectedRouter, notificationLogic)
	controllers.NewTelegramController(protectedRouter, telegramLogic)

	err = telegramClient.Start()
	if err != nil {
		log.Fatalf("Failed to start, because Telegram Bot could not start: %s\n", err)
	}

	port := "8080"
	portEnv := os.Getenv("PORT")
	if portEnv != "" {
		port = portEnv
	}
	err = router.Run(fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("Error running server: %s\n", err)
	}
}

func clearCookies(ctx *gin.Context) {
	ctx.SetCookie(CookieKeyUserID, "", -1, "", "", secureCookies, true)
	ctx.SetCookie(CookieKeySession, "", -1, "", "", secureCookies, true)
}

func HTML(ctx *gin.Context, status int, templateName string, obj gin.H) {
	if value := ctx.GetString("userID"); value != "" {
		obj["userID"] = value
	}
	ctx.HTML(status, templateName, obj)
}
