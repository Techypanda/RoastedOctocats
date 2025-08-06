package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/data/aztables"
	"techytechster.com/roastedoctocats/pkg/idb"
)

type database struct {
	serviceClient *aztables.Client
}

var errFailedSetAnswer = errors.New("failed to set answer")
var ErrFailedGetAnswerNoItem = errors.New("failed marshal")
var errFailedGetAnswer = errors.New("failed get answer")

func (db *database) GetAnswer(ctx context.Context, username string, idempotencyToken string) (*idb.JobTableSchema, error) {
	response, err := db.serviceClient.GetEntity(ctx, idempotencyToken, fmt.Sprintf("%s:%s", username, idempotencyToken), nil)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return nil, errors.Join(ErrFailedGetAnswerNoItem, err)
		}
		return nil, errors.Join(errFailedGetAnswer, err)
	}
	var entity aztables.EDMEntity
	err = json.Unmarshal(response.Value, &entity)
	if err != nil {
		return nil, errors.Join(ErrFailedGetAnswerNoItem, err)
	}
	var parsedResp idb.JobTableSchema
	status, ok := entity.Properties["Status"].(string)
	if !ok {
		return nil, ErrFailedGetAnswerNoItem
	}
	parsedResp.Status = idb.Status(status)
	errText, ok := entity.Properties["Error"].(string)
	if ok {
		parsedResp.Err = &errText
	}
	responseText, ok := entity.Properties["Response"].(string)
	if ok {
		parsedResp.Result = &responseText
	}
	return &parsedResp, nil
}

func (db *database) SetAnswer(ctx context.Context, username string, idempotencyToken string, status idb.Status, errStr *string, response *string) error {
	entity := aztables.EDMEntity{
		Entity: aztables.Entity{
			RowKey:       fmt.Sprintf("%s:%s", username, idempotencyToken),
			PartitionKey: idempotencyToken,
		},
		Properties: map[string]any{
			"Status":   status,
			"Error":    errStr,
			"Response": response,
		},
	}
	bytes, err := json.Marshal(&entity)
	if err != nil {
		return errors.Join(errFailedSetAnswer, err)
	}
	_, err = db.serviceClient.UpsertEntity(ctx, bytes, nil)
	if err != nil {
		return errors.Join(errFailedSetAnswer, err)
	}
	return nil
}

var errFailedConstructingDatabase = errors.New("failed to create db")

func New() (idb.Database, error) {
	credential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, errors.Join(errFailedConstructingDatabase, err)
	}
	client, err := aztables.NewServiceClient("https://octocatroasterdb.table.cosmos.azure.com:443/", credential, nil)
	if err != nil {
		return nil, errors.Join(errFailedConstructingDatabase, err)
	}
	table := client.NewClient("octoroastertable")
	return &database{
		serviceClient: table,
	}, nil
}
