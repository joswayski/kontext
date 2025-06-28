package main

import (
	"github.com/joswayski/kontext/apps/api/config"
	"github.com/joswayski/kontext/apps/api/router"
)

func main() {
	cfg := config.GetConfig()
	r := router.GetRouter()

	r.Run(":" + cfg.Port)
}
