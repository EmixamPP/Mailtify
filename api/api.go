package api

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"mailtify/database"
	"mailtify/message"
	"mailtify/model"
)

const BAD_REQUEST_MESSAGE = "bad request"
const NO_TITLE_MESSAGE = "no title specified"
const NO_MESSAGE_MESSAGE = "no title specified"
const NO_RECIPIENTS_MESSAGE = "no recipients specified"
const OK_MESSAGE = ""

// messageHandler allows to send a mail.
// It requieres a valid token in the context.
// It requieres a title string in the context.
// It requieres a message string in the context.
// It requieres a recipients string in the context.
func messageHandler(c *gin.Context, d *database.GormDB, m *message.Messenger) Response {
	tokenInterface, ok := c.Get("token")
	if !ok {
		panic("token missing from the gin context")
	}
	token := tokenInterface.(*model.Token)

	titleInterface, ok := c.Get("title")
	if !ok {
		return Response{Status: http.StatusBadRequest, Message: NO_TITLE_MESSAGE}
	}
	messageInterface, ok := c.Get("message")
	if !ok {
		return Response{Status: http.StatusBadRequest, Message: NO_MESSAGE_MESSAGE}
	}
	recipientsInterface, ok := c.Get("recipients")
	if !ok {
		return Response{Status: http.StatusBadRequest, Message: NO_RECIPIENTS_MESSAGE}
	}
	title, message_, recipients := titleInterface.(string), messageInterface.(string), recipientsInterface.(string)

	err := m.SendMessage(title, message_, strings.Split(recipients, ","))
	if err != nil {
		return Response{Status: http.StatusInternalServerError, Message: err.Error()}
	}

	var now = time.Now()
	token.LastUse = &now
	d.UpdateToken(token)

	return Response{Status: http.StatusOK, Message: OK_MESSAGE}
}

// newHandler allows to create a token.
// It requieres a valid user in the context, and responses a new token.
func newHandler(c *gin.Context, d *database.GormDB, _ *message.Messenger) Response {
	userInterface, ok := c.Get("user")
	if !ok {
		panic("user missing from the gin context")
	}
	user := userInterface.(*model.User)

	token, err := d.NewToken(user)
	if err != nil {
		return Response{Status: http.StatusInternalServerError, Message: err.Error()}
	}

	return Response{Status: http.StatusCreated, Message: token.Value}
}

// deleteHandler allows to delete a token.
// It requieres a valid token in the context.
func deleteHandler(c *gin.Context, d *database.GormDB, _ *message.Messenger) Response {
	tokenInterface, ok := c.Get("token")
	if !ok {
		panic("token missing from the gin context")
	}
	token := tokenInterface.(*model.Token)

	err := d.DelToken(token)
	if err != nil {
		return Response{Status: http.StatusInternalServerError, Message: err.Error()}
	}

	return Response{Status: http.StatusOK, Message: ""}
}

// tokensHandler responses all the tokens in the database.
func tokensHandler(c *gin.Context, d *database.GormDB, _ *message.Messenger) Response {
	tokens, err := d.GetTokens()
	if err != nil {
		return Response{Status: http.StatusInternalServerError, Message: err.Error()}
	}

	return Response{http.StatusOK, tokens}
}
