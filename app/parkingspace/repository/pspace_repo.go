package repository

import (
	"log"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"streetbox.id/app/parkingspace"
	"streetbox.id/entity"
	"streetbox.id/util"
)

// PSpaceRepo ...
type PSpaceRepo struct {
	DB *gorm.DB
}

// New ...
func New(db *gorm.DB) parkingspace.RepoInterface {
	return &PSpaceRepo{db}
}

// Create from recruit apps
func (r *PSpaceRepo) Create(pspace *entity.ParkingSpace) error {
	if err := r.DB.Create(pspace).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
		return err
	}
	log.Printf("INFO: Created ParkingSpace: %+v", pspace)
	return nil
}

// GetAll ...
func (r *PSpaceRepo) GetAll(
	limit, page int, sort []string) (*[]entity.ParkingSpace, int, int) {
	data := new([]entity.ParkingSpace)
	count := 0
	offset := util.Offset(page, limit)
	qry := r.DB.Limit(limit).Offset(offset)
	// sorting
	if len(sort) > 0 {
		for _, o := range sort {
			qry = qry.Order(o)
		}
	}
	qry = qry.Find(&data)
	qry = qry.Offset(0).Count(&count)
	return data, count, offset
}

// UpdateImagesMeta ...
func (r *PSpaceRepo) UpdateImagesMeta(img pq.StringArray, id int64) error {
	data := new(entity.ParkingSpace)
	data.ImagesMeta = img
	if err := r.DB.Model(&entity.ParkingSpace{ID: id}).Updates(data).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
		return err
	}
	log.Printf("INFO: Updated ImagesMeta: %+v", data)
	return nil
}

// Update ...
func (r *PSpaceRepo) Update(ps *entity.ParkingSpace, id int64) error {
	if err := r.DB.Model(&entity.ParkingSpace{ID: id}).
		Updates(ps).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
		return err
	}
	log.Printf("INFO: Updated ParkingSpace: %+v", ps)
	return nil
}

// GetOne ...
func (r *PSpaceRepo) GetOne(id int64) *entity.ParkingSpace {
	data := new(entity.ParkingSpace)
	r.DB.Find(&data, "id = ?", id)
	if data.ID == 0 {
		return nil
	}
	return data
}

// UpdateDocsMeta ...
func (r *PSpaceRepo) UpdateDocsMeta(doc pq.StringArray, id int64) error {
	data := new(entity.ParkingSpace)
	data.DocumentsMeta = doc
	if err := r.DB.Model(&entity.ParkingSpace{ID: id}).Updates(data).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
		return err
	}
	log.Printf("INFO: Updated DocsMeta: %+v", data)
	return nil
}

// DeleteByID ...
func (r *PSpaceRepo) DeleteByID(id int64) error {
	data := new(entity.ParkingSpace)
	if err := r.DB.Delete(&data, "id = ?", id).
		Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
		return err
	}
	log.Printf("INFO: Deleted ParkingSpace: %+v", data)
	return nil
}

// GetAllList ..
func (r *PSpaceRepo) GetAllList() *[]entity.ParkingSpace {
	data := new([]entity.ParkingSpace)
	if err := r.DB.Find(&data, "end_contract >= ?", time.Now()).
		Order("name asc").Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
		return nil
	}
	return data
}
