package idb

type Status string

const (
	InProgress Status = "inProgress"
	Failed     Status = "failed"
	Complete   Status = "complete"
)

type JobTableSchema struct {
	Status Status
	Err    *string
	Result *string
}
