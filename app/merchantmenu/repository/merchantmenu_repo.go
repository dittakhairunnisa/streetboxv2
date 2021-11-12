package merchantmenu

import (
	"log"

	"github.com/jinzhu/gorm"
	"streetbox.id/app/merchantmenu"
	"streetbox.id/entity"
	"streetbox.id/util"
)

// MerchantMenusRepo ...
type MerchantMenusRepo struct {
	DB *gorm.DB
}

// New ...
func New(db *gorm.DB) merchantmenu.RepoInterface {
	return &MerchantMenusRepo{db}
}

// CreateMenu ...
func (r *MerchantMenusRepo) CreateMenu(menu *entity.MerchantMenu) error {
	if err := r.DB.Create(&menu).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
		return err
	}
	log.Printf("INFO : Created Merchant Menu : %+v", menu)
	return nil
}

// GetAllMenu ...
func (r *MerchantMenusRepo) GetAllMenu(merchantID int64, limit int, page int, sort []string) (*[]entity.MerchantMenu, int, int) {
	data := new([]entity.MerchantMenu)
	count := 0
	offset := util.Offset(page, limit)
	qry := r.DB.Limit(limit).Offset(offset)
	// sorting
	if len(sort) > 0 {
		for _, o := range sort {
			qry = qry.Order(o)
		}
	}
	qry = qry.Where("merchant_id = ?", merchantID)
	qry = qry.Find(data)
	qry = qry.Offset(0).Count(&count)
	return data, count, offset
}

// GetListMenu ...
func (r *MerchantMenusRepo) GetListMenu(merchantID int64, nearby, visit bool) *[]entity.MerchantMenu {
	data := new([]entity.MerchantMenu)
	var filter *gorm.DB
	if nearby && visit {
		filter = r.DB.Where("is_nearby = ? OR is_visit = ?", nearby, visit)
	} else if nearby {
		filter = r.DB.Where("is_nearby = ?", nearby)
	} else if visit {
		filter = r.DB.Where("is_visit = ?", visit)
	}
	filter.Find(&data, "merchant_id = ? AND is_active = true", merchantID)
	return data
}

// GetMenuByID ...
func (r *MerchantMenusRepo) GetMenuByID(merchantID int64, ID int64) *entity.MerchantMenu {
	data := new(entity.MerchantMenu)
	r.DB.Find(&data, "merchant_id = ? AND id = ?", merchantID, ID)
	return data
}

// GetOne ...
func (r *MerchantMenusRepo) GetOne(merchantID int64, ID int64) *entity.MerchantMenu {
	data := new(entity.MerchantMenu)
	r.DB.Find(&data, "id = ? AND merchant_id = ?", ID, merchantID)
	if data.ID == 0 {
		return nil
	}
	return data
}

// Update ...
func (r *MerchantMenusRepo) Update(menu *entity.MerchantMenu, merchantID int64, ID int64, isUpload bool) error {
	data := new(entity.MerchantMenu)
	r.DB.Find(&data, "id = ?", ID)
	if data.ID == 0 {
		return nil
	}
	if menu.Description != "" {
		data.Description = menu.Description
	}
	// if menu.Discount != 0 {
	// 	data.Discount = menu.Discount
	// }
	if !isUpload {
		data.IsActive = menu.IsActive
		data.Discount = menu.Discount
		data.Qty = menu.Qty
	}

	// if err := r.DB.Model(&entity.MerchantMenu{ID: ID}).Update("is_active", menu.IsActive).Error; err != nil {
	// 	log.Printf("ERROR: %s", err.Error())
	// 	return err
	// }

	if menu.Name != "" {
		data.Name = menu.Name
	}
	if menu.Price != 0 {
		data.Price = menu.Price
	}
	if menu.Photo != "" {
		data.Photo = menu.Photo
	}

	if err := r.DB.Model(&entity.MerchantMenu{ID: ID, MerchantID: merchantID}).Save(data).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
		return err
	}
	return nil
}

// UpdateStock ...
func (r *MerchantMenusRepo) UpdateStock(stock int, ID int64) error {
	data := new(entity.MerchantMenu)
	r.DB.Find(&data, "id = ?", ID)
	if data.ID == 0 {
		return nil
	}
	newStock := data.Qty - stock
	if err := r.DB.Model(&entity.MerchantMenu{ID: ID}).Update("Qty", newStock).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
		return err
	}
	return nil
}

// CekStock ...
func (r *MerchantMenusRepo) CekStock(ID int64, Qty int) bool {
	data := new(entity.MerchantMenu)
	r.DB.Find(&data, "id = ?", ID)
	if data.ID == 0 {
		return false
	}
	if (data.Qty - int(Qty)) < 0 {
		return false
	}
	return true
}

// Delete ...
func (r *MerchantMenusRepo) Delete(ID int64) error {
	if err := r.DB.Where("id = ?", ID).Delete(&entity.MerchantMenu{}).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
		return err
	}
	return nil
}

// DeleteImageMenu ..
func (r *MerchantMenusRepo) DeleteImageMenu(menu *entity.MerchantMenu, ID int64) error {
	if err := r.DB.Model(&entity.MerchantMenu{}).
		Where("id = ?", ID).Update("photo", nil).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
		return err
	}
	return nil
}
