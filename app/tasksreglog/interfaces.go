package tasksreglog

import (
	"streetbox.id/entity"
)

// RepoInterface ..
type RepoInterface interface {
	Create(*entity.TasksRegularLog) error
}
