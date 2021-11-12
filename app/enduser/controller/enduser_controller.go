package controller

import (
	"math"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"streetbox.id/app/enduser"
	"streetbox.id/app/fcm"
	"streetbox.id/app/homevisitsales"
	"streetbox.id/app/merchant"
	"streetbox.id/app/trx"
	"streetbox.id/cfg"
	"streetbox.id/model"
	"streetbox.id/util"
)

// EndUserController ..
type EndUserController struct {
	Svc           enduser.ServiceInterface
	VisitSalesSvc homevisitsales.ServiceInterface
	MerchantSvc   merchant.ServiceInterface
	TrxSvc        trx.ServiceInterface
	FcmSvc        fcm.ServiceInterface
}

// Nearby  godoc
// @Summary Landing Page End User Nearby (permission = consumer)
// @Id Nearby
// @Tags End User
// @Security Token
// @Param lat path number true "Latitude"
// @Param lon path number true "Longitude"
// @Param limit query string false " " default(999999)
// @Param page query string false " " default(1)
// @Param distance query string false "KM" default(10)
// @Success 200 {object} model.Pagination "data: model.Pagination"
// @Router /consumer/home/nearby/{lat}/{lon} [get]
func (r *EndUserController) Nearby(c *gin.Context) {
	lat := util.ParamToFloat64(c.Param("lat"))
	lon := util.ParamToFloat64(c.Param("lon"))
	req := model.ReqMerchantNearby{
		Latitude:  lat,
		Longitude: lon,
	}
	distance := util.ParamToFloat64(c.DefaultQuery("distance", "100"))
	limit := util.ParamIDToInt(c.DefaultQuery("limit", "999999"))
	page := util.ParamIDToInt(c.DefaultQuery("page", "1"))
	data := r.Svc.GetNearby(limit, page, distance, &req)
	model.ResponsePagination(c, data)
	return
}

// LiveTracking  godoc
// @Summary Live Tracking Foodtruck Nearby on Maps  (permission = consumer)
// @Id LiveTracking
// @Tags End User
// @Security Token
// @Param lat path number true "Latitude"
// @Param lon path number true "Longitude"
// @Param distance query string false "KM" default(10)
// @Success 200 {object} []model.ResLiveTracking "data: []model.ResLiveTracking"
// @Router /consumer/home/map/livetracking/{lat}/{lon} [get]
func (r *EndUserController) LiveTracking(c *gin.Context) {
	lat := util.ParamToFloat64(c.Param("lat"))
	lon := util.ParamToFloat64(c.Param("lon"))
	distance := util.ParamToFloat64(c.DefaultQuery("distance", "10"))
	model.ResponseJSON(c, r.Svc.GetLiveTracking(lat, lon, distance))
	return
}

// MapParkingSpace  godoc
// @Summary Show Parking Space Nearby on Maps  (permission = consumer)
// @Id MapParkingSpace
// @Tags End User
// @Security Token
// @Param lat path number true "Latitude Consumer"
// @Param lon path number true "Longitude Consumer"
// @Param distance query string false "KM" default(10)
// @Success 200 {object} []model.ResParkingSpace "data: []model.ResParkingSpace"
// @Router /consumer/home/map/parking-space/{lat}/{lon} [get]
func (r *EndUserController) MapParkingSpace(c *gin.Context) {
	lat := util.ParamToFloat64(c.Param("lat"))
	lon := util.ParamToFloat64(c.Param("lon"))
	distance := util.ParamToFloat64(c.DefaultQuery("distance", "10"))
	model.ResponseJSON(c, r.Svc.MapParkingSpace(lat, lon, distance))
	return
}

// MapParkingSpaceDetail  godoc
// @Summary Show Schedule Parking Space Nearby on Maps  (permission = consumer)
// @Id MapParkingSpaceDetail
// @Tags End User
// @Security Token
// @Param id path integer true "Parking Space"
// @Success 200 {object} []model.ResParkingSpaceDetail "data: []model.ResParkingSpaceDetail"
// @Router /consumer/home/map/schedules/{id}/parking-space [get]
func (r *EndUserController) MapParkingSpaceDetail(c *gin.Context) {
	id := util.ParamIDToInt64(c.Param("id"))
	model.ResponseJSON(c, r.Svc.MapParkingSpaceDetail(id))
	return
}

