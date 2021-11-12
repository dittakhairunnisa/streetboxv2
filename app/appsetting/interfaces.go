package appsetting

import (
	"streetbox.id/model"
)

// RepoInterface ..
type RepoInterface interface {
	GetByKey(string) *model.AppSetting
	UpdateByKey(string, string) error
}

// ServiceInterface ..
type ServiceInterface interface {
	GetByKey(string) *model.AppSetting
	UpdateByKey(string, *model.ReqUpdateAppSettingByKey) error
}
