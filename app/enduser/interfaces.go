package enduser

import (
	"streetbox.id/entity"
	"streetbox.id/model"
)

// ServiceInterface ..
type ServiceInterface interface {
	GetNearby(int, int, float64, *model.ReqMerchantNearby) model.Pagination
	MapParkingSpace(userLat, userLon, distance float64) *[]model.ResParkingSpace
	MapParkingSpaceDetail(int64) *[]model.ResParkingSpaceDetail
	GetLiveTracking(lat, lon, distance float64) *[]model.ResLiveTracking
	GetPaymentMethod() *[]model.ResPaymentMethod
	RegistrationToken(string, int64) error
	GetProfile(int64) *entity.Users
	UpdateEndUser(model.ReqEndUserUpdate) error
	UpdateEndUserWithImage(model.ReqEndUserUpdateWithImage) error
	GetSchedulesByTypesID(int64) *[]model.Schedules
}
