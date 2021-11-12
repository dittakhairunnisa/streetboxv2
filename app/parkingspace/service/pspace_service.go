package service

import (
	"fmt"

	"github.com/jinzhu/copier"
	"github.com/lib/pq"
	"streetbox.id/app/logactivity"
	"streetbox.id/app/parkingspace"
	"streetbox.id/app/user"
	"streetbox.id/entity"
	"streetbox.id/model"
	"streetbox.id/util"
)

// ParkingSpaceService ...
type ParkingSpaceService struct {
	Repo     parkingspace.RepoInterface
	LogRepo  logactivity.RepoInterface
	UserRepo user.RepoInterface
}

// New ...
func New(s parkingspace.RepoInterface,
	log logactivity.RepoInterface, user user.RepoInterface) parkingspace.ServiceInterface {
	return &ParkingSpaceService{s, log, user}
}

// Create by recruit apps
func (s *ParkingSpaceService) Create(
	req *model.ReqParkingSpaceCreate, userID int64) (*entity.ParkingSpace, error) {
	data := new(entity.ParkingSpace)
	copier.Copy(&data, req)
	if err := s.Repo.Create(data); err != nil {
		return nil, err
	}
	userName := s.UserRepo.FindByID(userID).UserName
	msg := fmt.Sprintf("Add New Parking Space %s by %s", data.Name, userName)
	s.LogRepo.Create(msg)
	return data, nil
}

// UploadImage save 5 image meta data
func (s *ParkingSpaceService) UploadImage(img pq.StringArray, id int64) error {
	return s.Repo.UpdateImagesMeta(img, id)
}

// GetAll ...
func (s *ParkingSpaceService) GetAll(limit, page int, sort []string) model.Pagination {
	data, count, offset := s.Repo.GetAll(limit, page, sort)
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

// GetOne ...
func (s *ParkingSpaceService) GetOne(id int64) *entity.ParkingSpace {
	return s.Repo.GetOne(id)
}

// Update ...
func (s *ParkingSpaceService) Update(
	req *model.ReqParkingSpaceUpdate,
	id int64, userID int64) (*entity.ParkingSpace, error) {
	data := new(entity.ParkingSpace)
	copier.Copy(&data, req)
	if err := s.Repo.Update(data, id); err != nil {
		return nil, err
	}
	userName := s.UserRepo.FindByID(userID).UserName
	msg := fmt.Sprintf("Edit Parking Space %s by %s", data.Name, userName)
	s.LogRepo.Create(msg)
	return data, nil
}

// UploadDoc ...
func (s *ParkingSpaceService) UploadDoc(doc pq.StringArray, id int64) error {
	return s.Repo.UpdateDocsMeta(doc, id)
}

// DeleteByID ...
func (s *ParkingSpaceService) DeleteByID(id int64, usersID int64) error {
	var err error
	parkingSpace := s.GetOne(id).Name
	if err = s.Repo.DeleteByID(id); err != nil {
		return err
	}
	userName := s.UserRepo.FindByID(usersID).UserName
	msg := fmt.Sprintf("Delete Parking Space %s by %s", parkingSpace, userName)
	s.LogRepo.Create(msg)
	return nil
}

// GetAllList ..
func (s *ParkingSpaceService) GetAllList() *[]entity.ParkingSpace {
	return s.Repo.GetAllList()
}
