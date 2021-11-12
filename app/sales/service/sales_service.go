package service

import (
	"fmt"

	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
	"streetbox.id/app/logactivity"
	"streetbox.id/app/parkingspace"
	"streetbox.id/app/sales"
	"streetbox.id/app/user"
	"streetbox.id/entity"
	"streetbox.id/model"
	"streetbox.id/util"
)

// SalesService ...
type SalesService struct {
	SalesRepo        sales.RepoInterface
	LogRepo          logactivity.RepoInterface
	UsersRepo        user.RepoInterface
	ParkingSpaceRepo parkingspace.RepoInterface
}

// New ...
func New(repo sales.RepoInterface,
	log logactivity.RepoInterface,
	userRepo user.RepoInterface, parkingSpaceRepo parkingspace.RepoInterface) sales.ServiceInterface {
	return &SalesService{repo, log, userRepo, parkingSpaceRepo}
}

// CreateSales ...
func (s *SalesService) CreateSales(
	req *model.ReqSalesCreate, userID int64) (*entity.ParkingSpaceSales, error) {
	data := new(entity.ParkingSpaceSales)
	copier.Copy(&data, req)
	data.AvailableSlot = data.TotalSlot
	if err := s.SalesRepo.Create(data); err != nil {
		return nil, err
	}
	userName := s.UsersRepo.FindByID(userID).UserName
	parkingSpace := s.ParkingSpaceRepo.GetOne(data.ParkingSpaceID).Name
	msg := fmt.Sprintf("Add New Parking Space Sales %s by %s", parkingSpace, userName)
	s.LogRepo.Create(msg)
	return data, nil
}

// GetBySpaceID ...
func (s *SalesService) GetBySpaceID(id int64,
	limit, page int, sort []string) model.Pagination {
	data, count, offset := s.SalesRepo.FindBySpaceID(id, limit, page, sort)
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

// UpdateAvailableSlot ...
func (s *SalesService) UpdateAvailableSlot(
	db *gorm.DB, id int64, qty int) (*gorm.DB, error) {
	return s.SalesRepo.UpdateAvailableSlot(db, qty, id)
}

// GetByID ...
func (s *SalesService) GetByID(id int64) *entity.ParkingSpaceSales {
	return s.SalesRepo.GetOne(id)
}

// Update ...
func (s *SalesService) Update(
	req *model.ReqSalesUpdate, id int64, userID int64) (*entity.ParkingSpaceSales, error) {
	data := new(entity.ParkingSpaceSales)
	copier.Copy(&data, req)
	userName := s.UsersRepo.FindByID(userID).UserName
	parkingSpace := s.ParkingSpaceRepo.GetOne(data.ParkingSpaceID).Name
	msg := fmt.Sprintf("Edit Parking Space Sales %s by %s", parkingSpace, userName)
	s.LogRepo.Create(msg)
	return s.SalesRepo.Update(data, id)
}

// DeleteByID ...
func (s *SalesService) DeleteByID(id int64, userID int64) error {
	parkingSpace := s.ParkingSpaceRepo.GetOne(s.GetByID(id).ParkingSpaceID).Name
	userName := s.UsersRepo.FindByID(userID).UserName
	if err := s.SalesRepo.DeleteByID(id); err != nil {
		return err
	}
	msg := fmt.Sprintf("Delete Parking Space Sales %s by %s", parkingSpace, userName)
	s.LogRepo.Create(msg)
	return nil
}

// FindLikeName ..
func (s *SalesService) FindLikeName(name string,
	limit, page int, sort []string) model.Pagination {
	data, count, offset := s.SalesRepo.FindLikeName(name, limit, page, sort)
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

// FindLikeNameBackoffice ..
func (s *SalesService) FindLikeNameBackoffice(name string,
	limit, page int, sort []string) model.Pagination {
	data, count, offset := s.SalesRepo.FindLikeNameBackoffice(name, limit, page, sort)
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

// GetAll ...
func (s *SalesService) GetAll(limit, page int, sort []string) model.Pagination {
	data, count, offset := s.SalesRepo.GetAll(limit, page, sort)
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

// GetAllBackoffice ...
func (s *SalesService) GetAllBackoffice(limit, page int, sort []string) model.Pagination {
	data, count, offset := s.SalesRepo.GetAllBackoffice(limit, page, sort)
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

// GetAllListNonPaginate ...
func (s *SalesService) GetAllListNonPaginate(search string) *[]model.ResSales {
	return s.SalesRepo.GetAllBackofficeNonPaginate(search)
}

// GetAllList ...
func (s *SalesService) GetAllList() *[]entity.ParkingSpaceSales {
	return s.SalesRepo.GetAllList()
}

// GetSalesBySpace ...
func (s *SalesService) GetSalesBySpace(salesID int64, startDate string, endDate string) *[]entity.ParkingSpaceSales {
	return s.SalesRepo.GetSalesBySpace(salesID, startDate, endDate)
}
