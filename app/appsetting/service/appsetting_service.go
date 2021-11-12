package service

import (
	"github.com/jinzhu/copier"
	"streetbox.id/app/appsetting"
	"streetbox.id/model"
	"streetbox.id/entity"
)

// AppSettingService ..
type AppSettingService struct {
	AppSettingRepo    appsetting.RepoInterface
}

// New ..
func New(
	appSettingRepo appsetting.RepoInterface) appsetting.ServiceInterface {
	return &AppSettingService{appSettingRepo}
}

// GetByKey ..
func (s *AppSettingService) GetByKey(key string) *model.AppSetting {
	return s.AppSettingRepo.GetByKey(key)
}

// UpdateByKey ..
func (s *AppSettingService) UpdateByKey(key string, req *model.ReqUpdateAppSettingByKey) error {
	data := new(entity.AppSetting)
	copier.Copy(&data, req)

	return s.AppSettingRepo.UpdateByKey(key, data.Value)
}