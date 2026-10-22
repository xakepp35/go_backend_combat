package pkg

type CreateTaskRequest struct {
	ID       string   `json:"id"`
	Title    string   `json:"title"`
	Priority Priority `json:"priority"`
}

type ListTasksRequest struct{}
