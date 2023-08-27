package runner

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"mailtify/configuration"
)

func Run(router *gin.Engine, c *configuration.Configuration) error {
	corsConfig := cors.Config{
		AllowOrigins: c.Server.AllowOrigins,
	}
	router.Use(cors.New(corsConfig))

	addr := c.Server.ListenAddr + ":" + c.Server.Port
	return router.Run(addr)
}
