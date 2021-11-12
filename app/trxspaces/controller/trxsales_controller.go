package controller

import (
	"net/http"

	"streetbox.id/util"

	"github.com/gin-gonic/gin"
	"streetbox.id/app/merchant"
	"streetbox.id/app/sales"
	"streetbox.id/app/trxspaces"
	"streetbox.id/model"
)

// TrxSalesController ...
type TrxSalesController struct {
	Service      trxspaces.ServiceInterface
	ParkSalesSvc sales.ServiceInterface
	MerchantSvc  merchant.ServiceInterface
}

// Create godoc
// @Summary Create Transaction Parking Space Sales (permission = superadmin)
// @Id CreateTrxSales
// @Tags Transaction Sales
// @Security Token
// @Param trxSales body model.ReqCreateTrxSales true "New Transaction Sales"
// @Success 200 {object} model.ResponseSuccess "message: "Create Transaction Sales Succeed" "
// @Failure 422 {object} model.ResponseErrors "code: 422, message: "Invalid request" "
// @Failure 500 {object} model.ResponseErrors "code: 500, message: "Create Transaction Sales Succeed" "
// @Router /trxsales [post]
func (r *TrxSalesController) Create(c *gin.Context) {
	req := model.ReqCreateTrxSales{}
	if err := c.ShouldBindJSON(&req); err != nil {
		model.ResponseError(c, "Invalid request", http.StatusUnprocessableEntity)
		return
	}
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	parkingSpaceSales := req.ParkingSpaceSalesID
	getData := r.ParkSalesSvc.GetByID(parkingSpaceSales)
	checkExisted, errCheck := r.Service.GetByMerchantIDAndParkingSalesID(req.MerchantID, parkingSpaceSales)
	//compare total slot with available slot
	if req.TotalSlot > getData.AvailableSlot {
		model.ResponseError(c, "Inputted Slot Transaction already reached maximum slot", http.StatusInternalServerError)
		return
	}
	if errCheck != nil || checkExisted == nil {
		if err := r.Service.CreateTrx(&req, jwtModel.UserID); err != nil {
			model.ResponseError(c, "Create Transaction Failed", http.StatusInternalServerError)
			return
		}
	} else {
		req.TotalSlot = checkExisted.TotalSlot + req.TotalSlot
		if err := r.Service.UpdateTrx(&req, checkExisted.ID, jwtModel.UserID); err != nil {
			model.ResponseError(c, "Create Transaction Failed", http.StatusInternalServerError)
			return
		}
	}
	model.ResponseCreated(c, gin.H{"message": "Create Transaction Sales Succeed"})
	return
}

// GetMyParking godoc
// @Summary Get My Parking Space (permission = merchant)
// @Id GetMyParking
// @Tags Transaction Sales
// @Security Token
// @Success 200 {object} []model.ResMyParkingList "data: []model.ResMyParking"
// @Failure 500 {object} model.ResponseErrors "code: 500, message: "Failed" "
// @Router /trxsales/myparking [get]
func (r *TrxSalesController) GetMyParking(c *gin.Context) {
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	merchant := r.MerchantSvc.GetInfo(jwtModel.UserID)
	if merchant == nil {
		model.ResponseError(c, "Merchant Not Found", http.StatusUnprocessableEntity)
		return
	}
	data := r.Service.GetMyParking(jwtModel.UserID, merchant.ID)
	model.ResponseJSON(c, data)
	return
}

// GetSlotMyParking godoc
// @Summary Get Slot My Parking Space (permission = merchant)
// @Id GetSlotMyParking
// @Tags Transaction Sales
// @Security Token
// @Param id path integer true "Parking Space ID"
// @Success 200 {object} []model.ResSlotMyParking "data: []model.ResSlotMyParking "
// @Failure 500 {object} model.ResponseErrors "code: 500, message: "Failed" "
// @Router /trxsales/myparking/slot/{id} [get]
func (r *TrxSalesController) GetSlotMyParking(c *gin.Context) {
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	id := util.ParamIDToInt64(c.Param("id"))
	var (
		data *[]model.ResSlotMyParking
		err  error
	)
	if data, err = r.Service.GetSlotMyParking(id, jwtModel.UserID); err != nil {
		model.ResponseError(c, "Failed", http.StatusInternalServerError)
		return
	}
	model.ResponseJSON(c, data)
	return
}

// GetAll godoc
// @Summary Get All Transaction (permission = superadmin)
// @Id GetAllTrxSales
// @Tags Transaction Sales
// @Security Token
// @Success 200 {object} []model.ResTrxList "data: []model.ResTrxList"
// @Router /trxsales/all [get]
func (r *TrxSalesController) GetAll(c *gin.Context) {
	model.ResponseJSON(c, r.Service.GetAll())
}

// GetList godoc
// @Summary Get List Pagination Transaction (permission = superadmin)
// @Id GetList
// @Tags Transaction Sales
// @Security Token
// @Param limit query string false " " default(10)
// @Param page query string false " " default(1)
// @Param filter query string false " "
// @Param sort query string false "e.g.: id,desc / id,asc"
// @Success 200 {object} model.ResponseSuccess " "
// @Router /trxsales/list [get]
func (r *TrxSalesController) GetList(c *gin.Context) {
	limit := util.ParamIDToInt(c.DefaultQuery("limit", "10"))
	page := util.ParamIDToInt(c.DefaultQuery("page", "1"))
	filter := c.DefaultQuery("filter", "")
	sort := util.SortedBy(c.QueryArray("sort"))
	model.ResponseJSON(c, r.Service.GetList(limit, page, sort, filter))
}

// GetByID godoc
// @Summary Get Transaction by ID (permission = superadmin)
// @Id GetByIDTrxSales
// @Tags Transaction Sales
// @Security Token
// @Param id path integer true "Trx ID"
// @Success 200 {object} []model.ResTrxList "data: []model.ResTrxList"
// @Router /trxsales/info/{id} [get]
func (r *TrxSalesController) GetByID(c *gin.Context) {
	id := util.ParamIDToInt64(c.Param("id"))
	model.ResponseJSON(c, r.Service.GetByID(id))
	return
}
