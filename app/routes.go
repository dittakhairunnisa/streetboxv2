package app

import (
	"streetbox.id/app/canvassing"
	"streetbox.id/app/trxrefund"
	"streetbox.id/cfg"

	"github.com/gin-gonic/gin"
	"streetbox.id/middleware"

	"streetbox.id/app/appsetting"
	AppSetting "streetbox.id/app/appsetting/controller"
	Canvas "streetbox.id/app/canvassing/controller"
	"streetbox.id/app/enduser"
	EndUser "streetbox.id/app/enduser/controller"
	"streetbox.id/app/fcm"
	"streetbox.id/app/homevisitsales"
	HomeSales "streetbox.id/app/homevisitsales/controller"
	"streetbox.id/app/logactivity"
	LogActivity "streetbox.id/app/logactivity/controller"
	"streetbox.id/app/logactivitymerchant"
	LogMerch "streetbox.id/app/logactivitymerchant/controller"
	"streetbox.id/app/merchant"
	Merchant "streetbox.id/app/merchant/controller"
	"streetbox.id/app/parkingspace"
	PSpace "streetbox.id/app/parkingspace/controller"
	"streetbox.id/app/payment"
	Payment "streetbox.id/app/payment/controller"
	"streetbox.id/app/role"
	Role "streetbox.id/app/role/controller"
	"streetbox.id/app/sales"
	"streetbox.id/app/tasks"
	Tasks "streetbox.id/app/tasks/controller"
	"streetbox.id/app/trx"
	Trx "streetbox.id/app/trx/controller"
	"streetbox.id/app/trxspaces"
	TrxSales "streetbox.id/app/trxspaces/controller"
	"streetbox.id/app/user"
	User "streetbox.id/app/user/controller"
)

// AppSettingController ..
func AppSettingController(r *gin.Engine,
	svc appsetting.ServiceInterface) {
	appSetting := &AppSetting.AppSettingController{
		AppSettingSvc: svc,
	}
	routes := r.Group("appsetting")
	routes.GET("get-by-key/:key", middleware.AuthMiddleware("all"), appSetting.GetByKey)
	routes.POST("update-by-key/:key", middleware.AuthMiddleware("superadmin"), appSetting.UpdateByKey)
}

// UserController ...
func UserController(r *gin.Engine, service user.ServiceInterface) {
	user := &User.UserController{UserService: service}
	r.POST("login", user.Login)
	r.POST("forgotpassword", user.ForgotPassword)
	r.PUT("resetpassword", user.ResetPassword)
	r.GET("check", user.CheckToken) // validation email token
	r.POST("login/google", user.LoginGoogle)
	routes := r.Group("user")
	routes.POST("", middleware.AuthMiddleware("superadmin"), user.CreateUser)
	routes.GET("merchant", middleware.AuthMiddleware("superadmin"), user.GetUserMerchant)
	routes.GET("", middleware.AuthMiddleware("superadmin"), user.GetAll)
	routes.DELETE(":id/delete", middleware.AuthMiddleware("superadmin"), user.DeleteByID)
	routes.GET("info", middleware.AuthMiddleware("all"), user.GetUser)
	routes.PUT("update", middleware.AuthMiddleware("all"), user.UpdateUser)
	routes.PUT("role/:id/update", middleware.AuthMiddleware("superadmin"), user.UpdateUserRole)
	routes.PUT("changepassword", middleware.AuthMiddleware("all"), user.ChangePassword)
	routes.GET("info/:id", middleware.AuthMiddleware("superadmin"), user.GetUserByID)
	routes.POST("address", middleware.AuthMiddleware("consumer"), user.CreateAddress)
	routes.GET("address", middleware.AuthMiddleware("consumer"), user.GetAddressByUserID)
	routes.GET("address/primary", middleware.AuthMiddleware("consumer"), user.GetPrimaryAddressByUserID)
	routes.PUT("address", middleware.AuthMiddleware("consumer"), user.UpdateAddress)
	routes.DELETE(":id/address", middleware.AuthMiddleware("consumer"), user.DeleteAddress)
	routes.PUT("address/:id/switch", middleware.AuthMiddleware("consumer"), user.SwitchAddress)
	routes.PUT("config/:rad", middleware.AuthMiddleware("superadmin"), user.UpdateRadius)
	routes.GET("config", middleware.AuthMiddleware("all"), user.GetConfig)
}

