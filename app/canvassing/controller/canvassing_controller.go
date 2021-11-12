package controller

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"firebase.google.com/go/v4/messaging"
	"github.com/gin-gonic/gin"
	"streetbox.id/app/canvassing"
	"streetbox.id/app/fcm"
	"streetbox.id/app/merchant"
	"streetbox.id/app/user"
	"streetbox.id/entity"
	"streetbox.id/model"
	"streetbox.id/util"
)

type CanvassingController struct {
	Svc      canvassing.Service
	Fcm      fcm.ServiceInterface
	Merchant merchant.ServiceInterface
	User 		 user.ServiceInterface
}

// CreateCanvass godoc
// @Summary Create Canvassing Rule (permission = admin)
// @Id CreateCanvass
// @Tags Canvassing
// @Security Token
// @Param canvassing body entity.Canvassing true "merchant_id and last_auto_blast fields are not mandatory"
// @Success 200 {object} entity.Canvassing "{ "data": Model }"
// @Failure 422 {object} model.ResponseErrors "Invalid request"
// @Failure 500 {object} model.ResponseErrors "Create Canvassing Rule Failed"
// @Router /canvassing [post]
func (c *CanvassingController) CreateCanvas(ctx *gin.Context) {
	jwtModel, _ := util.ExtractTokenMetadata(ctx.Request)
	req := &entity.Canvassing{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		model.ResponseError(ctx, "Invalid Request", http.StatusUnprocessableEntity)
		return
	}
	merch := c.Merchant.GetMerchantUsersByUsersID(jwtModel.UserID)
	req.ID = merch.MerchantID
	if err := c.Svc.CreateCanvas(req); err != nil {
		model.ResponseError(ctx, "Create canvassing failed", http.StatusInternalServerError)
		return
	}
	model.ResponseJSON(ctx, req)
	return
}

// GetCanvas godoc
// @Summary Get Admin's Canvassing Rule (permission = admin)
// @Id GetCanvas
// @Tags Canvassing
// @Security Token
// @Success 200 {object} entity.Canvassing "{ "data": Model }"
// @Failure 500 {object} model.ResponseErrors "Inquiry Canvassing Rule Failed"
// @Router /canvassing [get]
func (c *CanvassingController) GetCanvas(ctx *gin.Context) {
	jwtModel, _ := util.ExtractTokenMetadata(ctx.Request)
	var canv entity.Canvassing
	var err error
	merch := c.Merchant.GetMerchantUsersByUsersID(jwtModel.UserID)
	if canv, err = c.Svc.GetCanvas(merch.MerchantID); err != nil {
		model.ResponseError(ctx, "Inquiry Canvassing Failed", http.StatusInternalServerError)
		return
	}
	model.ResponseJSON(ctx, canv)
	return
}

// GetCanvas godoc
// @Summary Get Foodtruck's Canvassing Rule (permission = merchant)
// @Id GetFoodtruckCanvas
// @Tags Canvassing
// @Security Token
// @Success 200 {object} model.RespCanvassing "{ "data": Model }"
// @Failure 500 {object} model.ResponseErrors "Inquiry Canvassing Rule Failed"
// @Router /canvassing/foodtruck [get]
func (c *CanvassingController) GetFoodtruckCanvas(ctx *gin.Context) {
	jwtModel, _ := util.ExtractTokenMetadata(ctx.Request)
	var canv model.RespCanvassing
	var err error
	merch := c.Merchant.GetMerchantUsersByUsersID(jwtModel.UserID)
	if canv, err = c.Svc.GetFoodtruckCanvas(merch.ID, merch.MerchantID); err != nil {
		model.ResponseError(ctx, "Inquiry Canvassing Failed", http.StatusInternalServerError)
		return
	}
	model.ResponseJSON(ctx, canv)
	return
}

