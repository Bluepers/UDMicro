package main

import (
	"context"
	"log"
	"logger-service/data"
	"time"
)

type RPCServer struct{}

type RPCPayload struct {
	Name string
	Data string
}

func (r *RPCServer) LogInfo(payload RPCPayload, resp *string) error {
	collection := client.Database("logs").Collection("logs")
	_, err := collection.InsertOne(context.TODO(), data.LogEntry{
		Name:      payload.Name,
		Data:      payload.Data,
		CreateAt:  time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		log.Println("error logging to MongoDB:", err)
		return err
	}

	*resp = "Process payload via RPC: " + payload.Name

	return nil
}