// RoleController ...
func RoleController(r *gin.Engine, service role.ServiceInterface) {
	role := &Role.RoleController{Service: service}
	routes := r.Group("role")
	routes.POST("", middleware.AuthMiddleware("superadmin"), role.Create)
	routes.GET("", middleware.AuthMiddleware("superadmin"), role.GetAll)
	routes.GET("exclude", middleware.AuthMiddleware("superadmin"), role.GetAllExclude)
}

// PSpaceController ...
func PSpaceController(r *gin.Engine, pspaceSrv parkingspace.ServiceInterface,
	salesSrv sales.ServiceInterface) {
	imagePath := cfg.Config.Path.Image
	docPath := cfg.Config.Path.Doc
	pspace := &PSpace.ParkingSpaceController{PSpaceSvc: pspaceSrv, SalesSvc: salesSrv}
	routes := r.Group("parkingspace")

	routes.GET("", middleware.AuthMiddleware("admin"), pspace.GetAll)
	routes.GET("list", middleware.AuthMiddleware("superadmin"), pspace.GetAllList)
	routes.GET("spacesalesdate/:id", middleware.AuthMiddleware("superadmin"), pspace.GetSpaceBySalesDate)
	routes.GET("show/:id", middleware.AuthMiddleware("superadmin"), pspace.GetByID)
	routes.DELETE(":id/delete", middleware.AuthMiddleware("superadmin"), pspace.DeleteByID)
	routes.GET("sales", middleware.AuthMiddleware("admin"), pspace.GetSales)
	routes.GET("list/sales", middleware.AuthMiddleware("superadmin"), pspace.GetSalesList)
	routes.GET("sales/:id", middleware.AuthMiddleware("admin"), pspace.GetSalesByPSpaceID)
	routes.GET("search/sales/:key", middleware.AuthMiddleware("admin"), pspace.SearchSalesByName)
	routes.GET("sales/:id/info", middleware.AuthMiddleware("superadmin"), pspace.GetSalesByID)
	routes.DELETE(":id/sales/delete", middleware.AuthMiddleware("superadmin"), pspace.DeleteSalesByID)
	routes.POST("", middleware.AuthMiddleware("superadmin"), pspace.CreateParkingSpace)
	routes.PUT(":id/sales/update", middleware.AuthMiddleware("superadmin"), pspace.UpdateSales)
	routes.POST("sales/create", middleware.AuthMiddleware("superadmin"), pspace.CreateSales)
	routes.PUT(":id/update", middleware.AuthMiddleware("superadmin"), pspace.Update)
	routes.PUT(":id/image/upload", middleware.AuthMiddleware("superadmin"), pspace.UploadImage)
	routes.PUT(":id/image/delete", middleware.AuthMiddleware("superadmin"), pspace.DeleteImage)
	routes.PUT(":id/doc/upload", middleware.AuthMiddleware("superadmin"), pspace.UploadDoc)
	routes.PUT(":id/doc/delete", middleware.AuthMiddleware("superadmin"), pspace.DeleteDoc)
	routesStatic := r.Group("static")
	routesStatic.Static("image", imagePath)
	routesStatic.Static("doc", docPath)
}

// TrxSalesController ...
func TrxSalesController(
	r *gin.Engine,
	srv trxspaces.ServiceInterface,
	pspaceSrv sales.ServiceInterface,
	merchantSvc merchant.ServiceInterface) {
	trxSales := &TrxSales.TrxSalesController{
		Service:      srv,
		ParkSalesSvc: pspaceSrv,
		MerchantSvc:  merchantSvc}
	routes := r.Group("trxsales")
	routes.POST("", middleware.AuthMiddleware("superadmin"), trxSales.Create)
	routes.GET("all", middleware.AuthMiddleware("superadmin"), trxSales.GetAll)
	routes.GET("list", middleware.AuthMiddleware("superadmin"), trxSales.GetList)
	routes.GET("myparking", middleware.AuthMiddleware("merchant"), trxSales.GetMyParking)
	routes.GET("myparking/slot/:id", middleware.AuthMiddleware("merchant"), trxSales.GetSlotMyParking)
	routes.GET("info/:id", middleware.AuthMiddleware("superadmin"), trxSales.GetByID)
}

