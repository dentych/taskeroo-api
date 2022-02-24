package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type PWAController struct {
	router gin.IRouter
}

func NewPWAController(router gin.IRouter) *PWAController {
	handler := &PWAController{router: router}

	router.GET("/manifest.json", handler.GetManifest())
	router.GET("/sw.js", handler.GetServiceWorkerJs())

	return handler
}

func (c *PWAController) GetManifest() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, map[string]interface{}{
			"background_color": "#ffffff",
			"display":          "fullscreen",
			"name":             "Taskeroo",
			"description":      "A small task app, which can help you keep track of recurring tasks that needs to be done in your household.",
			"start_url":        "/",
			"icons": []map[string]string{
				{
					"src":   "https://i.imgur.com/Ch1BU7E.png",
					"type":  "image/png",
					"sizes": "512x512",
				},
				{
					"src":   "https://i.imgur.com/RKeLCfC.png",
					"type":  "image/png",
					"sizes": "192x192",
				},
			},
		})
	}
}

const swjs = `
self.addEventListener('install', function(event) {
  event.waitUntil(
    caches.open('sw-cache').then(function(cache) {
      return cache.add('index');
    })
  );
});
 
self.addEventListener('fetch', function(event) {
  event.respondWith(
    caches.match(event.request).then(function(response) {
      return response || fetch(event.request);
    })
  );
});
`

func (c *PWAController) GetServiceWorkerJs() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("Content-Type", "text/javascript")
		ctx.String(http.StatusOK, swjs)
	}
}
