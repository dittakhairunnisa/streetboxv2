package service

import (
	"github.com/jinzhu/gorm"
	"streetbox.id/app/homevisitsales"
	"streetbox.id/app/trxvisitsales"
	"streetbox.id/entity"
	"streetbox.id/model"
	"streetbox.id/util"
)

// HomevisitService ..
type HomevisitService struct {
	HomevisitRepo     homevisitsales.RepoInterface
	TrxVisitSalesRepo trxvisitsales.RepoInterface
}

// New ..
func New(
	homeRepo homevisitsales.RepoInterface,
	trxVisitSalesRepo trxvisitsales.RepoInterface) homevisitsales.ServiceInterface {
	return &HomevisitService{homeRepo, trxVisitSalesRepo}
}

// GetAll ..
func (r *HomevisitService) GetAll(merchantID int64) *[]entity.HomevisitSales {
	return r.HomevisitRepo.GetAll(merchantID)
}

// Create ..
func (r *HomevisitService) Create(
	req *entity.HomevisitSales) (*entity.HomevisitSales, error) {
	return r.HomevisitRepo.Create(req)
}

// CheckDate ...
func (r *HomevisitService) CheckDate(date string, merchantID int64) int {
	return r.HomevisitRepo.CheckDate(date, merchantID)
}

// Update ..
func (r *HomevisitService) Update(data *entity.HomevisitSales) (*entity.HomevisitSales, error) {
	result, err := r.HomevisitRepo.Update(data)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetAllByMerchantAndDate ..
func (r *HomevisitService) GetAllByMerchantAndDate(
	merchantID int64, startDate string, endDate string) *[]entity.HomevisitSales {
	return r.HomevisitRepo.GetAllByMerchantAndDate(merchantID, startDate, endDate)
}

// GetInfoByDate ..
func (r *HomevisitService) GetInfoByDate(merchantID int64, date string) *model.ResHomeVisitGetInfo {
	return r.HomevisitRepo.GetInfoByDate(merchantID, date)
}

// GetByID ..
func (r *HomevisitService) GetByID(ID int64) *entity.HomevisitSales {
	return r.HomevisitRepo.GetByID(ID)
}

// DeleteByDate ..
func (r *HomevisitService) DeleteByDate(date string, merchantID int64) (*gorm.DB, error) {
	var (
		db  *gorm.DB
		err error
	)
	if db, err = r.HomevisitRepo.DeleteByDate(date, merchantID); err != nil {
		return nil, err
	}
	db.Commit()
	return db, nil
}

// DeleteByID ..
func (r *HomevisitService) DeleteByID(ID int64, merchantID int64) (*gorm.DB, error) {
	var (
		db  *gorm.DB
		err error
	)
	if db, err = r.HomevisitRepo.DeleteByID(ID, merchantID); err != nil {
		return nil, err
	}
	return db, nil
}

// GetAllEndUser ..
func (r *HomevisitService) GetAllEndUser(limit, page int) model.Pagination {
	data, count, offset := r.HomevisitRepo.GetAllList(limit, page)
	totalPages := util.TotalPages(count, limit)
	model := model.Pagination{
		Data:         data,
		Offset:       offset,
		TotalPages:   totalPages,
		Limit:        limit,
		NextPage:     util.NextPage(page, totalPages),
		Page:         page,
		PrevPage:     util.PrevPage(page),
		TotalRecords: count,
	}
	return model
}

// GetAvailableByMerchantID ...
func (r *HomevisitService) GetAvailableByMerchantID(id int64) *[]model.ResVisitSalesDetail {
	return r.HomevisitRepo.GetAvailableByMerchantID(id)
}

// UpdateByTrxID update available
func (r *HomevisitService) UpdateByTrxID(id string) error {
	trxVisitSales := r.TrxVisitSalesRepo.FindByTrxID(id)
	if len(*trxVisitSales) > 0 {
		for _, trxData := range *trxVisitSales {
			updateAvailable := trxData.Available - 1
			if err := r.HomevisitRepo.UpdateAvailableByID(updateAvailable,
				trxData.HomevisitSalesID); err != nil {
				return err
			}
		}
	}
	return nil
}