// MerchantController ...
func MerchantController(r *gin.Engine,
	srv merchant.ServiceInterface, userSvc user.ServiceInterface, trxSvc trx.ServiceInterface) {
	merchant := &Merchant.MerchantController{MerchantSvc: srv, UserSvc: userSvc, TrxSvc: trxSvc}
	routes := r.Group("merchant")
	routes.POST("category", middleware.AuthMiddleware("superadmin"), merchant.CreateCategory)
	routes.GET("category", middleware.AuthMiddleware("all"), merchant.GetAllCategory)
	routes.PUT("category", middleware.AuthMiddleware("superadmin"), merchant.UpdateCategory)
	routes.DELETE(":id/category", middleware.AuthMiddleware("superadmin"), merchant.DeleteCategory)
	routes.GET("info", middleware.AuthMiddleware("merchant"), merchant.GetInfo)
	routes.GET("all", middleware.AuthMiddleware("superadmin"), merchant.GetAll)
	routes.GET("info/:id", middleware.AuthMiddleware("superadmin"), merchant.GetByID)
	routes.POST("", middleware.AuthMiddleware("admin"), merchant.CreateMerchant)
	routes.POST("xendit-generate-subaccount", middleware.AuthMiddleware("superadmin"), merchant.XenditGenerateSubAccount)
	routes.DELETE(":id", middleware.AuthMiddleware("superadmin"), merchant.DeleteByID)
	routes.DELETE(":id/foodtruck", middleware.AuthMiddleware("admin"), merchant.DeleteFoodtruckByID)
	routes.DELETE(":id/menu/delete", middleware.AuthMiddleware("admin"), merchant.DeleteMenu)
	routes.PUT("", middleware.AuthMiddleware("admin"), merchant.UpdateMerchant)
	routes.GET("foodtruck/all", middleware.AuthMiddleware("admin"), merchant.GetAllFoodtruck)
	routes.POST("foodtruck", middleware.AuthMiddleware("admin"), merchant.CreateFoodtruck)
	routes.PUT("foodtruck/:id/update", middleware.AuthMiddleware("admin"), merchant.UpdateFoodtruck)
	routes.GET("info/:id/foodtruck", middleware.AuthMiddleware("admin"), merchant.GetFoodtruckByID)
	routes.PUT("upload-logo", middleware.AuthMiddleware("admin"), merchant.UploadLogo)
	routes.PUT("upload-banner", middleware.AuthMiddleware("admin"), merchant.UploadBanner)
	routes.PUT("foodtruck/:id/resetpassword", middleware.AuthMiddleware("admin"), merchant.ResetPasswordFoodTruck)
	routes.POST("menu", middleware.AuthMiddleware("admin"), merchant.CreateMenu)
	routes.GET("menu/all", middleware.AuthMiddleware("admin"), merchant.ListPaginateMenu)
	routes.PUT("menu/:id", middleware.AuthMiddleware("admin"), merchant.UpdateMenu)
	routes.PUT("upload-menu/:id", middleware.AuthMiddleware("admin"), merchant.UploadMenu)
	routes.GET("list/menu", middleware.AuthMiddleware("merchant"), merchant.ListMenu)
	routes.GET("count/foodtruck", middleware.AuthMiddleware("admin"), merchant.CountFoodtruck)
	routes.GET("menu/info/:id", middleware.AuthMiddleware("admin"), merchant.GetMenuByID)
	routes.POST("csv/uploadmenu", middleware.AuthMiddleware("admin"), merchant.ImportCSV)
	routes.GET("csv/templatemenu", merchant.DownloadMenuTemplateCSV)
	routes.POST("tax/menu", middleware.AuthMiddleware("admin"), merchant.SetMerchantTax)
	routes.GET("taxsetting/menu", middleware.AuthMiddleware("merchant"), merchant.GetTaxSetting)
	routes.GET("pos/gettransaction", middleware.AuthMiddleware("merchant"), merchant.GetPosTransaction)
	routes.POST("registration-token/:token", middleware.AuthMiddleware("merchant"), merchant.RegistrationToken)
	routes.PUT("remove-logo/:filename", middleware.AuthMiddleware("admin"), merchant.RemoveLogo)
	routes.PUT("remove-banner/:filename", middleware.AuthMiddleware("admin"), merchant.RemoveBanner)
	routes.PUT("delete-image-menu/:id", middleware.AuthMiddleware("admin"), merchant.DeleteImageMenu)
}