// MerchantSchedule  godoc
// @Summary Show Schedules Regular Merchants  (permission = consumer)
// @Id MerchantSchedule
// @Tags End User
// @Security Token
// @Param typesId path integer true "Tasks Regular ID"
// @Success 200 {object} []model.Schedules "data: []model.Schedules"
// @Router /consumer/home/schedules-regular/{typesId} [get]
func (r *EndUserController) MerchantSchedule(c *gin.Context) {
	id := util.ParamIDToInt64(c.Param("typesId"))
	model.ResponseJSON(c, r.Svc.GetSchedulesByTypesID(id))
	return
}

// VisitSales  godoc
// @Summary Show HomeVisit Sales  (permission = consumer)
// @Id VisitSales
// @Tags End User
// @Security Token
// @Param limit query string false " " default(999999)
// @Param page query string false " " default(1)
// @Success 200 {object} []model.ResVisitSales "data: []model.ResVisitSales"
// @Router /consumer/home/visit-sales [get]
func (r *EndUserController) VisitSales(c *gin.Context) {
	limit := util.ParamIDToInt(c.DefaultQuery("limit", "999999"))
	page := util.ParamIDToInt(c.DefaultQuery("page", "1"))
	data := r.VisitSalesSvc.GetAllEndUser(limit, page)
	model.ResponsePagination(c, data)
	return
}

// VisitSalesDetail  godoc
// @Summary Show HomeVisit Sales Available  (permission = consumer)
// @Id VisitSalesDetail
// @Tags End User
// @Security Token
// @Param merchantId path integer true "merchant id"
// @Success 200 {object} []model.ResVisitSalesDetail "data: []model.ResVisitSalesDetail"
// @Router /consumer/home/visit-sales/detail/{merchantId} [get]
func (r *EndUserController) VisitSalesDetail(c *gin.Context) {
	id := util.ParamIDToInt64(c.Param("merchantId"))
	data := r.VisitSalesSvc.GetAvailableByMerchantID(id)
	model.ResponseJSON(c, data)
	return
}

// GetMerchantTaxByMerchantID  godoc
// @Summary Get Merchant Tax (permission = consumer)
// @Id GetMerchantTax
// @Tags End User
// @Security Token
// @Param merchantId path integer true "Merchant ID"
// @Success 200 {object} entity.MerchantTax "data: entity.MerchantTax"
// @Router /consumer/merchant/tax/{merchantId} [get]
func (r *EndUserController) GetMerchantTaxByMerchantID(c *gin.Context) {
	merchantID := util.ParamIDToInt64(c.Param("merchantId"))
	model.ResponseJSON(c, r.MerchantSvc.GetTax(merchantID))
	return
}

// GetMerchantMenuByMerchantID  godoc
// @Summary Get Merchant Menu (permission = consumer)
// @Id GetMerchantMenu
// @Tags End User
// @Security Token
// @Param filter query string false "all, nearby, visit" default(all)
// @Param merchantId path integer true "Merchant ID"
// @Success 200 {object} []entity.MerchantMenu "data: []entity.MerchantMenu"
// @Router /consumer/merchant/menu/{merchantId} [get]
func (r *EndUserController) GetMerchantMenuByMerchantID(c *gin.Context) {
	var nearby, visit bool
	merchantID := util.ParamIDToInt64(c.Param("merchantId"))
	filter := c.DefaultQuery("filter", "all")
	if filter == "all" {
		nearby, visit = true, true
	} else if filter == "nearby" {
		nearby = true
	} else if filter == "visit" {
		visit = true
	}
	menus := r.MerchantSvc.GetMenuList(merchantID, nearby, visit)
	var (
		updatedMenus []model.ResMerchantMenuList
		updatedMenu  model.ResMerchantMenuList
		price        float64
	)
	for _, value := range *menus {
		copier.Copy(&updatedMenu, value)
		price = float64(value.Price)
		if price <= 100 {
			updatedMenu.PriceAfterDiscount = math.Floor(price * (float64(1) - float64(value.Discount/100)))
		} else {
			updatedMenu.PriceAfterDiscount = price - float64(value.Discount)
		}
		updatedMenus = append(updatedMenus, updatedMenu)
	}

	model.ResponseJSON(c, updatedMenus)
	return
}

