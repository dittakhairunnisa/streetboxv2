package repository

import (
	"log"

	"github.com/jinzhu/gorm"
	"streetbox.id/app/merchantusers"
	"streetbox.id/entity"
)

// MerchantUsersRepo ...
type MerchantUsersRepo struct {
	DB *gorm.DB
}

// New ...
func New(db *gorm.DB) merchantusers.RepoInterface {
	return &MerchantUsersRepo{db}
}

// Create ...
func (r *MerchantUsersRepo) Create(
	db *gorm.DB, merchantID, usersID int64) error {
	data := new(entity.MerchantUsers)
	data.MerchantID = merchantID
	data.UsersID = usersID
	if err := db.Create(&data).Error; err != nil {
		db.Rollback()
		log.Printf("ERROR: %s", err.Error())
		return err
	}
	log.Printf("INFO: Created MerchantUsers: %+v", data)
	return nil
}

// IsExist check is data valid (merchantID and usersID)
func (r *MerchantUsersRepo) IsExist(usersID int64) bool {
	data := new(entity.MerchantUsers)
	r.DB.Where("users_id = ?", usersID).Find(&data)
	if data.ID == 0 {
		return false
	}
	return true
}

// GetUserIdsByMerchantID ..
func (r *MerchantUsersRepo) GetUserIdsByMerchantID(merchantID int64) *[]int64 {
	data := new([]entity.MerchantUsers)
	var userIDs []int64
	if err := r.DB.Where("merchant_id = ?", merchantID).Find(&data).Pluck("users_id", &userIDs).Error; err != nil {
		return nil
	}
	return &userIDs
}

// GetByUsersID ...
func (r *MerchantUsersRepo) GetByUsersID(
	usersID int64) *entity.MerchantUsers {
	data := new(entity.MerchantUsers)
	r.DB.Find(&data, "users_id = ?", usersID)
	return data
}

// GetAdminByMerchantID ..
func (r *MerchantUsersRepo) GetAdminByMerchantID(merchantID int64) *entity.MerchantUsers {
	data := new(entity.MerchantUsers)
	r.DB.Select("mu.*").Joins("JOIN "+
		"users u on mu.users_id = u.id").Joins("JOIN "+
		"users_role ur on u.id = ur.users_id").Joins("JOIN "+
		"role r on ur.role_id = r.id").
		Where("mu.merchant_id = ? and r.name = ?", merchantID, "admin").
		Table("merchant_users mu").Scan(&data)
	if data.ID == 0 {
		return nil
	}
	return data
}

// DeleteByMerchantID ..
func (r *MerchantUsersRepo) DeleteByMerchantID(
	db *gorm.DB, id int64) error {
	if err := db.Delete(entity.MerchantUsers{},
		"merchant_id = ?", id).Error; err != nil {
		db.Rollback()
		log.Printf("ERROR: %s", err.Error())
		return err
	}
	return nil
}

// DeleteByFoodtruckID ..
func (r *MerchantUsersRepo) DeleteByFoodtruckID(id int64) error {
	if err := r.DB.Delete(entity.MerchantUsers{},
		"users_id = ?", id).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
		return err
	}
	return nil
}

// CountFoodtruckByMerchantID ..
func (r *MerchantUsersRepo) CountFoodtruckByMerchantID(merchantID int64) int {
	//merchantUsers := new(entity.MerchantUsers)
	var count int
	if err := r.DB.Table("merchant_users mu").
		Joins("LEFT JOIN users_role ur on mu.users_id = ur.users_id").
		Joins("LEFT JOIN role r on ur.role_id = r.id").
		Where("merchant_id = ? AND r.name = ? and mu.deleted_at is null",
			merchantID, "foodtruck").Count(&count).Error; err != nil {
		return 0
	}
	return count
}

// Update ...
func (r *MerchantUsersRepo) Update(data *entity.MerchantUsers, id int64) error {
	if err := r.DB.Model(&entity.MerchantUsers{ID: id}).Updates(data).Error; err != nil {
		return err
	}
	return nil
}

// GetOne ..
func (r *MerchantUsersRepo) GetOne(id int64) *entity.MerchantUsers {
	data := new(entity.MerchantUsers)
	r.DB.Select("mu.*").
		Where("mu.id = ? ", id).
		Table("merchant_users mu").Scan(&data)
	if data.ID == 0 {
		return nil
	}
	return data
}
