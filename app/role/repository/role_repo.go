package repository

import (
	"log"

	"github.com/jinzhu/gorm"
	"streetbox.id/app/role"
	"streetbox.id/entity"
)

// RoleRepo ...
type RoleRepo struct {
	DB *gorm.DB
}

// New ...
func New(db *gorm.DB) role.RepoInterface {
	return &RoleRepo{db}
}

// Create ...
func (r *RoleRepo) Create(role *entity.Role) error {
	if err := r.DB.Create(role).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
		return err
	}
	log.Printf("INFO: Created Role: %+v", role)
	return nil
}

// FindByName ...
func (r *RoleRepo) FindByName(name string) *entity.Role {
	role := new(entity.Role)
	r.DB.Where("name = ?", name).First(role)
	return role
}

// DeleteByID ...
func (r *RoleRepo) DeleteByID(id int64) error {
	role := new(entity.Role)
	role.ID = id
	if err := r.DB.Where("id = ?", id).Delete(role).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
		return err
	}
	log.Printf("INFO: Deleted Role: %+v", role)
	return nil
}

// GetAll ...
func (r *RoleRepo) GetAll() *[]entity.Role {
	role := new([]entity.Role)
	r.DB.Find(&role)
	return role
}

// GetAllExclude ...
func (r *RoleRepo) GetAllExclude() *[]entity.Role {
	role := new([]entity.Role)
	r.DB.Where("name <> ?", "foodtruck").Find(&role)
	return role
}

// GetOne ...
func (r *RoleRepo) GetOne(ID int64) *entity.Role {
	role := new(entity.Role)
	r.DB.Where("id = ?", ID).Find(&role)
	return role
}
