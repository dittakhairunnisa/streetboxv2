package repository

import (
	"log"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	"streetbox.id/app/trxspaces"
	"streetbox.id/entity"
	"streetbox.id/model"
	"streetbox.id/util"
)

// TrxSalesRepo ...
type TrxSalesRepo struct {
	DB *gorm.DB
}

// New ...
func New(db *gorm.DB) trxspaces.RepoInterface {
	return &TrxSalesRepo{db}
}

// Create ...
func (r *TrxSalesRepo) Create(data *entity.TrxParkingSpaceSales) (*gorm.DB, error) {
	trx := r.DB.Begin()
	if err := trx.Create(data).Error; err != nil {
		trx.Rollback()
		log.Printf("ERROR: %s", err.Error())
		return nil, err
	}
	log.Printf("INFO: Created TrxParkingSpaceSales: %+v", data)
	return trx, nil
}

// Update ...
func (r *TrxSalesRepo) Update(data *entity.TrxParkingSpaceSales, ID int64) (*gorm.DB, error) {
	trx := r.DB.Begin()
	if err := trx.Where("id = ?", ID).Table("trx_parking_space_sales").Update(data).Error; err != nil {
		trx.Rollback()
		log.Printf("ERROR: %s", err.Error())
		return nil, err
	}
	log.Printf("INFO: Created TrxParkingSpaceSales: %+v", data)
	return trx, nil
}

// GetByUserID ...
func (r *TrxSalesRepo) GetByUserID(id int64) (*[]entity.TrxParkingSpaceSales, error) {
	data := new([]entity.TrxParkingSpaceSales)
	if err := r.DB.Where("users_id = ?", id).Find(&data).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
		return nil, err
	}
	return data, nil
}

// GetMyParking ...
func (r *TrxSalesRepo) GetMyParking(id int64) (*[]model.ResMyParking, error) {
	data := new([]model.ResMyParking)
	if err := r.DB.Select(
		"distinct ps.id, ps.name, ps.address, ps.images_meta, "+
			"ps.rating, ps.description, ps.latitude, ps.longitude, "+
			"ps.start_time, ps.end_time").Joins("JOIN "+
		"parking_space_sales pss on ps.id = pss.parking_space_id").Joins("JOIN "+
		"trx_parking_space_sales t on pss.id = t.parking_space_sales_id").Joins("JOIN "+
		"merchant m on m.id = t.merchant_id").Joins("JOIN "+
		"merchant_users mu on m.id = mu.merchant_id").Where(
		"mu.users_id = ? and pss.end_date >= ? and ps.deleted_at is null and pss.deleted_at is null", id, time.Now()).
		Table("parking_space ps").Scan(&data).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
		return nil, err
	}
	return data, nil
}

// GetSlotMyParking ...
func (r *TrxSalesRepo) GetSlotMyParking(
	pspaceSalesID []int64, usersID int64) (*[]model.ResSlotMyParking, error) {
	data := new([]model.ResSlotMyParking)
	if err := r.DB.Select(
		"t.id, t.total_slot, pss.start_date, pss.end_date").Joins("JOIN "+
		"parking_space_sales pss on ps.id = pss.parking_space_id").Joins("JOIN "+
		"trx_parking_space_sales t on pss.id = t.parking_space_sales_id").Joins("JOIN "+
		"merchant m on m.id = t.merchant_id").Joins("JOIN "+
		"merchant_users mu on m.id = mu.merchant_id").Where(
		"t.deleted_at is null and t.parking_space_sales_id in (?) and mu.users_id = ?",
		pspaceSalesID, usersID).Table("parking_space ps").Scan(&data).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
		return nil, err
	}
	return data, nil
}

// GetAll ...
func (r *TrxSalesRepo) GetAll() *[]model.ResTrxList {
	data := new([]model.ResTrxList)
	r.DB.Raw("SELECT * from get_all_trx_space()").Scan(&data)
	return data
}

// DeleteByID ....
func (r *TrxSalesRepo) DeleteByID(id int64) error {
	if err := r.DB.Where("id = ?", id).
		Delete(new(entity.TrxParkingSpaceSales)).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
		return err
	}
	return nil
}

// GetByID ..
func (r *TrxSalesRepo) GetByID(id int64) *model.ResTrxList {
	data := new(model.ResTrxList)
	r.DB.Raw("SELECT * from get_all_trx_space() where id = ?", id).Scan(&data)
	if data.ID == 0 {
		return nil
	}
	return data
}

