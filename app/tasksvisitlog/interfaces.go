package tasksvisitlog

import (
	"streetbox.id/entity"
)

// RepoInterface ..
type RepoInterface interface {
	Create(*entity.TasksHomevisitLog) error
}
