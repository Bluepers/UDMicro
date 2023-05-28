package main

import (
	"fmt"
	"listener/event"
	"log"
	"math"
	"strconv"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	log.Println("Starting listener service on port")

	// Try to connect to RabbitMQ
	rabbitConn, err := connect()

	if err != nil {
		log.Fatalln(err)
	}

	defer rabbitConn.Close()

	// Listen to messages
	log.Println("Listen and Comsuming RabbitMQ messages")

	// Create a Consumer
	consumer, err := event.NewConsumer(rabbitConn)
	if err != nil {
		log.Fatalln(err)
	}

	// Watch the Queue and comsume event
	err = consumer.Listen([]string{"log.INFO", "log.WARNING", "log.ERROR"})
	if err != nil {
		log.Fatalln(err)
	}
}

func connect() (*amqp.Connection, error) {
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
