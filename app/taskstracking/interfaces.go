package taskstracking

import (
	"streetbox.id/entity"
	"streetbox.id/model"
)

// RepoInterface ..
type RepoInterface interface {
	Create(*entity.TasksTracking) error
	GetTrackingByID(int64) *model.ResTasksTracking
	DeleteByTasksID(int64)
	GetLiveTracking(lat, lon, distance float64) *[]model.ResLiveTracking
	DeleteAll()
}
