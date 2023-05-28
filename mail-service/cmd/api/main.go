package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

const webPort = 80

type Config struct {
	Mailer Mail
}

func main() {
	app := Config{
		Mailer: createMail(),
	}

	log.Println("Starting mail service on port " + strconv.Itoa(webPort))

	srv := &http.Server{
		Addr:    fmt.Sprint(":" + strconv.Itoa(webPort)),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()

	if errors.Is(err, http.ErrServerClosed) {
		log.Println("Server closed")
	} else if err != nil {
		log.Panicln("Error: server -", err)
	}
}

func createMail() Mail {
	port, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))
	m := Mail{
		Domain:      os.Getenv("MAIL_DOMAIN"),
		Host:        os.Getenv("MAIL_HOST"),
		Port:        port,
		Username:    os.Getenv("MAIL_USERNAME"),
		Password:    os.Getenv("MAIL_Password"),
		Encryption:  os.Getenv("MAIL_ENCRYPTION"),
		FromName:    os.Getenv("FROM_NAME"),
		FromAddress: os.Getenv("FROM_ADDRESS"),
	}

	return m
}