// UpdateCanvass godoc
// @Summary Update Foodtruck's Canvassing Rule (permission = admin)
// @Id UpdateCanvass
// @Tags Canvassing
// @Security Token
// @Param canvassing body entity.Canvassing true "merchant_id and last_auto_blast fields are not mandatory"
// @Success 200 {object} entity.Canvassing "{ "data": Model }"
// @Failure 422 {object} model.ResponseErrors "Invalid request"
// @Failure 500 {object} model.ResponseErrors "Update Canvassing Rule Failed"
// @Router /canvassing [put]
func (c *CanvassingController) UpdateCanvas(ctx *gin.Context) {
	jwtModel, _ := util.ExtractTokenMetadata(ctx.Request)
	req := entity.Canvassing{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		model.ResponseError(ctx, "Invalid Request", http.StatusUnprocessableEntity)
		return
	}
	merch := c.Merchant.GetMerchantUsersByUsersID(jwtModel.UserID)
	req.ID = merch.MerchantID
	if err := c.Svc.UpdateCanvas(req); err != nil {
		model.ResponseError(ctx, "Update Canvassing Rule Failed", http.StatusUnprocessableEntity)
		return
	}
	model.ResponseJSON(ctx, gin.H{"message": "Success"})
	return
}

// GetFoodtruckLocation godoc
// @Summary Get Foodtruck's Location (permission = all)
// @Id GetFoodtruckLocation
// @Tags Canvassing
// @Security Token
// @Param id path string true "1"
// @Success 200 {object} entity.FoodtruckLocation "{ "data": Model }"
// @Failure 500 {object} model.ResponseErrors "Inquiry Foodtruck Location Failed"
// @Router /canvassing/foodtruck/location/:id [get]
func (c *CanvassingController) GetFoodtruckLocation(ctx *gin.Context) {
	id := util.ParamIDToInt64(ctx.Param("id"))
	var loc entity.FoodtruckLocation
	var err error
	if loc, err = c.Svc.GetFoodtruckLocation(id); err != nil {
		model.ResponseError(ctx, "Inquiry Foodtruck Location Failed", http.StatusInternalServerError)
		return
	}
	model.ResponseJSON(ctx, loc)
	return
}

// GetUserLocation godoc
// @Summary Get User's Location (permission = consumer)
// @Id GetUserLocation
// @Tags Canvassing
// @Security Token
// @Success 200 {object} entity.UsersLocation "{ "data": Model }"
// @Failure 500 {object} model.ResponseErrors "Inquiry User's Location Failed"
// @Router /canvassing/users/location/ [get]
func (c *CanvassingController) GetUserLocation(ctx *gin.Context) {
	jwtModel, _ := util.ExtractTokenMetadata(ctx.Request)
	var loc entity.UsersLocation
	var err error
	if loc, err = c.Svc.GetUserLocation(jwtModel.UserID); err != nil {
		model.ResponseError(ctx, "Inquiry User's Location Failed", http.StatusInternalServerError)
		return
	}
	model.ResponseJSON(ctx, loc)
	return
}

// UpdateFoodtruckLocation godoc
// @Summary Update Foodtruck's Location (permission = merchant)
// @Id UpdateFoodtruckLocation
// @Tags Canvassing
// @Security Token
// @Param location body entity.FoodtruckLocation true "id fields are not mandatory"
// @Success 200 {object} model.ResponseSuccess "{"message": "Success"}"
// @Failure 422 {object} model.ResponseErrors "Invalid request"
// @Failure 500 {object} model.ResponseErrors "Update Foodtruck Location Failed"
// @Router /canvassing/foodtruck/location [put]
func (c *CanvassingController) UpdateFoodtruckLocation(ctx *gin.Context) {
	req := entity.FoodtruckLocation{}
	jwtModel, _ := util.ExtractTokenMetadata(ctx.Request)
	if err := ctx.ShouldBindJSON(&req); err != nil {
		model.ResponseError(ctx, "Invalid Request", http.StatusUnprocessableEntity)
		return
	}
	merch := c.Merchant.GetMerchantUsersByUsersID(jwtModel.UserID)
	req.ID = merch.ID
	if err := c.Svc.UpdateFoodtruckLocation(req); err != nil {
		model.ResponseError(ctx, "Update Merchant Location Failed", http.StatusInternalServerError)
		return
	}
	model.ResponseJSON(ctx, gin.H{"message": "Success"})
	return
}

