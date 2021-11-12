package controller

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"streetbox.id/app/parkingspace"
	"streetbox.id/app/sales"
	"streetbox.id/cfg"
	"streetbox.id/entity"
	"streetbox.id/model"
	"streetbox.id/util"
)

// ParkingSpaceController ...
type ParkingSpaceController struct {
	PSpaceSvc parkingspace.ServiceInterface
	SalesSvc  sales.ServiceInterface
}

// GetAll godoc
// @Summary Get All Parking Space Pagination (permission = admin)
// @Id GetAllPSpace
// @Tags Parking Space
// @Security Token
// @Param limit query string false " " default(10)
// @Param page query string false " " default(1)
// @Param sort query string false "e.g.: id,desc / id,asc"
// @Success 200 {object} model.ResponseSuccess "model.Pagination"
// @Router /parkingspace [get]
func (s *ParkingSpaceController) GetAll(c *gin.Context) {
	limit := util.ParamIDToInt(c.DefaultQuery("limit", "10"))
	page := util.ParamIDToInt(c.DefaultQuery("page", "1"))
	sorted := util.SortedBy(c.QueryArray("sort"))
	data := s.PSpaceSvc.GetAll(limit, page, sorted)
	model.ResponsePagination(c, data)
	return
}

// GetAllList godoc
// @Summary Get All Parking Space (permission = superadmin)
// @Id GetAllList
// @Tags Parking Space
// @Security Token
// @Success 200 {object} model.ResponseSuccess "model.Pagination"
// @Router /parkingspace/list [get]
func (s *ParkingSpaceController) GetAllList(c *gin.Context) {
	model.ResponseJSON(c, s.PSpaceSvc.GetAllList())
	return
}

// CreateParkingSpace godoc
// @Summary Create Parking Space (permission = superadmin)
// @Id CreateParkingSpace
// @Tags Parking Space
// @Security Token
// @Accept multipart/form-data
// @Produce json
// @Param name formData string true "Parking Space Name"
// @Param address formData string true "Address"
// @Param name formData string true "Parking Space Name"
// @Param lat formData number true "Latitude"
// @Param long formData number true "Longitude"
// @Param totalSpace formData string true "Total Space"
// @Param desc formData string true "Description"
// @Param landlordInfo formData string true "Landlord info"
// @Param rating formData number true "Rating"
// @Param startContract formData string true "Start Contract yyyy-MM-dd"
// @Param endContract formData string true "End Contract yyyy-MM-dd"
// @Param startTime formData string true "Start Operation yyyy-MM-dd HH:mm:SS"
// @Param endTime formData string true "End Operation yyyy-MM-dd HH:mm:SS"
// @Success 201 {object} entity.ParkingSpace "data: entity.ParkingSpace, message: "Create Parking Space Success" "
// @Failure 500 {object} model.ResponseErrors "code: 500, message: "Create Parking Space Failed" "
// @Router /parkingspace [post]
func (s *ParkingSpaceController) CreateParkingSpace(c *gin.Context) {
	req := model.ReqParkingSpaceCreate{}
	req.Address = c.PostForm("address")
	req.Description = c.PostForm("desc")
	req.EndContract = util.ParamToDate(c.PostForm("endContract"))
	req.EndTime = util.ParamToDatetime(c.PostForm("endTime"))
	req.LandlordInfo = c.PostForm("landlordInfo")
	req.Latitude = util.ParamToFloat64(c.PostForm("lat"))
	req.Longitude = util.ParamToFloat64(c.PostForm("long"))
	req.Name = c.PostForm("name")
	req.Rating = util.ParamToFloat32(c.PostForm("rating"))
	req.StartContract = util.ParamToDate(c.PostForm("startContract"))
	req.StartTime = util.ParamToDatetime(c.PostForm("startTime"))
	req.TotalSpace = util.ParamIDToInt(c.PostForm("totalSpace"))
	req.City = c.PostForm("city")
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	var (
		data *entity.ParkingSpace
		err  error
	)
	if data, err = s.PSpaceSvc.Create(&req, jwtModel.UserID); err != nil {
		model.ResponseError(c, "Create Parking Space Failed", http.StatusInternalServerError)
		return
	}
	model.ResponseCreated(c, gin.H{"data": data, "message": "Create Parking Space Success"})
	return
}

