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

	protectedRouter.POST("/task/:id/delete", handler.PostDelete())

	protectedRouter.GET("/task/:id/edit", handler.GetEditTask())
	protectedRouter.POST("/task/:id/edit", handler.PostEditTask())

	protectedRouter.POST("/task/:id/complete", handler.PostTaskComplete())

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
			"whole": func(number float64) int {
				return int(number * 100)
			},
		})
	}
}

type Member struct {
	ID   string
	Name string
}

func (c *TaskController) GetCreateTask() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := ctx.GetString(KeyUserID)
		user, err := c.userRepo.Get(ctx, userID)
		if err != nil {
			log.Printf("Failed to get information on user=%s: %s\n", userID, err)
			HTML(ctx, http.StatusInternalServerError, "pages/create-task", gin.H{
				"title": "Opret opgave",
				"error": "Kunne ikke hente brugerinformation. Prøv igen om lidt.",
			})
			return
		}
		if user.GroupID == nil {
			HTML(ctx, http.StatusInternalServerError, "pages/create-task", gin.H{
				"title": "Opret opgave",
				"error": "Bruger er ikke i en gruppe",
			})
			return
		}
		users, err := c.userRepo.GetByGroup(ctx, *user.GroupID)
		if err != nil {
			log.Printf("Failed to get information on user=%s: %s\n", userID, err)
			HTML(ctx, http.StatusInternalServerError, "pages/create-task", gin.H{
				"title": "Opret opgave",
				"error": "Kunne ikke hente information om medlemmer i gruppen. Prøv igen om lidt.",
			})
			return
		}

		var members []Member
		for _, member := range users {
			members = append(members, Member{
				ID:   member.ID,
				Name: member.Name,
			})
		}

		HTML(ctx, http.StatusOK, "pages/create-task", gin.H{
			"title":   "Opret opgave",
			"members": members,
		})
	}
}

func (c *TaskController) PostCreateTask() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		title := ctx.PostForm("title")
		description := ctx.PostForm("description")
		intervalSize := ctx.PostForm("intervalSize")
		intervalUnit := ctx.PostForm("intervalUnit")
		assignee := ctx.PostForm("assignee")
		rotatingAssignee := ctx.PostForm("rotatingAssignee")

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

		var assignedPerson *string
		if assignee != "" {
			assignedPerson = &assignee
		}

		formattedRotatingAssignee, _ := strconv.ParseBool(rotatingAssignee)

		userID := ctx.GetString(KeyUserID)

		var err error
		_, err = c.taskLogic.Create(ctx.Request.Context(), userID, app.NewTask{
			Title:            title,
			Description:      description,
			Assignee:         assignedPerson,
			RotatingAssignee: formattedRotatingAssignee,
			IntervalSize:     formattedIntervalSize,
			IntervalUnit:     intervalUnit,
		})
		if err != nil {
			log.Printf("Failed to create task: %s\n", err)
			ctx.Status(http.StatusInternalServerError)
			return
		}

		ctx.Redirect(http.StatusFound, "/")
	}
}

func (c *TaskController) PostDelete() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		taskID := ctx.Param("id")
		if taskID == "" {
			HTML(ctx, http.StatusBadRequest, "pages/index", gin.H{
				"title": "Taskeroo",
				"alert": "Der skete en fejl da opgaven skulle slettes. Prøv igen om lidt.",
			})
			return
		}

		userID := ctx.GetString(KeyUserID)
		err := c.taskLogic.Delete(ctx.Request.Context(), userID, taskID)
		if err != nil {
			log.Printf("Failed to delete task for user=%s: %s\n", userID, err)
		}

		ctx.Redirect(http.StatusFound, "/")
	}
}

