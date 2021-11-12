package useraddress

import "streetbox.id/entity"

type RepoInterface interface {
	Create(addr *entity.UsersAddress) (err error)
	GetPrimaryByUserID(userID int64) (addrs entity.UsersAddress, err error)
	GetByUserID(userID int64) (addrs []entity.UsersAddress, err error)
	Update(addr entity.UsersAddress) (err error)
	Delete(id, userID int64) (err error)
	Switch(id, userID int64) (err error)
}