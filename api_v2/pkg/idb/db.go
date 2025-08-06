package idb

import "context"

type Database interface {
	SetAnswer(ctx context.Context, username string, idempotencyToken string, status Status, errStr *string, response *string) error
	GetAnswer(ctx context.Context, username string, idempotencyToken string) (*JobTableSchema, error)
}
