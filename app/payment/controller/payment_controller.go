package controller

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"firebase.google.com/go/v4/messaging"
	"github.com/gin-gonic/gin"
	"streetbox.id/app/enduser"
	"streetbox.id/app/fcm"
	"streetbox.id/app/homevisitsales"
	"streetbox.id/app/merchant"
	"streetbox.id/app/payment"
	"streetbox.id/app/trx"
	"streetbox.id/cfg"
	"streetbox.id/entity"
	"streetbox.id/model"
	"streetbox.id/util"
)

// PaymentController ..
type PaymentController struct {
	PaymentSvc    payment.ServiceInterface
	FcmSvc        fcm.ServiceInterface
	EndUserSvc    enduser.ServiceInterface
	TrxSvc        trx.ServiceInterface
	MerchantSvc   merchant.ServiceInterface
	VisitSalesSvc homevisitsales.ServiceInterface
}

// CreateQRIS godoc
// @Summary Create QRIS Xendit Payment Method
// @Id CreateQRIS
// @Tags Payment
// @Accept application/x-www-form-urlencoded
// @Produce application/json
// @Security Token
// @Param merchant_id formData integer true "merchant ID"
// @Param amount formData integer true "total amount transaction"
// @Param types formData string true "types transaction" Enums(ORDER,VISIT)
// @Param address formData string true "Address Transaction"
// @Success 200 {object} model.ResCreateQRIS "data: model.ResCreateQRIS"
// @Router /payment/create-qrcode [post]
func (r *PaymentController) CreateQRIS(c *gin.Context) {
	merchantID := util.ParamIDToInt64(c.PostForm("merchant_id"))
	amount := util.ParamIDToInt64(c.PostForm("amount"))
	types := c.PostForm("types")
	address := c.PostForm("address")
	var err error
	// req := model.ReqCreateQRIS{}
	// merchantID := req.MerchantID
	// amount := req.Amount
	// types := req.Types
	// address := req.Address
	// productSales := req.Order.ProductSales

	// if productSales != nil {

	// 	if err = r.MerchantSvc.CheckStock(productSales); err != nil {
	// 		model.ResponseError(c, err, http.StatusUnprocessableEntity)
	// 		return
	// 	}
	// }
	if len(strings.TrimSpace(c.PostForm("amount"))) == 0 ||
		len(strings.TrimSpace(types)) == 0 ||
		len(strings.TrimSpace(address)) == 0 {
		model.ResponseError(c, "Invalid Request", http.StatusUnprocessableEntity)
		return
	}
	data := new(model.ResCreateQRIS)
	jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	var xenditID string
	var usersID int64
	if jwtModel.RoleName == util.RoleConsumer {
		if len(strings.TrimSpace(c.PostForm("merchant_id"))) == 0 {
			model.ResponseError(c, "Invalid Request", http.StatusUnprocessableEntity)
			return
		}
		merchant := r.MerchantSvc.GetByID(merchantID)
		xenditID = merchant.XenditID
		usersID = jwtModel.UserID
	} else if jwtModel.RoleName == util.RoleFoodtruck {
		merchant := r.MerchantSvc.GetInfo(jwtModel.UserID)
		xenditID = merchant.XenditID
		usersID = int64(-1)
	} else {
		model.ResponseError(c, "Unauthorized Access", http.StatusUnprocessableEntity)
		return
	}
	if data, err = r.PaymentSvc.CreateQRIS(amount, types, usersID, xenditID, address); err != nil {
		model.ResponseError(c, "Create QR Code failed", http.StatusInternalServerError)
		return
	}
	model.ResponseJSON(c, data)
	return
}

