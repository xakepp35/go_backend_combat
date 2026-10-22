package pkg

type Task struct {
	ID       string
	Title    string
	Priority int
	Done     bool
}

func (t *Task) Complete() { t.Done = true }