// UpdateUsersLocation godoc
// @Summary Update Users's Location (permission = consumer)
// @Id UpdateUsersLocation
// @Tags Canvassing
// @Security Token
// @Param location body entity.UsersLocation true "id fields are not mandatory"
// @Success 200 {object} model.ResponseSuccess "{"message": "Success"}"
// @Failure 422 {object} model.ResponseErrors "Invalid request"
// @Failure 500 {object} model.ResponseErrors "Update Users Location Failed"
// @Router /canvassing/users/location [put]
func (c *CanvassingController) UpdateUsersLocation(ctx *gin.Context) {
	jwtModel, _ := util.ExtractTokenMetadata(ctx.Request)
	req := entity.UsersLocation{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		model.ResponseError(ctx, "Invalid Request", http.StatusUnprocessableEntity)
		return
	}
	req.ID = jwtModel.UserID
	if err := c.Svc.UpdateUsersLocation(req); err != nil {
		model.ResponseError(ctx, "Update Users Location Failed", http.StatusInternalServerError)
		return
	}
	model.ResponseJSON(ctx, gin.H{"message": "Success"})
	return
}

// Blast godoc
// @Summary Blast Notification to Users Nearby Manually (permission = merchant)
// @Id Blast
// @Tags Canvassing
// @Security Token
// @Success 200 {object} model.ResponseSuccess "{"message": "Success"}"
// @Failure 422 {object} model.ResponseErrors "No users nearby"
// @Failure 500 {object} model.ResponseErrors "Blast Canvassing Failed"
// @Router /canvassing/blast [post]
func (c *CanvassingController) Blast(ctx *gin.Context) {
	jwtModel, _ := util.ExtractTokenMetadata(ctx.Request)
	merch := c.Merchant.GetMerchantUsersByUsersID(jwtModel.UserID)
	dataMerchant := c.Merchant.GetInfo(jwtModel.UserID)
	merchant := c.Merchant.GetByID(merch.MerchantID)
	users, err := c.Svc.Blast(merch.ID, merch.MerchantID, false)
	if err != nil {
		model.ResponseError(ctx, "Blast Canvassing Failed", http.StatusInternalServerError)
		return
	}
	if len(users) < 1 {
		model.ResponseError(ctx, "No users available", http.StatusUnprocessableEntity)
		return
	}
	msg := messaging.Notification{
		Title: "Foodtruck " + dataMerchant.Name,
		Body:  "Ada foodtruck keliling di dekat kamu, dipanggil yukk!",
		ImageURL: "https://api.streetbox.id/static/image/" + merchant.Logo,
	}
	data := map[string]string{
		"title": "Foodtruck " + dataMerchant.Name,
		"body":  "Ada foodtruck keliling di dekat kamu, dipanggil yukk!",
		// "image": "https://api.streetbox.id/static/image/" + merchant.Logo,
		"merchant_id": strconv.Itoa(int(merch.ID)),
	}
	for _, user := range users {
		c.Fcm.SendNotificationWithData(fmt.Sprintf("blast_%d", user.ID), &msg, &data, "id.streetbox.live")
		notif := entity.CanvassingNotif{
			FoodtruckID:   merch.ID,
			CustomerID:    user.ID,
			CustomerToken: user.RegistrationToken,
		}
		c.Svc.SaveNotification(notif)
	}

	model.ResponseJSON(ctx, gin.H{"message": "Success"})
	return
}

