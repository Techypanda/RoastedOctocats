package api

import (
	"context"
	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"techytechster.com/roastedoctocats/internal/github"
	"techytechster.com/roastedoctocats/pkg/igithub"
	"techytechster.com/roastedoctocats/pkg/proto"
)

func (o *octocatGrpcAPI) OAuth(ctx context.Context, req *proto.OAuthRequest) (*proto.OAuthResponse, error) {
	client := github.New()
	resp, err := client.Login(igithub.GithubLoginRequest{
		ClientId:      req.ClientId,
		RedirectURI:   req.RedirectUri,
		Code:          req.Code,
		CodeChallenge: req.CodeChallenge,
	})
	if err != nil {
		slog.Error("failed to do oauth login", "err", err)
		return nil, status.Errorf(codes.InvalidArgument, "Failed OAuth Login")
	}
	return &proto.OAuthResponse{
		AccessToken:        resp.AccessToken,
		AccessTokenExpiry:  int32(resp.AccessTokenExpiry),
		RefreshToken:       resp.RefreshToken,
		RefreshTokenExpiry: int32(resp.RefreshTokenExpiry),
	}, nil
}

func (o *octocatGrpcAPI) Refresh(ctx context.Context, req *proto.RefreshRequest) (*proto.OAuthResponse, error) {
	client := github.New()
	resp, err := client.Refresh(igithub.GithubRefreshRequest{
		ClientId:     req.ClientId,
		RefreshToken: req.RefreshToken,
	})
	if err != nil {
		slog.Error("failed to do refresh", "err", err)
		return nil, status.Errorf(codes.InvalidArgument, "Failed OAuth Login")
	}
	return &proto.OAuthResponse{
		AccessToken:        resp.AccessToken,
		AccessTokenExpiry:  int32(resp.AccessTokenExpiry),
		RefreshToken:       resp.RefreshToken,
		RefreshTokenExpiry: int32(resp.RefreshTokenExpiry),
	}, nil
}
