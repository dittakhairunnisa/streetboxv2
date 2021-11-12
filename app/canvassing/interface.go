package canvassing

import (
	"streetbox.id/entity"
	"streetbox.id/model"
)

type Repository interface {
	CreateCanvas(canv *entity.Canvassing) (err error)
	GetCanvas(id int64) (canv entity.Canvassing, err error)
	GetFoodtruckCanvas(id, merchID int64) (canv model.RespCanvassing, err error)
	GetTriggeredCanvas() (canv []model.AutoCanvassing, err error)
	UpdateCanvas(canv entity.Canvassing) (err error)
	GetNearby(id int64, long, lat, rad float64) (users []entity.UserTokens, err error)
	GetFoodtruckLocation(id int64) (loc entity.FoodtruckLocation, err error)
	GetUserLocation(id int64) (loc entity.UsersLocation, err error)
	UpdateFoodtruckLocation(loc entity.FoodtruckLocation) (err error)
	UpdateUsersLocation(loc entity.UsersLocation) (err error)
	SaveNotification(notif entity.CanvassingNotif) (err error)
	GetNotificationByID(id int64) (notif entity.CanvassingNotif, err error)
	GetNotificationsByUserID(id int64) (notifs []model.RespNotifByUserID, err error)
	UpdateNotification(id int64, status string) (err error)
	CreateCall(call *entity.CanvassingCall) (foodtruck entity.UserTokens, err error)
	UpdateCall(id int64, status string, queue int) (call entity.CanvassingCall, err error)
	UpdateStatusCall(id int64, status string) (call entity.CanvassingCall, err error)
	FinishCall(id int64) (err error)
	GetCallsByFoodtruckID(id int64, status string) (calls []model.RespCallsByFoodtruckID, err error)
	GetCallsByCustomerID(id int64) (calls []model.RespCallsByUserID, err error)
	UpdateFoodtruckBlast(id int64) (err error)
	UpdateAutoBlast(id int64)
	ToggleAutoBlast(id int64)
	ExpireNotification()
	ExpireCall()
	ResetQueue()
}

type Service interface {
	CreateCanvas(canv *entity.Canvassing) (err error)
	GetCanvas(id int64) (canv entity.Canvassing, err error)
	GetFoodtruckCanvas(id, merchID int64) (canv model.RespCanvassing, err error)
	UpdateCanvas(canv entity.Canvassing) (err error)
	GetFoodtruckLocation(id int64) (loc entity.FoodtruckLocation, err error)
	GetUserLocation(id int64) (loc entity.UsersLocation, err error)
	UpdateFoodtruckLocation(loc entity.FoodtruckLocation) (err error)
	UpdateUsersLocation(loc entity.UsersLocation) (err error)
	Blast(id, merchID int64, auto bool) (users []entity.UserTokens, err error)
	GetTriggeredCanvas() (canv []model.AutoCanvassing, err error)
	SaveNotification(notif entity.CanvassingNotif) (err error)
	CallFoodtruck(id int64) (call entity.CanvassingCall, foodtruck entity.UserTokens, err error)
	AnswerCall(callID, foodtruckID int64, status string) (call entity.CanvassingCall, token string, err error)
	UpdateStatusCall(callID, foodtruckID int64, status string) (call entity.CanvassingCall, token string, err error)
	FinishCall(id int64) (err error)
	GetNotifications(id int64) (notifs []model.RespNotifByUserID, err error)
	GetCallsByFoodtruckID(id int64, status string) (calls []model.RespCallsByFoodtruckID, err error)
	GetCallsByCustomerID(id int64) (calls []model.RespCallsByUserID, err error)
	ToggleAutoBlast(id int64)
	ExpireNotification()
	ExpireCall()
	ResetQueue()
}