package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"strconv"
	"time"

	"streetbox.id/app/trxrefund"
	"streetbox.id/cfg"

	//this images intended for excelize addpicture, without this, image will not loaded on xls file
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"streetbox.id/app/merchant"
	"streetbox.id/app/payment"
	"streetbox.id/app/trx"
	"streetbox.id/entity"
	"streetbox.id/model"
	"streetbox.id/util"
)

// TrxController ...
type TrxController struct {
	TrxSvc       trx.ServiceInterface
	MerchantSvc  merchant.ServiceInterface
	TrxRefundSvc trxrefund.ServiceInterface
	PaymentSvc   payment.ServiceInterface
}

// GetInfo godoc
// @Summary Get Transaction Info (permission = consumer)
// @Id GetInfotrx
// @Tags Transaction
// @Security token
// @Param trxId path string true "transaction ID"
// @Success 200 {object} model.ResTrx "data: model.ResTrx"
// @Router /trx/info/{trxId} [get]
func (r *TrxController) GetInfo(c *gin.Context) {
	model.ResponseJSON(c, model.ResTrx{})
	return
}

// CreateSyncTrx godoc
// @Summary Create Sync Trx (permission = merchant)
// @Id CreateSyncTrx
// @Tags Transaction
// @Security Token
// @Param req body model.ReqCreateSyncTrx true " "
// @Success 200 {object} entity.TrxSync "data: entity.TrxSync"
// @Router /trx/createsync [post]
func (r *TrxController) CreateSyncTrx(c *gin.Context) {
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	// merchant only
	if jwtModel.RoleName != "foodtruck" {
		model.ResponseError(
			c,
			"Sorry, Food Truck Role Only",
			http.StatusUnprocessableEntity)
		return
	}

	req := model.ReqCreateSyncTrx{}
	if err := c.ShouldBindJSON(&req); err != nil {
		model.ResponseError(c,
			"Invalid Request", http.StatusUnprocessableEntity)
		return
	}

	merchant := r.MerchantSvc.GetInfo(jwtModel.UserID)
	if merchant == nil {
		model.ResponseError(
			c,
			"Merchant ID not found",
			http.StatusUnprocessableEntity)
		return
	}

	data := r.TrxSvc.CreateSyncTrx(&req, merchant.ID)

	if data == nil {
		model.ResponseError(c,
			"Failed to Create Transaction", http.StatusInternalServerError)
		return
	}

	model.ResponseJSON(c, data)
	return
}

// CreateTrxVisit godoc
// @Summary Create Homevisit Trx (permission = consumer)
// @Id CreateTrxVisit
// @Tags Transaction
// @Security Token
// @Param req body model.ReqCreateVisitTrx true " "
// @Success 200 {object} entity.TrxVisit "data: entity.TrxVisit"
// @Router /trx/homevisit [post]
func (r *TrxController) CreateTrxVisit(c *gin.Context) {
	req := model.ReqCreateVisitTrx{}
	if err := c.ShouldBindJSON(&req); err != nil {
		model.ResponseError(c,
			err.Error(), http.StatusUnprocessableEntity)
		return
	}
	resCreateTrx := new(entity.TrxVisit)
	var err error
	if resCreateTrx, err = r.TrxSvc.CreateTrxVisit(&req); err != nil {
		model.ResponseError(c, "Create Transaction Visit Failed", http.StatusUnprocessableEntity)
		r.TrxSvc.DeleteTrxByID(req.TrxID)
		return
	}
	model.ResponseJSON(c, resCreateTrx)

	merchantID := r.TrxSvc.GetMerchantIDByTrxID(req.TrxID)
	merchant := r.MerchantSvc.GetInfo(merchantID)
	SMTPEmail := cfg.Config.Smpt.Email
	SMTPPassword := cfg.Config.Smpt.Password
	SMTPHost := cfg.Config.Smpt.Host
	SMTPPort := cfg.Config.Smpt.Port
	auth := smtp.PlainAuth("", SMTPEmail, SMTPPassword, SMTPHost)
	smtpAddr := fmt.Sprintf("%s:%s", SMTPHost, SMTPPort)
	to := []string{merchant.Email}
	msg := []byte("To: " + merchant.Email + "\r\n" +
		"Subject: Homevisit Notification\r\n" +
		"\r\n" +
		"Pesanan homevisit baru muncul. Silahkan proses pesanan melalui backoffice merchant.\r\n")
	err = smtp.SendMail(smtpAddr, auth, SMTPEmail, to, msg)
	if err != nil {
		log.Printf(fmt.Sprintf("ERROR: Sending homevisit email notification error: %s", err.Error()))
	}
	return
}