// CallFoodtruck godoc
// @Summary Call Foodtruck (permission = consumer)
// @Id CallFoodtruck
// @Tags Canvassing
// @Security Token
// @Param notif-id path string true "1"
// @Success 200 {object} model.ResponseSuccess "{"message": "Success"}"
// @Failure 500 {object} model.ResponseErrors "Calling Foodtruck Failed"
// @Router /canvassing/call/:notif-id [post]
func (c *CanvassingController) CallFoodtruck(ctx *gin.Context) {
	jwtModel, _ := util.ExtractTokenMetadata(ctx.Request)
	id := util.ParamIDToInt64(ctx.Param("notif-id"))
	call, foodtruck, err := c.Svc.CallFoodtruck(id)
	user := c.User.GetUserByID(jwtModel.UserID)
	if err != nil {
		model.ResponseError(ctx, "Calling Foodtruck Failed", http.StatusInternalServerError)
		return
	}
	msg := messaging.Notification{
		Title: "Incoming calll!",
		Body:  "There is customer wants you",
		ImageURL: "https://api.streetbox.id/static/image/" + user.ProfilePicture,
	}
	data := map[string]string{
		"title": "Incoming calll!",
		"body":  "There is customer wants you",
		// "image": "https://api.streetbox.id/static/image/" + user.ProfilePicture,
		"call_id": strconv.Itoa(int(call.ID)),
	}
	c.Fcm.SendNotificationWithData(fmt.Sprintf("blast_%d", foodtruck.ID), &msg, &data, "com.zeepos.streetbox")
	model.ResponseJSON(ctx, gin.H{"message": "Success"})
	return
}

// AnswerCall godoc
// @Summary Answer Call From Consumer (permission = merchant)
// @Id AnswerCall
// @Tags Canvassing
// @Security Token
// @Param call-id path string true "1"
// @Param status path string true "onprocess/accept/reject/expire"
// @Success 200 {object} entity.CanvassingCall "{ "data": Model }"
// @Failure 500 {object} model.ResponseErrors "Answer Consumer Call Failed"
// @Router /canvassing/call/:call-id/:status [put]
func (c *CanvassingController) AnswerCall(ctx *gin.Context) {
	callID := util.ParamIDToInt64(ctx.Param("call-id"))
	status := ctx.Param("status")
	jwtModel, _ := util.ExtractTokenMetadata(ctx.Request)
	foodtruck := c.Merchant.GetMerchantUsersByUsersID(jwtModel.UserID)
	merchant := c.Merchant.GetByID(foodtruck.MerchantID)
	call, token, err := c.Svc.AnswerCall(callID, foodtruck.ID, status)
	if err != nil {
		model.ResponseError(ctx, "Answer Consumer Call Failed", http.StatusInternalServerError)
		return
	}
	msg := messaging.Notification{
		Title: "Call Answer!",
		Body:  fmt.Sprintf("Your call has been %sed", status),
		ImageURL: "https://api.streetbox.id/static/image/" + merchant.Logo,
	}
	var tokens []string
	tokens = append(tokens, foodtruck.RegistrationToken)
	data := map[string]string{
		"title": "Call Answer!",
		"body":  fmt.Sprintf("Your call has been %sed", status),
		// "image": "https://api.streetbox.id/static/image/" + merchant.Logo,
		"call_id": strconv.Itoa(int(call.ID)),
	}
	c.Fcm.SendNotificationWithDataToken(token, &msg, &data, "id.streetbox.live")
	model.ResponseJSON(ctx, call)
	return
}

