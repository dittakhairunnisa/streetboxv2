package service

import (
	"fmt"

	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
	"streetbox.id/app/logactivity"
	"streetbox.id/app/parkingspace"
	"streetbox.id/app/sales"
	"streetbox.id/app/trxspaces"
	"streetbox.id/app/trxvisitsales"
	"streetbox.id/app/user"
	"streetbox.id/entity"
	"streetbox.id/model"
	"streetbox.id/util"
)

// TrxSalesService ...
type TrxSalesService struct {
	TrxSalesRepo     trxspaces.RepoInterface
	SalesRepo        sales.RepoInterface
	UserRepo         user.RepoInterface
	LogRepo          logactivity.RepoInterface
	ParkingSpaceRepo parkingspace.RepoInterface
	VisitSalesRepo   trxvisitsales.RepoInterface
}

// New ...
func New(trxSalesRepo trxspaces.RepoInterface,
	salesRepo sales.RepoInterface, userRepo user.RepoInterface,
	logRepo logactivity.RepoInterface,
	parkingSpaceRepo parkingspace.RepoInterface,
	visitSalesRepo trxvisitsales.RepoInterface) trxspaces.ServiceInterface {
	return &TrxSalesService{trxSalesRepo, salesRepo,
		userRepo, logRepo, parkingSpaceRepo, visitSalesRepo}
}

// CreateTrx ...
func (s *TrxSalesService) CreateTrx(req *model.ReqCreateTrxSales, UserID int64) error {
	data := new(entity.TrxParkingSpaceSales)
	copier.Copy(&data, req)
	var (
		trx *gorm.DB
		err error
	)
	if trx, err = s.TrxSalesRepo.Create(data); err != nil {
		return err
	}
	if trx, err = s.SalesRepo.UpdateAvailableSlot(trx, req.TotalSlot,
		req.ParkingSpaceSalesID); err != nil {
		return err
	}
	trx.Commit()
	parkingSpaceID := s.SalesRepo.GetOne(data.ParkingSpaceSalesID).ParkingSpaceID
	parkingSpace := s.ParkingSpaceRepo.GetOne(parkingSpaceID).Name
	userName := s.UserRepo.FindByID(UserID).UserName
	msg := fmt.Sprintf("Add New Transaction %s by %s", parkingSpace, userName)
	s.LogRepo.Create(msg)
	return nil
}

// UpdateTrx ...
func (s *TrxSalesService) UpdateTrx(req *model.ReqCreateTrxSales, id int64, UserID int64) error {
	data := new(entity.TrxParkingSpaceSales)
	copier.Copy(&data, req)
	var (
		trx *gorm.DB
		err error
	)
	if trx, err = s.TrxSalesRepo.Update(data, id); err != nil {
		return err
	}
	if trx, err = s.SalesRepo.UpdateAvailableSlot(trx, req.TotalSlot,
		req.ParkingSpaceSalesID); err != nil {
		return err
	}
	trx.Commit()
	parkingSpaceID := s.SalesRepo.GetOne(data.ParkingSpaceSalesID).ParkingSpaceID
	parkingSpace := s.ParkingSpaceRepo.GetOne(parkingSpaceID).Name
	userName := s.UserRepo.FindByID(UserID).UserName
	msg := fmt.Sprintf("Add New Transaction %s by %s", parkingSpace, userName)
	s.LogRepo.Create(msg)
	return nil
}

// GetByUserID ...
func (s *TrxSalesService) GetByUserID(id int64) (*[]entity.TrxParkingSpaceSales, error) {
	data := new([]entity.TrxParkingSpaceSales)
	var err error
	if data, err = s.TrxSalesRepo.GetByUserID(id); err != nil {
		return nil, err
	}
	return data, nil
}

// GetMyParking ...
func (s *TrxSalesService) GetMyParking(usersID, merchantID int64) *[]model.ResMyParkingList {
	data := make([]model.ResMyParkingList, 0)
	visit := s.VisitSalesRepo.FindByMerchantID(merchantID)
	for _, visitData := range *visit {
		dat := new(model.ResMyParkingList)
		dat.Address = visitData.Address
		dat.Description = visitData.Notes
		dat.EndTime = visitData.EndDate
		dat.ID = visitData.ID
		dat.ProfilePicture = visitData.ProfilePicture
		dat.Latitude = visitData.Latitude
		dat.Longitude = visitData.Longitude
		dat.Name = visitData.CustomerName
		dat.StartTime = visitData.StartDate
		dat.ImagesMeta = append([]string{}, visitData.ProfilePicture)
		dat.TrxVisitSalesID = visitData.TrxVisitSalesID
		data = append(data, *dat)
	}
	space, _ := s.TrxSalesRepo.GetMyParking(usersID)
	for _, spaceData := range *space {
		dat := new(model.ResMyParkingList)
		copier.Copy(&dat, spaceData)
		data = append(data, *dat)
	}
	return &data
}

// GetSlotMyParking ...
func (s *TrxSalesService) GetSlotMyParking(
	pspaceID, usersID int64) (*[]model.ResSlotMyParking, error) {
	salesID := s.SalesRepo.GetSalesIDByPSpaceID(pspaceID, usersID)
	return s.TrxSalesRepo.GetSlotMyParking(salesID, usersID)
}

// GetAll ...
func (s *TrxSalesService) GetAll() *[]model.ResTrxList {
	return s.TrxSalesRepo.GetAll()
}

// DeleteByID ...
func (s *TrxSalesService) DeleteByID(id int64) error {
	return s.TrxSalesRepo.DeleteByID(id)
}

// GetByID ..
func (s *TrxSalesService) GetByID(id int64) *model.ResTrxList {
	return s.TrxSalesRepo.GetByID(id)
}

// GetByMerchantIDAndParkingSalesID ..
func (s *TrxSalesService) GetByMerchantIDAndParkingSalesID(merchantID int64, parkingSpaceSalesID int64) (*entity.TrxParkingSpaceSales, error) {
	return s.TrxSalesRepo.GetByMerchantIDAndParkingSalesID(merchantID, parkingSpaceSalesID)
}

// GetList ...
func (s *TrxSalesService) GetList(limit, page int, sort []string, filter string) model.Pagination {
	data, count, offset := s.TrxSalesRepo.GetList(limit, page, sort, filter)
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
