package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"google.golang.org/api/option"
	"gopkg.in/natefinch/lumberjack.v2"
	"streetbox.id/cfg"
	"streetbox.id/jobs"

	firebase "firebase.google.com/go/v4"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
	"github.com/jinzhu/configor"
	"github.com/jinzhu/gorm"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "streetbox.id/docs"

	AppSettingRepo "streetbox.id/app/appsetting/repository"
	CanvasRepo "streetbox.id/app/canvassing/repository"
	HomevisitRepo "streetbox.id/app/homevisitsales/repository"
	LogRepo "streetbox.id/app/logactivity/repository"
	LogMerchRepo "streetbox.id/app/logactivitymerchant/repository"
	LogXenditReqRepo "streetbox.id/app/logxenditreq/repository"
	MerchantRepo "streetbox.id/app/merchant/repository"
	MerchantCategoryRepo "streetbox.id/app/merchantcategory/repository"
	MerchantMenusRepo "streetbox.id/app/merchantmenu/repository"
	MerchantTaxsRepo "streetbox.id/app/merchanttax/repository"
	MerchantUsersRepo "streetbox.id/app/merchantusers/repository"
	ShiftRepo "streetbox.id/app/merchantusersshift/repository"
	PSpaceRepo "streetbox.id/app/parkingspace/repository"
	PayMethodRepo "streetbox.id/app/paymentmethod/repository"
	RoleRepo "streetbox.id/app/role/repository"
	SalesRepo "streetbox.id/app/sales/repository"
	TasksRepo "streetbox.id/app/tasks/repository"
	TasksHomeRepo "streetbox.id/app/taskshomevisit/repository"
	TasksNonRegLogRepo "streetbox.id/app/tasksnonreglog/repository"
	TasksNonRegRepo "streetbox.id/app/tasksnonregular/repository"
	TasksRegLogRepo "streetbox.id/app/tasksreglog/repository"
	TasksRegRepo "streetbox.id/app/tasksregular/repository"
	TasksTrackingRepo "streetbox.id/app/taskstracking/repository"
	TasksHomeLogRepo "streetbox.id/app/tasksvisitlog/repository"
	TrxRepo "streetbox.id/app/trx/repository"
	TrxOrderRepo "streetbox.id/app/trxorder/repository"
	TrxOrderBillRepo "streetbox.id/app/trxorderbill/repository"
	TrxOrderPaymentSalesRepo "streetbox.id/app/trxorderpaymentsales/repository"
	TrxOrderProductSalesRepo "streetbox.id/app/trxorderproductsales/repository"
	TrxOrderTaxSalesRepo "streetbox.id/app/trxordertaxsales/repository"
	TrxRefundRepo "streetbox.id/app/trxrefund/repository"
	TrxSalesRepo "streetbox.id/app/trxspaces/repository"
	TrxVisitRepo "streetbox.id/app/trxvisit/repository"
	TrxVisitMenuSalesRepo "streetbox.id/app/trxvisitmenusales/repository"
	TrxVisitSalesRepo "streetbox.id/app/trxvisitsales/repository"
	UserRepo "streetbox.id/app/user/repository"
	UserAddressRepo "streetbox.id/app/useraddress/repository"
	UserAuthRepo "streetbox.id/app/userauth/repository"
	UserConfigRepo "streetbox.id/app/userconfig/repository"
	UserRoleRepo "streetbox.id/app/userrole/repository"

	routes "streetbox.id/app"
	AppSettingService "streetbox.id/app/appsetting/service"
	CanvasService "streetbox.id/app/canvassing/service"
	EndUserSvc "streetbox.id/app/enduser/service"
	FcmSvc "streetbox.id/app/fcm/service"
	HomevisitSvc "streetbox.id/app/homevisitsales/service"
	LogService "streetbox.id/app/logactivity/service"
	LogMerchSvc "streetbox.id/app/logactivitymerchant/service"
	MerchantService "streetbox.id/app/merchant/service"
	PSpaceService "streetbox.id/app/parkingspace/service"
	PaymentSvc "streetbox.id/app/payment/service"
	RoleService "streetbox.id/app/role/service"
	SalesService "streetbox.id/app/sales/service"
	TasksService "streetbox.id/app/tasks/service"
	TrxService "streetbox.id/app/trx/service"
	TrxRefundService "streetbox.id/app/trxrefund/service"
	TrxSalesService "streetbox.id/app/trxspaces/service"
	UserService "streetbox.id/app/user/service"
)