// UpdateStatusCall godoc
// @Summary Update Status Call From Consumer (permission = merchant)
// @Id UpdateStatusCall
// @Tags Canvassing
// @Security Token
// @Param call-id path string true "1"
// @Param status path string true "onprocess/accept/reject/expire"
// @Success 200 {object} entity.CanvassingCall "{ "data": Model }"
// @Failure 500 {object} model.ResponseErrors "Call Status Updated Failed"
// @Router /canvassing/call-status/:call-id/:status [put]
func (c *CanvassingController) UpdateStatusCall(ctx *gin.Context) {
	callID := util.ParamIDToInt64(ctx.Param("call-id"))
	status := ctx.Param("status")
	jwtModel, _ := util.ExtractTokenMetadata(ctx.Request)
	foodtruck := c.Merchant.GetMerchantUsersByUsersID(jwtModel.UserID)
	merchant := c.Merchant.GetByID(foodtruck.MerchantID)
	call, token, err := c.Svc.UpdateStatusCall(callID, foodtruck.ID, status)
	if err != nil {
		model.ResponseError(ctx, "Call Status Updated Failed", http.StatusInternalServerError)
		return
	}
	msg := messaging.Notification{
		Title: "Call Status Updated!",
		Body:  fmt.Sprintf("Your call has been %sed", status),
		ImageURL: "https://api.streetbox.id/static/image/" + merchant.Logo,
	}
	var tokens []string
	tokens = append(tokens, foodtruck.RegistrationToken)
	data := map[string]string{
		"title": "Call Status Updated!",
		"body":  fmt.Sprintf("Your call has been %sed", status),
		// "image": "https://api.streetbox.id/static/image/" + merchant.Logo,
		"call_id": strconv.Itoa(int(call.ID)),
	}
	c.Fcm.SendNotificationWithDataToken(token, &msg, &data, "id.streetbox.live")
	model.ResponseJSON(ctx, call)
	return
}

// FinishCall godoc
// @Summary Finish Call From Consumer (permission = merchant)
// @Id FinishCall
// @Tags Canvassing
// @Security Token
// @Param call-id path string true "1"
// @Success 200 {object} model.ResponseSuccess "{"message": "Success"}"
// @Failure 500 {object} model.ResponseErrors "Finish Consumer Call Failed"
// @Router /canvassing/finish/:call-id [put]
func (c *CanvassingController) FinishCall(ctx *gin.Context) {
	callID := util.ParamIDToInt64(ctx.Param("call-id"))
	jwtModel, _ := util.ExtractTokenMetadata(ctx.Request)
	foodtruck := c.Merchant.GetMerchantUsersByUsersID(jwtModel.UserID)
	merchant := c.Merchant.GetByID(foodtruck.MerchantID)
	call, token, err := c.Svc.UpdateStatusCall(callID, foodtruck.ID, "FINISH")
	if err != nil {
		model.ResponseError(ctx, "Finish Consumer Call Failed", http.StatusInternalServerError)
		return
	}
	msg := messaging.Notification{
		Title: "Call Finished!",
		Body:  "Your call has been finish",
		ImageURL: "https://api.streetbox.id/static/image/" + merchant.Logo,
	}
	data := map[string]string{
		"title": "Call Finished!",
		"body":  "Your call has been finish",
		// "image": "https://api.streetbox.id/static/image/" + merchant.Logo,
		"call_id": strconv.Itoa(int(call.ID)),
	}
	c.Fcm.SendNotificationWithDataToken(token, &msg, &data, "id.streetbox.live")
	model.ResponseJSON(ctx, gin.H{"message": "Success"})
	return
}

// GetNotifications godoc
// @Summary Inquiry Users' Notifications (permission = consumer)
// @Id GetNotifications
// @Tags Canvassing
// @Security Token
// @Success 200 {object} []model.RespNotifByUserID "{ "data": Model }"
// @Failure 500 {object} model.ResponseErrors "Inquiry Notifications Failed"
// @Router /canvassing/notifications [get]
func (c *CanvassingController) GetNotifications(ctx *gin.Context) {
	jwtModel, _ := util.ExtractTokenMetadata(ctx.Request)
	notifs, err := c.Svc.GetNotifications(jwtModel.UserID)
	if err != nil {
		model.ResponseError(ctx, "Inquiry Notifications Failed", http.StatusInternalServerError)
		return
	}
	model.ResponseJSON(ctx, notifs)
	return
}

