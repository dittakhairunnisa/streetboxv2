package controller

import (
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"streetbox.id/entity"

	"github.com/gin-gonic/gin"
	"streetbox.id/app/homevisitsales"
	"streetbox.id/app/merchant"
	"streetbox.id/model"
	"streetbox.id/util"
)

// HomevisitSalesController ..
type HomevisitSalesController struct {
	HomeSalesSvc homevisitsales.ServiceInterface
	MerchantSvc  merchant.ServiceInterface
}

// GetAll godoc
// @Summary Get All Home Visit Pagination (permission = admin)
// @Id GetAllHV
// @Tags Home Visit
// @Security Token
// @Param limit query string false " " default(10)
// @Param page query string false " " default(1)
// @Param sort query string false "e.g.: id,desc / id,asc"
// @Success 200 {object} model.ResponseSuccess "model.Pagination"
// @Router /homevisit [get]
func (r *HomevisitSalesController) GetAll(c *gin.Context) {
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	merchantID := r.MerchantSvc.GetInfo(jwtModel.UserID).ID
	data := r.HomeSalesSvc.GetAll(merchantID)
	model.ResponseJSON(c, data)
	return
}

// Create godoc
// @Summary Batch Create New Home Visit (permission = admin)
// @Id BatchCreateHV
// @Tags Home Visit
// @Security Token
// @Param batchHomeVisit body model.ReqBatchCreateHomevisitSales true "all fields mandatory"
// @Success 200 {object} []entity.HomevisitSales "data: []entity.HomevisitSales"
// @Router /homevisit/batch [post]
func (r *HomevisitSalesController) BatchCreate(c *gin.Context) {
	request := &model.ReqBatchCreateHomevisitSales{}
	var datas []entity.HomevisitSales
	if err := c.ShouldBindJSON(&request); err != nil {
		model.ResponseError(c, "Invalid Request", http.StatusUnprocessableEntity)
		return
	}
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	merchantID := r.MerchantSvc.GetInfo(jwtModel.UserID).ID
	for _, req := range request.Request {
		checkDate := r.HomeSalesSvc.CheckDate(req.Date, merchantID)
		if checkDate > 0 {
			model.ResponseError(c, "Cannot Created Existing Date", http.StatusNotAcceptable)
			return
		}
		Foodtrucks := r.MerchantSvc.CountFoodtruckByMerchantID(merchantID)
		var (
			homeVisitSales []entity.HomevisitSales
			homeVisitSale  entity.HomevisitSales
			startDate      string
			endDate        string
			startDateTime  time.Time
			endDateTime    time.Time
			err            error
			data           *entity.HomevisitSales
			date           string
			dates          []string
		)
	
		for _, value := range req.Summary {
			if value.NumberOfFoodtruck > Foodtrucks {
				model.ResponseError(c, "Invalid Request", http.StatusNotAcceptable)
				return
			}
	
			startDate = req.Date + " " + value.StartTime
			endDate = req.Date + " " + value.EndTime
			startDateTime, err = time.Parse("2006-01-02 15:04:05", startDate)
			if err != nil {
				model.ResponseError(c, "Invalid Date Format", http.StatusNotAcceptable)
				return
			}
			endDateTime, err = time.Parse("2006-01-02 15:04:05", endDate)
			if err != nil {
				model.ResponseError(c, "Invalid Date Format", http.StatusNotAcceptable)
				return
			}
			homeVisitSale = entity.HomevisitSales{
				MerchantID: merchantID,
				StartDate:  startDateTime,
				EndDate:    endDateTime,
				Deposit:    req.Deposit,
				Total:      value.NumberOfFoodtruck,
				Available:  value.NumberOfFoodtruck,
			}
			homeVisitSales = append(homeVisitSales, homeVisitSale)
	
			date = startDate + "#" + endDate + "#" + strconv.Itoa(value.NumberOfFoodtruck)
			dates = append(dates, date)
	
		}
	
		sort.Strings(dates)
		var results int
		for _, value := range dates {
			split := strings.Split(value, "#")
			results += util.ScheduleCompare(split[0], split[1], dates, "#")
			if results > Foodtrucks {
				model.ResponseError(c, "Number of Foodtruck Exceeded Available Foodtruck!", http.StatusNotAcceptable)
				return
			}
		}
	
		for _, value := range homeVisitSales {
			if data, err = r.HomeSalesSvc.Create(&value); err != nil {
				model.ResponseError(c, "Invalid Request", http.StatusNotAcceptable)
				return
			}
			datas = append(datas, *data)
			data = nil
		}
	}
	model.ResponseJSON(c, datas)
	return
}

