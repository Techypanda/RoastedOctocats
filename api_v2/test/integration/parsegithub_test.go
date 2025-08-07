package integration

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"techytechster.com/roastedoctocats/internal/config"
	"techytechster.com/roastedoctocats/pkg/proto"
)

func TestParseGithub(t *testing.T) {
	NewIntegTest("ParseGithub", func(t *testing.T, client proto.OctoRoasterAPIClient) {
		jwt := config.New().ReadConfiguration("local").GetPAT()
		token := uuid.Must(uuid.NewV7()).String()
		promptTypes := []proto.ModelPromptType{
			proto.ModelPromptType_ModelPromptType_EARLY2000s,
			proto.ModelPromptType_ModelPromptType_UWUIFIED,
		}
		for _, promptType := range promptTypes {
			_, err := client.ParseGithub(t.Context(), &proto.ParseGithubRequest{
				GithubToken:      jwt,
				IdempotencyToken: token,
				PromptType:       promptType,
			})
			if err != nil {
				t.Fatalf("expected no err: %s", err.Error())
			}
			for {
				resp, err := client.GetParsedGithubResult(t.Context(), &proto.GetParsedGithubResultRequest{
					GithubToken:      jwt,
					IdempotencyToken: token,
				})
				if err != nil {
					t.Fatalf("expected no err: %s", err.Error())
				}
				if resp.Status == "inProgress" {
					fmt.Println("still getting ai response... sleeping and retrying")
					time.Sleep(2 * time.Second) // Pause for 2 seconds
				} else {
					if resp.Status != "complete" {
						t.Fatalf("expected complete status: %s", resp.Status)
					}
					if len(*resp.Result) == 0 {
						t.Fatal("response is 0 characters")
					}
					break
				}
			}
		}
	})(t)
}
