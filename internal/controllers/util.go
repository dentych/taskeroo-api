package controllers

import (
	"github.com/gin-gonic/gin"
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

func HTML(ctx *gin.Context, status int, templateName string, obj gin.H) {
	if value := ctx.GetString("userID"); value != "" {
		obj["userID"] = value
	}
	ctx.HTML(status, templateName, obj)
}

func clearCookies(ctx *gin.Context) {
	ctx.SetCookie(CookieKeyUserID, "", -1, "", "", true, true)
	ctx.SetCookie(CookieKeySession, "", -1, "", "", true, true)
}
