package usersconfig

import "streetbox.id/entity"

type RepoInterface interface {
	UpdateRadius(rad int) (err error)
	GetConfig() (cfg entity.UsersConfig, err error)
}
