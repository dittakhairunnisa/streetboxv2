package jobs

import (
	"github.com/go-co-op/gocron"
	"streetbox.id/app/canvassing"
	Ctrl "streetbox.id/app/canvassing/controller"
	"streetbox.id/app/fcm"
	"streetbox.id/app/merchant"
	"streetbox.id/app/taskstracking"
	"streetbox.id/app/trx"

	//Trx "streetbox.id/app/trx/controller"
	"streetbox.id/app/user"
)

// CronJob ..
func CronJob(
	schedulerApp *gocron.Scheduler,
	trxSvc trx.ServiceInterface,
	merchantSvc merchant.ServiceInterface,
	trackingRepo taskstracking.RepoInterface,
	Svc canvassing.Service,
	Fcm fcm.ServiceInterface,
	Merchant merchant.ServiceInterface,
	User user.ServiceInterface) {
	// trx := &Trx.TrxController{
	// 	TrxSvc:      trxSvc,
	// 	MerchantSvc: merchantSvc,
	// }
	ctrl := Ctrl.CanvassingController{
		Svc:      Svc,
		Fcm:      Fcm,
		Merchant: Merchant,
		User:     User,
	}
	//schedulerApp.Every(5).Seconds().Do(trx.ParseTrxJobs)
	schedulerApp.Every(1).Day().At("00:01").Do(trackingRepo.DeleteAll)
	schedulerApp.Every(10).Seconds().Do(Svc.ExpireNotification)
	schedulerApp.Every(30).Seconds().Do(ctrl.AutoBlast)
	// schedulerApp.Cron("*/5 8-22 * * *").Do(ctrl.AutoBlast)
	schedulerApp.Cron("0 0 * * *").Do(Svc.ResetQueue)
	schedulerApp.Every(10).Seconds().Do(Svc.ExpireCall)
	schedulerApp.StartAsync()
}
