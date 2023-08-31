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
func newHandler(c *gin.Context, d *database.GormDB) Response {
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

// deleteTokenHandler allows to delete a token.
// It requieres a valid token in the context.
func deleteTokenHandler(c *gin.Context, d *database.GormDB) Response {
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

// tokensHandler responses all the tokens of an user in the database.
// It requieres a valid user in the context.
func tokensHandler(c *gin.Context, d *database.GormDB) Response {
	userInterface, ok := c.Get("user")
	if !ok {
		panic("user missing from the gin context")
	}
	user := userInterface.(*model.User)

	tokens, err := d.GetUserTokens(user)
	if err != nil {
		return Response{Status: http.StatusInternalServerError, Message: err.Error()}
	}

	return Response{http.StatusOK, tokens}
}

// usersHandler responses all the tokens in the database.
func usersHandler(c *gin.Context, d *database.GormDB) Response {
	users, err := d.GetUsers()
	if err != nil {
		return Response{Status: http.StatusInternalServerError, Message: err.Error()}
	}

	return Response{http.StatusOK, users}
}

// deleteUserHandler allows to delete a token.
// It requieres a valid token in the context.
func deleteUserHandler(c *gin.Context, d *database.GormDB) Response {
	userInterface, ok := c.Get("user")
	if !ok {
		panic("user missing from the gin context")
	}
	user := userInterface.(*model.User)

	err := d.DelUser(user)
	if err != nil {
		return Response{Status: http.StatusInternalServerError, Message: err.Error()}
	}

	return Response{http.StatusOK, OK_MESSAGE}
}

// createUserHandler allows to create a token.
// It requieres a valid user in the context, and responses a new token.
func createUserHandler(c *gin.Context, d *database.GormDB) Response {
	username := c.Query("username")
	password := c.Query("password")
	admin := c.Query("admin")
	if username == "" || password == "" || admin == "" {
		return Response{Status: http.StatusBadRequest, Message: BAD_REQUEST_MESSAGE}
	}

	err := d.CreateUser(username, password, admin == "1" || admin == "true")
	if err != nil {
		return Response{Status: http.StatusInternalServerError, Message: err.Error()}
	}

	return Response{Status: http.StatusCreated, Message: OK_MESSAGE}
}
