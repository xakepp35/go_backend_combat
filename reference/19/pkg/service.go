package pkg

import "context"

// TaskService is consumed by HTTP (and a future gRPC adapter).
// On error the response is not a successful result.
type TaskService interface {
	Create(context.Context, CreateTaskRequest) (CreateTaskResponse, error)
	Get(context.Context, GetTaskRequest) (GetTaskResponse, error)
	List(context.Context, ListTasksRequest) (ListTasksResponse, error)
	Update(context.Context, UpdateTaskRequest) (UpdateTaskResponse, error)
	Delete(context.Context, DeleteTaskRequest) (DeleteTaskResponse, error)
	Complete(context.Context, CompleteTaskRequest) (CompleteTaskResponse, error)
}
