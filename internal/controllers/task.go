package controllers

import (
	"github.com/dentych/taskeroo/internal/app"
	"github.com/dentych/taskeroo/internal/database"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type TaskController struct {
	userRepo  *database.UserRepo
	taskLogic *app.TaskLogic
}

func NewTaskController(
	protectedRouter gin.IRouter,
	userRepo *database.UserRepo,
	taskLogic *app.TaskLogic,
) *TaskController {
	handler := &TaskController{userRepo: userRepo, taskLogic: taskLogic}

	protectedRouter.GET("/", handler.GetIndex())

	protectedRouter.GET("/task/create", handler.GetCreateTask())
	protectedRouter.POST("/task/create", handler.PostCreateTask())

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
				"alert": "",
			})
			return
		}
		if user.GroupID == nil {
			log.Printf("User=%s is not in a group, so can't retrieve tasks\n", userID)
			HTML(ctx, http.StatusBadRequest, "pages/index", gin.H{
				"title": "Taskeroo",
			})
			return
		}

		tasks, err := c.taskLogic.GetForGroup(ctx.Request.Context(), userID, *user.GroupID)

		HTML(ctx, http.StatusOK, "pages/index", gin.H{
			"title":   "Taskeroo",
			"groupID": user.GroupID,
			"tasks":   tasks,
		})
	}
}

func (c *TaskController) GetCreateTask() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		HTML(ctx, http.StatusOK, "pages/create-task", gin.H{
			"title": "Opret opgave",
		})
	}
}

func (c *TaskController) PostCreateTask() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		title := ctx.PostForm("title")
		description := ctx.PostForm("description")
		intervalSize := ctx.PostForm("intervalSize")
		intervalUnit := ctx.PostForm("intervalUnit")

		if title == "" {
			HTML(ctx, http.StatusBadRequest, "pages/create-task", gin.H{
				"title": "Opret opgave",
				"error": "Titel skal udfyldes",
			})
			return
		}
		if description == "" {
			HTML(ctx, http.StatusBadRequest, "pages/create-task", gin.H{
				"title": "Opret opgave",
				"error": "Beskrivelse skal udfyldes",
			})
			return
		}
		if intervalUnit == "" {
			HTML(ctx, http.StatusBadRequest, "pages/create-task", gin.H{
				"title": "Opret opgave",
				"error": "Opgavens hyppighed skal udfyldes",
			})
			return
		}

		formattedIntervalSize := 0
		if intervalSize != "" {
			var err error
			formattedIntervalSize, err = strconv.Atoi(intervalSize)
			if err != nil {
				HTML(ctx, http.StatusBadRequest, "pages/create-test", gin.H{
					"title": "Opret opgave",
					"error": "Opgavens hyppighed skal defineres i tal og engangsopgave, dag, uge, måned.",
				})
				return
			}
		}

		userID := ctx.GetString(KeyUserID)

		var err error
		_, err = c.taskLogic.Create(ctx.Request.Context(), userID, app.NewTask{
			Title:        title,
			Description:  description,
			IntervalSize: formattedIntervalSize,
			IntervalUnit: intervalUnit,
		})
		if err != nil {
			log.Printf("Failed to create task: %s\n", err)
			ctx.Status(http.StatusInternalServerError)
			return
		}

		ctx.Redirect(http.StatusFound, "/")
	}
}