// Create godoc
// @Summary Create New Home Visit (permission = admin)
// @Id CreateHV
// @Tags Home Visit
// @Security Token
// @Param homeVisit body model.ReqCreateHomevisitSales true "all fields mandatory"
// @Success 200 {object} []entity.HomevisitSales "data: []entity.HomevisitSales"
// @Router /homevisit [post]
func (r *HomevisitSalesController) Create(c *gin.Context) {
	req := &model.ReqCreateHomevisitSales{}
	if err := c.ShouldBindJSON(&req); err != nil {
		model.ResponseError(c, "Invalid Request", http.StatusUnprocessableEntity)
		return
	}

	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	merchantID := r.MerchantSvc.GetInfo(jwtModel.UserID).ID
	checkDate := r.HomeSalesSvc.CheckDate(req.Date, merchantID)
	if checkDate > 0 {
		model.ResponseError(c, "Cannot Created Existing Date", http.StatusNotAcceptable)
		return
	}
	Foodtrucks := r.MerchantSvc.CountFoodtruckByMerchantID(merchantID)
	var (
		homeVisitSales []entity.HomevisitSales
		homeVisitSale  entity.HomevisitSales
		startDate      string
		endDate        string
		startDateTime  time.Time
		endDateTime    time.Time
		err            error
		data           *entity.HomevisitSales
		datas          []entity.HomevisitSales
		date           string
		dates          []string
	)

	for _, value := range req.Summary {
		if value.NumberOfFoodtruck > Foodtrucks {
			model.ResponseError(c, "Invalid Request", http.StatusNotAcceptable)
			return
		}

		startDate = req.Date + " " + value.StartTime
		endDate = req.Date + " " + value.EndTime
		startDateTime, err = time.Parse("2006-01-02 15:04:05", startDate)
		if err != nil {
			model.ResponseError(c, "Invalid Date Format", http.StatusNotAcceptable)
			return
		}
		endDateTime, err = time.Parse("2006-01-02 15:04:05", endDate)
		if err != nil {
			model.ResponseError(c, "Invalid Date Format", http.StatusNotAcceptable)
			return
		}
		homeVisitSale = entity.HomevisitSales{
			MerchantID: merchantID,
			StartDate:  startDateTime,
			EndDate:    endDateTime,
			Deposit:    req.Deposit,
			Total:      value.NumberOfFoodtruck,
			Available:  value.NumberOfFoodtruck,
		}
		homeVisitSales = append(homeVisitSales, homeVisitSale)

		date = startDate + "#" + endDate + "#" + strconv.Itoa(value.NumberOfFoodtruck)
		dates = append(dates, date)

	}

	sort.Strings(dates)
	var results int
	for _, value := range dates {
		split := strings.Split(value, "#")
		results += util.ScheduleCompare(split[0], split[1], dates, "#")
		if results > Foodtrucks {
			model.ResponseError(c, "Number of Foodtruck Exceeded Available Foodtruck!", http.StatusNotAcceptable)
			return
		}
	}

	for _, value := range homeVisitSales {
		if data, err = r.HomeSalesSvc.Create(&value); err != nil {
			model.ResponseError(c, "Invalid Request", http.StatusNotAcceptable)
			return
		}
		datas = append(datas, *data)
		data = nil
	}
	model.ResponseJSON(c, datas)
	return
}

// GetInfo godoc
// @Summary Get Info Home Visit (permission = admin)
// @Id GetInfoHV
// @Tags Home Visit
// @Security Token
// @Param date path string true "yyyy-mm-dd"
// @Success 200 {object} model.ResHomeVisitGetInfo "data: model.ResHomeVisitGetInfo"
// @Router /homevisit/info/{date} [GET]
func (r *HomevisitSalesController) GetInfo(c *gin.Context) {
	date := c.Param("date")
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	merchantID := r.MerchantSvc.GetInfo(jwtModel.UserID).ID

	resultInfo := r.HomeSalesSvc.GetInfoByDate(merchantID, date)

	model.ResponseJSON(c, resultInfo)
	return
}

// Delete godoc
// @Summary Delete Home Visit (permission = admin)
// @Id DeleteHomeVisit
// @Tags Home Visit
// @Security Token
// @Param date path string true "yyyy-mm-dd"
// @Success 200 {object} model.ResponseSuccess "message: "successfully deleted home visit!" "
// @Failure 500 {object} model.ResponseErrors  "message: "Failed to Delete Menu""
// @Router /homevisit/deletebydate/{date} [DELETE]
func (r *HomevisitSalesController) Delete(c *gin.Context) {
	date := c.Param("date")
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	merchantID := r.MerchantSvc.GetInfo(jwtModel.UserID).ID

	if _, err := r.HomeSalesSvc.DeleteByDate(date, merchantID); err != nil {
		model.ResponseError(c, "Problem when trying to Delete Home Visit", http.StatusNotAcceptable)
		return
	}
	model.ResponseJSON(c, gin.H{"message": "successfully deleted home visit!"})
	return
}

