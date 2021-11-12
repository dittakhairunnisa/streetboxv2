package parkingspace

import (
	"github.com/lib/pq"
	"streetbox.id/entity"
	"streetbox.id/model"
)

// RepoInterface ...
type RepoInterface interface {
	Create(*entity.ParkingSpace) error
	GetAll(limit, page int, sort []string) (*[]entity.ParkingSpace, int, int)
	GetOne(int64) *entity.ParkingSpace
	UpdateImagesMeta(pq.StringArray, int64) error
	UpdateDocsMeta(pq.StringArray, int64) error
	Update(*entity.ParkingSpace, int64) error
	DeleteByID(id int64) error
	GetAllList() *[]entity.ParkingSpace
}

// ServiceInterface ...
type ServiceInterface interface {
	GetAll(limit, page int, sort []string) model.Pagination
	GetAllList() *[]entity.ParkingSpace
	GetOne(int64) *entity.ParkingSpace
	Create(*model.ReqParkingSpaceCreate, int64) (*entity.ParkingSpace, error)
	UploadImage(pq.StringArray, int64) error
	UploadDoc(pq.StringArray, int64) error
	Update(*model.ReqParkingSpaceUpdate, int64, int64) (*entity.ParkingSpace, error)
	DeleteByID(id int64, usersID int64) error
}
