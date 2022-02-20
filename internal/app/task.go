package app

type TaskLogic struct {
}

type Task struct {
	ID          string
	GroupID     string
	Title       string
	Description string
}

func NewTask() *TaskLogic {
	return &TaskLogic{}
}

func (t *TaskLogic) Create() {

}
