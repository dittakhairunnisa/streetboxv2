package repository

import (
	"log"
	"time"

	"github.com/jinzhu/gorm"
	"streetbox.id/app/homevisitsales"
	"streetbox.id/entity"
	"streetbox.id/model"
	"streetbox.id/util"
)

// HomevisitRepo ..
type HomevisitRepo struct {
	DB *gorm.DB
}

// New ..
func New(db *gorm.DB) homevisitsales.RepoInterface {
	return &HomevisitRepo{db}
}

// GetAll ..
func (r *HomevisitRepo) GetAll(merchantID int64) *[]entity.HomevisitSales {
	data := new([]entity.HomevisitSales)
	r.DB.Where("merchant_id = ?", merchantID).Find(&data)
	if len(*data) <= 0 || data == nil {
		return nil
	}
	return data
}

// GetAllByMerchantAndDate ..
func (r *HomevisitRepo) GetAllByMerchantAndDate(merchantID int64,
	startDate string, endDate string) *[]entity.HomevisitSales {
	data := new([]entity.HomevisitSales)
	r.DB.Where("merchant_id = ? and startDate >= ? and endDate <= ? and deleted_at is null",
		merchantID, startDate, endDate).Find(&data)
	if len(*data) <= 0 || data == nil {
		return nil
	}
	return data
}

// Create ..
func (r *HomevisitRepo) Create(data *entity.HomevisitSales) (*entity.HomevisitSales, error) {
	if err := r.DB.Create(&data).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
		return nil, err
	}
	return data, nil
}

// Update ..
func (r *HomevisitRepo) Update(data *entity.HomevisitSales) (*entity.HomevisitSales, error) {
	if err := r.DB.Save(&data).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
		return nil, err
	}
	return data, nil
}

// CheckDate ...
func (r *HomevisitRepo) CheckDate(date string, merchantID int64) int {
	var countDate int
	r.DB.Where("deleted_at is null AND merchant_id = ? AND TO_CHAR(start_date, 'yyyy-mm-dd') = ?", merchantID, date).
		Select("COUNT(*)").Table("homevisit_sales").Count(&countDate)
	return countDate
}

// GetInfoByDate ..
func (r *HomevisitRepo) GetInfoByDate(merchantID int64, date string) *model.ResHomeVisitGetInfo {
	data := new(model.ResHomeVisitGetInfo)
	r.DB.Table("homevisit_sales").Where("merchant_id = ? and TO_CHAR(start_date, 'yyyy-mm-dd') = ?",
		merchantID, date).
		Group("TO_CHAR(start_date, 'yyyy-mm-dd'), deposit").
		Select("TO_CHAR(start_date, 'yyyy-mm-dd') as Date, deposit").Scan(&data)
	if data == nil {
		return nil
	}
	summary := new([]model.ResHomeVisitDetails)
	r.DB.Table("homevisit_sales").Where("merchant_id = ? AND deleted_at IS NULL and TO_CHAR(start_date, 'yyyy-mm-dd') = ?",
		merchantID, date).
		Select("id, TO_CHAR(start_date, 'HH24:MI:SS') as start_time, TO_CHAR(end_date, 'HH24:MI:SS') " +
			"as end_time, total as number_of_foodtruck").Scan(&summary)
	data.Summary = summary
	return data
}

// GetByID ..
func (r *HomevisitRepo) GetByID(ID int64) *entity.HomevisitSales {
	data := new(entity.HomevisitSales)
	r.DB.Where("id = ?", ID).Find(&data)
	return data
}

func (r *HomevisitRepo) GetMenuByTrxVisitSalesID(ID int64) (menus []model.ResMenu) {
	r.DB.Raw("SELECT m.name, t.quantity FROM merchant_menu m, trx_homevisit_menu_sales t WHERE t.trx_homevisit_sales_id = ? AND t.menu_id = m.id", ID).Scan(&menus)
	return
}

// GetAllList Get Visit Sales by End User
func (r *HomevisitRepo) GetAllList(limit, page int) (*[]model.ResVisitSales, int, int) {
	data := new([]model.ResVisitSales)
	count := 0
	offset := util.Offset(page, limit)
	qry := r.DB.Select("distinct m.id, m.banner, m.name, m.address, m.city, m.logo, m.ig_account, m.terms, mc.category as category, mc.hexcode as category_color").Joins("JOIN "+
		"homevisit_sales s on m.id = s.merchant_id").Joins("LEFT JOIN "+
		"merchant_category mc on mc.id = m.category_id").
		Where("m.deleted_at is null and s.deleted_at is null and s.end_date >= ? and s.available > 0",
			time.Now()).Table("merchant m").Offset(offset).Limit(limit)
	qry = qry.Order("m.city asc").Order("m.name asc")
	qry = qry.Scan(&data)
	qry = qry.Offset(0).Count(&count)
	return data, count, offset
}

// DeleteByDate ..
func (r *HomevisitRepo) DeleteByDate(date string, merchantID int64) (*gorm.DB, error) {
	db := r.DB.Begin()
	if err := db.Where("TO_CHAR(start_date, 'yyyy-mm-dd') = ? AND merchant_id = ?",
		date, merchantID).Delete(new([]entity.HomevisitSales)).Error; err != nil {
		db.Rollback()
		log.Printf("ERROR: %s", err.Error())
		return nil, err
	}
	return db, nil
}

// DeleteByID ..
func (r *HomevisitRepo) DeleteByID(ID int64, merchantID int64) (*gorm.DB, error) {
	db := r.DB.Begin()
	if err := db.Where("id = ? AND merchant_id = ?",
		ID, merchantID).Delete(new(entity.HomevisitSales)).Error; err != nil {
		db.Rollback()
		log.Printf("ERROR: %s", err.Error())
		return nil, err
	}
	db.Commit()
	return db, nil
}

// GetAvailableByMerchantID ...
func (r *HomevisitRepo) GetAvailableByMerchantID(id int64) *[]model.ResVisitSalesDetail {
	data := new([]model.ResVisitSalesDetail)
	r.DB.Select("id,start_date, end_date,deposit").Model(entity.HomevisitSales{}).
		Where("deleted_at is null and end_date > ? and merchant_id = ? and available > 0",
			time.Now(), id).Order("start_date").Scan(&data)
	return data
}

// UpdateAvailableByID ..
func (r *HomevisitRepo) UpdateAvailableByID(available int, id int64) error {
	if err := r.DB.Model(&entity.HomevisitSales{}).
		Where("id = ?", id).Update("available", available).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
		return err
	}
	return nil
}

// UpdateByID ..
func (r *HomevisitRepo) UpdateByID(data *entity.TrxHomevisitSales, id int64) error {
	qry := entity.TrxHomevisitSales{ID: id}
	if err := r.DB.Where(qry).Updates(data).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
		return err
	}
	return nil
}
