package main

import (
	"fmt"
	"github.com/foolin/goview"
	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
)

func main() {
	router := gin.Default()

	goviewConfig := goview.DefaultConfig
	if os.Getenv("ENVIRONMENT") != "prod" {
		goviewConfig.DisableCache = true
	}
	router.HTMLRender = ginview.New(goviewConfig)

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

	router.GET("/page", func(ctx *gin.Context) {
		//render only file, must full name with extension
		ctx.HTML(http.StatusOK, "page.html", gin.H{"title": "Page file title!!"})
	})

	port := "8080"
	portEnv := os.Getenv("PORT")
	if portEnv != "" {
		port = portEnv
	}
	err := router.Run(fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("Error running server: %s\n", err)
	}
}