// TasksController new Module Tasks
func TasksController(
	r *gin.Engine,
	tasksSvc tasks.ServiceInterface,
	merchantSvc merchant.ServiceInterface,
	logSvc logactivitymerchant.ServiceInterface,
	userSvc user.ServiceInterface) {
	tasks := &Tasks.TasksController{
		TasksSvc:    tasksSvc,
		MerchantSvc: merchantSvc,
		LogSvc:      logSvc,
		UserSvc:     userSvc}
	routes := r.Group("tasks")
	routes.POST("regular", middleware.AuthMiddleware("admin"), tasks.CreateTasksRegular)
	routes.POST("homevisit", middleware.AuthMiddleware("admin"), tasks.CreateTasksHomevisit)
	routes.POST("nonregular", middleware.AuthMiddleware("merchant"), tasks.CreateTasksNonRegular)
	routes.POST("regular/check-in", middleware.AuthMiddleware("merchant"), tasks.CreateRegularCheckIn)
	routes.POST("regular/check-out", middleware.AuthMiddleware("merchant"), tasks.CreateRegularCheckOut)
	routes.POST("homevisit/check-in", middleware.AuthMiddleware("merchant"), tasks.CreateHomeVisitCheckIn)
	routes.POST("homevisit/check-out", middleware.AuthMiddleware("merchant"), tasks.CreateHomeVisitCheckOut)
	routes.POST("nonregular/check-in", middleware.AuthMiddleware("merchant"), tasks.CreateNonRegCheckIn)
	routes.POST("nonregular/check-out", middleware.AuthMiddleware("merchant"), tasks.CreateNonRegCheckOut)
	routes.GET("regular/list", middleware.AuthMiddleware("merchant"), tasks.MyTaskRegular)
	routes.GET("nonregular/list", middleware.AuthMiddleware("merchant"), tasks.MyTaskNonRegular)
	routes.GET("nonregular/status", middleware.AuthMiddleware("merchant"), tasks.NonRegStatus)
	routes.PUT("undo/:id", middleware.AuthMiddleware("merchant"), tasks.ChangeToOpen)
	routes.POST("shift-in", middleware.AuthMiddleware("merchant"), tasks.ShiftIn)
	routes.GET("shift-in/status", middleware.AuthMiddleware("merchant"), tasks.ShiftInStatus)
	routes.POST("shift-out", middleware.AuthMiddleware("merchant"), tasks.ShiftOut)
	routes.POST("tracking", middleware.AuthMiddleware("merchant"), tasks.CreateTracking)
	routes.GET("tracking/:id/foodtruck", middleware.AuthMiddleware("merchant"), tasks.GetTracking)
	routes.PUT("ongoing/:id", middleware.AuthMiddleware("merchant"), tasks.ChangeToOngoing)
	routes.GET("all", middleware.AuthMiddleware("admin"), tasks.GetAll)
}

// LogMerchantController ..
func LogMerchantController(r *gin.Engine,
	svc logactivitymerchant.ServiceInterface,
	merchantSvc merchant.ServiceInterface) {
	log := &LogMerch.LogMerchantController{
		Svc:         svc,
		MerchantSvc: merchantSvc,
	}
	routes := r.Group("log-merchant")
	routes.GET("", middleware.AuthMiddleware("admin"), log.GetAll)
}

// LogController ..
func LogController(r *gin.Engine, svc logactivity.ServiceInterface) {
	log := &LogActivity.LogController{
		Svc: svc,
	}
	routes := r.Group("log")
	routes.GET("", middleware.AuthMiddleware("superadmin"), log.GetAll)
	routes.GET("generatecsv", log.GenerateCSV)
}

