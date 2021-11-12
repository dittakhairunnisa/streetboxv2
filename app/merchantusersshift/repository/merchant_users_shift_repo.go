package repository

import (
	"log"
	"time"

	"github.com/jinzhu/gorm"
	"streetbox.id/app/merchantusersshift"
	"streetbox.id/entity"
)

// MerchantUsersShiftRepo ...
type MerchantUsersShiftRepo struct {
	DB *gorm.DB
}

// New ...
func New(db *gorm.DB) merchantusersshift.RepoInterface {
	return &MerchantUsersShiftRepo{db}
}

// Create ...
func (r *MerchantUsersShiftRepo) Create(
	data *entity.MerchantUsersShift) (*entity.MerchantUsersShift, error) {
	if err := r.DB.Create(data).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
		return nil, err
	}
	log.Printf("INFO: Created MerchantUsersShift: %+v", data)
	return data, nil
}

// IsUsersShiftIn ...
func (r *MerchantUsersShiftRepo) IsUsersShiftIn(
	usersID int64) bool {
	data := new(entity.MerchantUsersShift)
	timeNow := time.Now().Format("2006-01-02")
	r.DB.Joins("JOIN merchant_users on "+
		"merchant_users_shift.merchant_users_id = merchant_users.id").
		Where("merchant_users.users_id = ? and "+
			"TO_CHAR(merchant_users_shift.updated_at,'yyyy-mm-dd') like ?",
			usersID, timeNow).Model(data).Find(&data)
	if data.ID > 0 && data.Shift == "IN" {
		return true
	}
	return false
}
