package repository

import (
	"errors"
	"log"
	"time"

	"streetbox.id/app/sales"
	"streetbox.id/util"

	"streetbox.id/model"

	"github.com/jinzhu/gorm"
	"streetbox.id/entity"
)

// SalesRepo ...
type SalesRepo struct {
	DB *gorm.DB
}

// New Sales Repo
func New(db *gorm.DB) sales.RepoInterface {
	return &SalesRepo{db}
}

// Create ...
func (r *SalesRepo) Create(pss *entity.ParkingSpaceSales) error {
	if err := r.DB.Create(&pss).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
		return err
	}
	log.Printf("INFO: Created ParkingSpaceSales: %+v", pss)
	return nil
}

// FindBySpaceID ...
func (r *SalesRepo) FindBySpaceID(id int64, limit,
	page int, sort []string) (*[]model.ResSales, int, int) {
	data := new([]model.ResSales)
	count := 0
	offset := util.Offset(page, limit)
	qry := r.DB.Select("s.id,p.name,s.start_date,s.end_date,"+
		"s.total_slot,s.available_slot,s.point,s.parking_space_id,"+
		"s.created_at,s.updated_at").Joins("JOIN "+
		"parking_space_sales s on p.id = s.parking_space_id").
		Where("p.deleted_at is null and s.end_date >= ?"+
			"and s.parking_space_id = ?", time.Now(), id).
		Table("parking_space p").Offset(offset).Limit(limit)
	// sorting
	if len(sort) > 0 {
		for _, o := range sort {
			qry = qry.Order(o)
		}
	}
	qry = qry.Scan(&data)
	qry = qry.Offset(0).Count(&count)
	return data, count, offset
}

// UpdateAvailableSlot ...
func (r *SalesRepo) UpdateAvailableSlot(db *gorm.DB, qty int, id int64) (*gorm.DB, error) {
	currentSlot := r.GetOne(id).AvailableSlot
	availableSlot := currentSlot - qty
	if availableSlot < 0 {
		return nil, errors.New("Available Slot Not Enough")
	}
	if err := db.Model(new(entity.ParkingSpaceSales)).Where("id = ?", id).
		Update("available_slot", availableSlot).Error; err != nil {
		db.Rollback()
		log.Printf("ERROR: %s", err.Error())
		return nil, err
	}
	log.Printf("INFO: Update Available Slot With Qty: %d, ID: %d", qty, id)
	return db, nil
}

// GetOne ...
func (r *SalesRepo) GetOne(id int64) *entity.ParkingSpaceSales {
	data := new(entity.ParkingSpaceSales)
	r.DB.Find(&data, "id = ?", id)
	if data.ID == 0 {
		return nil
	}
	return data
}

// Update ...
func (r *SalesRepo) Update(
	data *entity.ParkingSpaceSales, id int64) (*entity.ParkingSpaceSales, error) {
	if err := r.DB.Model(&entity.ParkingSpaceSales{ID: id}).Updates(&data).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
		return nil, err
	}
	log.Printf("INFO: Updated ParkingSpaceSales: %+v", data)
	return data, nil
}

// DeleteByID ...
func (r *SalesRepo) DeleteByID(id int64) error {
	data := new(entity.ParkingSpaceSales)
	if err := r.DB.Delete(&data, "id = ?", id).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
		return err
	}
	log.Printf("INFO: Deleted ParkingSpaceSales: %+v", data)
	return nil
}

// FindLikeName ..
func (r *SalesRepo) FindLikeName(name string,
	limit, page int, sort []string) (*[]model.ResSales, int, int) {
	data := new([]model.ResSales)
	count := 0
	offset := util.Offset(page, limit)
	qry := r.DB.Select("s.id,p.name,s.start_date,s.end_date,"+
		"s.total_slot,s.available_slot,s.point,s.parking_space_id,"+
		"s.created_at,s.updated_at").Joins("JOIN "+
		"parking_space_sales s on p.id = s.parking_space_id").
		Where("p.deleted_at is null and s.deleted_at is null and s.end_date >= ?"+
			"and (p.name ilike ? or p.address ilike ?)", time.Now(), "%"+name+"%", "%"+name+"%").
		Table("parking_space p").Offset(offset).Limit(limit)
	// sorting
	if len(sort) > 0 {
		for _, o := range sort {
			qry = qry.Order(o)
		}
	}
	qry = qry.Scan(&data)
	qry = qry.Offset(0).Count(&count)
	return data, count, offset
}

// FindLikeNameBackoffice ..
func (r *SalesRepo) FindLikeNameBackoffice(name string,
	limit, page int, sort []string) (*[]model.ResSales, int, int) {
	data := new([]model.ResSales)
	count := 0
	offset := util.Offset(page, limit)
	qry := r.DB.Select("s.id,p.name,s.start_date,s.end_date,"+
		"s.total_slot,s.available_slot,s.point,s.parking_space_id,"+
		"s.created_at,s.updated_at").Joins("JOIN "+
		"parking_space_sales s on p.id = s.parking_space_id").
		Where("p.deleted_at is null and s.deleted_at is null "+
			"and (p.name ilike ? or p.address ilike ?)", "%"+name+"%", "%"+name+"%").
		Table("parking_space p").Offset(offset).Limit(limit)
	// sorting
	if len(sort) > 0 {
		for _, o := range sort {
			qry = qry.Order(o)
		}
	}
	qry = qry.Scan(&data)
	qry = qry.Offset(0).Count(&count)
	return data, count, offset
}

