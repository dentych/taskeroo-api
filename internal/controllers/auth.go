package controllers

import (
	"errors"
	"github.com/dentych/taskeroo/internal/auth"
	internalerrors "github.com/dentych/taskeroo/internal/errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type AuthController struct {
	authService   *auth.Auth
	secureCookies bool
}

func NewAuthController(
	router gin.IRouter,
	protectedRouter gin.IRouter,
	authService *auth.Auth,
	secureCookies bool,
) *AuthController {
	handler := &AuthController{authService: authService, secureCookies: secureCookies}
	router.GET("/login", handler.GetLogin())
	router.POST("/login", handler.PostLogin())

	router.GET("/register", handler.GetRegister())
	router.POST("/register", handler.PostRegister())

	protectedRouter.GET("/logout", handler.GetLogout())

	protectedRouter.GET("/profile", handler.GetProfile())

	return handler
}

func AuthMiddleware(authService *auth.Auth) gin.HandlerFunc {
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

func (c *AuthController) GetLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		HTML(ctx, http.StatusOK, "pages/login", gin.H{
			"title": "Login",
		})
	}
}

func (c *AuthController) PostLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		email := ctx.PostForm("email")
		password := ctx.PostForm("password")

		userSession, err := c.authService.Login(ctx.Request.Context(), email, password)
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

		ctx.SetCookie(CookieKeyUserID, userSession.UserID, int(Time31Days.Seconds()), "", "", c.secureCookies, true)
		ctx.SetCookie(CookieKeySession, userSession.Session, int(Time31Days.Seconds()), "", "", c.secureCookies, true)
		ctx.Redirect(http.StatusFound, "/")
	}
}

func (c *AuthController) GetRegister() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		HTML(ctx, http.StatusOK, "pages/register", gin.H{
			"title": "Register",
		})
	}
}

func (c *AuthController) PostRegister() gin.HandlerFunc {
	return func(ctx *gin.Context) {
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

		err := c.authService.Register(ctx.Request.Context(), email, password)
		if err != nil {
			HTML(ctx, http.StatusInternalServerError, "pages/index", nil)
			return
		}

		ctx.Redirect(http.StatusFound, "/login")
	}
}

func (c *AuthController) GetLogout() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		clearCookies(ctx)
		ctx.Redirect(http.StatusFound, "/login")
	}
}

func (c *AuthController) GetProfile() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		HTML(ctx, http.StatusOK, "pages/profile", gin.H{
			"title": "Profil",
		})
	}
}