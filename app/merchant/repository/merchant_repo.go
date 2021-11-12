package repository

import (
	"log"
	"sort"

	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
	"streetbox.id/app/merchant"
	"streetbox.id/entity"
	"streetbox.id/model"
	"streetbox.id/util"
)

// MerchantRepo ...
type MerchantRepo struct {
	DB *gorm.DB
}

// New ...
func New(db *gorm.DB) merchant.RepoInterface {
	return &MerchantRepo{db}
}

// Create ...
func (r *MerchantRepo) Create(data *entity.Merchant) (*gorm.DB, error) {
	db := r.DB.Begin()
	if err := db.Create(&data).Error; err != nil {
		db.Rollback()
		log.Printf("ERROR: %s", err.Error())
		return nil, err
	}
	log.Printf("INFO: Created Merchant: %+v", data)
	return db, nil
}

// GetAllFoodtruck ...
func (r *MerchantRepo) GetAllFoodtruck(
	merchantID int64) *[]entity.Users {
	data := new([]entity.Users)
	if err := r.DB.Select("u.*").Joins("JOIN "+
		"users u on mu.users_id = u.id").Joins("JOIN "+
		"users_role ur on u.id = ur.users_id").Joins("JOIN "+
		"role r on ur.role_id = r.id").Where(
		"mu.merchant_id = ? and r.name = ? and "+
			"mu.deleted_at is null", merchantID,
		"foodtruck").Table("merchant_users mu").
		Scan(&data).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
	}

	return data
}

// GetInfo ... merchant
func (r *MerchantRepo) GetInfo(usersID int64) *model.Merchant {
	data := new(model.Merchant)
	if err := r.DB.Select("m.*,mu.id as merchant_users_id, mc.category, mc.hexcode").Joins("LEFT JOIN merchant_category mc ON m.category_id = mc.id").Joins("JOIN merchant_users mu on m.id = mu.merchant_id").Where("mu.users_id = ? and m.deleted_at is null", usersID).Table("merchant m").Scan(&data).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
	}
	if data.ID == 0 {
		return nil
	}
	return data
}

// Update ..
func (r *MerchantRepo) Update(data *entity.Merchant, id int64) error {
	if err := r.DB.Model(&entity.Merchant{ID: id}).Updates(data).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
		return err
	}
	return nil
}

// GetAll ...
func (r *MerchantRepo) GetAll() *[]model.Merchant {
	data := new([]model.Merchant)
	// r.DB.Find(&data).Order("name")
	r.DB.Table("merchant m").Select("m.*, mc.category, mc.hexcode").Joins("LEFT JOIN merchant_category mc ON m.category_id = mc.id;").Scan(&data)
	return data
}

// DeleteByID ...
func (r *MerchantRepo) DeleteByID(id int64) (*gorm.DB, error) {
	db := r.DB.Begin()
	if err := db.Where("id = ?", id).
		Delete(new(entity.Merchant)).Error; err != nil {
		db.Rollback()
		log.Printf("ERROR: %s", err.Error())
		return nil, err
	}
	return db, nil
}

// GetByID ..
func (r *MerchantRepo) GetByID(id int64) *entity.Merchant {
	data := new(entity.Merchant)
	r.DB.Find(&data, "id = ?", id)
	if data.ID == 0 {
		return nil
	}
	return data
}

// GetFoodtruckByID ..
func (r *MerchantRepo) GetFoodtruckByID(id int64) *entity.Users {
	data := new(entity.Users)
	r.DB.Find(&data, "id = ?", id)
	if data.ID == 0 {
		return nil
	}
	return data
}

// GetFoodtruckTasks ..
func (r *MerchantRepo) GetFoodtruckTasks(merchantID int64) *[]model.ResGetFoodtruckTasks {
	data := make([]model.ResGetFoodtruckTasks, 0)
	foodtruck := r.GetAllFoodtruck(merchantID)
	foodtruckTasks := new([]model.ResGetFoodtruckTasks)
	r.DB.Select("distinct u.*").Joins("JOIN "+
		"merchant_users mu on u.id = mu.users_id").Joins("JOIN "+
		"tasks t on mu.id = t.merchant_users_id").
		Where("mu.merchant_id = ? and t.status < 4", merchantID).Table("users u").
		Scan(&foodtruckTasks)
	if len(*foodtruckTasks) > 0 {
		for _, u1 := range *foodtruck {
			for _, u2 := range *foodtruckTasks {
				if u1.ID == u2.ID {
					data = append(data, u2)
				}
			}
		}
	}
	return &data
}

// GetNearby .. Landing Page Nearby End User
func (r *MerchantRepo) GetNearby(
	limit int, page int, lat float64, lon float64,
	distance float64) (
	*model.NearbySorted, int, int) {
	data := make([]model.ResMerchantNearby, 0)
	offset := util.Offset(page, limit)
	rows, _ := r.DB.Raw("select * from merchant_nearby(?,?,?,?,?) order by updated_at desc, nearby asc,  status desc",
		limit, offset, lat, lon, distance).Rows()
	for rows.Next() {
		var dat model.ResMerchantNearby
		r.DB.ScanRows(rows, &dat)
		data = append(data, dat)
	}
	distinct := make(model.NearbySorted, 0)
	merchantUsersIds := make([]int64, 0)
	merchantUsersCheckout := make(model.NearbySorted, 0)
	for _, dataMerchant := range data {
		if util.IsExistSlicesInt64(&merchantUsersIds, dataMerchant.MerchantUsersID) == true {
			continue
		}
		merchantUsersIds = append(merchantUsersIds, dataMerchant.MerchantUsersID)
		if dataMerchant.Status == 4 {
			merchantUsersCheckout = append(merchantUsersCheckout, dataMerchant)
			continue
		}
		copier.Copy(&distinct, dataMerchant)
	}
	distinct = append(distinct, merchantUsersCheckout...)
	sort.Sort(distinct)
	return &distinct, len(distinct), offset
}

// RemoveImage method by types: logo/banner
func (r *MerchantRepo) RemoveImage(types string, merchantID int64) error {
	if err := r.DB.Model(&entity.Merchant{}).Where("id = ?", merchantID).
		Update(types, "").Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
		return err
	}
	log.Printf("INFO: Remove Image %s from merchantID %d", types, merchantID)
	return nil
}