// CreateTrxRefundSpace godoc
// @Summary Create Refund for Parking Space Sales(permission = superadmin)
// @Id CreateTrxSpaceRefund
// @Tags Transaction
// @Security Token
// @Param req body model.ReqRefundParkingSpaceSales true " "
// @Success 200 {object} model.ResponseSuccess "message: Refund Space Transaction Success!"
// @Router /trx/refund/space [post]
func (r *TrxController) CreateTrxRefundSpace(c *gin.Context) {
	req := model.ReqRefundParkingSpaceSales{}

	if err := c.ShouldBindJSON(&req); err != nil {
		model.ResponseError(c,
			"Invalid Request", http.StatusUnprocessableEntity)
		return
	}

	if err := r.TrxRefundSvc.CreateRefundParkingSpace(&req); err != nil {
		model.ResponseError(c, "Create Refund Parking Space Failed", http.StatusUnprocessableEntity)
		return
	}

	model.ResponseJSON(c, "Refund Space Transaction Success!")
	return
}

// CreateTrxRefundVisit godoc
// @Summary Create Refund for Home Visit (permission = merchant)
// @Id CreateTrxRefundVisit
// @Tags Transaction
// @Security Token
// @Param req body model.ReqRefundHomeVisit true " "
// @Success 200 {object} model.ResponseSuccess "message: Refund Space Transaction Success!"
// @Router /trx/refund/visit [post]
func (r *TrxController) CreateTrxRefundVisit(c *gin.Context) {
	req := model.ReqRefundHomeVisit{}

	if err := c.ShouldBindJSON(&req); err != nil {
		model.ResponseError(c,
			"Invalid Request", http.StatusUnprocessableEntity)
		return
	}

	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	merchant := r.MerchantSvc.GetInfo(jwtModel.UserID)

	if check := r.TrxSvc.CheckTrxStatusPendingByTrxHomeVisitSalesID(req.ID); check > 0 {
		model.ResponseError(c, "Cannot Refund Pending Status!", http.StatusUnprocessableEntity)
		return
	}

	if err := r.TrxRefundSvc.CreateRefundHomeVisit(&req, merchant.ID); err != nil {
		model.ResponseError(c, "Create Refund Parking Space Failed", http.StatusUnprocessableEntity)
		return
	}

	model.ResponseJSON(c, "Refund Visit Transaction Success!")
	return
}

// ParseTrxJobs ..
func (r *TrxController) ParseTrxJobs() {
	getActiveSyncData := r.TrxSvc.GetSyncTrx()
	var (
		data            string
		body            []byte
		reqTrxOrderList model.ReqTrxOrderList
		db              *gorm.DB
	)

	for _, value := range *getActiveSyncData {
		data = value.Data
		body = []byte(data)
		json.Unmarshal(body, &reqTrxOrderList)
		r.TrxSvc.CountTrx()
		if err := r.TrxSvc.CreateTrxOrder(&reqTrxOrderList, reqTrxOrderList.Order.UserID, value.MerchantID, value.UniqueID); err != nil {
			db = db.Begin()
			log.Printf("ERROR: Create Sync Transaction Failed")
			db, _ = r.TrxSvc.UpdateStatusSyncTrx(value.UniqueID, value.MerchantID, -1, db)
			continue
		}
	}
	return
}

// OnlineOrder godoc
// @Summary Create Online Order Trx (permission = consumer)
// @Id OnlineOrder
// @Tags Transaction
// @Param req body model.ReqTrxOrderOnline true " "
// @Success 200 {object} entity.Trx "data: entity.Trx"
// @Router /trx/order [post]
func (r *TrxController) OnlineOrder(c *gin.Context) {
	req := model.ReqTrxOrderOnline{}

	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.Error(err.Error())
		model.ResponseError(c,
			"Invalid Request", http.StatusUnprocessableEntity)
		return
	}

	data := new(entity.Trx)
	var err error
	if data, err = r.TrxSvc.CreateTrxOrderOnline(&req); err != nil {
		model.ResponseError(c, err.Error(), http.StatusUnprocessableEntity)
		r.TrxSvc.DeleteTrxByID(req.TrxID)
		return
	}
	model.ResponseJSON(c, data)
	return
}

// VoidOrder godoc
// @Summary Void Order Trx (permission = merchant)
// @Id VoidOrder
// @Tags Transaction
// @Param trxId path string true "Transaction ID"
// @Success 200 {object} entity.Trx "data: entity.Trx"
// @Router /trx/online-order/closed/{trxId} [put]
func (r *TrxController) VoidOrder(c *gin.Context) {
	trxID := c.Param("trxId")
	r.TrxSvc.VoidTrxByID(trxID)
	model.ResponseJSON(c, "Void Order Transaction Success!")
	return
}

