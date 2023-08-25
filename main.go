package main

import (
	"mailtify/config"
	"mailtify/database"
	"mailtify/message"
	"mailtify/router"
	"mailtify/runner"
)

func main() {
	c := config.Get()

	d, err := database.New(c.Database.Dialect, c.Database.Connection, c.Security.TokenSize)
	if err != nil {
		panic(err)
	}
	defer d.Close()

	m := message.Create(c.SMTP.From, c.SMTP.Username, c.SMTP.Password, c.SMTP.Host, c.SMTP.Port)

	r := router.Create(d, m)
	
	err = runner.Run(r, c)
	if err != nil {
		panic(err)
	}
	
}
