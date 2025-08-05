package main

import (
	"fmt"
	"log/slog"
	"net"
	"os"
	"strings"

	"google.golang.org/grpc"
	"techytechster.com/roastedoctocats/internal/api"
	"techytechster.com/roastedoctocats/pkg/proto"
)

const defaultPort string = "8080"

func getPort() string {
	if port, found := os.LookupEnv("PORT"); found {
		return port
	}
	return defaultPort
}

func getLogLevel() slog.Level {
	if val, found := os.LookupEnv("ENABLE_DEBUG"); found && strings.ToUpper(val) == "TRUE" {
		return slog.LevelDebug
	}
	return slog.LevelInfo
}

const defaultAddress string = "localhost"

func getAddress() string {
	if address, found := os.LookupEnv("ADDRESS"); found {
		return address
	}
	return defaultAddress
}

func main() {
	logLevel := getLogLevel()
	logOpts := &slog.HandlerOptions{
		Level: logLevel,
	}
	port := getPort()
	address := getAddress()
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, logOpts)))
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", address, port))
	if err != nil {
		slog.Error("Failed to bind to tcp port", "port", port, "err", err)
		return
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	apiInstance, cancel := api.New()
	defer cancel()
	proto.RegisterOctoRoasterAPIServer(grpcServer, apiInstance)
	slog.Info("Server started", "port", port, "logLevel", logLevel, "address", address)
	err = grpcServer.Serve(lis)
	if err != nil {
		slog.Error("Failed to start server", "err", err)
	}
}
