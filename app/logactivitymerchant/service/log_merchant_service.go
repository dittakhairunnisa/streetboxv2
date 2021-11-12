package service

import (
	"log"
	"time"

	"streetbox.id/app/logactivitymerchant"
	"streetbox.id/entity"
	"streetbox.id/model"
	"streetbox.id/util"
)

// LogMerchantService ..
type LogMerchantService struct {
	LogRepo logactivitymerchant.RepoInterface
}

// New ..
func New(logRepo logactivitymerchant.RepoInterface) logactivitymerchant.ServiceInterface {
	return &LogMerchantService{logRepo}
}

// Create ..
func (r *LogMerchantService) Create(merchantID int64, activity string) {
	data := new(entity.LogActivityMerchant)
	data.MerchantID = merchantID
	data.Activity = activity
	data.LogTime = time.Now()
	if err := r.LogRepo.Create(data); err != nil {
		log.Printf("ERROR: %s", err.Error())
	}
}

// GetAll ..
func (r *LogMerchantService) GetAll(
	limit, page int, sort []string, usersID int64) model.Pagination {
	data, count, offset := r.LogRepo.GetAll(limit, page, sort, usersID)
	totalPages := util.TotalPages(count, limit)
	model := model.Pagination{
		Data:         data,
		Limit:        limit,
		NextPage:     util.NextPage(page, totalPages),
		Offset:       offset,
		Page:         page,
		PrevPage:     util.PrevPage(page),
		TotalPages:   totalPages,
		TotalRecords: count,
	}
	return model
}
