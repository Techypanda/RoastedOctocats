package integration

import (
	"testing"

	"github.com/google/uuid"
	"techytechster.com/roastedoctocats/pkg/api"
	"techytechster.com/roastedoctocats/pkg/proto"
)

func TestPingEndpoint(t *testing.T) {
	NewIntegTest("TestPingEndpoint", func(t *testing.T, client proto.OctoRoasterAPIClient) {
		ping, err := client.Ping(t.Context(), &proto.PingRequest{IdempotencyToken: uuid.Must(uuid.NewV7()).String()})
		if err != nil {
			t.Fatalf("expected no err: %s", err.Error())
		}
		if ping.Message != api.PingResponseMessage {
			t.Fatalf("api response is not correct: %s", api.PingResponseMessage)
		}
	})(t)
}