// HomevisitSalesController ..
func HomevisitSalesController(
	r *gin.Engine,
	homeSvc homevisitsales.ServiceInterface,
	merchantSvc merchant.ServiceInterface) {
	homevisit := &HomeSales.HomevisitSalesController{
		HomeSalesSvc: homeSvc,
		MerchantSvc:  merchantSvc,
	}
	routes := r.Group("homevisit")
	routes.GET("", middleware.AuthMiddleware("admin"), homevisit.GetAll)
	routes.POST("", middleware.AuthMiddleware("admin"), homevisit.Create)
	routes.POST("/batch", middleware.AuthMiddleware("admin"), homevisit.BatchCreate)
	routes.GET("info/:date", middleware.AuthMiddleware("admin"), homevisit.GetInfo)
	routes.DELETE("deletebyid/:id", middleware.AuthMiddleware("admin"), homevisit.DeleteByID)
	routes.PUT("", middleware.AuthMiddleware("admin"), homevisit.Update)
	routes.DELETE("deletebydate/:date", middleware.AuthMiddleware("admin"), homevisit.Delete)
}

// EndUserController ..
func EndUserController(r *gin.Engine,
	svc enduser.ServiceInterface,
	visitSales homevisitsales.ServiceInterface,
	merchant merchant.ServiceInterface,
	trx trx.ServiceInterface,
	fcmSvc fcm.ServiceInterface) {
	endUser := &EndUser.EndUserController{
		Svc:           svc,
		VisitSalesSvc: visitSales,
		MerchantSvc:   merchant,
		TrxSvc:        trx,
		FcmSvc:        fcmSvc,
	}
	routes := r.Group("consumer")
	routes.GET("home/nearby/:lat/:lon", middleware.AuthMiddleware("consumer"), endUser.Nearby)
	routes.GET("home/map/livetracking/:lat/:lon", middleware.AuthMiddleware("consumer"), endUser.LiveTracking)
	routes.GET("home/map/parking-space/:lat/:lon", middleware.AuthMiddleware("consumer"), endUser.MapParkingSpace)
	routes.GET("home/map/schedules/:id/parking-space", middleware.AuthMiddleware("consumer"), endUser.MapParkingSpaceDetail)
	routes.GET("home/schedules-regular/:typesId", middleware.AuthMiddleware("consumer"), endUser.MerchantSchedule)
	routes.GET("home/visit-sales", middleware.AuthMiddleware("consumer"), endUser.VisitSales)
	routes.GET("home/visit-sales/detail/:merchantId", middleware.AuthMiddleware("consumer"), endUser.VisitSalesDetail)
	routes.GET("merchant/tax/:merchantId", middleware.AuthMiddleware("consumer"), endUser.GetMerchantTaxByMerchantID)
	routes.GET("merchant/menu/:merchantId", middleware.AuthMiddleware("consumer"), endUser.GetMerchantMenuByMerchantID)
	routes.GET("order/history", middleware.AuthMiddleware("consumer"), endUser.OrderHistory)
	routes.GET("payment-method", middleware.AuthMiddleware("consumer"), endUser.GetPaymentMethod)
	routes.POST("registration-token/:token", middleware.AuthMiddleware("consumer"), endUser.RegistrationToken)
	routes.PUT("update/userprofile", middleware.AuthMiddleware("consumer"), endUser.UpdateEndUser)
}

// PaymentController ...
func PaymentController(
	r *gin.Engine,
	paymentSvc payment.ServiceInterface,
	fcmSvc fcm.ServiceInterface,
	endSvc enduser.ServiceInterface,
	trxSvc trx.ServiceInterface,
	merchantSvc merchant.ServiceInterface,
	homeVisitSvc homevisitsales.ServiceInterface) {
	payment := &Payment.PaymentController{
		PaymentSvc:    paymentSvc,
		FcmSvc:        fcmSvc,
		EndUserSvc:    endSvc,
		TrxSvc:        trxSvc,
		MerchantSvc:   merchantSvc,
		VisitSalesSvc: homeVisitSvc,
	}
	routes := r.Group("payment")
	routes.POST("create-qrcode", middleware.AuthMiddleware("all"), payment.CreateQRIS)
	routes.POST("xendit/qris/callback", payment.XenditCallback)
	routes.GET("qrcode/:trxId", middleware.AuthMiddleware("consumer"), payment.GetQRByTrxID)
	routes.POST("simulate-qrcode/:trxId", middleware.AuthMiddleware("consumer"), payment.SimulateQR)
}

