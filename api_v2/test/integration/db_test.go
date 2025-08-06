package integration

import (
	"testing"
)

func TestDatabaseOperations(t *testing.T) {
	// dbInstance, err := db.New()
	// if err != nil {
	// 	t.Fatalf("failed new: %s", err)
	// }
	// uuidVal := uuid.Must(uuid.NewV7()).String()
	// _, err = dbInstance.GetAnswer(t.Context(), "testuser", uuidVal)
	// if err == nil || !errors.Is(err, db.ErrFailedGetAnswerNoItem) {
	// 	t.Fatal("failed get item, expected ErrFailedGetAnswerNoItem")
	// }
	// err = dbInstance.SetAnswer(t.Context(), "testuser", uuidVal, "inProgress", nil, nil)
	// if err != nil {
	// 	t.Fatalf("failed set in progress: %s", err)
	// }
	// val, err := dbInstance.GetAnswer(t.Context(), "testuser", uuidVal)
	// if err != nil {
	// 	t.Fatalf("failed get in progress: %s", err)
	// }
	// if val.Status != idb.InProgress {
	// 	t.Fatal("get inprogress val is not inprogress")
	// }
	// errText := "some mock internal failure"
	// err = dbInstance.SetAnswer(t.Context(), "testuser", uuidVal, "failed", &errText, nil)
	// if err != nil {
	// 	t.Fatalf("failed set failed: %s", err)
	// }
	// val, err = dbInstance.GetAnswer(t.Context(), "testuser", uuidVal)
	// if err != nil {
	// 	t.Fatalf("failed get failed: %s", err)
	// }
	// if val.Status != idb.Failed {
	// 	t.Fatal("get failed val is not failed")
	// }
	// if val.Err == nil || *val.Err != errText {
	// 	t.Fatal("get failed err is not expected")
	// }
	// respText := "some ai text"
	// err = dbInstance.SetAnswer(t.Context(), "testuser", uuidVal, "complete", nil, &respText)
	// if err != nil {
	// 	t.Fatalf("failed set complete: %s", err)
	// }
	// val, err = dbInstance.GetAnswer(t.Context(), "testuser", uuidVal)
	// if err != nil {
	// 	t.Fatalf("failed get complete: %s", err)
	// }
	// if val.Status != idb.Complete {
	// 	t.Fatal("get complete val is not complete")
	// }
	// if val.Result == nil || *val.Result != respText {
	// 	t.Fatal("get complete resp is not complete")
	// }
}
