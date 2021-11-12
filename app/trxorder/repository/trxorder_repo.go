package repository

import (
	"fmt"
	"log"
	"strconv"

	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
	"streetbox.id/app/trxorder"
	"streetbox.id/entity"
	"streetbox.id/model"
	"streetbox.id/util"
)

// TrxOrderRepo ...
type TrxOrderRepo struct {
	DB *gorm.DB
}

// New ...
func New(db *gorm.DB) trxorder.RepoInterface {
	return &TrxOrderRepo{db}
}

// FindByTrxID ...
func (r *TrxOrderRepo) FindByTrxID(trxID string) *model.TrxOrderMerchant {
	data := new(model.TrxOrderMerchant)
	r.DB.Select("tro.*,m.name as merchant_name ,m.logo as merchant_logo, t.status, t.address").Joins("JOIN "+
		"trx_order tro on t.id = tro.trx_id").Joins("JOIN "+
		"merchant_users mu on tro.merchant_users_id = mu.id").Joins("JOIN "+
		"merchant m on mu.merchant_id = m.id").
		Table("trx t").Where("tro.deleted_at is null and "+
		"tro.trx_id = ?", trxID).Scan(&data)
	return data
}

// Create ..
func (r *TrxOrderRepo) Create(trxID string, data *model.TrxOrder, db *gorm.DB) (int64, *gorm.DB, error) {
	trxOrder := new(entity.TrxOrder)
	copier.Copy(&trxOrder, data)
	updatedAt := util.MillisToTime(data.UpdatedAt).Local()
	trxOrder.ID = 0
	trxOrder.IsClose = true
	trxOrder.BusinessDate = util.MillisToTime(data.BusinessDate).Local()
	trxOrder.CreatedAt = util.MillisToTime(data.CreatedAt).Local()
	trxOrder.UpdatedAt = &updatedAt
	trxOrder.TrxID = trxID
	if err := db.Create(&trxOrder).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
		db.Rollback()
		return 0, nil, err
	}
	return trxOrder.ID, db, nil
}

// CreateOnline from end user apps
func (r *TrxOrderRepo) CreateOnline(req *model.ReqTrxOrderOnline) (*gorm.DB, int64, error) {
	orderno := "0001"
	trxOrder := new(entity.TrxOrder)
	db := r.DB.Begin()
	copier.Copy(&trxOrder, req.Order)
	trxOrder.BusinessDate = util.MillisToTime(req.Order.BusinessDate)

	fmt.Println("current date : ", trxOrder.BusinessDate)

	if trxOrder.TypeOrder == "Online" {
		result := new(entity.TrxOrder)
		if err := db.Select("t.*").Table("trx_order t").Joins("JOIN merchant_users ms on ms.id = t.merchant_users_id").Where("t.business_date = ? "+
			" AND t.type_order = ? AND ms.merchant_id = ?", trxOrder.BusinessDate, "Online", req.Order.MerchantID).Order("t.id desc").Limit("1").Scan(&result).Error; err != nil {
			log.Printf("ERROR: %s", err.Error())
		}

		fmt.Println(*result)
		if result.ID > 0 {
			previousOrderNo, _ := strconv.Atoi(result.OrderNo)
			previousOrderNo++

			if previousOrderNo < 10 {
				orderno = "000" + strconv.Itoa(previousOrderNo)
			}

			if previousOrderNo < 100 && previousOrderNo >= 10 {
				orderno = "00" + strconv.Itoa(previousOrderNo)
			}

			if previousOrderNo < 1000 && previousOrderNo >= 100 {
				orderno = "0" + strconv.Itoa(previousOrderNo)
			}

			if previousOrderNo < 10000 && previousOrderNo >= 1000 {
				orderno = strconv.Itoa(previousOrderNo)
			}
		}
	}

	copier.Copy(&trxOrder, req.Order)
	trxOrder.BusinessDate = util.MillisToTime(req.Order.BusinessDate)
	updatedAt := util.MillisToTime(req.Order.UpdatedAt)
	trxOrder.CreatedAt = util.MillisToTime(req.Order.CreatedAt)
	trxOrder.UpdatedAt = &updatedAt
	trxOrder.TrxID = req.TrxID
	trxOrder.IsClose = false
	trxOrder.ID = 0
	trxOrder.OrderNo = orderno

	if err := db.Create(trxOrder).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
		db.Rollback()
		return nil, 0, err
	}
	log.Printf("INFO: Created TrxOrder: %+v", trxOrder)
	return db, trxOrder.ID, nil
}

//

// UpdateByTrxID ..
func (r *TrxOrderRepo) UpdateByTrxID(data *entity.TrxOrder, trxID string) *entity.TrxOrder {
	if err := r.DB.Model(&entity.TrxOrder{}).Where("trx_id = ?", trxID).
		Updates(data).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
		return nil
	}
	resp := new(entity.TrxOrder)
	r.DB.Find(&resp, "trx_id = ?", trxID)
	return resp
}

// FindAll ..
func (r *TrxOrderRepo) FindAll(model *entity.TrxOrder) *[]entity.TrxOrder {
	data := new([]entity.TrxOrder)
	r.DB.Where(model).Find(&data)
	return data
}

// FindOpenByMerchantUsersID order online
func (r *TrxOrderRepo) FindOpenByMerchantUsersID(merchantUserID int64) *[]entity.TrxOrder {
	data := new([]entity.TrxOrder)
	r.DB.Select("tro.*").Joins("JOIN "+
		"trx_order tro on t.id = tro.trx_id").Table("trx t").
		Where("t.deleted_at is null and tro.deleted_at is null and t.status = ? "+
			"and merchant_users_id = ? and is_close is false and tro.types = 1",
			util.TrxStatusSuccess, merchantUserID).Scan(&data)
	return data
}
