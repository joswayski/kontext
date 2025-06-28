package router

import (
	"github.com/gin-gonic/gin"
	"github.com/joswayski/kontext/apps/api/handlers"
)

func GetRouter() *gin.Engine {
	r := gin.Default()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("", handlers.RootHandler)

	return r
}
