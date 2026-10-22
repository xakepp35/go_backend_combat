package pkg

import "context"

// TaskRepository is the application's outgoing port, implemented by adapters.
// Insert enforces unique IDs atomically. List is sorted by ID and never nil.
// Invalid, conflicting or missing writes leave existing records unchanged.
type TaskRepository interface {
	Insert(context.Context, InsertTaskRequest) (InsertTaskResponse, error)
	Get(context.Context, GetTaskRequest) (GetTaskResponse, error)
	List(context.Context, ListTasksRequest) (ListTasksResponse, error)
	Update(context.Context, UpdateTaskRequest) (UpdateTaskResponse, error)
	Delete(context.Context, DeleteTaskRequest) (DeleteTaskResponse, error)
	Complete(context.Context, CompleteTaskRequest) (CompleteTaskResponse, error)
}
