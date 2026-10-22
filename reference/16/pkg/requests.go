package pkg

type CreateTaskRequest struct {
	ID       string
	Title    string
	Priority int
	Done     bool
}
type GetTaskRequest struct{ ID string }
type ListTasksRequest struct{}
type UpdateTaskRequest struct{ Task Task }
type DeleteTaskRequest struct{ ID string }
type CompleteTaskRequest struct{ ID string }
type InsertTaskRequest struct{ Task Task }
