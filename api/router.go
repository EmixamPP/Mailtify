package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"mailtify/database"
	"mailtify/message"
	"mailtify/model"
)

const UNAUTHORIZED_MESSAGE = "unauthorized access"
const INVALID_TOKEN_MESSAGE = "invalid token"

// Response represents the reponse returned to the client.
type Response struct {
	Status  int `json:"status"`
	Message any `json:"message"`
}

func Create(d *database.GormDB, m *message.Messenger) *gin.Engine {
	routeMessenger := func(handler func(*gin.Context, *database.GormDB, *message.Messenger) Response) func(*gin.Context) {
		return func(c *gin.Context) {
			res := handler(c, d, m)
			c.JSON(res.Status, res)
		}
	}

	route := func(handler func(*gin.Context, *database.GormDB) Response) func(*gin.Context) {
		return func(c *gin.Context) {
			res := handler(c, d)
			c.JSON(res.Status, res)
		}
	}

	router := gin.Default()
	router.Use(parseMultipartForm())

	router.Match([]string{"POST", "PUT"}, "/msg", authenticateToken(d),
		routeMessenger(messageHandler))
	router.Match([]string{"POST", "PUT"}, "/message", authenticateToken(d),
		routeMessenger(messageHandler))

	router.GET("/new", authenticateUser(d), route(newHandler))

	router.DELETE("/delete", authenticateUser(d), authenticateToken(d),
		tokenOwnerOrAdmin(), route(deleteHandler))

	router.GET("/tokens", authenticateUser(d), admin(), route(tokensHandler))

	return router
}

// authenticate stores the token in the context for unified access from the API.
// Abord if the token is invalid, i.e. is not in the database,
// or if the token is missing.
func authenticateToken(d *database.GormDB) gin.HandlerFunc {
	return func(c *gin.Context) {
		value := c.Query("token")

		if value == "" {
			res := Response{Status: http.StatusBadRequest, Message: BAD_REQUEST_MESSAGE}
			c.JSON(res.Status, res)
			c.Abort()
			return
		}

		token, err := d.GetToken(value)
		if err != nil {
			res := Response{Status: http.StatusInternalServerError, Message: err.Error()}
			c.JSON(res.Status, res)
			c.Abort()
			return
		} else if token == nil {
			res := Response{Status: http.StatusUnauthorized, Message: INVALID_TOKEN_MESSAGE}
			c.JSON(res.Status, res)
			c.Abort()
			return
		}

		c.Set("token", token)
		c.Next()
	}
}

// authenticate stores the user in the context for unified access from the API.
// Abord if the user is unauthorized, i.e. is not in the database,
// or if the authentication is missing.
func authenticateUser(d *database.GormDB) gin.HandlerFunc {
	return func(c *gin.Context) {
		username, password, ok := c.Request.BasicAuth()

		if !ok {
			c.Header("WWW-Authenticate", "Basic realm=Authorization Required")
			res := Response{Status: http.StatusUnauthorized, Message: UNAUTHORIZED_MESSAGE}
			c.JSON(res.Status, res)
			c.Abort()
			return
		}

		user, err := d.GetUser(username, password)
		if err != nil {
			res := Response{Status: http.StatusInternalServerError, Message: err.Error()}
			c.JSON(res.Status, res)
			c.Abort()
			return
		} else if user == nil {
			res := Response{Status: http.StatusUnauthorized, Message: UNAUTHORIZED_MESSAGE}
			c.JSON(res.Status, res)
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}

// parseMultipartForm stores each fields of the mutlipart form
// in the context for unified access from the API.
// Does not abord if the request does not contains multipart form.
func parseMultipartForm() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.HasPrefix(c.GetHeader("Content-Type"), "multipart/form-data") {
			err := c.Request.ParseMultipartForm(16_000_000)
			if err != nil {
				res := Response{Status: http.StatusBadRequest, Message: BAD_REQUEST_MESSAGE}
				c.JSON(res.Status, res)
				c.Abort()
				return
			}

			for field, values := range c.Request.MultipartForm.Value {
				if len(values) > 0 {
					c.Set(field, values[0])
				}
			}
		}
		c.Next()
	}
}

// admin abord if the user stored in the conext is not an admin.
func admin() gin.HandlerFunc {
	return func(c *gin.Context) {
		userInterface, ok := c.Get("user")
		if !ok {
			panic("user missing from the gin context")
		}
		user := userInterface.(*model.User)

		if !user.Admin {
			res := Response{Status: http.StatusUnauthorized, Message: UNAUTHORIZED_MESSAGE}
			c.JSON(res.Status, res)
			c.Abort()
			return
		}

		c.Next()
	}
}

// tokenOwnerOrAdmin abord if the token stored in the conext has not been
// created by the user in the context, except if he is an admin.
func tokenOwnerOrAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		userInterface, ok := c.Get("user")
		if !ok {
			panic("user missing from the gin context")
		}
		user := userInterface.(*model.User)

		tokenInterface, ok := c.Get("token")
		if !ok {
			panic("token missing from the gin context")
		}
		token := tokenInterface.(*model.Token)

		if token.UserID != user.ID && !user.Admin {
			res := Response{Status: http.StatusUnauthorized, Message: UNAUTHORIZED_MESSAGE}
			c.JSON(res.Status, res)
			c.Abort()
			return
		}

		c.Next()
	}
}