func main() {
	// init config
	if err := configor.New(&configor.Config{ErrorOnUnmatchedKeys: true}).Load(&cfg.Config, "config.yml"); err != nil {
		log.Fatalf("ERROR: %s", err.Error())
	}
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetOutput(&lumberjack.Logger{
		Filename: cfg.Config.Log.Filename,
		MaxSize:  cfg.Config.Log.MaxSize,
		MaxAge:   cfg.Config.Log.MaxAge,
	})

	// init db connection
	conString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable TimeZone=Asia/Jakarta",
		cfg.Config.Postgres.Host,
		cfg.Config.Postgres.Port,
		cfg.Config.Postgres.User,
		cfg.Config.Postgres.Name,
		cfg.Config.Postgres.Password)
	dbConn, err := gorm.Open("postgres", conString)
	defer dbConn.Close()
	if err != nil {
		panic(err.Error())
	}
	err = dbConn.DB().Ping()
	if err != nil {
		panic(err.Error)
	}
	if cfg.Config.Env != "development" {
		dbConn.LogMode(false)
	} else {
		dbConn.LogMode(true)
	}
	dbConn.SingularTable(true)
	log.Printf("Connected to db with host %s", cfg.Config.Postgres.Host)

	// init server with default config
	app := gin.Default()

	// init middleware
	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "CLIENT_ID"},
	}))

	// init swagger
	var url func(c *ginSwagger.Config)
	url = ginSwagger.URL(fmt.Sprintf("http://%s:%s/swagger/doc.json", cfg.Config.Api.Host, cfg.Config.Api.Port))
	if cfg.Config.Env != "development" {
		url = ginSwagger.URL(fmt.Sprintf("https://%s/swagger/doc.json", cfg.Config.Api.Host))
	}
	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	// init repositories
	appSettingRepo := AppSettingRepo.New(dbConn)
	userRepo := UserRepo.New(dbConn)
	roleRepo := RoleRepo.New(dbConn)
	userAuthRepo := UserAuthRepo.New(dbConn)
	userRoleRepo := UserRoleRepo.New(dbConn)
	userAddressRepo := UserAddressRepo.New(dbConn)
	userConfigRepo := UserConfigRepo.New(dbConn)
	pSpaceRepo := PSpaceRepo.New(dbConn)
	salesRepo := SalesRepo.New(dbConn)
	trxSalesRepo := TrxSalesRepo.New(dbConn)
	merchantRepo := MerchantRepo.New(dbConn)
	merchantUsersRepo := MerchantUsersRepo.New(dbConn)
	shiftRepo := ShiftRepo.New(dbConn)
	logRepo := LogRepo.New(dbConn)
	tasksRepo := TasksRepo.New(dbConn)
	tasksRegRepo := TasksRegRepo.New(dbConn)
	tasksNonRegRepo := TasksNonRegRepo.New(dbConn)
	tasksHomeRepo := TasksHomeRepo.New(dbConn)
	logMerchRepo := LogMerchRepo.New(dbConn)
	homevisitRepo := HomevisitRepo.New(dbConn)
	tasksRegLogRepo := TasksRegLogRepo.New(dbConn)
	tasksNonRegLogRepo := TasksNonRegLogRepo.New(dbConn)
	tasksHomeLogRepo := TasksHomeLogRepo.New(dbConn)
	tasksTrackingRepo := TasksTrackingRepo.New(dbConn)
	merchantMenuRepo := MerchantMenusRepo.New(dbConn)
	merchantTaxRepo := MerchantTaxsRepo.New(dbConn)
	merchantCategoryRepo := MerchantCategoryRepo.New(dbConn)
	trxRepo := TrxRepo.New(dbConn)
	trxOrderRepo := TrxOrderRepo.New(dbConn)
	trxOrderBillRepo := TrxOrderBillRepo.New(dbConn)
	trxOrderProductSalesRepo := TrxOrderProductSalesRepo.New(dbConn)
	trxOrderTaxSalesRepo := TrxOrderTaxSalesRepo.New(dbConn)
	trxOrderPaymentSalesRepo := TrxOrderPaymentSalesRepo.New(dbConn)
	trxVisitSalesRepo := TrxVisitSalesRepo.New(dbConn)
	trxVisitMenuSalesRepo := TrxVisitMenuSalesRepo.New(dbConn)
	payMethodRepo := PayMethodRepo.New(dbConn)
	logXenditReqRepo := LogXenditReqRepo.New(dbConn)
	trxRefundRepo := TrxRefundRepo.New(dbConn)
	trxVisitRepo := TrxVisitRepo.New(dbConn)
	canvasRepo := CanvasRepo.New(dbConn)

	// init fcm admin SDK
	ctx := context.Background()
	opt := option.WithCredentialsFile(cfg.Config.GoogleApplicationCredentials)
	firebaseApp, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Printf("error initializing firebase app: %v", err)
	}
	log.Printf("Connected to firebase: app -> %+v", firebaseApp)
	// init fcm service
	fcmSvc := FcmSvc.New(&ctx, firebaseApp)

	// init services
	appSettingService := AppSettingService.New(appSettingRepo)
	userService := UserService.New(userRepo, roleRepo, userRoleRepo, userAddressRepo, userAuthRepo, userConfigRepo, logRepo)
	roleService := RoleService.New(roleRepo)
	pSpaceService := PSpaceService.New(pSpaceRepo, logRepo, userRepo)
	salesService := SalesService.New(salesRepo, logRepo, userRepo, pSpaceRepo)
	trxSalesService := TrxSalesService.New(trxSalesRepo,
		salesRepo, userRepo, logRepo, pSpaceRepo, trxVisitSalesRepo)

	merchantSvc := MerchantService.New(
		merchantRepo, merchantUsersRepo, merchantTaxRepo, merchantCategoryRepo, userRepo, roleRepo,
		userRoleRepo, userAuthRepo, shiftRepo, tasksRepo, merchantMenuRepo, logMerchRepo, logRepo, cfg.Config.Xendit.ApiKey, cfg.Config.Xendit.ApiHost)
	tasksSvc := TasksService.New(merchantUsersRepo,
		tasksRegRepo, tasksHomeRepo, tasksNonRegRepo, tasksRepo,
		tasksRegLogRepo, tasksHomeLogRepo, tasksNonRegLogRepo,
		tasksTrackingRepo, homevisitRepo, trxVisitSalesRepo)
	logMerchSvc := LogMerchSvc.New(logMerchRepo)
	homevisitSvc := HomevisitSvc.New(homevisitRepo, trxVisitSalesRepo)
	logActivitySvc := LogService.New(logRepo)
	endUserSvc := EndUserSvc.New(merchantRepo, salesRepo,
		tasksTrackingRepo, homevisitRepo, payMethodRepo, trxSalesRepo, userRepo, tasksRegRepo)
	trxSvc := TrxService.New(trxRepo, trxVisitSalesRepo,
		trxOrderBillRepo, trxOrderRepo,
		trxOrderProductSalesRepo, trxOrderTaxSalesRepo,
		trxOrderPaymentSalesRepo, trxVisitRepo, trxVisitMenuSalesRepo, homevisitRepo, merchantTaxRepo, merchantMenuRepo, userRepo, merchantUsersRepo)
	paymentSvc := PaymentSvc.New(
		trxRepo, trxOrderRepo, trxOrderBillRepo, trxOrderProductSalesRepo, merchantMenuRepo, logXenditReqRepo,
		cfg.Config.Xendit.ApiHost, cfg.Config.Xendit.ApiKey)
	trxRefundService := TrxRefundService.New(trxRefundRepo, trxVisitSalesRepo)
	var queue sync.Map
	canvasService := CanvasService.New(canvasRepo, &queue)
	// init routes
	app.MaxMultipartMemory = 8 << 20 // 8 mb limit upload
	routes.AppSettingController(app, appSettingService)
	routes.UserController(app, userService)
	routes.RoleController(app, roleService)
	routes.PSpaceController(app, pSpaceService, salesService)
	routes.TrxSalesController(app, trxSalesService, salesService, merchantSvc)
	routes.MerchantController(app, merchantSvc, userService, trxSvc)
	routes.TasksController(app, tasksSvc, merchantSvc,
		logMerchSvc, userService)
	routes.HomevisitSalesController(app, homevisitSvc, merchantSvc)
	routes.LogMerchantController(app, logMerchSvc, merchantSvc)
	routes.LogController(app, logActivitySvc)
	routes.EndUserController(app, endUserSvc, homevisitSvc,
		merchantSvc, trxSvc, fcmSvc)
	routes.PaymentController(app, paymentSvc, fcmSvc, endUserSvc,
		trxSvc, merchantSvc, homevisitSvc)
	routes.TrxController(app, trxSvc, merchantSvc, trxRefundService)
	routes.CanvassingController(app, canvasService, fcmSvc, merchantSvc, userService)

	cron := gocron.NewScheduler(time.Local)
	jobs.CronJob(cron, trxSvc, merchantSvc, tasksTrackingRepo, canvasService, fcmSvc, merchantSvc, userService)
	// init tls config
	// certMan := autocert.Manager{
	// 	Prompt:     autocert.AcceptTOS,
	// 	HostPolicy: autocert.HostWhitelist(apiHost),
	// 	Cache:      autocert.DirCache(certsDir),
	// }
	// tlsCfg := certMan.TLSConfig()
	// tlsCfg.GetCertificate = util.GetSelfSignedOrLetsEncryptCert(&certMan)
	// config tls
	// server := http.Server{
	// 	Addr:      fmt.Sprintf(":%s", apiPort),
	// 	Handler:   app,
	// 	TLSConfig: tlsCfg,
	// }
	// fmt.Println("Server listening on", server.Addr)
	// if err := server.ListenAndServeTLS("", ""); err != nil {
	// 	log.Fatal(err.Error())
	// }
	if err := app.Run(fmt.Sprintf(":%s", cfg.Config.Api.Port)); err != nil {
		log.Fatal(err)
	}
}
