package api

import (
	"context"
	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"techytechster.com/roastedoctocats/internal/db"
	"techytechster.com/roastedoctocats/internal/github"
	pkgapi "techytechster.com/roastedoctocats/pkg/api"
	"techytechster.com/roastedoctocats/pkg/idb"
	"techytechster.com/roastedoctocats/pkg/proto"
)

type octocatGrpcParseGithubJobWorkerParam struct {
	Username         string
	Bio              string
	PromptType       proto.ModelPromptType
	idempotencyToken string
	githubOAuthToken string
}

// dummy worker pool
type octocatGrpcWorkerPool struct {
	// define
	parseGithubJobs chan octocatGrpcParseGithubJobWorkerParam
}

type octocatGrpcAPI struct {
	dbTable    idb.Database
	workerPool octocatGrpcWorkerPool
	proto.UnimplementedOctoRoasterAPIServer
}

func (o *octocatGrpcAPI) WhoAmI(ctx context.Context, req *proto.WhoAmIRequest) (*proto.WhoAmIResponse, error) {
	slog.Debug("New WhoAmI Request")
	client := github.New()
	resp, err := client.WhoAmI(req.GithubToken)
	if err != nil {
		slog.Error("failed to do whoami", "err", err)
		return nil, status.Errorf(codes.Unauthenticated, "Unauthorized")
	}
	return &proto.WhoAmIResponse{
		Username: resp.Username,
		Bio:      resp.Bio,
	}, nil
}

func (o *octocatGrpcAPI) Ping(ctx context.Context, request *proto.PingRequest) (*proto.PingResponse, error) {
	slog.Debug("New Ping Request", "idempotencyToken", request.IdempotencyToken)
	return &proto.PingResponse{Message: pkgapi.PingResponseMessage}, nil
}

const poolSize int = 10

func New() (proto.OctoRoasterAPIServer, context.CancelFunc) {
	dbInstance, err := db.New()
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	parseGithubJobs := make(chan octocatGrpcParseGithubJobWorkerParam, poolSize)
	for i := range poolSize {
		go parseGithubWorker(ctx, i, parseGithubJobs, dbInstance)
	}
	return &octocatGrpcAPI{
		workerPool: octocatGrpcWorkerPool{
			parseGithubJobs: parseGithubJobs,
		},
		dbTable: dbInstance,
	}, cancel
}
