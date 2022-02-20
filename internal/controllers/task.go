package controllers

import (
	"github.com/dentych/taskeroo/internal/app"
	"github.com/dentych/taskeroo/internal/database"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type TaskController struct {
	userRepo *database.UserRepo
}

func NewTaskController(protectedRouter gin.IRouter, userRepo *database.UserRepo) *TaskController {
	handler := &TaskController{userRepo: userRepo}

	protectedRouter.GET("/", handler.GetIndex())

	protectedRouter.GET("/task/create", handler.GetCreateTask())

	return handler
}

func (c *TaskController) GetIndex() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := ctx.GetString(KeyUserID)
		user, err := c.userRepo.Get(ctx.Request.Context(), userID)
		if err != nil {
			log.Printf("Failed to get user with ID '%s': %s\n", userID, err)
			HTML(ctx, http.StatusInternalServerError, "pages/index", gin.H{
				"title": "Taskeroo",
			})
			return
		}
		HTML(ctx, http.StatusOK, "pages/index", gin.H{
			"title":   "Taskeroo",
			"groupID": user.GroupID,
			"tasks": []app.Task{
				{
					ID:          "1234",
					GroupID:     "1234",
					Title:       "Vask 30 grader",
					Description: "Vask tøj på 30 grader, hver uge",
				},
				{
					ID:          "1235",
					GroupID:     "1235",
					Title:       "Vask 60 grader",
					Description: "Vask tøj på 60 grader, hver uge",
				},
				{
					ID:          "1236",
					GroupID:     "1236",
					Title:       "Støvsug ovenpå",
					Description: "Støvsug alle værelser ovenpå: Badeværelse, soveværelse, gæsteværelse, kontor.",
				},
			},
		})
	}
}

func (c *TaskController) GetCreateTask() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		HTML(ctx, http.StatusOK, "pages/index", gin.H{
			"title": "Opret opgave",
		})
	}
}
