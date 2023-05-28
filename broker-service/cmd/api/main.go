package main

import (
	"errors"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Config struct {
	Rabbit *amqp.Connection
}

const webPort = 8080

func main() {
	// Try to connect to RabbitMQ
	rabbitConn, err := connectRabbitMQ()

	if err != nil {
		log.Fatalln(err)
	}

	defer rabbitConn.Close()

	app := Config{Rabbit: rabbitConn}

	log.Println("Starting broker service on port " + strconv.Itoa(webPort))

	srv := &http.Server{
		Addr:    fmt.Sprint(":" + strconv.Itoa(webPort)),
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()

	if errors.Is(err, http.ErrServerClosed) {
		log.Println("Server closed")
	} else if err != nil {
		log.Panicln("Error: server -", err)
	}
}

func connectRabbitMQ() (*amqp.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	// loop until connected
	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq")

		if err != nil {
			fmt.Println("RabbitMQ not ready to connect...")
			counts++
		} else {
			connection = c
			break
		}

		if counts > 5 {
			fmt.Println(err)
			return nil, err
		}

		backOff = time.Duration(math.Pow(float64(counts), 2) * float64(time.Second))
		log.Println("backing off for " + strconv.Itoa(int(backOff)) + " seconds")
		time.Sleep(backOff)
		continue
	}
	log.Println("Connected to RabbitMQ")
	return connection, nil
}
