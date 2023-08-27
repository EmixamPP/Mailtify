package main

import (
	"mailtify/configuration"
	"mailtify/database"
	"mailtify/message"
	"mailtify/api"
	"mailtify/runner"
)

func main() {
	config := configuration.Get()

	db, err := database.New(config.Database.Dialect, config.Database.Connection, config.Security.TokenSize)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	messenger := message.Create(config.SMTP.From, config.SMTP.Username, config.SMTP.Password, config.SMTP.Host, config.SMTP.Port)

	router := api.Create(db, messenger)
	
	err = runner.Run(router, config)
	if err != nil {
		panic(err)
	}
	
}
