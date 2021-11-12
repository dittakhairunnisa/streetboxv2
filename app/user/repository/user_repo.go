package repository

import (
	"log"

	"github.com/jinzhu/gorm"
	"streetbox.id/app/user"
	"streetbox.id/entity"
	"streetbox.id/model"
)

// UserRepo ...
type UserRepo struct {
	DB *gorm.DB
}

// New ... init UserRepo
func New(db *gorm.DB) user.RepoInterface {
	return &UserRepo{db}
}

// Create ...
func (repo *UserRepo) Create(user *entity.Users) (*gorm.DB, error) {
	trx := repo.DB.Begin()
	if err := trx.Create(user).Error; err != nil {
		trx.Rollback()
		log.Printf("ERROR: %s", err.Error())
		return trx, err
	}
	log.Printf("INFO: Created User: %+v ", user)
	return trx, nil
}

// GetAllPagination ...
// func (repo *UserRepo) GetAllPagination(offset, perPage int, sort string) ([]*entity.Users, int, error) {

// }

// FindByID ...
func (repo *UserRepo) FindByID(id int64) *model.ResUserAll {
	user := new(model.ResUserAll)
	if err := repo.DB.Select("u.id, u.user_name, "+
		"u.name, u.phone, "+
		"u.address, r.name as role_name, ur.role_id, u.created_at, "+
		"u.updated_at, u.plat_no, u.profile_picture").Joins(
		"JOIN users_role ur on u.id = ur.users_id").
		Joins("JOIN role r on ur.role_id = r.id").
		Where("u.deleted_at is null and u.id = ?", id).
		Table("users u").Scan(&user).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
		return nil
	}
	return user
}

// FindUsernameByMultipleID ...
func (repo *UserRepo) FindUsernameByMultipleID(id []int64) *[]string {
	var username []string
	user := new(entity.Users)
	if err := repo.DB.Where("deleted_at is null and id IN (?)", id).
		Find(&user).Pluck("user_name", &username).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
		return nil
	}
	return &username
}

// GetAll ...
func (repo *UserRepo) GetAll(filter string) *[]model.ResUserAll {
	user := new([]model.ResUserAll)
	filterQuery := ""
	if filter != "" {
		filterQuery = " AND r.name = ?"
	}
	var err error
	if filterQuery == "" {
		err = repo.DB.Select("u.id, u.user_name, " +
			"u.name, u.phone, " +
			"u.address, r.name as role_name, ur.role_id, u.created_at, " +
			"u.updated_at").Joins(
			"JOIN users_role ur on u.id = ur.users_id").
			Joins("JOIN role r on ur.role_id = r.id").
			Where("u.deleted_at is null").
			Table("users u").Scan(&user).Error
	} else {
		err = repo.DB.Select("u.id, u.user_name, "+
			"u.name, u.phone, "+
			"u.address, r.name as role_name, ur.role_id, u.created_at, "+
			"u.updated_at").Joins(
			"JOIN users_role ur on u.id = ur.users_id").
			Joins("JOIN role r on ur.role_id = r.id").
			Where("u.deleted_at is null"+filterQuery, filter).
			Table("users u").Scan(&user).Error
	}
	if err != nil {
		log.Printf("ERROR: %s", err.Error())
	}
	return user
}

// DeleteByID ...
func (repo *UserRepo) DeleteByID(id int64) (*gorm.DB, error) {
	user := new(entity.Users)
	user.ID = id
	trx := repo.DB.Begin()
	if err := trx.Delete(user).Error; err != nil {
		trx.Rollback()
		log.Printf("ERROR: %s", err.Error())
		return nil, err
	}
	log.Printf("INFO: Deleted Users: %+v", user)
	return trx, nil
}

// DeleteByMultipleID ...
func (repo *UserRepo) DeleteByMultipleID(trx *gorm.DB, id []int64) (*gorm.DB, error) {
	user := new([]entity.Users)
	if err := trx.Where(id).Delete(&user).Error; err != nil {
		trx.Rollback()
		log.Printf("ERROR: %s", err.Error())
		return nil, err
	}
	log.Printf("INFO: Deleted Users: %+v", user)
	return trx, nil
}

// FindByUsername ...
func (repo *UserRepo) FindByUsername(userName string) *entity.Users {
	user := new(entity.Users)
	repo.DB.Find(&user, "user_name = ?", userName)
	if user.ID == 0 {
		return nil
	}
	return user
}

// Update ...
func (repo *UserRepo) Update(user *entity.Users, id int64) error {
	if err := repo.DB.Model(&entity.Users{ID: id}).Updates(&user).
		Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
		return err
	}
	log.Printf("INFO: Updated Users: %+v", user)
	return nil
}

// GetUserAdmin ...
func (repo *UserRepo) GetUserAdmin() *[]model.ResUserMerchant {
	users := new([]model.ResUserMerchant)
	if err := repo.DB.Select(
		"u.user_name, m.name, mu.merchant_id, mu.users_id").Joins("JOIN "+
		"merchant_users mu on m.id = mu.merchant_id").Joins("JOIN "+
		"users u on mu.users_id = u.id").Joins("JOIN "+
		"users_role ur on u.id = ur.users_id").Joins("JOIN "+
		"role r on ur.role_id = r.id").
		Where("r.name = ? and m.deleted_at is null", "admin").
		Table("merchant m").Scan(&users).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
	}
	return users
}

// FindEndUserByID ..
func (repo *UserRepo) FindEndUserByID(id int64) *entity.Users {
	data := new(entity.Users)
	repo.DB.Find(&data, "id = ?", id)
	return data
}

// GetByMerchantUsersID ..
func (repo *UserRepo) GetByMerchantUsersID(id int64) *entity.Users {
	data := new(entity.Users)
	if err := repo.DB.Select("u.*").Joins("JOIN "+
		"merchant_users mu on u.id = mu.users_id").Table("users u").
		Where("mu.id = ? ", id).Scan(&data).Error; err != nil {
		return nil
	}
	return data
}

// GetFoodtruckByTrxVisitID ..
func (repo *UserRepo) GetFoodtruckByTrxVisitID(id int64) *entity.Users {
	data := new(entity.Users)
	if err := repo.DB.Select("u.*").Joins("JOIN "+
		"trx_homevisit_sales ths on tv.id = ths.trx_visit_id").Joins("JOIN "+
		"tasks_homevisit tkh on ths.id = tkh.trx_homevisit_sales_id").Joins("JOIN "+
		"tasks tk on tkh.tasks_id = tk.id").Joins("JOIN "+
		"merchant_users mu on tk.merchant_users_id = mu.id").Joins("JOIN "+
		"users u on mu.users_id = u.id").
		Table("trx_visit tv").Where("tv.id = ?", id).Scan(&data).Error; err != nil {
		return nil
	}
	return data
}
