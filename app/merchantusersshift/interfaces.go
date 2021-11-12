package merchantusersshift

import (
	"streetbox.id/entity"
)

// RepoInterface ...
type RepoInterface interface {
	Create(*entity.MerchantUsersShift) (*entity.MerchantUsersShift, error)
	IsUsersShiftIn(usersID int64) bool
}