// GetOnlineOrder godoc
// @Summary Get Online Order Trx for POS (permission = merchant)
// @Id GetOnlineOrder
// @Tags Transaction
// @Security Token
// @Success 200 {object} model.ResTrxOrderList "data: model.ResTrxOrderList"
// @Router /trx/online-order [get]
func (r *TrxController) GetOnlineOrder(c *gin.Context) {
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	merchantUsers := r.MerchantSvc.GetMerchantUsersByUsersID(jwtModel.UserID)
	if merchantUsers == nil {
		model.ResponseError(c, "Foodtruck Not Register Yet to Any Merchant", http.StatusUnprocessableEntity)
		return
	}
	data := r.TrxSvc.GetOnlineOrder(*merchantUsers)
	model.ResponseJSON(c, data)
	return
}

// ClosedOnlineOrder godoc
// @Summary Closed Online Order Trx for POS (permission = merchant)
// @Id ClosedOnlineOrder
// @Tags Transaction
// @Security Token
// @Param trxId path string true "Transaction ID"
// @Success 200 {object} model.ResTrxOnlineClosed "data: model.ResTrxOnlineClosed"
// @Router /trx/online-order/closed/{trxId} [put]
func (r *TrxController) ClosedOnlineOrder(c *gin.Context) {
	trxID := c.Param("trxId")
	data := r.TrxSvc.ClosedOnlineOrderByTrxID(trxID)
	model.ResponseJSON(c, data)
	return
}

// TrxVisitBookingList godoc
// @Summary List Home Visit Booking List (permission = merchant)
// @Id TrxVisitBookingList
// @Tags Transaction
// @Security Token
// @Param limit query string false " " default(10)
// @Param page query string false " " default(1)
// @Param sort query string false "e.g.: id,desc / id,asc"
// @Param filter query string false " "
// @Success 200 {object} model.ResponseSuccess  "data: model.Pagination"
// @Router /trx/visit/bookingall [get]
func (r *TrxController) TrxVisitBookingList(c *gin.Context) {
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	merchant := r.MerchantSvc.GetInfo(jwtModel.UserID)
	limit := util.ParamIDToInt(c.DefaultQuery("limit", "10"))
	page := util.ParamIDToInt(c.DefaultQuery("page", "1"))
	sorted := util.SortedBy(c.QueryArray("sort"))
	filter := c.DefaultQuery("filter", "")
	model.ResponsePagination(c, r.TrxSvc.ListBookingTrxVisitSale(merchant.ID, limit, page, sorted, filter))
	return
}

// TrxVisitBookingByID godoc
// @Summary List Home Visit Booking By Date (permission = merchant)
// @Id TrxVisitBookingByDate
// @Tags Transaction
// @Security Token
// @Param id path string true "0"
// @Success 200 {object} model.ResHomeVisitBookingDetailTimeNew "data: model.ResHomeVisitBookingDetailTimeNew"
// @Router /trx/visit/booking/{id} [get]
func (r *TrxController) TrxVisitBookingByID(c *gin.Context) {
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	id := util.ParamIDToInt64(c.Param("id"))
	merchant := r.MerchantSvc.GetInfo(jwtModel.UserID)
	model.ResponseJSON(c, r.TrxSvc.ListBookingTrxVisitSalesByID(id, merchant.ID))
	return
}

// TrxReportAll godoc
// @Summary TransactionReportAll (permission = merchant)
// @Id TrxReportAll
// @Tags Transaction
// @Security Token
// @Param month query string false "format 01 to 12"
// @Param year query string false "2015-now"
// @Param limit query string false "10"
// @Param page query string false "1"
// @Param sort query string false "id,desc"
// @Success 200 {object} model.Pagination "data: model.Pagination"
// @Router /trx/report-all [get]
func (r *TrxController) TrxReportAll(c *gin.Context) {
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	merchant := r.MerchantSvc.GetInfo(jwtModel.UserID)

	now := time.Now()
	monthNow := now.Format("01")
	yearNow := now.Format("2006")
	month := c.DefaultQuery("month", monthNow)
	if len(month) < 2 {
		monthConvert := util.ParamIDToInt(month)
		month = fmt.Sprintf("%02d", monthConvert)
	}
	year := c.DefaultQuery("year", yearNow)
	limit := util.ParamIDToInt(c.DefaultQuery("limit", "10"))
	page := util.ParamIDToInt(c.DefaultQuery("page", "1"))
	sorted := util.SortedBy(c.QueryArray("sort"))
	model.ResponsePagination(c, r.TrxSvc.TrxReportPagination(merchant.ID, month, year, limit, page, sorted))
	return
}

