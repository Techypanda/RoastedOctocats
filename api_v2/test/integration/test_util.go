package integration

import (
	"fmt"
	"log/slog"
	"os"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"techytechster.com/roastedoctocats/pkg/proto"
)

const defaultPort string = "8080"

// var ghPem []byte = []byte(config.New().ReadConfiguration("local").GetPAT())

// func generateGithubJWT() string {
// 	parsedPrivateKey, err := jwt.ParseRSAPrivateKeyFromPEM(ghPem)
// 	if err != nil {
// 		panic(err)
// 	}
// 	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
// 		"iat": time.Now().Unix() - 60,
// 		"exp": time.Now().Unix() + 180,
// 		"iss": "Iv23ligun1uyOZYdvxnq",
// 	})
// 	tokenString, err := token.SignedString(parsedPrivateKey)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return tokenString
// }

func getPort() string {
	if port, found := os.LookupEnv("PORT"); found {
		return port
	}
	return defaultPort
}

const defaultAddress string = "localhost"

func getAddress() string {
	if add, found := os.LookupEnv("ADDRESS"); found {
		return add
	}
	return defaultAddress
}

func BuildClient() (proto.OctoRoasterAPIClient, *grpc.ClientConn) {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%s", getAddress(), getPort()), opts...)
	if err != nil {
		slog.Error("failed to create client", "err", err)
		panic(err)
	}
	client := proto.NewOctoRoasterAPIClient(conn)
	return client, conn
}

// Utility Function To Give Each Test A Client That Is Cleaned Up Correctly
func NewIntegTest(Name string, Test func(t *testing.T, client proto.OctoRoasterAPIClient)) func(t *testing.T) {
	client, conn := BuildClient()
	return func(t *testing.T) {
		Test(t, client)
		err := conn.Close()
		if err != nil {
			t.Fatalf("failed to close client: %s", err.Error())
		}
	}
}
