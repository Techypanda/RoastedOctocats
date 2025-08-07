package api

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/azure"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"techytechster.com/roastedoctocats/internal/config"
	"techytechster.com/roastedoctocats/internal/github"
	"techytechster.com/roastedoctocats/pkg/idb"
	"techytechster.com/roastedoctocats/pkg/igithub"
	"techytechster.com/roastedoctocats/pkg/proto"
)

type CommitDetails struct {
	RepoName       string   `json:"repositoryName"`
	CommitMessages []string `json:"commitMessages"`
}
type GenAIJsonInput struct {
	Username   string          `json:"username"`
	Bio        string          `json:"bio"`
	CommitDets []CommitDetails `json:"commitDetails"`
}

var errInternal = errors.New("an internal server error occured, please retry")

const maxPagesOfIteration = 10 // at most 10 pages for now
func callRelevantGithubAPIS(j octocatGrpcParseGithubJobWorkerParam) (*string, error) {
	client := github.New()
	page := 1
	pageSize := 100
	var commitDetails []CommitDetails
	for {
		resp, err := client.ListEvents(igithub.GithubListEventsRequest{
			Username:   j.Username,
			OAuthToken: j.githubOAuthToken,
			Page:       &page,
			PageSize:   &pageSize,
		})
		if err != nil {
			slog.Error("failed to do listevents", "err", err)
			return nil, errInternal
		}
		if resp != nil {
			for _, item := range *resp {
				curr := CommitDetails{
					RepoName:       item.Repo.Name,
					CommitMessages: []string{},
				}
				if item.Type == igithub.PushEventType {
					for _, commit := range item.Payload.Commits {
						curr.CommitMessages = append(curr.CommitMessages, commit.Message)
					}
					commitDetails = append(commitDetails, curr)
				} else {
					slog.Debug("filtering out event from list event.", "item", item)
				}
			}
			if len(*resp) == 0 || page == maxPagesOfIteration {
				break
			}
			page++
		} else {
			break
		}
	}
	jsonified, _ := json.Marshal(&GenAIJsonInput{
		Username:   j.Username,
		Bio:        j.Bio,
		CommitDets: commitDetails,
	})
	strVers := string(jsonified)
	return &strVers, nil
}

func callAzureOpenAi(promptType proto.ModelPromptType, input string) (*string, error) {
	credential, err := azidentity.NewDefaultAzureCredential(&azidentity.DefaultAzureCredentialOptions{
		TenantID: "5f45a518-36c3-406f-92fa-6f74ac0ea00f",
	})
	if err != nil {
		slog.Error("failed to get azure creds", "err", err)
		return nil, errInternal
	}
	client := openai.NewClient(
		azure.WithEndpoint("https://jonathan-testing-ai.openai.azure.com/", "2024-08-01-preview"),
		azure.WithTokenCredential(credential),
	)
	return askOpenAiPrompt(promptType, &client, "gpt-4o", input)
}

func askOpenAiPrompt(promptType proto.ModelPromptType, client *openai.Client, model string, input string) (*string, error) {
	chatParams := openai.ChatCompletionNewParams{
		Model: openai.ChatModel(model),
		// MaxTokens:           openai.Int(25000),
		MaxCompletionTokens: openai.Int(16384),
		Messages: []openai.ChatCompletionMessageParamUnion{
			{
				OfAssistant: &openai.ChatCompletionAssistantMessageParam{
					Content: openai.ChatCompletionAssistantMessageParamContentUnion{
						OfString: openai.String(config.New().GetModelPrompt(promptType)),
					},
				},
			},
			{
				OfUser: &openai.ChatCompletionUserMessageParam{
					Content: openai.ChatCompletionUserMessageParamContentUnion{
						OfString: openai.String(input),
					},
				},
			},
		},
	}
	resp, err := client.Chat.Completions.New(
		context.TODO(),
		chatParams,
	)
	if err != nil {
		slog.Error("failed to call openai", "err", err)
		return nil, errInternal
	}
	slog.Debug("OpenAI Response Generated", "resp", resp)
	return &resp.Choices[0].Message.Content, nil
}

