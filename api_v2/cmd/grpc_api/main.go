package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"strings"

	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
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

var rateLimitedUsersMap map[string]*rate.Limiter = map[string]*rate.Limiter{}

const rateLimit = 4
const burstRateLimit = 8

func rateLimiterInterceptor(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	if p, ok := peer.FromContext(ctx); ok {
		ipAddress := strings.Split(p.Addr.String(), ":")[0]
		slog.Debug("Checking rate limit for user...", "ip", ipAddress)
		if _, exists := rateLimitedUsersMap[ipAddress]; !exists {
			rateLimitedUsersMap[ipAddress] = rate.NewLimiter(rateLimit, burstRateLimit)
		}
		if !rateLimitedUsersMap[ipAddress].Allow() {
			return nil, status.Errorf(codes.ResourceExhausted, "slowdown... you are being rate limited")
		}
		return handler(ctx, req)
	} else {
		return nil, status.Errorf(codes.Aborted, "please retry, we failed to read your peer information")
	}
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
	opts = append(opts, grpc.UnaryInterceptor(rateLimiterInterceptor))
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
