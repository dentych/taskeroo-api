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
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
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

	err = db.AutoMigrate(&database.User{}, &database.Session{}, &database.Group{})
	if err != nil {
		log.Fatalf("Failed to migrate database models: %s\n", err)
	}

	userRepo := database.NewUserRepo(db)
	sessionRepo := database.NewSessionRepo(db)
	teamRepo := database.NewTeamRepo(db)

	auth := auth.New(sessionRepo, userRepo)

	goviewConfig := goview.DefaultConfig
	if os.Getenv("ENVIRONMENT") != "prod" {
		goviewConfig.DisableCache = true
	}
	router.HTMLRender = ginview.New(goviewConfig)

	protectedRouter := router.Group("")
	protectedRouter.Use(authMiddleware(auth))
	protectedRouter.GET("/", func(ctx *gin.Context) {
		userID := ctx.GetString(KeyUserID)
		user, err := userRepo.Get(ctx.Request.Context(), userID)
		if err != nil {
			log.Printf("Failed to get user with ID '%s': %s\n", userID, err)
			HTML(ctx, http.StatusInternalServerError, "pages/index", gin.H{
				"title": "Taskeroo",
			})
			return
		}
		HTML(ctx, http.StatusOK, "pages/index", gin.H{
			"title":  "Taskeroo",
			"teamID": user.GroupID,
		})
	})

	protectedRouter.GET("/profile", func(ctx *gin.Context) {
		HTML(ctx, http.StatusOK, "pages/profile", gin.H{
			"title": "Profil",
		})
	})

	protectedRouter.GET("/group/create", func(ctx *gin.Context) {
		HTML(ctx, http.StatusOK, "pages/create-group", gin.H{
			"title": "Opret gruppe",
		})
	})

	protectedRouter.POST("/group/create", func(ctx *gin.Context) {
		name := ctx.PostForm("name")
		if name == "" {
			HTML(ctx, http.StatusBadRequest, "pages/create-group", gin.H{
				"title": "Opret gruppe",
				"error": "Gruppens navn skal udfyldes",
			})
			return
		}

		userID := ctx.GetString(KeyUserID)
		if userID == "" {
			ctx.Status(http.StatusInternalServerError)
			return
		}

		teamID := uuid.NewString()
		err := teamRepo.Create(ctx.Request.Context(), database.Group{ID: teamID, Name: name, OwnerUserID: userID})
		if err != nil {
			log.Printf("Failed to create team: %s\n", err)
			ctx.Status(http.StatusInternalServerError)
			return
		}

		err = userRepo.SetGroup(ctx.Request.Context(), userID, teamID)
		if err != nil {
			log.Printf("Failed to create team: %s\n", err)
			ctx.Status(http.StatusInternalServerError)
			return
		}

		ctx.Redirect(http.StatusFound, "/")
	})

	router.GET("/login", func(ctx *gin.Context) {
		HTML(ctx, http.StatusOK, "pages/login", gin.H{
			"title": "Login",
		})
	})

	router.POST("/login", func(ctx *gin.Context) {
		email := ctx.PostForm("email")
		password := ctx.PostForm("password")

		userSession, err := auth.Login(ctx.Request.Context(), email, password)
		if err != nil {
			if errors.Is(err, internalerrors.InvalidEmailOrPassword) {
				HTML(ctx, http.StatusOK, "pages/login", gin.H{
					"title": "Login",
					"error": "Email eller password ugyldig",
				})
				return
			}
			HTML(ctx, http.StatusInternalServerError, "pages/index", nil)
		}

		ctx.SetCookie(CookieKeyUserID, userSession.UserID, int(Time31Days.Seconds()), "", "", secureCookies, true)
		ctx.SetCookie(CookieKeySession, userSession.Session, int(Time31Days.Seconds()), "", "", secureCookies, true)
		ctx.Redirect(http.StatusFound, "/")
	})

	router.GET("/register", func(ctx *gin.Context) {
		HTML(ctx, http.StatusOK, "pages/register", gin.H{
			"title": "Register",
		})
	})

	router.POST("/register", func(ctx *gin.Context) {
		email := ctx.PostForm("email")
		password := ctx.PostForm("password")
		repeatedPassword := ctx.PostForm("repeated-password")

		if password != repeatedPassword {
			HTML(ctx, http.StatusBadRequest, "pages/register", gin.H{
				"title": "Login",
				"error": "De to passwords matcher ikke",
			})
			return
		}

		if email == "" {
			HTML(ctx, http.StatusBadRequest, "pages/register", gin.H{
				"title": "Login",
				"error": "Email felt skal udfyldes",
			})
			return
		}
		if password == "" {
			HTML(ctx, http.StatusBadRequest, "pages/register", gin.H{
				"title": "Login",
				"error": "Password felt skal udfyldes",
			})
			return
		}

		err := auth.Register(ctx.Request.Context(), email, password)
		if err != nil {
			HTML(ctx, http.StatusInternalServerError, "pages/index", nil)
			return
		}

		ctx.Redirect(http.StatusFound, "/login")
	})

	protectedRouter.GET("/logout", func(ctx *gin.Context) {
		clearCookies(ctx)
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
		userID, err := ctx.Cookie(CookieKeyUserID)
		if err != nil || userID == "" {
			clearCookies(ctx)
			ctx.Redirect(http.StatusFound, "/login")
			return
		}
		session, err := ctx.Cookie(CookieKeySession)
		if err != nil || session == "" {
			clearCookies(ctx)
			ctx.Redirect(http.StatusFound, "/login")
			return
		}

		authenticated, err := authService.IsAuthenticated(ctx.Request.Context(), userID, session)
		if err != nil {
			log.Printf("Failed to check if user is authenticated: %s\n", err)
			HTML(ctx, http.StatusInternalServerError, "pages/index", nil)
			return
		}

		if !authenticated {
			clearCookies(ctx)
			ctx.Redirect(http.StatusFound, "/login")
			return
		}

		ctx.Set("userID", userID)
		ctx.Set("session", session)
		ctx.Next()
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
