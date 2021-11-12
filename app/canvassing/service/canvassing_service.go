package service

import (
	"strings"
	"sync"

	"streetbox.id/app/canvassing"
	"streetbox.id/entity"
	"streetbox.id/model"
)

type CanvassingService struct {
	repo  canvassing.Repository
	queue *sync.Map
}

func (c *CanvassingService) CreateCanvas(canv *entity.Canvassing) (err error) {
	return c.repo.CreateCanvas(canv)
}

func (c *CanvassingService) GetCanvas(id int64) (canv entity.Canvassing, err error) {
	return c.repo.GetCanvas(id)
}

func (c *CanvassingService) GetFoodtruckCanvas(id, merchID int64) (canv model.RespCanvassing, err error) {
	return c.repo.GetFoodtruckCanvas(id, merchID)
}

func (c *CanvassingService) UpdateCanvas(canv entity.Canvassing) (err error) {
	return c.repo.UpdateCanvas(canv)
}

func (c *CanvassingService) GetFoodtruckLocation(id int64) (loc entity.FoodtruckLocation, err error) {
	return c.repo.GetFoodtruckLocation(id)
}

func (c *CanvassingService) GetUserLocation(id int64) (loc entity.UsersLocation, err error) {
	return c.repo.GetUserLocation(id)
}

func (c *CanvassingService) UpdateFoodtruckLocation(loc entity.FoodtruckLocation) (err error) {
	return c.repo.UpdateFoodtruckLocation(loc)
}

func (c *CanvassingService) UpdateUsersLocation(loc entity.UsersLocation) (err error) {
	return c.repo.UpdateUsersLocation(loc)
}

func (c *CanvassingService) Blast(id, merchID int64, auto bool) (users []entity.UserTokens, err error) {
	merchLoc, err := c.repo.GetFoodtruckLocation(id)
	if err != nil {
		return
	}
	canv, err := c.repo.GetCanvas(merchID)
	if err != nil {
		return
	}
	users, err = c.repo.GetNearby(id, merchLoc.Longitude, merchLoc.Latitude, canv.Radius)
	if err != nil {
		return
	}
	if auto {
		c.repo.UpdateAutoBlast(canv.ID)
	} else {
		c.repo.UpdateFoodtruckBlast(id)
	}
	return

}

func (c *CanvassingService) GetTriggeredCanvas() (canv []model.AutoCanvassing, err error) {
	return c.repo.GetTriggeredCanvas()
}

func (c *CanvassingService) SaveNotification(notif entity.CanvassingNotif) (err error) {
	return c.repo.SaveNotification(notif)
}

func (c *CanvassingService) GetNotifications(id int64) (notifs []model.RespNotifByUserID, err error) {
	return c.repo.GetNotificationsByUserID(id)
}

func (c *CanvassingService) CallFoodtruck(id int64) (call entity.CanvassingCall, foodtruck entity.UserTokens, err error) {
	notif, err := c.repo.GetNotificationByID(id)
	if err != nil {
		return
	}
	err = c.repo.UpdateNotification(id, "CALLING")
	if err != nil {
		return
	}
	call.NotifID = id
	call.CustomerID = notif.CustomerID
	call.FoodtruckID = notif.FoodtruckID
	foodtruck, err = c.repo.CreateCall(&call)
	if err != nil {
		return
	}
	return
}

func (c *CanvassingService) AnswerCall(callID, foodtruckID int64, status string) (call entity.CanvassingCall, token string, err error) {
	var queueNo int
	if status == "accept" {
		value, ok := c.queue.Load(foodtruckID)
		if !ok {
			queueNo = 1
		} else {
			queueNo = value.(int)
		}
		c.queue.Store(foodtruckID, queueNo+1)
	}

	status = strings.ToUpper(status)
	call, err = c.repo.UpdateCall(callID, status, queueNo)
	if err != nil {
		return
	}
	err = c.repo.UpdateNotification(call.NotifID, status)
	if err != nil {
		return
	}
	notif, err := c.repo.GetNotificationByID(call.NotifID)
	token = notif.CustomerToken
	return
}

func (c *CanvassingService) UpdateStatusCall(callID, foodtruckID int64, status string) (call entity.CanvassingCall, token string, err error) {
	status = strings.ToUpper(status)
	call, err = c.repo.UpdateStatusCall(callID, status)
	if err != nil {
		return
	}
	err = c.repo.UpdateNotification(call.NotifID, status)
	if err != nil {
		return
	}
	notif, err := c.repo.GetNotificationByID(call.NotifID)
	token = notif.CustomerToken
	return
}

func (c *CanvassingService) FinishCall(id int64) (err error) {
	return c.repo.FinishCall(id)
}

func (c *CanvassingService) GetCallsByFoodtruckID(id int64, status string) (calls []model.RespCallsByFoodtruckID, err error) {
	status = strings.ToUpper(status)
	return c.repo.GetCallsByFoodtruckID(id, status)
}

func (c *CanvassingService) GetCallsByCustomerID(id int64) (calls []model.RespCallsByUserID, err error) {
	return c.repo.GetCallsByCustomerID(id)
}

func (c *CanvassingService) ToggleAutoBlast(id int64) {
	c.repo.ToggleAutoBlast(id)
}

func (c *CanvassingService) ExpireNotification() {
	c.repo.ExpireNotification()
}

func (c *CanvassingService) ExpireCall() {
	c.repo.ExpireCall()
}

func (c *CanvassingService) ResetQueue() {
	c.queue.Range(func(key interface{}, value interface{}) bool {
		c.queue.Delete(key)
		return true
	})
	c.repo.ResetQueue()
}

func New(repo canvassing.Repository, queue *sync.Map) *CanvassingService {
	return &CanvassingService{
		repo:  repo,
		queue: queue,
	}
}
