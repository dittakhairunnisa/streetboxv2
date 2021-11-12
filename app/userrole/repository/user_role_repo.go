package repository

import (
	"log"

	"github.com/jinzhu/gorm"
	"streetbox.id/app/userrole"
	"streetbox.id/entity"
)

// UserRoleRepo ...
type UserRoleRepo struct {
	DB *gorm.DB
}

// New ...
func New(db *gorm.DB) userrole.RepoInterface {
	return &UserRoleRepo{db}
}

// Create ...
func (repo *UserRoleRepo) Create(trx *gorm.DB, userRole *entity.UsersRole) error {
	if err := trx.Create(userRole).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
		trx.Rollback()
		return err
	}
	log.Printf("INFO: Created UsersRole: %+v", userRole)
	return nil
}

// DeleteByID ...
func (repo *UserRoleRepo) DeleteByID(trx *gorm.DB, id int64) error {
	userRole := new(entity.UsersRole)
	userRole.ID = id
	if err := trx.Delete(&userRole).Error; err != nil {
		trx.Rollback()
		log.Printf("ERROR: %s", err.Error())
		return err
	}
	log.Printf("INFO: Deleted UsersRole: %+v", userRole)
	return nil
}

// DeleteByMultipleID ...
func (repo *UserRoleRepo) DeleteByMultipleID(trx *gorm.DB, id []int64) error {
	userRole := new([]entity.UsersRole)
	if err := trx.Where("users_id IN (?)", id).Delete(&userRole).Error; err != nil {
		trx.Rollback()
		log.Printf("ERROR: %s", err.Error())
		return err
	}
	log.Printf("INFO: Deleted UsersRole: %+v", userRole)
	return nil
}

// GetNameByUserID ...
func (repo *UserRoleRepo) GetNameByUserID(id int64) string {
	role := new(entity.Role)
	if err := repo.DB.Joins("JOIN users_role on role.id = users_role.role_id").
		Where("users_role.users_id = ?", id).First(&role).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
		return ""
	}
	return role.Name
}

// Update ...
func (repo *UserRoleRepo) Update(usersID, roleID int64) error {
	if err := repo.DB.Model(new(entity.UsersRole)).
		Where("users_id = ?", usersID).Update("role_id", roleID).
		Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
		return err
	}
	log.Printf("INFO: Updated UserRole: "+
		"usersID = %d and roleID = %d", usersID, roleID)
	return nil
}