// DeleteByID godoc
// @Summary Delete Home Visit By ID (permission = admin)
// @Id DeleteHomeVisitByID
// @Tags Home Visit
// @Security Token
// @Param id path string true "number"
// @Success 200 {object} model.ResponseSuccess "message: "successfully deleted home visit!" "
// @Failure 500 {object} model.ResponseErrors  "message: "Failed to Delete Menu""
// @Router /homevisit/deletebyid/{id} [DELETE]
func (r *HomevisitSalesController) DeleteByID(c *gin.Context) {
	id := util.ParamIDToInt64(c.Param("id"))
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	merchantID := r.MerchantSvc.GetInfo(jwtModel.UserID).ID

	if _, err := r.HomeSalesSvc.DeleteByID(id, merchantID); err != nil {
		model.ResponseError(c, "Problem when trying to Delete Home Visit", http.StatusNotAcceptable)
		return
	}
	model.ResponseJSON(c, gin.H{"message": "successfully deleted home visit!"})
	return
}

// Update godoc
// @Summary Update Home Visit (permission = admin)
// @Id UpdateHomeVisit
// @Tags Home Visit
// @Security Token
// @Param homeVisit body model.ReqUpdateHomevisitSales true "all fields mandatory"
// @Success 200 {object} []entity.HomevisitSales "data: []entity.HomevisitSales"
// @Router /homevisit [PUT]
func (r *HomevisitSalesController) Update(c *gin.Context) {
	req := &model.ReqUpdateHomevisitSales{}
	if err := c.ShouldBindJSON(&req); err != nil {
		model.ResponseError(c, "Invalid Request", http.StatusUnprocessableEntity)
		return
	}

	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	merchantID := r.MerchantSvc.GetInfo(jwtModel.UserID).ID
	Foodtrucks := r.MerchantSvc.CountFoodtruckByMerchantID(merchantID)
	var (
		homeVisitSales []entity.HomevisitSales
		homeVisitSale  entity.HomevisitSales
		startDate      string
		endDate        string
		startDateTime  time.Time
		endDateTime    time.Time
		err            error
		data           *entity.HomevisitSales
		datas          []entity.HomevisitSales
		date           string
		dates          []string
	)
	for _, value := range req.Summary {
		if value.NumberOfFoodtruck > Foodtrucks {
			model.ResponseError(c, "Invalid Request", http.StatusNotAcceptable)
			return
		}

		startDate = req.Date + " " + value.StartTime
		endDate = req.Date + " " + value.EndTime
		startDateTime, err = time.Parse("2006-01-02 15:04:05", startDate)
		if err != nil {
			model.ResponseError(c, "Invalid Date Format", http.StatusNotAcceptable)
			return
		}
		endDateTime, err = time.Parse("2006-01-02 15:04:05", endDate)
		if err != nil {
			model.ResponseError(c, "Invalid Date Format", http.StatusNotAcceptable)
			return
		}
		homeVisitSale = entity.HomevisitSales{
			ID:         value.ID,
			MerchantID: merchantID,
			StartDate:  startDateTime,
			EndDate:    endDateTime,
			Deposit:    req.Deposit,
			Total:      value.NumberOfFoodtruck,
			Available:  value.NumberOfFoodtruck,
		}
		homeVisitSales = append(homeVisitSales, homeVisitSale)

		date = startDate + "#" + endDate + "#" + strconv.Itoa(value.NumberOfFoodtruck)
		dates = append(dates, date)
	}

	sort.Strings(dates)
	var results int
	for _, value := range dates {
		split := strings.Split(value, "#")
		results += util.ScheduleCompare(split[0], split[1], dates, "#")
		if results > Foodtrucks {
			model.ResponseError(c, "Number of Foodtruck Exceeded Available Foodtruck!", http.StatusNotAcceptable)
			return
		}
	}

	for _, value := range homeVisitSales {
		if value.ID == 0 {
			if data, err = r.HomeSalesSvc.Create(&value); err != nil {
				model.ResponseError(c, "Invalid Request", http.StatusNotAcceptable)
				return
			}
		} else {
			if data, err = r.HomeSalesSvc.Update(&value); err != nil {
				model.ResponseError(c, "Invalid Request", http.StatusNotAcceptable)
				return
			}
		}
		datas = append(datas, *data)
		data = nil
	}
	model.ResponseJSON(c, datas)
	return
}
