package router

import (
	_ "bias/docs"
	"bias/http/handlers"
	"bias/store"
	"context"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"github.com/swaggo/echo-swagger"
	"gorm.io/gorm"
)

func InitializeStreamRoutes(g *echo.Group, redisClient *redis.Client) {
	streamStore := store.NewStreamStore(nil, redisClient, context.Background())
	streamHandler := handlers.NewStreamHandler(*streamStore)

	g.GET("/streams", streamHandler.ListStreams)
	g.POST("/streams", streamHandler.CreateStream)
}

func InitializeTagRoutes(g *echo.Group, db *gorm.DB) {
	tagStore := store.NewTagStore(db, nil, context.Background())
	tagsHandler := handlers.NewTagHandler(tagStore)

	g.GET("/tags", tagsHandler.ListTags)
	g.GET("/tags/:id", tagsHandler.GetTagByID)
	g.POST("/tags", tagsHandler.CreateTag)
}

func SetupSwaggerUI(g *echo.Group) {
	g.GET("/docs/*", echoSwagger.WrapHandler)
}

func SetupRoutes(e *echo.Echo, redisClient *redis.Client, db *gorm.DB) {
	apiGroup := e.Group("/api")

	InitializeStreamRoutes(apiGroup, redisClient)
	InitializeTagRoutes(apiGroup, db)
	SetupSwaggerUI(apiGroup)
}
