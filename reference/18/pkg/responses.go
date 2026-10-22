package pkg

type CreateTaskResponse struct{ Task Task }
type GetTaskResponse struct{ Task Task }
type ListTasksResponse struct{ Tasks []Task }
type UpdateTaskResponse struct{ Task Task }
type DeleteTaskResponse struct{}
type CompleteTaskResponse struct{ Task Task }
type InsertTaskResponse struct{ Task Task }
type ErrorResponse struct {
	Error string `json:"error"`
}
type HealthResponse struct {
	Status string `json:"status"`
}
