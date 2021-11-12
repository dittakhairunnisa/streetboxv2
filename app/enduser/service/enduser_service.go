package service

import (
	"fmt"

	"github.com/jinzhu/copier"
	"streetbox.id/app/enduser"
	"streetbox.id/app/homevisitsales"
	"streetbox.id/app/merchant"
	"streetbox.id/app/paymentmethod"
	"streetbox.id/app/sales"
	"streetbox.id/app/tasksregular"
	"streetbox.id/app/taskstracking"
	"streetbox.id/app/trxspaces"
	"streetbox.id/app/user"
	"streetbox.id/entity"
	"streetbox.id/model"
	"streetbox.id/util"
)

// EndUserService ..
type EndUserService struct {
	MerchantRepo      merchant.RepoInterface
	SpaceSalesRepo    sales.RepoInterface
	TrackingRepo      taskstracking.RepoInterface
	VisitSalesRepo    homevisitsales.RepoInterface
	PaymentMethodRepo paymentmethod.RepoInterface
	TrxSalesRepo      trxspaces.RepoInterface
	UsersRepo         user.RepoInterface
	TasksRegRepo      tasksregular.RepoInterface
}

// New ..
func New(merchantRepo merchant.RepoInterface,
	spaceSalesRepo sales.RepoInterface,
	tracking taskstracking.RepoInterface,
	visitSalesRepo homevisitsales.RepoInterface,
	payMethodRepo paymentmethod.RepoInterface,
	trxSalesRepo trxspaces.RepoInterface,
	usersRepo user.RepoInterface,
	tasksRegRepo tasksregular.RepoInterface,
) enduser.ServiceInterface {
	return &EndUserService{merchantRepo, spaceSalesRepo,
		tracking, visitSalesRepo,
		payMethodRepo, trxSalesRepo, usersRepo, tasksRegRepo}
}

// GetNearby ..
func (s *EndUserService) GetNearby(
	limit, page int, distance float64, req *model.ReqMerchantNearby) model.Pagination {
	lat := req.Latitude
	lon := req.Longitude
	data, count, offset := s.MerchantRepo.GetNearby(limit, page, lat, lon, distance)
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

// MapParkingSpace get nearby parking space which have sales
func (s *EndUserService) MapParkingSpace(userLat, userLon, distance float64) *[]model.ResParkingSpace {
	return s.SpaceSalesRepo.GetSalesNearby(userLat, userLon, distance)
}

// GetLiveTracking ..
func (s *EndUserService) GetLiveTracking(lat, lon, distance float64) *[]model.ResLiveTracking {
	return s.TrackingRepo.GetLiveTracking(lat, lon, distance)
}

// VisitSalesList Show Home Visit Sales List
func (s *EndUserService) VisitSalesList(limit, page int) model.Pagination {
	data, count, offset := s.VisitSalesRepo.GetAllList(limit, page)
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

// GetPaymentMethod ...
func (s *EndUserService) GetPaymentMethod() *[]model.ResPaymentMethod {
	return s.PaymentMethodRepo.FindByActive()
}

// MapParkingSpaceDetail ..
func (s *EndUserService) MapParkingSpaceDetail(id int64) *[]model.ResParkingSpaceDetail {
	data := make([]model.ResParkingSpaceDetail, 0)
	merchant := s.TrxSalesRepo.GetMerchantBySpaceID(id)
	for _, merchantChild := range *merchant {
		detail := new(model.ResParkingSpaceDetail)
		copier.Copy(&detail, merchantChild)
		schedules := s.TrxSalesRepo.GetMerchantSchedules(merchantChild.MerchantID, id)
		for _, scheduleChild := range *schedules {
			schedule := new(model.Schedules)
			copier.Copy(&schedule, scheduleChild)
			detail.Schedules = append(detail.Schedules, *schedule)
		}
		if len(detail.Schedules) > 0 {
			salesID := detail.Schedules[0].ID
			if count := s.TasksRegRepo.CountBySalesID(salesID); count > 0 {
				detail.IsCheckin = true
			}
		}
		data = append(data, *detail)
	}
	return &data
}

// RegistrationToken stored registration token fcm sdk
func (s *EndUserService) RegistrationToken(token string, usersID int64) error {
	data := new(entity.Users)
	data.RegistrationToken = token
	return s.UsersRepo.Update(data, usersID)
}

// GetProfile ..
func (s *EndUserService) GetProfile(id int64) *entity.Users {
	return s.UsersRepo.FindEndUserByID(id)
}

// UpdateEndUser ..
func (s *EndUserService) UpdateEndUser(data model.ReqEndUserUpdate) error {
	findByUsername := s.UsersRepo.FindByUsername(data.Email)
	if findByUsername == nil {
		return fmt.Errorf("User not found for %s", data.Email)
	}
	user := new(entity.Users)
	copier.Copy(user, data)
	user.UserName = data.Email

	return s.UsersRepo.Update(user, findByUsername.ID)
}

// UpdateEndUserWithImage ..
func (s *EndUserService) UpdateEndUserWithImage(data model.ReqEndUserUpdateWithImage) error {
	findByUsername := s.UsersRepo.FindByUsername(data.Email)
	if findByUsername == nil {
		return fmt.Errorf("User not found for %s", data.Email)
	}
	user := new(entity.Users)
	copier.Copy(user, data)
	user.UserName = data.Email
	user.ProfilePicture = data.PhotoProfile
	return s.UsersRepo.Update(user, findByUsername.ID)
}

// GetSchedulesByTypesID method show schedules by tasks regular ID
func (s *EndUserService) GetSchedulesByTypesID(id int64) *[]model.Schedules {
	return s.TrxSalesRepo.GetSchedulesByTypesID(id)
}
