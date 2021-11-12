package merchanttax

import (
	"github.com/jinzhu/gorm"
	"streetbox.id/app/merchanttax"
	"streetbox.id/entity"
)

// MerchantTaxsRepo ..
type MerchantTaxsRepo struct {
	DB *gorm.DB
}

// New ..
func New(db *gorm.DB) merchanttax.RepoInterface {
	return &MerchantTaxsRepo{db}
}

// Create ..
func (r *MerchantTaxsRepo) Create(data *entity.MerchantTax) (*entity.MerchantTax, error) {
	if err := r.DB.Create(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

// Update ..
func (r *MerchantTaxsRepo) Update(data *entity.MerchantTax, merchantID int64, ID int64) (*entity.MerchantTax, error) {
	if err := r.DB.Model(&entity.MerchantTax{ID: ID, MerchantID: merchantID}).Update(data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

// GetTax ..
func (r *MerchantTaxsRepo) GetTax(merchantID int64) *entity.MerchantTax {
	data := new(entity.MerchantTax)
	if err := r.DB.Find(&data, "merchant_id = ? AND is_active = ?", merchantID, true).Error; err != nil {
		if err := r.DB.Find(&data, "merchant_id = ?", merchantID).Order(data.UpdatedAt).Limit(1).Error; err != nil {
			data.MerchantID = merchantID
			data.IsActive = new(bool)
			data.Amount = new(float32)
			data.Type = new(int)
			if err := r.DB.Create(&data).Error; err != nil {
				return nil
			}
			return data
		}
		return data
	}
	return data
}

// Find ..
func (r *MerchantTaxsRepo) Find(model *entity.MerchantTax) *entity.MerchantTax {
	data := new(entity.MerchantTax)
	r.DB.Where(model).Find(&data)
	return data
}
