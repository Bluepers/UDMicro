package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"logger-service/data"
	"net"
	"net/http"
	"net/rpc"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	webPort  = 80
	rpcPort  = 5001
	mongoURL = "mongodb://mongo:27017"
	gRpcPort = 50001
)

var client *mongo.Client

type Config struct {
	Models data.Models
}

func main() {
	// Connect to mongodb
	mongoClient, err := connectToMongoDB()
	if err != nil {
		log.Panic(err)
	}

	client = mongoClient

	// Create a context in order to disconnect
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)

	defer cancel()

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Panic(err)
		}
	}()

	app := Config{
		Models: data.New(client),
	}

	// Register RPC Server
	err = rpc.Register(new(RPCServer))
	if err != nil {
		log.Panic(err)
	}
	go app.rpcListen()

	go app.gRPCListen()

	srv := &http.Server{
		Addr:    fmt.Sprint(":" + strconv.Itoa(webPort)),
		Handler: app.routes(),
	}

	fmt.Println("Starting logger service on port " + strconv.Itoa(webPort))
	err = srv.ListenAndServe()

	if errors.Is(err, http.ErrServerClosed) {
		log.Println("Server closed")
	} else if err != nil {
		log.Panicln("Error: server -", err)
	}
}

func connectToMongoDB() (*mongo.Client, error) {
	// Create the connection options
	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	// Connect to MongoDB
	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Error connecting:", err)
		return nil, err
	}

	log.Println("Connected to MongoDB")

	return c, nil
}

func (app *Config) rpcListen() error {
	log.Println("Starting rpc server on port " + strconv.Itoa(rpcPort))

	listen, err := net.Listen("tcp", fmt.Sprint("0.0.0.0:"+strconv.Itoa(rpcPort)))
	if err != nil {
		return err
	}

	defer listen.Close()

	for {
		rpcConn, err := listen.Accept()
		if err != nil {
			continue
		}

		go rpc.ServeConn(rpcConn)
	}
}