// OrderHistory  godoc
// @Summary Order History (permission = consumer)
// @Id OrderHistory
// @Tags End User
// @Security Token
// @Param limit query string false " " default(10)
// @Param page query string false " " default(1)
// @Param filter query string false " " default(all)
// @Success 200 {object} model.Pagination "data: model.Pagination"
// @Router /consumer/order/history [get]
func (r *EndUserController) OrderHistory(c *gin.Context) {
	limit := util.ParamIDToInt(c.DefaultQuery("limit", "10"))
	page := util.ParamIDToInt(c.DefaultQuery("page", "1"))
	filter := c.DefaultQuery("filter", "all")
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	data := r.TrxSvc.GetOrderHistoryByUsersID(limit, page, jwtModel.UserID, filter)
	model.ResponsePagination(c, data)
	return
}

// GetPaymentMethod  godoc
// @Summary Get Payment Method (permission = consumer)
// @Id GetPaymentMethod
// @Tags End User
// @Security Token
// @Success 200 {object} []model.ResPaymentMethod "data: []model.ResPaymentMethod"
// @Router /consumer/payment-method [get]
func (r *EndUserController) GetPaymentMethod(c *gin.Context) {
	data := r.Svc.GetPaymentMethod()
	model.ResponseJSON(c, data)
	return
}

// RegistrationToken godoc
// @Summary Send Registration Token FCM, Its should be done to archive push notification (permission = merchant)
// @Id RegistrationTokenEndUser
// @Tags End User
// @Security Token
// @Param token path string true "fcm token from client SDK"
// @Success 200 {object} model.ResTrxOrderList "data: model.ResTrxOrderList"
// @Router /consumer/registration-token/{token} [post]
func (r *EndUserController) RegistrationToken(c *gin.Context) {
	regisToken := c.Param("token")
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	if err := r.Svc.RegistrationToken(regisToken, jwtModel.UserID); err != nil {
		model.ResponseError(c, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	model.ResponseJSON(c, "Storing Registration Token Successfully")
	return
}

// UpdateEndUser godoc
// @Summary Update End User Profile (permission = consumer)
// @Id UpdateEndUser
// @Tags End User
// @Security Token
// @Accept multipart/form-data
// @Produce json
// @Param name formData string true "End User Name"
// @Param address formData string true "End User Address"
// @Param phone formData string true "End User Phone"
// @Param email formData string true "End User Email"
// @Param image formData file false "Profile Picture"
// @Success 200 {object} model.ResponseSuccess "message: "User Successfully Updated" "
// @Failure 422 {object} model.ResponseErrors "code: 422, message: "Error When Update User Profile" "
// @Router /consumer/update/userprofile [put]
func (r *EndUserController) UpdateEndUser(c *gin.Context) {
	name := c.PostForm("name")
	address := c.PostForm("address")
	phone := c.PostForm("phone")
	email := c.PostForm("email")
	image, _ := c.FormFile("image")
	path := cfg.Config.Path.Image
	var req interface{}
	if image == nil {
		req = model.ReqEndUserUpdate{
			Name:    name,
			Phone:   phone,
			Address: address,
			Email:   email,
		}
		if err := r.Svc.UpdateEndUser(req.(model.ReqEndUserUpdate)); err != nil {
			model.ResponseError(c, "Error When Update User Profile", http.StatusUnprocessableEntity)
			return
		}
		model.ResponseJSON(c, "User Successfully updated")
		return
	}
	filename := util.GeneratedUUID(filepath.Base(image.Filename))
	req = model.ReqEndUserUpdateWithImage{
		Name:         name,
		Phone:        phone,
		Address:      address,
		Email:        email,
		PhotoProfile: filename,
	}
	pathFile := path + filename
	if err := c.SaveUploadedFile(image, pathFile); err != nil {
		model.ResponseError(c, "Upload Photo Profile error", http.StatusBadRequest)
		return
	}
	if err := r.Svc.UpdateEndUserWithImage(req.(model.ReqEndUserUpdateWithImage)); err != nil {
		model.ResponseError(c, "Error When Update User Profile", http.StatusUnprocessableEntity)
		return
	}
	model.ResponseJSON(c, "User Successfully updated")
	return
}
