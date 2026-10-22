package pkg

type CreateTaskResponse struct {
	Task Task `json:"task"`
}

type ListTasksResponse struct {
	Tasks []Task `json:"tasks"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