// Update godoc
// @Summary Update Parking Space except uploading files/images (permission = superadmin)
// @Id Update
// @Tags Parking Space
// @Security Token
// @Accept multipart/form-data
// @Produce json
// @Param id path integer true "parking_space_id"
// @Param name formData string false " "
// @Param address formData string false " "
// @Param latitude formData number false " "
// @Param longitude formData number false " "
// @Param total formData integer false "space"
// @Param description formData string false " "
// @Param latitude formData string false " "
// @Param rating formData number false " "
// @Param startContract formData string false " "
// @Param endContract formData string false " "
// @Param startTime formData string false " "
// @Param endTime formData string false " "
// @Success 200 {object} entity.ParkingSpace "data: entity.ParkingSpace"
// @Failure 422 {object} model.ResponseErrors "code: 422, message: "Invalid request" "
// @Failure 500 {object} model.ResponseErrors "code: 500, message: "Update Parking Space Failed" "
// @Router /parkingspace/{id}/update [put]
func (s *ParkingSpaceController) Update(c *gin.Context) {
	id := util.ParamIDToInt64(c.Param("id"))
	req := model.ReqParkingSpaceUpdate{}
	if err := c.ShouldBind(&req); err != nil {
		model.ResponseError(c, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	var (
		data *entity.ParkingSpace
		err  error
	)
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	if data, err = s.PSpaceSvc.Update(&req, id, jwtModel.UserID); err != nil {
		model.ResponseError(c, "Update Parking Space Failed", http.StatusInternalServerError)
		return
	}
	model.ResponseJSON(c, data)
	return
}

// UploadImage godoc
// @Summary Upload Image (permission = superadmin)
// @Id UploadImage
// @Tags Parking Space
// @Security Token
// @Accept multipart/form-data
// @Produce json
// @Param id path integer true "parking_space_id"
// @Param image formData file true "image parking space"
// @Success 201 {object} model.ResponseSuccess "message: "Upload Image Parking Space Success" "
// @Failure 400 {object} model.ResponseErrors "code: 400, message: "Form/Upoad image error" "
// @Failure 500 {object} model.ResponseErrors "code: 500, message: "Update Image Meta data failed" "
// @Router /parkingspace/{id}/image/upload [put]
func (s *ParkingSpaceController) UploadImage(c *gin.Context) {
	id := util.ParamIDToInt64(c.Param("id"))
	imageMeta := s.PSpaceSvc.GetOne(id).ImagesMeta
	path := cfg.Config.Path.Image
	// images upload when create success
	// image 1
	image, err := c.FormFile("image")
	if err != nil {
		model.ResponseError(c, "Form image error", http.StatusBadRequest)
		return
	}
	filename := util.GeneratedUUID(filepath.Base(image.Filename))
	pathImg := path + filename
	if err := c.SaveUploadedFile(image, pathImg); err != nil {
		model.ResponseError(c, "Upload image error", http.StatusBadRequest)
		return
	}
	imageMeta = append(imageMeta, filename)
	if err := s.PSpaceSvc.UploadImage(imageMeta, id); err != nil {
		model.ResponseError(c, "Update Image Meta data failed", http.StatusInternalServerError)
		return
	}
	model.ResponseUpdated(c, gin.H{"message": "Upload Image Parking Space Success"})
	return
}

// UploadDoc godoc
// @Summary Upload Documents (permission = superadmin)
// @Id UploadDoc
// @Tags Parking Space
// @Security Token
// @Accept multipart/form-data
// @Produce json
// @Param id path integer true "parking_space_id"
// @Param file formData file true "document parking space"
// @Success 201 {object} model.ResponseSuccess "{ "message": "Upload Document Parking Space Success" }"
// @Failure 400 {object} model.ResponseErrors "code: 400, message: "Form/Upoad file error" "
// @Failure 500 {object} model.ResponseErrors "code: 500, message: "Update Document Meta data failed" "
// @Router /parkingspace/{id}/doc/upload [put]
func (s *ParkingSpaceController) UploadDoc(c *gin.Context) {
	id := util.ParamIDToInt64(c.Param("id"))
	docMeta := s.PSpaceSvc.GetOne(id).DocumentsMeta
	path := cfg.Config.Path.Doc
	// images upload when create success
	// image 1
	file, err := c.FormFile("file")
	if err != nil {
		model.ResponseError(c, "Form file error", http.StatusBadRequest)
		return
	}
	filename := util.GeneratedUUID(filepath.Base(file.Filename))
	pathFile := path + filename
	if err := c.SaveUploadedFile(file, pathFile); err != nil {
		model.ResponseError(c, "Upload file error", http.StatusBadRequest)
		return
	}
	docMeta = append(docMeta, filename)
	if err := s.PSpaceSvc.UploadDoc(docMeta, id); err != nil {
		model.ResponseError(c, "Update Document Meta data failed", http.StatusInternalServerError)
		return
	}
	model.ResponseUpdated(c, gin.H{"message": "Upload Document Parking Space Success"})
	return
}

// GetSalesByPSpaceID godoc
// @Summary Get Parking Space Sales Pagination by ParkingSpaceID (permission = admin)
// @Id GetSalesByPSpaceID
// @Tags Parking Space
// @Security Token
// @Param id path integer true "ParkingSpaceID"
// @Param limit query string false " " default(10)
// @Param page query string false " " default(1)
// @Param sort query string false "e.g.: id,desc / id,asc"
// @Success 200 {object} model.ResponseSuccess "model.Pagination"
// @Router /parkingspace/sales/{id} [get]
func (s *ParkingSpaceController) GetSalesByPSpaceID(c *gin.Context) {
	id := util.ParamIDToInt64(c.Param("id"))
	limit := util.ParamIDToInt(c.DefaultQuery("limit", "10"))
	page := util.ParamIDToInt(c.DefaultQuery("page", "1"))
	sort := util.SortedBy(c.QueryArray("sort"))
	data := s.SalesSvc.GetBySpaceID(id, limit, page, sort)
	model.ResponsePagination(c, data)
	return
}

// CreateSales godoc
// @Summary Create Parking Space Sales (permission = superadmin)
// @Id CreateSales
// @Tags Parking Space
// @Security Token
// @Accept multipart/form-data
// @Produce json
// @Param parkingSpaceId formData integer true "Parking Space ID"
// @Param startDate formData string true "Start Date yyyy-MM-dd"
// @Param endDate formData string true "End Date yyyy-MM-dd"
// @Param totalSlot formData integer true "Total Slot"
// @Param point formData integer true "Point"
// @Success 200 {object} entity.ParkingSpaceSales "data: entity.ParkingSpaceSales, message: "Create Sales Succeed" "
// @Failure 422 {object} model.ResponseErrors "code: 422, message: "Invalid request" "
// @Failure 500 {object} model.ResponseErrors "code: 500, message: "Create Sales Failed" "
// @Router /parkingspace/sales/create [post]
func (s *ParkingSpaceController) CreateSales(c *gin.Context) {
	req := model.ReqSalesCreate{}
	if err := c.ShouldBind(&req); err != nil {
		model.ResponseError(c, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	startDate := strings.Split(req.StartDate.Format("2006-01-02 15:04:05"), " ")
	endDate := strings.Split(req.EndDate.Format("2006-01-02 15:04:05"), " ")
	parkingSpaceStartTime := time.Now().Format("2006-01-02") + " " + startDate[1]
	parkingSpaceEndTime := time.Now().Format("2006-01-02") + " " + endDate[1]

	parkingSpaceStartdate := startDate[0] + " " + "00:00:00"
	parkingSpaceEnddate := endDate[0] + " " + "00:00:00"
	var (
		data *entity.ParkingSpaceSales
		err  error
	)
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)

	spaceID := req.ParkingSpaceID
	getParkingSpace := s.PSpaceSvc.GetOne(spaceID)

	parkingSpaceStart := util.DateTimeSwap(getParkingSpace.StartContract, getParkingSpace.StartTime)
	parkingSpaceEnd := util.DateTimeSwap(getParkingSpace.EndContract, getParkingSpace.EndTime)

	startMilliseconds := util.DateStringToMilliSeconds(parkingSpaceStartTime)
	endMilliseconds := util.DateStringToMilliSeconds(parkingSpaceEndTime)
	startDateMilliseconds := util.DateStringToMilliSeconds(parkingSpaceStartdate)
	endDateMilliseconds := util.DateStringToMilliSeconds(parkingSpaceEnddate)

	getPreviousSlot := s.SalesSvc.GetSalesBySpace(spaceID, parkingSpaceStart, parkingSpaceEnd)
	var sumSlot int
	sumSlot += req.TotalSlot
	if len(*getPreviousSlot) > 0 {
		var (
			parseParkingSpaceStartMilliseconds     int64
			parseParkingSpaceEndMilliseconds       int64
			parseParkingSpaceStartDateMilliseconds int64
			parseParkingSpaceEndDateMilliseconds   int64
			startRecordDateString                  string
			endRecordDateString                    string
		)
		for _, value := range *getPreviousSlot {
			startRecordDateString = time.Now().Format("2006-01-02") + " " + value.StartDate.Format("15:04:05")
			endRecordDateString = time.Now().Format("2006-01-02") + " " + value.EndDate.Format("15:04:05")

			parseParkingSpaceStartMilliseconds = util.DateStringToMilliSeconds(startRecordDateString)
			parseParkingSpaceEndMilliseconds = util.DateStringToMilliSeconds(endRecordDateString)
			parseParkingSpaceStartDateMilliseconds = util.DateStringToMilliSeconds(value.StartDate.Format("2006-01-02") + " " + "00:00:00")
			parseParkingSpaceEndDateMilliseconds = util.DateStringToMilliSeconds(value.EndDate.Format("2006-01-02") + " " + "00:00:00")

			if (startDateMilliseconds < parseParkingSpaceStartDateMilliseconds && endDateMilliseconds > parseParkingSpaceEndDateMilliseconds) ||
				(startDateMilliseconds >= parseParkingSpaceStartDateMilliseconds && endDateMilliseconds <= parseParkingSpaceEndDateMilliseconds) ||
				(startDateMilliseconds <= parseParkingSpaceStartDateMilliseconds && endDateMilliseconds >= parseParkingSpaceStartDateMilliseconds &&
					endDateMilliseconds <= parseParkingSpaceEndDateMilliseconds) || (startDateMilliseconds >= parseParkingSpaceStartDateMilliseconds &&
				endDateMilliseconds >= parseParkingSpaceEndDateMilliseconds && startDateMilliseconds <= parseParkingSpaceEndDateMilliseconds) {
				if startMilliseconds >= parseParkingSpaceStartMilliseconds && endMilliseconds <= parseParkingSpaceEndMilliseconds {
					sumSlot += value.TotalSlot
				} else if startMilliseconds <= parseParkingSpaceStartMilliseconds && endMilliseconds > parseParkingSpaceStartMilliseconds &&
					endMilliseconds <= parseParkingSpaceEndMilliseconds {
					sumSlot += value.TotalSlot
				} else if startMilliseconds <= parseParkingSpaceStartMilliseconds && startMilliseconds < parseParkingSpaceEndMilliseconds &&
					endMilliseconds >= parseParkingSpaceEndMilliseconds {
					sumSlot += value.TotalSlot
				} else if startMilliseconds >= parseParkingSpaceStartMilliseconds && endMilliseconds >= parseParkingSpaceEndMilliseconds &&
					startMilliseconds < parseParkingSpaceEndMilliseconds {
					sumSlot += value.TotalSlot
				}
			}
		}
	}

	getOne := s.PSpaceSvc.GetOne(spaceID)
	if sumSlot > getOne.TotalSpace {
		model.ResponseError(c, "Total Slot You entered already reached maximum space capacity!", http.StatusNotAcceptable)
		return
	}
	if data, err = s.SalesSvc.CreateSales(&req, jwtModel.UserID); err != nil {
		model.ResponseError(c, "Create Sales Failed", http.StatusInternalServerError)
		return
	}
	model.ResponseCreated(c, gin.H{"data": data, "message": "Create Sales Succeed"})
	return
}

// GetByID godoc
// @Summary Get Parking Space by ID (permission = superadmin)
// @Id GetByIDPspace
// @Tags Parking Space
// @Security Token
// @Param id path integer true "id"
// @Success 200 {object} entity.ParkingSpace "data: entity.ParkingSpace"
// @Router /parkingspace/show/{id} [get]
func (s *ParkingSpaceController) GetByID(c *gin.Context) {
	id := util.ParamIDToInt64(c.Param("id"))
	data := s.PSpaceSvc.GetOne(id)
	model.ResponseJSON(c, data)
	return
}

// DeleteByID godoc
// @Summary Delete Parking Space by ID (permission = superadmin)
// @Id DeleteByIDPSpace
// @Tags Parking Space
// @Security Token
// @Param id path integer true "parking space"
// @Success 200 {object} model.ResponseSuccess "message: "Success" "
// @Failure 500 {object} model.ResponseErrors "code: 500, message: "Delete Parking Space Failed" "
// @Router /parkingspace/{id}/delete [delete]
func (s *ParkingSpaceController) DeleteByID(c *gin.Context) {
	id := util.ParamIDToInt64(c.Param("id"))
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	if err := s.PSpaceSvc.DeleteByID(id, jwtModel.UserID); err != nil {
		model.ResponseError(c, "Delete Parking Space Failed", http.StatusInternalServerError)
		return
	}
	model.ResponseJSON(c, gin.H{"message": "Success"})
	return
}

// GetSalesByID godoc
// @Summary Get Parking Space Sales by ID (permission = superadmin)
// @Id GetSalesByID
// @Tags Parking Space
// @Security Token
// @Param id path integer true "Parking Space"
// @Success 200 {object} entity.ParkingSpaceSales "data: entity.ParkingSpaceSales"
// @Router /parkingspace/sales/{id}/info [get]
func (s *ParkingSpaceController) GetSalesByID(c *gin.Context) {
	id := util.ParamIDToInt64(c.Param("id"))
	model.ResponseJSON(c, s.SalesSvc.GetByID(id))
	return
}

// GetSpaceBySalesDate godoc
// @Summary Get Parking Space by Sales Date (permission = superadmin)
// @Id GetSpaceBySalesDate
// @Tags Parking Space
// @Security Token
// @Param id path integer true "Parking Space"
// @Param startDate query string true "yyyy-mm-dd hh:mm:ss"
// @Param endDate query string true "yyyy-mm-dd hh:mm:ss"
// @Success 200 {object} entity.ParkingSpace "data: entity.ParkingSpace"
// @Router /parkingspace/spacesalesdate/{id} [get]
func (s *ParkingSpaceController) GetSpaceBySalesDate(c *gin.Context) {
	id := util.ParamIDToInt64(c.Param("id"))

	startDate := strings.Split(c.Query("startDate"), " ")
	endDate := strings.Split(c.Query("endDate"), " ")

	// get total space
	getParkingSpace := s.PSpaceSvc.GetOne(id)
	parkingSpaceStartTime := time.Now().Format("2006-01-02") + " " + startDate[1]
	parkingSpaceEndTime := time.Now().Format("2006-01-02") + " " + endDate[1]

	parkingSpaceStartdate := startDate[0] + " " + "00:00:00"
	parkingSpaceEnddate := endDate[0] + " " + "00:00:00"
	parkingSpaceStart := util.DateTimeSwap(getParkingSpace.StartContract, getParkingSpace.StartTime)
	parkingSpaceEnd := util.DateTimeSwap(getParkingSpace.EndContract, getParkingSpace.EndTime)
	startMilliseconds := util.DateStringToMilliSeconds(parkingSpaceStartTime)
	endMilliseconds := util.DateStringToMilliSeconds(parkingSpaceEndTime)
	startDateMilliseconds := util.DateStringToMilliSeconds(parkingSpaceStartdate)
	endDateMilliseconds := util.DateStringToMilliSeconds(parkingSpaceEnddate)
	// sum all related space slot by previous sales with range date given and substract with give total space
	getSales := s.SalesSvc.GetSalesBySpace(id, parkingSpaceStart, parkingSpaceEnd)
	if len(*getSales) > 0 {
		var (
			sumSlot                                int
			parseParkingSpaceStartMilliseconds     int64
			parseParkingSpaceEndMilliseconds       int64
			parseParkingSpaceStartDateMilliseconds int64
			parseParkingSpaceEndDateMilliseconds   int64
			startRecordDateString                  string
			endRecordDateString                    string
		)
		for _, value := range *getSales {
			startRecordDateString = time.Now().Format("2006-01-02") + " " + value.StartDate.Format("15:04:05")
			endRecordDateString = time.Now().Format("2006-01-02") + " " + value.EndDate.Format("15:04:05")
			parseParkingSpaceStartMilliseconds = util.DateStringToMilliSeconds(startRecordDateString)
			parseParkingSpaceEndMilliseconds = util.DateStringToMilliSeconds(endRecordDateString)
			parseParkingSpaceStartDateMilliseconds = util.DateStringToMilliSeconds(value.StartDate.Format("2006-01-02") + " " + "00:00:00")
			parseParkingSpaceEndDateMilliseconds = util.DateStringToMilliSeconds(value.EndDate.Format("2006-01-02") + " " + "00:00:00")
			if (startDateMilliseconds < parseParkingSpaceStartDateMilliseconds && endDateMilliseconds > parseParkingSpaceEndDateMilliseconds) ||
				(startDateMilliseconds >= parseParkingSpaceStartDateMilliseconds && endDateMilliseconds <= parseParkingSpaceEndDateMilliseconds) ||
				(startDateMilliseconds <= parseParkingSpaceStartDateMilliseconds && endDateMilliseconds >= parseParkingSpaceStartDateMilliseconds &&
					endDateMilliseconds <= parseParkingSpaceEndDateMilliseconds) || (startDateMilliseconds >= parseParkingSpaceStartDateMilliseconds &&
				endDateMilliseconds >= parseParkingSpaceEndDateMilliseconds && startDateMilliseconds <= parseParkingSpaceEndDateMilliseconds) {
				if startMilliseconds >= parseParkingSpaceStartMilliseconds && endMilliseconds <= parseParkingSpaceEndMilliseconds {
					sumSlot += value.TotalSlot
				} else if startMilliseconds <= parseParkingSpaceStartMilliseconds && endMilliseconds > parseParkingSpaceStartMilliseconds &&
					endMilliseconds <= parseParkingSpaceEndMilliseconds {
					sumSlot += value.TotalSlot
				} else if startMilliseconds <= parseParkingSpaceStartMilliseconds && startMilliseconds < parseParkingSpaceEndMilliseconds &&
					endMilliseconds >= parseParkingSpaceEndMilliseconds {
					sumSlot += value.TotalSlot
				} else if startMilliseconds >= parseParkingSpaceStartMilliseconds && endMilliseconds >= parseParkingSpaceEndMilliseconds &&
					startMilliseconds < parseParkingSpaceEndMilliseconds {
					sumSlot += value.TotalSlot
				}
			}
		}
		totalSpace := getParkingSpace.TotalSpace - sumSlot
		if totalSpace < 0 {
			totalSpace = 0
		}
		getParkingSpace.TotalSpace = totalSpace
	}
	model.ResponseJSON(c, getParkingSpace)
	return
}

// UpdateSales godoc
// @Summary Update Parking Space Sales (permission = superadmin)
// @Id UpdUpdateSalesate
// @Tags Parking Space
// @Security Token
// @Param id path integer true "Parking Space Sales"
// @Param req body model.ReqSalesUpdate true "Update Sales"
// @Success 200 {object} entity.ParkingSpaceSales "{ "data": Model }"
// @Failure 422 {object} model.ResponseErrors "code: 422, message: "Invalid request" "
// @Failure 500 {object} model.ResponseErrors "code: 500, message: "Update Sales Failed" "
// @Router /parkingspace/{id}/sales/update [put]
func (s *ParkingSpaceController) UpdateSales(c *gin.Context) {
	req := model.ReqSalesUpdate{}
	id := util.ParamIDToInt64(c.Param("id"))
	if err := c.ShouldBindJSON(&req); err != nil {
		model.ResponseError(c, "Invalid Request", http.StatusUnprocessableEntity)
		return
	}

	var (
		data *entity.ParkingSpaceSales
		err  error
	)
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)

	//check current data
	currentData := s.SalesSvc.GetByID(id)
	startDate := currentData.StartDate.Format("2006-01-02") + " " + "00:00:00"
	endDate := currentData.EndDate.Format("2006-01-02") + " " + "00:00:00"
	startTime := time.Now().Format("2006-01-02") + " " + currentData.StartDate.Format("15:04:05")
	endTime := time.Now().Format("2006-01-02") + " " + currentData.EndDate.Format("15:04:05")

	startDateMilliseconds := util.DateStringToMilliSeconds(startDate)
	endDateMilliseconds := util.DateStringToMilliSeconds(endDate)
	startTimeMilliseconds := util.DateStringToMilliSeconds(startTime)
	endTimeMilliseconds := util.DateStringToMilliSeconds(endTime)
	//check if the current parkingSpaceId is same as user inputted or not
	//this logic is for check the inputted  slot reached maximum space or not
	if int64(req.ParkingSpaceID) == currentData.ParkingSpaceID {
		var newTotalSlot int
		//differentiate old totalSlot with new totalSlot and do calculation
		sub := currentData.TotalSlot - req.TotalSlot
		newTotalSlot = currentData.TotalSlot - sub
		spaceID := int64(req.ParkingSpaceID)
		getParkingSpace := s.PSpaceSvc.GetOne(spaceID)
		parkingSpaceStart := util.DateTimeSwap(getParkingSpace.StartContract, getParkingSpace.StartTime)
		parkingSpaceEnd := util.DateTimeSwap(getParkingSpace.EndContract, getParkingSpace.EndTime)

		getPreviousSlot := s.SalesSvc.GetSalesBySpace(spaceID, parkingSpaceStart, parkingSpaceEnd)
		var (
			sumSlot                                int
			parseParkingSpaceStartMilliseconds     int64
			parseParkingSpaceEndMilliseconds       int64
			parseParkingSpaceStartDateMilliseconds int64
			parseParkingSpaceEndDateMilliseconds   int64
			startRecordDateString                  string
			endRecordDateString                    string
			startRecordTimeString                  string
			endRecordTimeString                    string
		)
		sumSlot += newTotalSlot

		for _, value := range *getPreviousSlot {
			startRecordDateString = value.StartDate.Format("2006-01-02") + " " + "00:00:00"
			endRecordDateString = value.EndDate.Format("2006-01-02") + " " + "00:00:00"
			startRecordTimeString = time.Now().Format("2006-01-02") + " " + value.StartDate.Format("15:04:05")
			endRecordTimeString = time.Now().Format("2006-01-02") + " " + value.EndDate.Format("15:04:05")
			parseParkingSpaceStartDateMilliseconds = util.DateStringToMilliSeconds(startRecordDateString)
			parseParkingSpaceEndDateMilliseconds = util.DateStringToMilliSeconds(endRecordDateString)
			parseParkingSpaceStartMilliseconds = util.DateStringToMilliSeconds(startRecordTimeString)
			parseParkingSpaceEndMilliseconds = util.DateStringToMilliSeconds(endRecordTimeString)
			if value.ID != id {
				if (startDateMilliseconds < parseParkingSpaceStartDateMilliseconds && endDateMilliseconds > parseParkingSpaceEndDateMilliseconds) ||
					(startDateMilliseconds >= parseParkingSpaceStartDateMilliseconds && endDateMilliseconds <= parseParkingSpaceEndDateMilliseconds) ||
					(startDateMilliseconds <= parseParkingSpaceStartDateMilliseconds && endDateMilliseconds >= parseParkingSpaceStartDateMilliseconds &&
						endDateMilliseconds <= parseParkingSpaceEndDateMilliseconds) || (startDateMilliseconds >= parseParkingSpaceStartDateMilliseconds &&
					endDateMilliseconds >= parseParkingSpaceEndDateMilliseconds && startDateMilliseconds <= parseParkingSpaceEndDateMilliseconds) {
					if startTimeMilliseconds >= parseParkingSpaceStartMilliseconds && endTimeMilliseconds <= parseParkingSpaceEndMilliseconds {
						sumSlot += value.TotalSlot
					} else if startTimeMilliseconds <= parseParkingSpaceStartMilliseconds && endTimeMilliseconds > parseParkingSpaceStartMilliseconds &&
						endTimeMilliseconds <= parseParkingSpaceEndMilliseconds {
						sumSlot += value.TotalSlot
					} else if startTimeMilliseconds <= parseParkingSpaceStartMilliseconds && startTimeMilliseconds < parseParkingSpaceEndMilliseconds &&
						endTimeMilliseconds >= parseParkingSpaceEndMilliseconds {
						sumSlot += value.TotalSlot
					} else if startTimeMilliseconds >= parseParkingSpaceStartMilliseconds && endTimeMilliseconds >= parseParkingSpaceEndMilliseconds &&
						startTimeMilliseconds < parseParkingSpaceEndMilliseconds {
						sumSlot += value.TotalSlot
					}
				}
			}
		}
		if sumSlot > getParkingSpace.TotalSpace {
			model.ResponseError(c, "Total Slot You entered already reached maximum space capacity!", http.StatusNotAcceptable)
			return
		}
	} else {
		spaceID := int64(req.ParkingSpaceID)
		getParkingSpace := s.PSpaceSvc.GetOne(spaceID)
		parkingSpaceStart := util.DateTimeSwap(getParkingSpace.StartContract, getParkingSpace.StartTime)
		parkingSpaceEnd := util.DateTimeSwap(getParkingSpace.EndContract, getParkingSpace.EndTime)

		getPreviousSlot := s.SalesSvc.GetSalesBySpace(spaceID, parkingSpaceStart, parkingSpaceEnd)
		var (
			sumSlot                                int
			parseParkingSpaceStartMilliseconds     int64
			parseParkingSpaceEndMilliseconds       int64
			parseParkingSpaceStartDateMilliseconds int64
			parseParkingSpaceEndDateMilliseconds   int64
			startRecordDateString                  string
			endRecordDateString                    string
			startRecordTimeString                  string
			endRecordTimeString                    string
		)
		sumSlot += req.TotalSlot
		if len(*getPreviousSlot) > 0 {
			for _, value := range *getPreviousSlot {
				startRecordDateString = value.StartDate.Format("2006-01-02") + " " + "00:00:00"
				endRecordDateString = value.EndDate.Format("2006-01-02") + " " + "00:00:00"
				startRecordTimeString = time.Now().Format("2006-01-02") + " " + value.StartDate.Format("15:04:05")
				endRecordTimeString = time.Now().Format("2006-01-02") + " " + value.EndDate.Format("15:04:05")
				parseParkingSpaceStartDateMilliseconds = util.DateStringToMilliSeconds(startRecordDateString)
				parseParkingSpaceEndDateMilliseconds = util.DateStringToMilliSeconds(endRecordDateString)
				parseParkingSpaceStartMilliseconds = util.DateStringToMilliSeconds(startRecordTimeString)
				parseParkingSpaceEndMilliseconds = util.DateStringToMilliSeconds(endRecordTimeString)

				if (startDateMilliseconds < parseParkingSpaceStartDateMilliseconds && endDateMilliseconds > parseParkingSpaceEndDateMilliseconds) ||
					(startDateMilliseconds >= parseParkingSpaceStartDateMilliseconds && endDateMilliseconds <= parseParkingSpaceEndDateMilliseconds) ||
					(startDateMilliseconds <= parseParkingSpaceStartDateMilliseconds && endDateMilliseconds >= parseParkingSpaceStartDateMilliseconds &&
						endDateMilliseconds <= parseParkingSpaceEndDateMilliseconds) || (startDateMilliseconds >= parseParkingSpaceStartDateMilliseconds &&
					endDateMilliseconds >= parseParkingSpaceEndDateMilliseconds && startDateMilliseconds <= parseParkingSpaceEndDateMilliseconds) {
					if startTimeMilliseconds >= parseParkingSpaceStartMilliseconds && endTimeMilliseconds <= parseParkingSpaceEndMilliseconds {
						sumSlot += value.TotalSlot
					} else if startTimeMilliseconds <= parseParkingSpaceStartMilliseconds && endTimeMilliseconds > parseParkingSpaceStartMilliseconds &&
						endTimeMilliseconds <= parseParkingSpaceEndMilliseconds {
						sumSlot += value.TotalSlot
					} else if startTimeMilliseconds <= parseParkingSpaceStartMilliseconds && startTimeMilliseconds < parseParkingSpaceEndMilliseconds &&
						endTimeMilliseconds >= parseParkingSpaceEndMilliseconds {
						sumSlot += value.TotalSlot
					} else if startTimeMilliseconds >= parseParkingSpaceStartMilliseconds && endTimeMilliseconds >= parseParkingSpaceEndMilliseconds &&
						startTimeMilliseconds < parseParkingSpaceEndMilliseconds {
						sumSlot += value.TotalSlot
					}
				}
			}
		}
		if sumSlot > getParkingSpace.TotalSpace {
			model.ResponseError(c, "Total Slot You entered already reached maximum space capacity!", http.StatusNotAcceptable)
			return
		}
	}
	if data, err = s.SalesSvc.Update(&req, id, jwtModel.UserID); err != nil {
		model.ResponseError(c, "Update Sales Failed", http.StatusInternalServerError)
		return
	}
	model.ResponseJSON(c, data)
	return
}

// DeleteImage godoc
// @Summary Delete Image by Image Name (permission = superadmin)
// @Id DeleteImage
// @Tags Parking Space
// @Security Token
// @Accept multipart/form-data
// @Produce json
// @Param id path integer true "parking space"
// @Param req body model.ReqDeleteAsset true "Delete Image"
// @Success 200 {object} model.ResponseSuccess "message: "Success""
// @Failure 422 {object} model.ResponseErrors "code: 422, message: "Invalid request" "
// @Failure 500 {object} model.ResponseErrors "code: 500, message: "Delete Image Failed" "
// @Router /parkingspace/{id}/image/delete [put]
func (s *ParkingSpaceController) DeleteImage(c *gin.Context) {
	req := model.ReqDeleteAsset{}
	if err := c.ShouldBindJSON(&req); err != nil {
		model.ResponseError(c, "Invalid Request", http.StatusUnprocessableEntity)
		return
	}
	pSpaceID := util.ParamIDToInt64(c.Param("id"))
	imgMeta := s.PSpaceSvc.GetOne(pSpaceID).ImagesMeta
	imgName := req.Filename
	imagePath := cfg.Config.Path.Image
	updateImgMeta := make(pq.StringArray, 0)
	for _, v := range imgMeta {
		if v != imgName {
			updateImgMeta = append(updateImgMeta, v)
		}
	}
	if err := s.PSpaceSvc.UploadImage(
		updateImgMeta, pSpaceID); err != nil {
		model.ResponseError(c, "Delete Image Failed", http.StatusInternalServerError)
		return
	}
	os.Remove(imagePath + imgName)
	model.ResponseJSON(c, gin.H{"message": "Success"})
	return
}

// DeleteDoc godoc
// @Summary Delete Document Parking Space by Document Name (permission = superadmin)
// @Id DeleteDoc
// @Tags Parking Space
// @Security Token
// @Param id path integer true "parking space"
// @Param req body model.ReqDeleteAsset true "Delete Document"
// @Success 200 {object} model.ResponseSuccess "message: "Success" "
// @Failure 422 {object} model.ResponseErrors "code: 422, message: "Invalid request" "
// @Failure 500 {object} model.ResponseErrors "code: 500, message: "Failed to Delete Doc" "
// @Router /parkingspace/{id}/doc/delete [put]
func (s *ParkingSpaceController) DeleteDoc(c *gin.Context) {
	req := model.ReqDeleteAsset{}
	if err := c.ShouldBindJSON(&req); err != nil {
		model.ResponseError(c, "Invalid Request", http.StatusUnprocessableEntity)
		return
	}
	pSpaceID := util.ParamIDToInt64(c.Param("id"))
	docName := req.Filename
	docMeta := s.PSpaceSvc.GetOne(pSpaceID).DocumentsMeta
	docPath := cfg.Config.Path.Doc
	updateDocMeta := make(pq.StringArray, 0)
	for _, v := range docMeta {
		if v != docName {
			updateDocMeta = append(updateDocMeta, v)
		}
	}
	if err := s.PSpaceSvc.UploadDoc(
		updateDocMeta, pSpaceID); err != nil {
		model.ResponseError(c, "Failed to Delete Doc", http.StatusInternalServerError)
		return
	}
	os.Remove(docPath + docName)
	model.ResponseJSON(c, gin.H{"message": "Success"})
	return
}

// DeleteSalesByID godoc
// @Summary Get Parking Space by ID (permission = superadmin)
// @Id DeleteSalesByID
// @Tags Parking Space
// @Security Token
// @Param id path integer true "parking space"
// @Success 200 {object} model.ResponseSuccess "message: "Success" "
// @Failure 500 {object} model.ResponseErrors "code: 500, message: "Failed to Delete ParkingSpaceSales" "
// @Router /parkingspace/{id}/sales/delete [delete]
func (s *ParkingSpaceController) DeleteSalesByID(c *gin.Context) {
	id := util.ParamIDToInt64(c.Param("id"))
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	if err := s.SalesSvc.DeleteByID(id, jwtModel.UserID); err != nil {
		model.ResponseError(c,
			"Failed to Delete ParkingSpaceSales",
			http.StatusInternalServerError)
		return
	}
	model.ResponseJSON(c, gin.H{"message": "Success"})
	return
}

// SearchSalesByName godoc
// @Summary Search Parking Space Sales by Name or Address (permission = admin)
// @Id SearchSalesByName
// @Tags Parking Space
// @Security Token
// @Param key path string true "keyword"
// @Param limit query string false " " default(10)
// @Param page query string false " " default(1)
// @Param sort query string false "e.g.: id,desc / id,asc"
// @Success 200 {object} model.ResponseSuccess " "
// @Router /parkingspace/search/sales/{key} [get]
func (s *ParkingSpaceController) SearchSalesByName(c *gin.Context) {
	limit := util.ParamIDToInt(c.DefaultQuery("limit", "10"))
	page := util.ParamIDToInt(c.DefaultQuery("page", "1"))
	sort := util.SortedBy(c.QueryArray("sortSales"))
	key := c.Param("key")
	data := s.SalesSvc.FindLikeNameBackoffice(key, limit, page, sort)
	model.ResponsePagination(c, data)
	return
}

// GetSales godoc
// @Summary Get All Parking Space Sales Pagination (permission = admin)
// @Id GetSales
// @Tags Parking Space
// @Security Token
// @Param limit query string false " " default(10)
// @Param page query string false " " default(1)
// @Param sort query string false "e.g.: id,desc / id,asc"
// @Success 200 {object} model.ResponseSuccess "model.Pagination"
// @Router /parkingspace/sales [get]
func (s *ParkingSpaceController) GetSales(c *gin.Context) {
	limit := util.ParamIDToInt(c.DefaultQuery("limit", "10"))
	page := util.ParamIDToInt(c.DefaultQuery("page", "1"))
	sort := util.SortedBy(c.QueryArray("sort"))
	data := s.SalesSvc.GetAllBackoffice(limit, page, sort)
	model.ResponsePagination(c, data)
	return
}

// GetSalesList godoc
// @Summary Get All Parking Space Sales List no pagination (permission = superadmin)
// @Tags Parking Space
// @Security Token
// @Success 200 {object} model.ResponseSuccess "entity.ParkingSpaceSales"
// @Router /parkingspace/list/sales [get]
func (s *ParkingSpaceController) GetSalesList(c *gin.Context) {
	model.ResponseJSON(c, s.SalesSvc.GetAllList())
	return
}