// TrxReportSingleAll godoc
// @Summary TransactionReportAll (permission = merchant)
// @Id TrxReportSingleAll
// @Tags Transaction
// @Security Token
// @Param month query string false "format 01 to 12"
// @Param year query string false "2015-now"
// @Param limit query string false "10"
// @Param page query string false "1"
// @Param sort query string false "id,desc"
// @Success 200 {object} model.Pagination "data: model.Pagination"
// @Router /trx/report-all [get]
func (r *TrxController) TrxReportSingleAll(c *gin.Context) {
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	merchant := r.MerchantSvc.GetInfo(jwtModel.UserID)

	now := time.Now()
	monthNow := now.Format("01")
	yearNow := now.Format("2006")
	month := c.DefaultQuery("month", monthNow)
	if len(month) < 2 {
		monthConvert := util.ParamIDToInt(month)
		month = fmt.Sprintf("%02d", monthConvert)
	}
	year := c.DefaultQuery("year", yearNow)
	limit := util.ParamIDToInt(c.DefaultQuery("limit", "10"))
	page := util.ParamIDToInt(c.DefaultQuery("page", "1"))
	sorted := util.SortedBy(c.QueryArray("sort"))
	model.ResponsePagination(c, r.TrxSvc.TrxReportSinglePagination(merchant.ID, month, year, limit, page, sorted))
	return
}

// TrxReport godoc
// @Summary TransactionReport (permission = merchant)
// @Id TrxReport
// @Tags Transaction
// @Security Token
// @Param month query string false "format 01 to 12"
// @Param year query string false "2015-now"
// @Success 200 {object} model.ResHomeVisitBookingList "data: model.ResHomeVisitBookingList"
// @Router /trx/report [get]
func (r *TrxController) TrxReport(c *gin.Context) {
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	merchant := r.MerchantSvc.GetInfo(jwtModel.UserID)

	now := time.Now()
	monthNow := now.Format("01")
	yearNow := now.Format("2006")
	month := c.DefaultQuery("month", monthNow)
	if len(month) < 2 {
		monthConvert := util.ParamIDToInt(month)
		month = fmt.Sprintf("%02d", monthConvert)
	}
	year := c.DefaultQuery("year", yearNow)
	f := excelize.NewFile()
	headerTable := map[string]string{
		"D3": "Order No",
		"E3": "Dates",
		"F3": "Times",
		"G3": "Transaction No",
		"H3": "Type",
		"I3": "Payment Method",
		"J3": "Product Name",
		"K3": "Employee",
		"L3": "Qty",
		"M3": "Total Invoice Tax",
		"N3": "Grand Total",
		"O3": "Status",
	}
	// path := cfg.Config.Path.Image
	// Insert a picture offset in the cell with external hyperlink, printing and positioning support.
	// if merchant.Logo != "" {
	// 	f.AddPicture("Sheet1", "D1", path+merchant.Logo, `{
	// 		"x_scale": 0.3,
	// 		"y_scale": 0.3,
	// 		"print_obj": true,
	// 		"lock_aspect_ratio": false,
	// 		"locked": false,
	// 		"positioning": "oneCell"
	// 	}`)
	// }

	// f.MergeCell("Sheet1", "D1", "N12")
	f.MergeCell("Sheet1", "D1", "O1")
	f.MergeCell("Sheet1", "D2", "O2")
	f.SetCellValue("Sheet1", "D1", "Merchant Name : "+merchant.Name)
	f.SetCellValue("Sheet1", "D2", "Merchant Address : "+merchant.Address)

	for key, value := range headerTable {
		f.SetCellValue("Sheet1", key, value)
	}

	valueTable := make(map[string]string)

	reports := r.TrxSvc.TrxReport(merchant.ID, month, year)
	if len(reports) > 0 {
		index := 4
		for _, value := range reports {
			indexConv := strconv.Itoa(index)
			convertTime, _ := time.Parse("2006-01-02 15:04:05", value.Dates+" 00:00:00")
			value.Dates = convertTime.Format("02/01/2006")
			valueTable["D"+indexConv] = value.OrderNo
			valueTable["E"+indexConv] = value.Dates
			valueTable["F"+indexConv] = value.Times
			valueTable["G"+indexConv] = value.TrxID
			valueTable["H"+indexConv] = value.TypeOrder
			valueTable["I"+indexConv] = value.ExternalID
			valueTable["J"+indexConv] = value.ProductName
			valueTable["K"+indexConv] = value.UserName
			valueTable["L"+indexConv] = strconv.Itoa(value.Qty)
			valueTable["M"+indexConv] = strconv.FormatFloat(value.TotalTax, 'f', 0, 64)
			valueTable["N"+indexConv] = strconv.FormatFloat(value.GrandTotal, 'f', 0, 64)
			valueTable["O"+indexConv] = value.Status

			index++
		}
	}

	for key, value := range valueTable {
		f.SetCellValue("Sheet1", key, value)
	}

	today := time.Now().Format("Mon 02 Jan 2006")
	fileName := "Transaction Report by Menu " + today + ".xlsx"

	// Set the headers necessary to get browsers to interpret the downloadable file
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment;filename=\""+fileName+"\"")
	c.Header("Content-Transfer-Encoding", "binary")
	_ = f.Write(c.Writer)
}

