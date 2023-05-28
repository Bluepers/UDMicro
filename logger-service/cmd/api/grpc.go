package main

import (
	"context"
	"fmt"
	"log"
	"logger-service/data"
	"logger-service/logs"
	"net"
	"strconv"

	"google.golang.org/grpc"
)

type LogServer struct {
	logs.UnimplementedLogServiceServer
	Models data.Models
}

func (l *LogServer) WriteLog(ctx context.Context, req *logs.LogRequest) (*logs.LogResponse, error) {
	input := req.GetLogEntry()

	// Writel the log
	logEntry := data.LogEntry{
		Name: input.Name,
		Data: input.Data,
	}

	err := l.Models.LogEntry.Insert(logEntry)
	if err != nil {
		res := &logs.LogResponse{
			Result: "failed",
		}
		return res, err
	}

	// Return response
	res := &logs.LogResponse{
		Result: "logged",
	}

	return res, nil
}

func (app *Config) gRPCListen() {
	lis, err := net.Listen("tcp", fmt.Sprint(":"+strconv.Itoa(gRpcPort)))
	if err != nil {
		log.Fatalln("Failed to listen to gRPC:", err)
	}

	srv := grpc.NewServer()
	logs.RegisterLogServiceServer(srv, &LogServer{Models: app.Models})

	log.Println("gRPC Server started on port " + strconv.Itoa(gRpcPort))
	if err := srv.Serve(lis); err != nil {
		log.Fatalln("Failed to listen to gRPC:", err)
	}
}