// GetCallsByFoodtruckID godoc
// @Summary Inquiry Foodtruck's Calls (permission = merchant)
// @Id GetCallsByFoodtruckID
// @Tags Canvassing
// @Security Token
// @Param status query string true "onprocess/accept/request/history"
// @Success 200 {object} []model.RespCallsByFoodtruckID "{ "data": Model }"
// @Failure 500 {object} model.ResponseErrors "Inquiry Calls Failed"
// @Router /canvassing/foodtruck/calls [get]
func (c *CanvassingController) GetCallsByFoodtruckID(ctx *gin.Context) {
	jwtModel, _ := util.ExtractTokenMetadata(ctx.Request)
	status := ctx.DefaultQuery("status", "request")
	foodtruck := c.Merchant.GetMerchantUsersByUsersID(jwtModel.UserID)
	calls, err := c.Svc.GetCallsByFoodtruckID(foodtruck.ID, status)
	if err != nil {
		model.ResponseError(ctx, "Inquiry Calls Failed", http.StatusInternalServerError)
		return
	}
	model.ResponseJSON(ctx, calls)
	return
}

// GetCallsByCustomerID godoc
// @Summary Inquiry Customer's Calls (permission = consumer)
// @Id GetCallsByCustomerID
// @Tags Canvassing
// @Security Token
// @Param status query string true "accept/request/history"
// @Success 200 {object} []model.RespCallsByUserID "{ "data": Model }"
// @Failure 500 {object} model.ResponseErrors "Inquiry Calls Failed"
// @Router /canvassing/calls [get]
func (c *CanvassingController) GetCallsByCustomerID(ctx *gin.Context) {
	jwtModel, _ := util.ExtractTokenMetadata(ctx.Request)
	calls, err := c.Svc.GetCallsByCustomerID(jwtModel.UserID)
	if err != nil {
		model.ResponseError(ctx, "Inquiry Calls Failed", http.StatusInternalServerError)
		return
	}
	model.ResponseJSON(ctx, calls)
	return
}


// ToggleAutoBlast godoc
// @Summary Toggle Foodtruck Auto Blast (permission = merchant)
// @Id ToggleAutoBlast
// @Tags Canvassing
// @Security Token
// @Success 200 {object} model.ResponseSuccess "{"message": "Success"}"
// @Router /canvassing/toggle/auto [put]
func (c *CanvassingController) ToggleAutoBlast(ctx *gin.Context) {
	jwtModel, _ := util.ExtractTokenMetadata(ctx.Request)
	foodtruck := c.Merchant.GetMerchantUsersByUsersID(jwtModel.UserID)
	c.Svc.ToggleAutoBlast(foodtruck.ID)
	model.ResponseJSON(ctx, gin.H{"message": "Success"})
	return
}

func (c *CanvassingController) AutoBlast() {
	canvs, err := c.Svc.GetTriggeredCanvas()
	if err != nil {
		return
	}
	for _, canv := range canvs {
		if canv.IsAutoBlast {
			users, err := c.Svc.Blast(canv.FoodtruckID, canv.ID, true)
			if err != nil {
				log.Printf("ERROR: %s | CANVASSING-ID: %d | FOODTRUCK-ID: %d", err.Error(), canv.ID, canv.FoodtruckID)
				continue
			}
			msg := messaging.Notification{
				Title: "Foodtruck didekat sini",
				Body:  "Ada foodtruck keliling di dekat kamu, dipanggil yukk!",
				ImageURL: "https://api.streetbox.id/static/image/" + canv.MerchantLogo,
			}
			data := map[string]string{
				"title": "Foodtruck didekat sini",
				"body":  "Ada foodtruck keliling di dekat kamu, dipanggil yukk!",
				// "image": "https://api.streetbox.id/static/image/" + canv.MerchantLogo,
				"merchant_id": strconv.Itoa(int(canv.FoodtruckID)),
			}
			for _, user := range users {
				c.Fcm.SendNotificationWithData(fmt.Sprintf("blast_%d", user.ID), &msg, &data, "id.streetbox.live")
				notif := entity.CanvassingNotif{
					FoodtruckID:   canv.FoodtruckID,
					CustomerID:    user.ID,
					CustomerToken: user.RegistrationToken,
				}
				c.Svc.SaveNotification(notif)
			}
		}
	}
}
