package pkg

import "context"

// TaskService is the application contract for incoming adapters.
// An error means the response is not a successful result.
type TaskService interface {
	Create(ctx context.Context, req CreateTaskRequest) (CreateTaskResponse, error)
	List(ctx context.Context, req ListTasksRequest) (ListTasksResponse, error)
}
