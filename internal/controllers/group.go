package controllers

import (
	"github.com/dentych/taskeroo/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
