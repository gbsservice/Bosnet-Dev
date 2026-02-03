package routes

import (
	dashboard "api_kino/app/controllers/dashboard"
	"api_kino/app/controllers/version"
	"api_kino/app/jobs"
	"api_kino/app/middleware"
	"api_kino/config/app"
	"api_kino/service/ws"
	"context"
	"time"

	api "github.com/appleboy/gin-status-api"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

var (
	c = context.Background()
)

func Router() *gin.Engine {
	router := gin.New()
	if app.Config().EnableJob {
		jobs.HandleJobs()
	}

	hub := ws.Hub
	go hub.Run()

	router.Use(gzip.Gzip(gzip.DefaultCompression))

	router.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "PUT", "POST", "PATCH", "OPTIONS", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "User-Agent", "Referrer", "Host", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowAllOrigins:  false,
		AllowOriginFunc:  func(origin string) bool { return true },
		MaxAge:           12 * time.Hour,
	}))
	router.Use(middleware.Api())

	router.GET("/ws/:auth", func(c *gin.Context) {
		ws.ServeWs(c)
	})

	v1 := router.Group("v1")
	v1.GET("/time", version.GetTime)
	v1.GET("/version", version.GetApi)
	v1.POST("/post-draftso", dashboard.PostDraftSO)
	v1.POST("/post-draftso-manual", dashboard.PostDraftSOManual)
	v1.GET("/parameter-check", dashboard.ParamCheck)

	authorize := v1.Group("/")
	authorize.Use(middleware.Auth())

	authorizeKey := v1.Group("/")
	authorizeKey.Use(middleware.AuthKey())

	router.GET("/api/status", api.GinHandler)
	v1.Static("/assets", "./assets")

	pprof.Register(router)
	return router
}
