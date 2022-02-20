package internal

import (
	"fmt"
	"github.com/dentych/taskeroo/internal/auth"
	"github.com/dentych/taskeroo/internal/database"
	"github.com/foolin/goview"
	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
)

func Run() {
	router := gin.Default()

	var dsn string
	if os.Getenv("ENVIRONMENT") == "prod" {
		dsn = os.Getenv("DATABASE_URL")
		if dsn == "" {
			log.Fatalf("DATABASE_URL is not set, but is required.")
		}
	} else {
		dsn = "postgres://postgres:postgres@localhost/postgres"
	}
	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		log.Fatalf("Failed to establish database connection: %s\n", err)
	}

	err = db.AutoMigrate(&database.User{}, &database.Session{})
	if err != nil {
		log.Fatalf("Failed to migrate database models: %s\n", err)
	}

	userRepo := database.NewUserRepo(db)
	sessionRepo := database.NewSessionRepo(db)

	auth := auth.New(sessionRepo, userRepo)

	goviewConfig := goview.DefaultConfig
	if os.Getenv("ENVIRONMENT") != "prod" {
		goviewConfig.DisableCache = true
	}
	router.HTMLRender = ginview.New(goviewConfig)

	protectedRouter := router.Group("")
	protectedRouter.Use(authMiddleware(auth))
	router.GET("/", func(ctx *gin.Context) {
		//render with master
		ctx.HTML(http.StatusOK, "pages/index", gin.H{
			"title": "Index title!",
			"add": func(a int, b int) int {
				return a + b
			},
		})
	})

	router.GET("/login", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "pages/login", gin.H{
			"title": "Login",
		})
	})

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

func authMiddleware(authService *auth.Auth) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session, err := ctx.Request.Cookie("auth_session")
		if err != nil || session == nil || !session.HttpOnly || !session.Secure {
			ctx.SetCookie("auth_session", "", -1, "", "", true, true)
			ctx.SetCookie("auth_userid", "", -1, "", "", true, true)
			ctx.Redirect(http.StatusFound, "/login")
			return
		}
		userID, err := ctx.Request.Cookie("auth_userid")
		if err != nil || userID == nil || !userID.HttpOnly || !userID.Secure {
			ctx.SetCookie("auth_session", "", -1, "", "", true, true)
			ctx.SetCookie("auth_userid", "", -1, "", "", true, true)
			ctx.Redirect(http.StatusFound, "/login")
			return
		}

		authenticated, err := authService.IsAuthenticated(userID.Value, session.Value)
		if err != nil {
			log.Printf("Failed to check if user is authenticated: %s\n", err)
			ctx.HTML(http.StatusInternalServerError, "", nil)
			return
		}

		if !authenticated {
			ctx.SetCookie("auth_session", "", -1, "", "", true, true)
			ctx.SetCookie("auth_userid", "", -1, "", "", true, true)
			ctx.Redirect(http.StatusFound, "/login")
			return
		}

		ctx.Next()
	}
}
