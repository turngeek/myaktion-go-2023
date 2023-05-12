package main

import (
	"fmt"
	"net"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/turngeek/myaktion-go-2023/src/banktransfer/grpc/banktransfer"
	"github.com/turngeek/myaktion-go-2023/src/banktransfer/service"
	"google.golang.org/grpc"
)

func init() {
	// init logger
	log.SetFormatter(&log.TextFormatter{})
	log.SetReportCaller(true)
	level, err := log.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		log.Info("Log level not specified, using default log level: INFO")
		log.SetLevel(log.InfoLevel)
		return
	}
	log.SetLevel(level)
}

var grpcPort = 9111

func main() {
	log.Info("Starting Banktransfer service")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen on grpc port %d: %v", grpcPort, err)
	}
	grpcServer := grpc.NewServer()
	banktransfer.RegisterBankTransferServer(grpcServer, service.NewBankTransferService())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve gRPC server over port %d: %v", grpcPort, err)
	}
}
