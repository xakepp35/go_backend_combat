package pkg

import "context"

type InsertTaskRequest struct {
	Task Task
}

type InsertTaskResponse struct {
	Task Task
}

// TaskRepository is the outgoing port required by the application.
// Insert enforces unique IDs atomically. List returns tasks sorted by ID.
type TaskRepository interface {
	Insert(ctx context.Context, req InsertTaskRequest) (InsertTaskResponse, error)
	List(ctx context.Context, req ListTasksRequest) (ListTasksResponse, error)
}