// TrxController ...
func TrxController(r *gin.Engine, trxSvc trx.ServiceInterface, merchantSvc merchant.ServiceInterface,
	trxrefund trxrefund.ServiceInterface) {
	trx := &Trx.TrxController{
		TrxSvc:       trxSvc,
		MerchantSvc:  merchantSvc,
		TrxRefundSvc: trxrefund,
	}
	routes := r.Group("trx")
	routes.GET("info/:trxId", middleware.AuthMiddleware("consumer"), trx.GetInfo)
	routes.POST("createsync", middleware.AuthMiddleware("merchant"), trx.CreateSyncTrx)
	routes.POST("homevisit", middleware.AuthMiddleware("consumer"), trx.CreateTrxVisit)
	routes.POST("refund/space", middleware.AuthMiddleware("superadmin"), trx.CreateTrxRefundSpace)
	routes.POST("refund/visit", middleware.AuthMiddleware("merchant"), trx.CreateTrxRefundVisit)
	routes.POST("order", middleware.AuthMiddleware("all"), trx.OnlineOrder)
	routes.PUT("order/void/:trxId", middleware.AuthMiddleware("merchant"), trx.VoidOrder)
	routes.GET("online-order", middleware.AuthMiddleware("merchant"), trx.GetOnlineOrder)
	routes.PUT("online-order/closed/:trxId", middleware.AuthMiddleware("merchant"), trx.ClosedOnlineOrder)
	routes.GET("visit/bookingall", middleware.AuthMiddleware("merchant"), trx.TrxVisitBookingList)
	routes.GET("visit/booking/:id", middleware.AuthMiddleware("merchant"), trx.TrxVisitBookingByID)
	routes.GET("report-all", middleware.AuthMiddleware("gmerchant"), trx.TrxReportAll)
	routes.GET("report-single-all", middleware.AuthMiddleware("merchant"), trx.TrxReportSingleAll)
	routes.GET("report", middleware.AuthMiddleware("merchant"), trx.TrxReport)
	routes.GET("report-single", middleware.AuthMiddleware("merchant"), trx.TrxReportSingle)
}

func CanvassingController(r *gin.Engine, can canvassing.Service, fcm fcm.ServiceInterface, mcv merchant.ServiceInterface, usr user.ServiceInterface) {
	canv := &Canvas.CanvassingController{
		Svc:      can,
		Fcm:      fcm,
		Merchant: mcv,
		User:     usr,
	}
	routes := r.Group("canvassing")
	routes.POST("", middleware.AuthMiddleware("admin"), canv.CreateCanvas)
	routes.GET("", middleware.AuthMiddleware("admin"), canv.GetCanvas)
	routes.GET("/foodtruck", middleware.AuthMiddleware("merchant"), canv.GetFoodtruckCanvas)
	routes.PUT("", middleware.AuthMiddleware("admin"), canv.UpdateCanvas)
	routes.GET("/foodtruck/location/:id", middleware.AuthMiddleware("all"), canv.GetFoodtruckLocation)
	routes.PUT("/foodtruck/location", middleware.AuthMiddleware("merchant"), canv.UpdateFoodtruckLocation)
	routes.GET("/users/location", middleware.AuthMiddleware("consumer"), canv.GetUserLocation)
	routes.PUT("/users/location", middleware.AuthMiddleware("consumer"), canv.UpdateUsersLocation)
	routes.POST("/blast/", middleware.AuthMiddleware("merchant"), canv.Blast)
	routes.POST("/call/:notif-id", middleware.AuthMiddleware("consumer"), canv.CallFoodtruck)
	routes.PUT("/call/:call-id/:status", middleware.AuthMiddleware("merchant"), canv.AnswerCall)
	routes.PUT("/call-status/:call-id/:status", middleware.AuthMiddleware("merchant"), canv.UpdateStatusCall)
	routes.PUT("/finish/:call-id", middleware.AuthMiddleware("merchant"), canv.FinishCall)
	routes.GET("/notifications", middleware.AuthMiddleware("consumer"), canv.GetNotifications)
	routes.GET("/foodtruck/calls", middleware.AuthMiddleware("merchant"), canv.GetCallsByFoodtruckID)
	routes.GET("/calls", middleware.AuthMiddleware("consumer"), canv.GetCallsByCustomerID)
	routes.PUT("/toggle/auto", middleware.AuthMiddleware("merchant"), canv.ToggleAutoBlast)
}
