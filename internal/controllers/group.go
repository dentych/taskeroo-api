package controllers

import (
	"errors"
	"github.com/dentych/taskeroo/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type GroupController struct {
	groupRepo *database.GroupRepo
	userRepo  *database.UserRepo
}

func NewGroupController(protectedRouter gin.IRouter, groupRepo *database.GroupRepo, userRepo *database.UserRepo) *GroupController {
	handler := &GroupController{groupRepo: groupRepo, userRepo: userRepo}

	protectedRouter.GET("/group/create", handler.GetCreateGroup())

	protectedRouter.POST("/group/create", handler.PostCreateGroup())

	protectedRouter.GET("/group/members/add", handler.GetAddGroupMember())
	protectedRouter.POST("/group/members/add", handler.PostAddGroupMember())

	return handler
}

func (c *GroupController) GetCreateGroup() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		HTML(ctx, http.StatusOK, "pages/create-group", gin.H{
			"title": "Opret gruppe",
		})
	}
}

func (c *GroupController) PostCreateGroup() gin.HandlerFunc {
	return func(ctx *gin.Context) {
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

		groupID := uuid.NewString()
		err := c.groupRepo.Create(ctx.Request.Context(), database.Group{ID: groupID, Name: name, OwnerUserID: userID})
		if err != nil {
			log.Printf("Failed to create team: %s\n", err)
			ctx.Status(http.StatusInternalServerError)
			return
		}

		err = c.userRepo.SetGroup(ctx.Request.Context(), userID, &groupID)
		if err != nil {
			log.Printf("Failed to create team: %s\n", err)
			ctx.Status(http.StatusInternalServerError)
			return
		}

		ctx.Redirect(http.StatusFound, "/")
	}
}

func (c *GroupController) GetAddGroupMember() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		HTML(ctx, http.StatusOK, "pages/add-member", gin.H{
			"title": "Tilføj medlem",
		})
	}
}

func (c *GroupController) PostAddGroupMember() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		email := ctx.PostForm("email")
		if email == "" {
			HTML(ctx, http.StatusOK, "pages/add-member", gin.H{
				"title": "Tilføj medlem",
				"error": "Email adresse skal udfyldes",
			})
			return
		}

		userID := ctx.GetString(KeyUserID)
		inviter, err := c.userRepo.Get(ctx, userID)
		if err != nil {
			log.Printf("Failed to get inviting user=%s: %s\n", userID, err)
			HTML(ctx, http.StatusInternalServerError, "pages/add-member", gin.H{
				"title": "Tilføj medlem",
				"error": "Der var en fejl da brugeren skulle tilføjes. Prøv igen senere, eller kontakt support hvis problemet bliver ved.",
			})
			return
		}

		invitee, err := c.userRepo.GetByEmail(ctx, email)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				HTML(ctx, http.StatusOK, "pages/add-member", gin.H{
					"title": "Tilføj medlem",
					"error": "Brugeren findes ikke.",
				})
				return
			}
			log.Printf("Failed to get invitee by email=%s: %s\n", email, err)
			HTML(ctx, http.StatusInternalServerError, "pages/add-member", gin.H{
				"title": "Tilføj medlem",
				"error": "Der var en fejl da brugeren skulle tilføjes. Prøv igen senere, eller kontakt support hvis problemet bliver ved.",
			})
			return
		}

		if invitee.GroupID != nil {
			HTML(ctx, http.StatusOK, "pages/add-member", gin.H{
				"title": "Tilføj medlem",
				"error": "Brugeren er allerede i en gruppe.",
			})
			return
		}

		err = c.userRepo.SetGroup(ctx, invitee.ID, inviter.GroupID)
		if err != nil {
			log.Printf("Failed to get invitee by email=%s: %s\n", email, err)
			HTML(ctx, http.StatusInternalServerError, "pages/add-member", gin.H{
				"title": "Tilføj medlem",
				"error": "Der var en fejl da brugeren skulle tilføjes. Prøv igen senere, eller kontakt support hvis problemet bliver ved.",
			})
			return
		}

		ctx.Redirect(http.StatusFound, "/profile")
	}
}
