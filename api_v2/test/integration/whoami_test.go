package integration

import (
	"testing"

	"techytechster.com/roastedoctocats/internal/config"
	"techytechster.com/roastedoctocats/pkg/proto"
)

func TestWhoAmI(t *testing.T) {
	NewIntegTest("WhoAmI", func(t *testing.T, client proto.OctoRoasterAPIClient) {
		jwt := config.New().ReadConfiguration("local").GetPAT()
		resp, err := client.WhoAmI(t.Context(), &proto.WhoAmIRequest{
			GithubToken: jwt,
		})
		if err != nil {
			t.Fatalf("expected no err: %s", err.Error())
		}
		if resp.Username != "Techypanda" {
			t.Fatalf("expected username to be Techypanda for test: %s", resp.Username)
		}
	})(t)
}