// GetByMerchantIDAndParkingSalesID ..
func (r *TrxSalesRepo) GetByMerchantIDAndParkingSalesID(
	merchantID int64, parkingSpaceSalesID int64) (*entity.TrxParkingSpaceSales, error) {
	data := new(entity.TrxParkingSpaceSales)
	if err := r.DB.Where("merchant_id = ? AND parking_space_sales_id = ?",
		merchantID, parkingSpaceSalesID).Find(&data).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
		return nil, err
	}
	return data, nil
}

// GetMerchantBySpaceID ..
func (r *TrxSalesRepo) GetMerchantBySpaceID(id int64) *[]model.MerchantSpace {
	data := new([]model.MerchantSpace)
	r.DB.Select("distinct m.id as merchant_id, m.name as merchant_name , "+
		"m.logo, m.banner, m.ig_account, ps.address").Joins("JOIN "+
		"merchant m on t.merchant_id = m.id").Joins("JOIN "+
		"parking_space_sales pss on t.parking_space_sales_id = pss.id").Joins("JOIN "+
		"parking_space ps on pss.parking_space_id = ps.id").
		Table("trx_parking_space_sales t").Where("ps.deleted_at is null and "+
		"pss.deleted_at is null and t.deleted_at is null and "+
		"pss.end_date >= ? and pss.parking_space_id = ?", time.Now(), id).
		Order("merchant_name asc").Scan(&data)
	return data
}

// GetMerchantSchedules ..
func (r *TrxSalesRepo) GetMerchantSchedules(id int64, parkingSpaceID int64) *[]model.Schedules {
	data := new([]model.Schedules)
	r.DB.Select("distinct pss.id, pss.start_date, pss.end_date").Joins("JOIN "+
		"trx_parking_space_sales t on pss.id = t.parking_space_sales_id").Joins("JOIN "+
		"parking_space ps on pss.parking_space_id = ps.id").
		Table("parking_space_sales pss").
		Where("pss.deleted_at is null and t.merchant_id = ? and "+
			"pss.end_date >= ? and ps.id = ?", id, time.Now(), parkingSpaceID).Order("pss.start_date asc").Scan(&data)
	return data
}

// GetSchedulesByTypesID method get schedules merchant
func (r *TrxSalesRepo) GetSchedulesByTypesID(id int64) *[]model.Schedules {
	data := new([]model.Schedules)
	r.DB.Select("distinct pss.id, pss.start_date, pss.end_date").Joins("JOIN "+
		"trx_parking_space_sales tpss on tr.trx_parking_space_sales_id = tpss.id").Joins("JOIN "+
		"parking_space_sales pss on tpss.parking_space_sales_id = pss.id").
		Table("tasks_regular tr").
		Where("pss.deleted_at is null and tr.id = ? and "+
			"pss.end_date >= ?", id, time.Now()).Order("pss.start_date asc").Scan(&data)
	return data
}

// GetList ...
func (r *TrxSalesRepo) GetList(
	limit, page int, sort []string, filter string) (*[]model.ResTrxList, int, int) {
	data := new([]model.ResTrxList)
	count := 0
	offset := util.Offset(page, limit)
	order := " "
	if len(sort) > 0 {
		flag := 0
		for _, value := range sort {
			if flag != 0 {
				order = ", " + value
			} else {
				order = "ORDER BY " + value
				flag = 1
			}

		}
	}
	filterQuery := ""
	if filter != "" {
		filterQuery = "WHERE status = ?"
	}
	var queries *gorm.DB

	if filter == "" {
		queries = r.DB.Raw("SELECT * from get_all_trx_space()" + order + " LIMIT " + strconv.Itoa(limit) + " OFFSET " + strconv.Itoa(offset))
	} else {
		queries = r.DB.Raw("SELECT * from get_all_trx_space() "+filterQuery+order+" LIMIT "+strconv.Itoa(limit)+" OFFSET "+strconv.Itoa(offset), filter)
	}
	queries.Scan(&data)
	var counts *gorm.DB
	if filter == "" {
		counts = r.DB.Raw("SELECT COUNT(*) from get_all_trx_space() ")
	} else {
		counts = r.DB.Raw("SELECT COUNT(*) from get_all_trx_space() "+filterQuery, filter)
	}
	counts.Count(&count)

	return data, count, offset
}
