package router

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gopkg.in/validator.v2"

	"mailtify/database"
	"mailtify/message"
	"mailtify/model"
)

const BAD_REQUEST_MESSAGE = "bad request"
const INVALID_TOKEN_MESSAGE = "invalid token"
const NO_TITLE_MESSAGE = "no title specified"
const NO_MESSAGE_MESSAGE = "no title specified"
const NO_RECIPIENTS_MESSAGE = "no recipients specified"
const OK_MESSAGE = ""

// response represents the reponse returned to the client.
type response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// messageHandler allows to send a mail.
// It requieres the query parameter "token=",
// where the value is a valid token.
// It requieres the form "title=" in the payload,
// where the value is the object of the mail.
// It requieres the form "message=" in the payload,
// where the value is the body of the mail.
// It requieres the form "recipents=" in the payload,
// where the value is each recipent separated by a comma.
// It sends the message to the recipients.
// The token must be valid.
func messageHandler(c *gin.Context, d *database.GormDB, m *message.Messenger) response {
	token := model.Token{Value: c.Query("token")}
	if (validator.Validate(token) != nil) {
		return response{http.StatusBadRequest, BAD_REQUEST_MESSAGE}
	}

	exists, err := d.IsTokenExists(&token)
	if err != nil {
		return response{http.StatusInternalServerError, err.Error()}
	} else if !exists {
		return response{http.StatusUnauthorized, INVALID_TOKEN_MESSAGE}
	}

	err = c.Request.ParseMultipartForm(16_000_000)
	if err != nil {
		return response{http.StatusBadRequest, BAD_REQUEST_MESSAGE}
	}

	title := c.Request.PostFormValue("title")
	if title == "" {
		return response{http.StatusBadRequest, NO_TITLE_MESSAGE}
	}
	message_ := c.Request.PostFormValue("message")
	if title == "" {
		return response{http.StatusBadRequest, NO_MESSAGE_MESSAGE}
	}
	recipients := c.Request.PostFormValue("recipients")
	if recipients == "" {
		return response{http.StatusBadRequest, NO_RECIPIENTS_MESSAGE}
	}

	err = m.SendMessage(title, message_, strings.Split(recipients, ","))
	if err != nil {
		return response{http.StatusInternalServerError, err.Error()}
	}

	var now = time.Now()
	token.LastUsed = &now
	d.DB.Save(token)

	return response{http.StatusOK, ""}
}

// newHandler allows to create a token.
// It responses a new token.
func newHandler(c *gin.Context, d *database.GormDB, _ *message.Messenger) response {
	token, err := d.NewToken()
	if err != nil {
		return response{http.StatusInternalServerError, err.Error()}
	}
	return response{http.StatusCreated, token.Value}
}

// deleteHandler allows to delete a token.
// It requieres the query parameter "token=",
// where the value is a valid token.
// It deletes the token which must the valid.
func deleteHandler(c *gin.Context, d *database.GormDB, _ *message.Messenger) response {
	token := model.Token{Value: c.Query("token")}
	if (validator.Validate(token) != nil) {
		return response{http.StatusBadRequest, BAD_REQUEST_MESSAGE}
	}

	exists, err := d.IsTokenExists(&token)
	if err != nil {
		return response{http.StatusInternalServerError, err.Error()}
	} else if !exists {
		return response{http.StatusUnauthorized, INVALID_TOKEN_MESSAGE}
	}

	err = d.DelToken(&token)
	if err != nil {
		return response{http.StatusInternalServerError, err.Error()}
	}

	return response{http.StatusOK, ""}
}

func Create(listenAddr, port string, allowOrigins []string, d *database.GormDB, m *message.Messenger) error {
	route := func(handler func(*gin.Context, *database.GormDB, *message.Messenger) response) func(*gin.Context) {
		return func(c *gin.Context) {
			res := handler(c, d, m)
			c.JSON(res.Status, res)
		}
	}

	router := gin.Default()
	router.POST("/msg", route(messageHandler))
	router.POST("/message", route(messageHandler))
	router.GET("/new", route(newHandler))
	router.DELETE("/del", route(deleteHandler))
	router.DELETE("/delete", route(deleteHandler))

	corsConfig := cors.Config{AllowOrigins: allowOrigins}
	router.Use(cors.New(corsConfig))

	return router.Run(listenAddr + ":" + port)
}