// GetAll ..
func (r *SalesRepo) GetAll(limit,
	page int, sort []string) (*[]model.ResSales, int, int) {
	data := new([]model.ResSales)
	count := 0
	offset := util.Offset(page, limit)
	qry := r.DB.Select("s.id,p.name,s.start_date,s.end_date,"+
		"s.total_slot,s.available_slot,s.point,s.parking_space_id,"+
		"s.created_at,s.updated_at").Joins("JOIN "+
		"parking_space_sales s on p.id = s.parking_space_id").
		Where("p.deleted_at is null and s.deleted_at is null and s.end_date >= ?", time.Now()).
		Table("parking_space p").Offset(offset).Limit(limit)
	// sorting
	if len(sort) > 0 {
		for _, o := range sort {
			qry = qry.Order(o)
		}
	} else {
		qry = qry.Order("p.name asc")
	}
	qry = qry.Scan(&data)
	qry = qry.Offset(0).Count(&count)
	return data, count, offset
}

// GetAllBackoffice ..
func (r *SalesRepo) GetAllBackoffice(limit,
	page int, sort []string) (*[]model.ResSales, int, int) {
	data := new([]model.ResSales)
	count := 0
	offset := util.Offset(page, limit)
	qry := r.DB.Select("s.id,p.name,s.start_date,s.end_date," +
		"s.total_slot,s.available_slot,s.point,s.parking_space_id," +
		"s.created_at,s.updated_at").Joins("JOIN " +
		"parking_space_sales s on p.id = s.parking_space_id").
		Where("p.deleted_at is null and s.deleted_at is null").
		Table("parking_space p").Offset(offset).Limit(limit)
	// sorting
	if len(sort) > 0 {
		for _, o := range sort {
			qry = qry.Order(o)
		}
	} else {
		qry = qry.Order("p.name asc")
	}
	qry = qry.Scan(&data)
	qry = qry.Offset(0).Count(&count)
	return data, count, offset
}

// GetAllBackofficeNonPaginate ..
func (r *SalesRepo) GetAllBackofficeNonPaginate(search string) *[]model.ResSales {
	data := new([]model.ResSales)
	qry := r.DB.Select("s.id,p.name,s.start_date,s.end_date," +
		"s.total_slot,s.available_slot,s.point,s.parking_space_id," +
		"s.created_at,s.updated_at").Joins("JOIN " +
		"parking_space_sales s on p.id = s.parking_space_id")
	if search != "" {
		qry = qry.Where("p.deleted_at is null and s.deleted_at is null "+
			"and (p.name ilike ? or p.address ilike ?)", "%"+search+"%", "%"+search+"%").
			Table("parking_space p")
	} else {
		qry = qry.Where("p.deleted_at is null and s.deleted_at is null").
			Table("parking_space p")
	}

	qry = qry.Scan(&data)
	return data
}

// GetAllList ..
func (r *SalesRepo) GetAllList() *[]entity.ParkingSpaceSales {
	data := new([]entity.ParkingSpaceSales)
	r.DB.Find(&data)
	return data
}

// GetSalesBySpace ...
func (r *SalesRepo) GetSalesBySpace(salesID int64, startDate string, endDate string) *[]entity.ParkingSpaceSales {
	data := new([]entity.ParkingSpaceSales)
	r.DB.Where("parking_space_id = ? AND start_date >= ? AND end_date <= ?", salesID, startDate, endDate).Find(&data)
	return data
}

// GetSalesIDByPSpaceID ..
func (r *SalesRepo) GetSalesIDByPSpaceID(id int64, usersID int64) []int64 {
	data := new([]entity.ParkingSpaceSales)
	var salesID []int64
	r.DB.Select("distinct pss.*").Joins("JOIN "+
		"parking_space_sales pss on ps.id = pss.parking_space_id").Joins("JOIN "+
		"trx_parking_space_sales trx on pss.id = trx.parking_space_sales_id").Joins("JOIN "+
		"merchant m on trx.merchant_id = m.id").Joins("JOIN "+
		"merchant_users mu on m.id = mu.merchant_id").Table("parking_space ps").
		Where("pss.parking_space_id = ? and ps.deleted_at is null and "+
			"pss.end_date >= ? and mu.users_id = ?", id, time.Now(), usersID).Scan(&data)
	if len(*data) > 0 {
		for _, v := range *data {
			salesID = append(salesID, v.ID)
		}
	}
	return salesID
}

// GetSalesNearby end user get parking space nearby
func (r *SalesRepo) GetSalesNearby(userLat, userLon, distance float64) *[]model.ResParkingSpace {
	data := make([]model.ResParkingSpace, 0)
	rows, _ := r.DB.Raw("select * from parkingspace_nearby(?,?,?) order by nearby asc", userLat, userLon, distance).Rows()
	for rows.Next() {
		var dat model.ResParkingSpace
		r.DB.ScanRows(rows, &dat)
		data = append(data, dat)
	}
	return &data
}