func parseGithubWorker(ctx context.Context, id int, jobs chan octocatGrpcParseGithubJobWorkerParam, table idb.Database) {
	for {
		select {
		case <-ctx.Done():
			slog.Info("worker received done, finishing", "id", id)
			return
		case job := <-jobs:
			slog.Info("worker processing", "id", id, "job", job)
			value, err := table.GetAnswer(ctx, job.Username, job.idempotencyToken)
			if err != nil {
				slog.Error("failed to get answer", "err", err)
				return
			}
			if value.Status != idb.InProgress {
				return
			}
			jsonData, err := callRelevantGithubAPIS(job) // get json
			if err != nil {
				errText := err.Error()
				err := table.SetAnswer(ctx, job.Username, job.idempotencyToken, idb.Failed, &errText, nil)
				if err != nil {
					slog.Error("failed to set answer", "err", err)
					return
				}
				return
			}
			aiResponse, err := callAzureOpenAi(job.PromptType, *jsonData) // get ai response
			if err != nil {
				errText := err.Error()
				err := table.SetAnswer(ctx, job.Username, job.idempotencyToken, idb.Failed, &errText, nil)
				if err != nil {
					slog.Error("failed to set answer", "err", err)
					return
				}
				return
			}
			slog.Info("finished ai response saving to db", "id", id, "job", job)
			err = table.SetAnswer(ctx, job.Username, job.idempotencyToken, idb.Complete, nil, aiResponse)
			if err != nil {
				slog.Error("failed to set answer", "err", err)
				return
			}
		}
	}
}

func (o *octocatGrpcAPI) ParseGithub(ctx context.Context, req *proto.ParseGithubRequest) (*proto.ParseGithubResponse, error) {
	resp, err := o.WhoAmI(ctx, &proto.WhoAmIRequest{
		GithubToken: req.GithubToken,
	})
	if err != nil {
		return nil, err
	}
	slog.Info("Spawning new ParseGithub Goroutine in background", "user", resp.Username, "idempotencyToken", req.IdempotencyToken)
	err = o.dbTable.SetAnswer(ctx, resp.Username, req.IdempotencyToken, idb.InProgress, nil, nil)
	if err != nil {
		slog.Error("failed to set answer", "err", err)
		return nil, err
	}
	o.workerPool.parseGithubJobs <- octocatGrpcParseGithubJobWorkerParam{
		Username:   resp.Username,
		Bio:        resp.Bio,
		PromptType: req.PromptType,
		// sensitive do not log
		idempotencyToken: req.IdempotencyToken,
		githubOAuthToken: req.GithubToken,
	}
	return &proto.ParseGithubResponse{
		IdempotencyToken: req.IdempotencyToken,
	}, nil
}

func (o *octocatGrpcAPI) GetParsedGithubResult(ctx context.Context, req *proto.GetParsedGithubResultRequest) (*proto.GetParsedGithubResultResponse, error) {
	resp, err := o.WhoAmI(ctx, &proto.WhoAmIRequest{
		GithubToken: req.GithubToken,
	})
	if err != nil {
		return nil, err
	}
	value, err := o.dbTable.GetAnswer(ctx, resp.Username, req.IdempotencyToken)
	if err != nil {
		slog.Error("no parsed github result for this user and token", "user", resp.Username, "token", req.IdempotencyToken)
		return nil, status.Errorf(codes.NotFound, "No result")
	}
	if value.Status == idb.InProgress {
		slog.Info("still in progress for user and token", "user", resp.Username, "token", req.IdempotencyToken)
		return &proto.GetParsedGithubResultResponse{
			Status: string(idb.InProgress),
		}, nil
	}
	if value.Status == idb.Failed {
		slog.Info("failed for user and token, returning failure", "user", resp.Username, "token", req.IdempotencyToken)
		return &proto.GetParsedGithubResultResponse{
			Status: string(idb.Failed),
			Error:  value.Err,
		}, nil
	}
	return &proto.GetParsedGithubResultResponse{
		Status: string(idb.Complete),
		Result: value.Result,
	}, nil
}
