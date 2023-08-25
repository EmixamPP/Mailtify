package runner

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"mailtify/config"
	"mailtify/router"
)

func Run(router *gin.Engine, c *config.Configuration) error {

	corsConfig := cors.Config{
		AllowOrigins: c.Server.AllowOrigins,
	}
	router.Use(cors.New(corsConfig))

	router.Use(multipartForm())

	addr := c.Server.ListenAddr + ":" + c.Server.Port
	return router.Run(addr)
}

func multipartForm() gin.HandlerFunc {
	return func(c *gin.Context) {
		contentType := c.GetHeader("Content-Type")
		if contentType == "multipart/form-data" {
			err := c.Request.ParseMultipartForm(16_000_000)
			if err != nil {
				res := router.Response{
					Status:  http.StatusBadRequest,
					Message: router.BAD_REQUEST_MESSAGE,
				}
				c.JSON(res.Status, res)
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