// TrxReportSingle godoc
// @Summary TransactionReport (permission = merchant)
// @Id TrxReportSingle
// @Tags Transaction
// @Security Token
// @Param month query string false "format 01 to 12"
// @Param year query string false "2015-now"
// @Success 200 {object} model.ResHomeVisitBookingList "data: model.ResHomeVisitBookingList"
// @Router /trx/report [get]
func (r *TrxController) TrxReportSingle(c *gin.Context) {
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	merchant := r.MerchantSvc.GetInfo(jwtModel.UserID)

	now := time.Now()
	monthNow := now.Format("01")
	yearNow := now.Format("2006")
	month := c.DefaultQuery("month", monthNow)
	if len(month) < 2 {
		monthConvert := util.ParamIDToInt(month)
		month = fmt.Sprintf("%02d", monthConvert)
	}
	year := c.DefaultQuery("year", yearNow)
	f := excelize.NewFile()
	headerTable := map[string]string{
		"D3": "Order No",
		"E3": "Dates",
		"F3": "Times",
		"G3": "Transaction No",
		"H3": "Type",
		"I3": "Payment Method",
		"J3": "Employee",
		"K3": "Total Tax",
		"L3": "Grand Total",
		"M3": "Status",
	}
	// path := cfg.Config.Path.Image
	// Insert a picture offset in the cell with external hyperlink, printing and positioning support.
	// if merchant.Logo != "" {
	// 	f.AddPicture("Sheet1", "D1", path+merchant.Logo, `{
	// 		"x_scale": 0.3,
	// 		"y_scale": 0.3,
	// 		"print_obj": true,
	// 		"lock_aspect_ratio": false,
	// 		"locked": false,
	// 		"positioning": "oneCell"
	// 	}`)
	// }

	// f.MergeCell("Sheet1", "D1", "N12")
	f.MergeCell("Sheet1", "D1", "M1")
	f.MergeCell("Sheet1", "D2", "M2")
	f.SetCellValue("Sheet1", "D1", "Merchant Name : "+merchant.Name)
	f.SetCellValue("Sheet1", "D2", "Merchant Address : "+merchant.Address)

	for key, value := range headerTable {
		f.SetCellValue("Sheet1", key, value)
	}

	valueTable := make(map[string]string)

	reports := r.TrxSvc.TrxReportSingle(merchant.ID, month, year)
	if len(reports) > 0 {
		index := 4
		for _, value := range reports {
			indexConv := strconv.Itoa(index)
			convertTime, _ := time.Parse("2006-01-02 15:04:05", value.Dates+" 00:00:00")
			value.Dates = convertTime.Format("02/01/2006")
			valueTable["D"+indexConv] = value.OrderNo
			valueTable["E"+indexConv] = value.Dates
			valueTable["F"+indexConv] = value.Times
			valueTable["G"+indexConv] = value.TrxID
			valueTable["H"+indexConv] = value.TypeOrder
			valueTable["I"+indexConv] = value.ExternalID
			valueTable["J"+indexConv] = value.UserName
			valueTable["K"+indexConv] = strconv.FormatFloat(value.TotalTax, 'f', 0, 64)
			valueTable["L"+indexConv] = strconv.FormatFloat(value.GrandTotal, 'f', 0, 64)
			valueTable["M"+indexConv] = value.Status

			index++
		}
	}

	for key, value := range valueTable {
		f.SetCellValue("Sheet1", key, value)
	}

	today := time.Now().Format("Mon 02 Jan 2006")
	fileName := "Transaction Report by Transaction " + today + ".xlsx"

	// Set the headers necessary to get browsers to interpret the downloadable file
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment;filename=\""+fileName+"\"")
	c.Header("Content-Transfer-Encoding", "binary")
	_ = f.Write(c.Writer)
}
