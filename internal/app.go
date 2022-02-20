package internal

import (
	"errors"
	"fmt"
	"github.com/dentych/taskeroo/internal/auth"
	"github.com/dentych/taskeroo/internal/database"
	internalerrors "github.com/dentych/taskeroo/internal/errors"
	"github.com/foolin/goview"
	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	CookieKeyEmail   = "auth_email"
	CookieKeySession = "auth_session"
)

var (
	Time31Days = 31 * 24 * time.Hour
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
	protectedRouter.GET("/", func(ctx *gin.Context) {
		//render with master
		ctx.HTML(http.StatusOK, "pages/index", gin.H{
			"title": "Taskeroo",
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

	router.POST("/login", func(ctx *gin.Context) {
		email := ctx.PostForm("email")
		password := ctx.PostForm("password")

		session, err := auth.Login(email, password)
		if err != nil {
			if errors.Is(err, internalerrors.InvalidEmailOrPassword) {
				ctx.HTML(http.StatusOK, "pages/login", gin.H{
					"title": "Login",
					"error": "Email eller password ugyldig.",
				})
				return
			}
			ctx.HTML(http.StatusInternalServerError, "", nil)
		}

		ctx.SetCookie(CookieKeyEmail, email, int(Time31Days.Seconds()), "", "", true, true)
		ctx.SetCookie(CookieKeySession, session, int(Time31Days.Seconds()), "", "", true, true)
		ctx.Redirect(http.StatusFound, "/")
	})

	router.GET("/register", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "pages/register", gin.H{
			"title": "Register",
		})
	})

	router.POST("/register", func(ctx *gin.Context) {
		email := ctx.PostForm("email")
		password := ctx.PostForm("password")

		if email == "" {
			ctx.HTML(http.StatusOK, "pages/login", gin.H{
				"title": "Login",
				"error": "Email felt skal udfyldes.",
			})
			return
		}
		if password == "" {
			ctx.HTML(http.StatusOK, "pages/login", gin.H{
				"title": "Login",
				"error": "Password felt skal udfyldes.",
			})
			return
		}

		err := auth.Register(email, password)
		if err != nil {
			ctx.HTML(http.StatusInternalServerError, "", nil)
			return
		}

		ctx.Redirect(http.StatusFound, "/login")
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
			clearCookies(ctx)
			ctx.Redirect(http.StatusFound, "/login")
			return
		}
		userID, err := ctx.Request.Cookie("auth_userid")
		if err != nil || userID == nil || !userID.HttpOnly || !userID.Secure {
			clearCookies(ctx)
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
			clearCookies(ctx)
			ctx.Redirect(http.StatusFound, "/login")
			return
		}

		ctx.Next()
	}
}

func clearCookies(ctx *gin.Context) {
	ctx.SetCookie(CookieKeyEmail, "", -1, "", "", true, true)
	ctx.SetCookie(CookieKeySession, "", -1, "", "", true, true)
}
