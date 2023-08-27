package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"mailtify/database"
	"mailtify/message"
)


// Response represents the reponse returned to the client.
type Response struct {
	Status  int `json:"status"`
	Message any `json:"message"`
}

func Create(d *database.GormDB, m *message.Messenger) *gin.Engine {
	route := func(handler func(*gin.Context, *database.GormDB, *message.Messenger) Response) func(*gin.Context) {
		return func(c *gin.Context) {
			res := handler(c, d, m)
			c.JSON(res.Status, res)
		}
	}

	router := gin.Default()
	router.POST("/msg", route(messageHandler))
	router.POST("/message", route(messageHandler))
	router.GET("/new", authenticate(d), route(newHandler))
	router.DELETE("/delete", route(deleteHandler))
	router.GET("/tokens", route(tokensHandler))

	router.Use(parseMultipartForm())

	return router
}

func authenticate(d *database.GormDB) gin.HandlerFunc {
	return func(c *gin.Context) {
		/*user, password, ok := c.Request.BasicAuth()
		
		if !ok || validCredentials[user] != password {
			c.Header("WWW-Authenticate", "Basic realm=Authorization Required")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}*/

		c.Next()
	}
}

func parseMultipartForm() gin.HandlerFunc {
	return func(c *gin.Context) {
		contentType := c.GetHeader("Content-Type")
		if contentType == "multipart/form-data" {
			err := c.Request.ParseMultipartForm(16_000_000)
			if err != nil {
				res := Response{Status: http.StatusBadRequest, Message: BAD_REQUEST_MESSAGE}
				c.JSON(res.Status, res)
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
