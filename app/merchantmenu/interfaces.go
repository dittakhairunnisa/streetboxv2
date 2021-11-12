package merchantmenu

import (
	"streetbox.id/entity"
)

// RepoInterface ...
type RepoInterface interface {
	CreateMenu(*entity.MerchantMenu) error
	GetAllMenu(merchantID int64, limit int, page int, sort []string) (*[]entity.MerchantMenu, int, int)
	GetListMenu(merchantID int64, nearby, visit bool) *[]entity.MerchantMenu
	GetMenuByID(merchantID int64, ID int64) *entity.MerchantMenu
	GetOne(int64, int64) *entity.MerchantMenu
	Update(*entity.MerchantMenu, int64, int64, bool) error
	UpdateStock(int, int64) error
	Delete(ID int64) error
	DeleteImageMenu(*entity.MerchantMenu, int64) error
	CekStock(ID int64, Qty int) bool
}
