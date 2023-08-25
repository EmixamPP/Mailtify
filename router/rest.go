package router

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

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

// Response represents the reponse returned to the client.
type Response struct {
	Status  int `json:"status"`
	Message any `json:"message"`
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
func messageHandler(c *gin.Context, d *database.GormDB, m *message.Messenger) Response {
	token := model.Token{Value: c.Query("token")}
	if validator.New().Struct(token) != nil {
		return Response{http.StatusBadRequest, BAD_REQUEST_MESSAGE}
	}

	exists, err := d.IsTokenExists(&token)
	if err != nil {
		return Response{http.StatusInternalServerError, err.Error()}
	} else if !exists {
		return Response{http.StatusUnauthorized, INVALID_TOKEN_MESSAGE}
	}

	title := c.Request.PostFormValue("title")
	if title == "" {
		return Response{http.StatusBadRequest, NO_TITLE_MESSAGE}
	}
	message_ := c.Request.PostFormValue("message")
	if title == "" {
		return Response{http.StatusBadRequest, NO_MESSAGE_MESSAGE}
	}
	recipients := c.Request.PostFormValue("recipients")
	if recipients == "" {
		return Response{http.StatusBadRequest, NO_RECIPIENTS_MESSAGE}
	}

	err = m.SendMessage(title, message_, strings.Split(recipients, ","))
	if err != nil {
		return Response{http.StatusInternalServerError, err.Error()}
	}

	var now = time.Now()
	token.LastUsed = &now
	d.DB.Save(token)

	return Response{http.StatusOK, ""}
}

// newHandler allows to create a token.
// It responses a new token.
func newHandler(c *gin.Context, d *database.GormDB, _ *message.Messenger) Response {
	token, err := d.NewToken()
	if err != nil {
		return Response{http.StatusInternalServerError, err.Error()}
	}
	return Response{http.StatusCreated, token.Value}
}

// deleteHandler allows to delete a token.
// It requieres the query parameter "token=",
// where the value is a valid token.
// It deletes the token which must the valid.
func deleteHandler(c *gin.Context, d *database.GormDB, _ *message.Messenger) Response {
	token := model.Token{Value: c.Query("token")}
	if validator.New().Struct(token) != nil {
		return Response{http.StatusBadRequest, BAD_REQUEST_MESSAGE}
	}

	exists, err := d.IsTokenExists(&token)
	if err != nil {
		return Response{http.StatusInternalServerError, err.Error()}
	} else if !exists {
		return Response{http.StatusUnauthorized, INVALID_TOKEN_MESSAGE}
	}

	err = d.DelToken(&token)
	if err != nil {
		return Response{http.StatusInternalServerError, err.Error()}
	}

	return Response{http.StatusOK, ""}
}

func tokensHandler(c *gin.Context, d *database.GormDB, _ *message.Messenger) Response {
	tokens, err := d.GetTokens()
	if err != nil {
		return Response{http.StatusInternalServerError, err.Error()}
	}

	return Response{http.StatusOK, tokens}
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
	router.GET("/new", route(newHandler))
	router.DELETE("/delete", route(deleteHandler))
	router.GET("/tokens", route(tokensHandler))

	return router
}