func (c *TaskController) GetEditTask() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		taskID := ctx.Param("id")
		userID := ctx.GetString(KeyUserID)
		task, err := c.taskLogic.Get(ctx, userID, taskID)
		if err != nil {
			log.Printf("Failed to get task=%s for user=%s: %s\n", taskID, userID, err)
			ctx.Status(http.StatusInternalServerError)
			return
		}

		user, err := c.userRepo.Get(ctx, userID)
		if err != nil {
			log.Printf("Failed to get information on user=%s: %s\n", userID, err)
			HTML(ctx, http.StatusInternalServerError, "pages/create-task", gin.H{
				"title": "Opret opgave",
				"error": "Kunne ikke hente brugerinformation. Prøv igen om lidt.",
			})
			return
		}
		if user.GroupID == nil {
			HTML(ctx, http.StatusInternalServerError, "pages/create-task", gin.H{
				"title": "Opret opgave",
				"error": "Bruger er ikke i en gruppe",
			})
			return
		}
		users, err := c.userRepo.GetByGroup(ctx, *user.GroupID)
		if err != nil {
			log.Printf("Failed to get information on user=%s: %s\n", userID, err)
			HTML(ctx, http.StatusInternalServerError, "pages/create-task", gin.H{
				"title": "Opret opgave",
				"error": "Kunne ikke hente information om medlemmer i gruppen. Prøv igen om lidt.",
			})
			return
		}

		var members []Member
		for _, member := range users {
			members = append(members, Member{
				ID:   member.ID,
				Name: member.Name,
			})
		}

		HTML(ctx, http.StatusOK, "pages/edit-task", gin.H{
			"title":            "Opdatere opgave",
			"task":             task,
			"members":          members,
			"assignee":         task.Assignee,
			"rotatingAssignee": task.RotatingAssignee,
			"compare": func(a *string, b string) bool {
				if a == nil {
					return false
				}
				return *a == b
			},
		})
	}
}

func (c *TaskController) PostEditTask() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		taskID := ctx.Param("id")
		userID := ctx.GetString(KeyUserID)

		title := ctx.PostForm("title")
		description := ctx.PostForm("description")
		intervalSize := ctx.PostForm("intervalSize")
		intervalUnit := ctx.PostForm("intervalUnit")
		assignee := ctx.PostForm("assignee")
		rotatingAssignee := ctx.PostForm("rotatingAssignee")

		formattedIntervalSize, err := strconv.Atoi(intervalSize)
		if err != nil {
			formattedIntervalSize = 1
		}

		task := app.Task{
			ID:           taskID,
			Title:        title,
			Description:  description,
			IntervalSize: formattedIntervalSize,
			IntervalUnit: intervalUnit,
		}

		if title == "" {
			HTML(ctx, http.StatusBadRequest, "pages/edit-task", gin.H{
				"title": "Opdatere opgave",
				"error": "Titel skal udfyldes",
				"task":  task,
			})
			return
		}
		if description == "" {
			HTML(ctx, http.StatusBadRequest, "pages/edit-task", gin.H{
				"title": "Opdatere opgave",
				"error": "Description skal udfyldes",
				"task":  task,
			})
			return
		}
		if intervalUnit == "" || intervalSize == "" {
			HTML(ctx, http.StatusBadRequest, "pages/edit-task", gin.H{
				"title": "Opdatere opgave",
				"error": "Hyppighed skal udfyldes",
				"task":  task,
			})
			return
		}

		var assignedPerson *string
		if assignee != "" {
			assignedPerson = &assignee
		}

		formattedRotatingAssignee, _ := strconv.ParseBool(rotatingAssignee)

		err = c.taskLogic.Update(ctx, userID, taskID, app.NewTask{
			Title:            title,
			Description:      description,
			Assignee:         assignedPerson,
			RotatingAssignee: formattedRotatingAssignee,
			IntervalSize:     formattedIntervalSize,
			IntervalUnit:     intervalUnit,
		})
		if err != nil {
			log.Printf("Failed to get task=%s for user=%s: %s\n", taskID, userID, err)
			ctx.Status(http.StatusInternalServerError)
			return
		}

		ctx.Redirect(http.StatusFound, "/")
	}
}

func (c *TaskController) PostTaskComplete() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		taskID := ctx.Param("id")
		if taskID == "" {
			HTML(ctx, http.StatusBadRequest, "pages/index", gin.H{
				"title": "Taskeroo",
				"alert": "Der skete en fejl da opgaven skulle udføres. Prøv igen om lidt.",
			})
			return
		}

		userID := ctx.GetString(KeyUserID)
		err := c.taskLogic.Complete(ctx.Request.Context(), userID, taskID)
		if err != nil {
			log.Printf("Failed to complete task for user=%s: %s\n", userID, err)
		}

		ctx.Redirect(http.StatusFound, "/")
	}
}