// XenditCallback godoc
// @Summary Callback executed by Xendit for Status Payment QRIS (Production Only)
// @Id XenditCallback
// @Tags Payment
// @Accept application/x-www-form-urlencoded
// @Produce application/json
// @Param x-callback-token header string true "Token Callback for Verify"
// @Param req body model.ReqXndQRISCallback true "Callback from Xendit"
// @Success 200 {object} model.ResponseSuccess "data: model.ResponseSuccess"
// @Router /payment/xendit/qris/callback [post]
func (r *PaymentController) XenditCallback(c *gin.Context) {
	req := model.ReqXndQRISCallback{}
	tokenTestMode := cfg.Config.Xendit.CallbackQrisToken
	token := c.GetHeader("x-callback-token")
	if strings.TrimSpace(token) != strings.TrimSpace(tokenTestMode) {
		log.Printf("ERROR: Token Callback Not Valid")
		model.ResponseError(c, "Token Callback Not Valid", http.StatusUnprocessableEntity)
		return
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		model.ResponseError(c, err.Error(), http.StatusUnprocessableEntity)
		log.Printf("ERROR: %s", err.Error())
		return
	}
	trx, err := r.PaymentSvc.XenditCallback(&req)
	if err != nil {
		log.Printf("ERROR: %s", err.Error())
		model.ResponseError(c, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	// send to end user
	// topic := os.Getenv("TOPIC_TRX_PUSH_NOTIF")
	trxOnly := r.TrxSvc.GetTrxByID(req.QRCode.ExternalID)
	endUser := r.EndUserSvc.GetProfile(trxOnly.UsersID)
	if endUser != nil && strings.TrimSpace(endUser.RegistrationToken) != "" {
		if trxOnly.Types == "ORDER" {
			msg := messaging.Notification{
				Title: "Online Order Nearby",
				Body:  fmt.Sprintf("Transaksi %s berhasil", trx.ID),
			}
			data := map[string]string{
				"title": "Online Order Nearby",
				"body":  fmt.Sprintf("Transaksi %s berhasil", trx.ID),
			}
			// r.FcmSvc.SendNotificationWithData(fmt.Sprintf("blast_%d", endUser.ID), &msg, &data, "id.streetbox.live")
			r.FcmSvc.SendNotificationWithDataToken(endUser.RegistrationToken, &msg, &data, "id.streetbox.live")
		} else {
			msg := messaging.Notification{
				Title: "Home Visit",
				Body:  fmt.Sprintf("Transaksi %s berhasil", trx.ID),
			}
			data := map[string]string{
				"title": "Home Visit",
				"body":  fmt.Sprintf("Transaksi %s berhasil", trx.ID),
			}
			// r.FcmSvc.SendNotificationWithData(fmt.Sprintf("blast_%d", endUser.ID), &msg, &data, "id.streetbox.live")
			r.FcmSvc.SendNotificationWithDataToken(endUser.RegistrationToken, &msg, &data, "id.streetbox.live")
		}
	}

	if trx.Types == util.TrxOrder {
		// send to POS
		trxOrder := r.TrxSvc.GetOneTrxOrderByTrxID(trx.ID)
		foodtruck := r.MerchantSvc.GetMerchantUsersByID(trxOrder.MerchantUsersID)
		if foodtruck != nil && strings.TrimSpace(foodtruck.RegistrationToken) != "" {
			msg := messaging.Notification{
				Title: "New Online Order",
				Body:  trxOrder.TrxID,
			}
			data := map[string]string{
				"title": "New Online Order",
				"body":  trxOrder.TrxID,
			}
			// r.FcmSvc.SendNotificationWithData(fmt.Sprintf("blast_%d", foodtruck.ID), &msg, &data, "com.streetbox.pos")
			r.FcmSvc.SendNotificationWithDataToken(foodtruck.RegistrationToken, &msg, &data, "com.streetbox.pos")
		}
		model.ResponseJSON(c, trx)
		return

	}

	// send to merchant admin
	merchantID := r.TrxSvc.GetMerchantIDByTrxID(trx.ID)
	admin := r.MerchantSvc.GetMerchantUsersAdminByMerchantID(merchantID)
	if admin != nil && strings.TrimSpace(admin.RegistrationToken) != "" {
		notifVisit := &messaging.Notification{
			Title: "New Home Visit Order",
			Body:  trx.ID,
		}
		log.Printf("INFO: FCM Merchant Admin Visit Order -> %s",
			r.FcmSvc.SendNotificationToOne(notifVisit, admin.RegistrationToken))
	}
	// update available home visit sales
	if err := r.VisitSalesSvc.UpdateByTrxID(trx.ID); err != nil {
		model.ResponseError(c, "Update Available Visit Sales failed",
			http.StatusUnprocessableEntity)
		log.Printf("ERROR: %s", err.Error())
		return
	}
	model.ResponseJSON(c, trx)
	return
}

// GetQRByTrxID godoc
// @Summary Get QRIS Xendit Payment Method by Transaction ID
// @Id GetQRByTrxID
// @Tags Payment
// @Security Token
// @Param trxId path string true "Transaction ID"
// @Success 200 {object} model.ResCreateQRIS "data: model.ResCreateQRIS"
// @Router /payment/qrcode/{trxId} [get]
func (r *PaymentController) GetQRByTrxID(c *gin.Context) {
	trxID := c.Param("trxId")
	data := new(model.ResCreateQRIS)
	var err error
	if data, err = r.PaymentSvc.GetQrCodeByTrxID(trxID); err != nil {
		model.ResponseError(c, "Get QRCodes Failed", http.StatusUnprocessableEntity)
		return
	}
	model.ResponseJSON(c, data)
	return
}

// SimulateQR godoc
// @Summary Simulate QRIS Xendit Payment Method (testmode only)
// @Id SimulateQR
// @Tags Payment
// @Security Token
// @Accept application/x-www-form-urlencoded
// @Produce application/json
// @Param trxId path string true "Transaction ID"
// @Success 200 {object} model.ResSimulateQR "data: model.ResSimulateQR"
// @Router /payment/simulate-qrcode/{trxId} [post]
// func (r *PaymentController) SimulateQR(c *gin.Context) {
// 	trxID := c.Param("trxId")
// 	data := new(model.ResSimulateQR)
// 	trx := new(entity.Trx)
// 	var err error
// 	if data, trx, err = r.PaymentSvc.SimulateQR(trxID); err != nil {
// 		model.ResponseError(c, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	// send to end user
// 	// topic := os.Getenv("TOPIC_TRX_PUSH_NOTIF")
// 	endUser := r.EndUserSvc.GetProfile(trx.UsersID)
// 	if endUser != nil && strings.TrimSpace(endUser.RegistrationToken) != "" {
// 		msg := messaging.Notification{
// 			Title: "Online Order",
// 			Body:  fmt.Sprintf("Transaksi %s berhasil", trx.ID),
// 		}
// 		data := map[string]string{
// 			"title": "Online Order",
// 			"body":  fmt.Sprintf("Transaksi %s berhasil", trx.ID),
// 		}
// 		r.FcmSvc.SendNotificationWithData(fmt.Sprintf("blast_%d", endUser.ID), &msg, &data, "id.streetbox.live")
// 	}

// 	if trx.Types == util.TrxOrder {
// 		// send to POS
// 		trxOrder := r.TrxSvc.GetOneTrxOrderByTrxID(trx.ID)
// 		foodtruck := r.MerchantSvc.GetMerchantUsersByID(trxOrder.MerchantUsersID)
// 		if foodtruck != nil && strings.TrimSpace(foodtruck.RegistrationToken) != "" {
// 			msg := messaging.Notification{
// 				Title: "New Online Order",
// 				Body:  trxOrder.TrxID,
// 			}
// 			data := map[string]string{
// 				"title": "New Online Order",
// 				"body":  trxOrder.TrxID,
// 			}
// 			r.FcmSvc.SendNotificationWithData(fmt.Sprintf("blast_%d", foodtruck.ID), &msg, &data, "com.streetbox.pos")
// 		}
// 		model.ResponseJSON(c, data)
// 		return

// 	}

// 	// send to merchant admin
// 	merchantID := r.TrxSvc.GetMerchantIDByTrxID(trx.ID)
// 	admin := r.MerchantSvc.GetMerchantUsersAdminByMerchantID(merchantID)
// 	if admin != nil && strings.TrimSpace(admin.RegistrationToken) != "" {
// 		notifVisit := &messaging.Notification{
// 			Title: "New Visit Order",
// 			Body:  trx.ID,
// 		}
// 		log.Printf("INFO: FCM Merchant Admin Visit Order -> %s",
// 			r.FcmSvc.SendNotificationToOne(notifVisit, admin.RegistrationToken))
// 	}
// 	// update available home visit sales
// 	if err := r.VisitSalesSvc.UpdateByTrxID(trxID); err != nil {
// 		model.ResponseError(c, "Update Available Visit Sales failed",
// 			http.StatusUnprocessableEntity)
// 		return
// 	}
// 	model.ResponseJSON(c, data)
// 	return

// }
func (r *PaymentController) SimulateQR(c *gin.Context) {
	trxID := c.Param("trxId")
	data := new(model.ResSimulateQR)
	trx := new(entity.Trx)
	var err error
	if data, trx, err = r.PaymentSvc.SimulateQR(trxID); err != nil {
		model.ResponseError(c, err.Error(), http.StatusInternalServerError)
		return
	}

	// send to end user
	// topic := os.Getenv("TOPIC_TRX_PUSH_NOTIF")
	// jwtModel, _ := util.ExtractTokenMetadata(c.Request)
	trxOnly := r.TrxSvc.GetTrxByID(trxID)
	endUser := r.EndUserSvc.GetProfile(trx.UsersID)
	if endUser != nil && strings.TrimSpace(endUser.RegistrationToken) != "" {
		if trxOnly.Types == "ORDER" {
			msg := messaging.Notification{
				Title: "Online Order Nearby",
				Body:  fmt.Sprintf("Transaksi %s berhasil", trx.ID),
			}
			data := map[string]string{
				"title": "Online Order Nearby",
				"body":  fmt.Sprintf("Transaksi %s berhasil", trx.ID),
			}
			// r.FcmSvc.SendNotificationWithData(fmt.Sprintf("blast_%d", endUser.ID), &msg, &data, "id.streetbox.live")
			r.FcmSvc.SendNotificationWithDataToken(endUser.RegistrationToken, &msg, &data, "id.streetbox.live")
		} else {
			msg := messaging.Notification{
				Title: "Home Visit",
				Body:  fmt.Sprintf("Transaksi %s berhasil", trx.ID),
			}
			data := map[string]string{
				"title": "Home Visit",
				"body":  fmt.Sprintf("Transaksi %s berhasil", trx.ID),
			}
			// r.FcmSvc.SendNotificationWithData(fmt.Sprintf("blast_%d", endUser.ID), &msg, &data, "id.streetbox.live")
			r.FcmSvc.SendNotificationWithDataToken(endUser.RegistrationToken, &msg, &data, "id.streetbox.live")
		}
	}

	if trx.Types == util.TrxOrder {
		// send to POS
		trxOrder := r.TrxSvc.GetOneTrxOrderByTrxID(trxID)
		foodtruck := r.MerchantSvc.GetMerchantUsersByID(trxOrder.MerchantUsersID)
		if foodtruck != nil {
			msg := messaging.Notification{
				Title: "New Online Order",
				Body:  trxOrder.TrxID,
			}
			data := map[string]string{
				"title": "New Online Order",
				"body":  trxOrder.TrxID,
			}
			// r.FcmSvc.SendNotificationWithData(fmt.Sprintf("blast_%d", foodtruck.ID), &msg, &data, "com.streetbox.pos")
			r.FcmSvc.SendNotificationWithDataToken(foodtruck.RegistrationToken, &msg, &data, "com.streetbox.pos")
		}
		model.ResponseJSON(c, data)
		return

	}

	// send to merchant admin
	merchantID := r.TrxSvc.GetMerchantIDByTrxID(trx.ID)
	admin := r.MerchantSvc.GetMerchantUsersAdminByMerchantID(merchantID)
	if admin != nil && strings.TrimSpace(admin.RegistrationToken) != "" {
		notifVisit := &messaging.Notification{
			Title: "New Home Visit Order",
			Body:  trxID,
		}
		log.Printf("INFO: FCM Merchant Admin Visit Order -> %s",
			r.FcmSvc.SendNotificationToOne(notifVisit, admin.RegistrationToken))
	}
	// update available home visit sales
	if err := r.VisitSalesSvc.UpdateByTrxID(trxID); err != nil {
		model.ResponseError(c, "Update Available Visit Sales failed",
			http.StatusUnprocessableEntity)
		return
	}
	model.ResponseJSON(c, data)
	return

}
