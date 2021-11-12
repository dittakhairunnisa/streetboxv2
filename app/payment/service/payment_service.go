package service

import (
	"bytes"
	bs64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/jinzhu/copier"
	"streetbox.id/app/logxenditreq"
	"streetbox.id/app/merchantmenu"
	"streetbox.id/app/payment"
	"streetbox.id/app/trx"
	"streetbox.id/app/trxorder"
	"streetbox.id/app/trxorderbill"
	"streetbox.id/app/trxorderproductsales"
	"streetbox.id/cfg"
	"streetbox.id/entity"
	"streetbox.id/model"
	"streetbox.id/util"
)

var (
	xenditHeader http.Header
	xenditHost   string
)

// PaymentService ...
type PaymentService struct {
	TrxRepo                  trx.RepoInterface
	TrxOrderRepo             trxorder.RepoInterface
	TrxOrderBillRepo         trxorderbill.RepoInterface
	TrxOrderProductSalesRepo trxorderproductsales.RepoInterface
	MerchantMenusRepo        merchantmenu.RepoInterface
	LogXenditReqRepo         logxenditreq.RepoInterface
}

// New ...
func New(trxRepo trx.RepoInterface,
	trxOrderRepo trxorder.RepoInterface,
	trxOrderBillRepo trxorderbill.RepoInterface,
	trxOrderProductSales trxorderproductsales.RepoInterface,
	merchantMenusRepo merchantmenu.RepoInterface,
	logXenditReqRepo logxenditreq.RepoInterface,
	xenditHosts, xenditAPIKey string) payment.ServiceInterface {
	data := xenditAPIKey + ":"
	xenditAuthHeader := "Basic " + bs64.StdEncoding.EncodeToString([]byte(data))
	xenditHeader = http.Header{
		"Authorization": []string{xenditAuthHeader},
		"Content-Type":  []string{"application/x-www-form-urlencoded"},
	}
	xenditHost = xenditHosts
	return &PaymentService{trxRepo, trxOrderRepo, trxOrderBillRepo, trxOrderProductSales, merchantMenusRepo, logXenditReqRepo}
}

// CreateQRIS ...
func (s *PaymentService) CreateQRIS(
	amount int64, types string, usersID int64, xenditID, address string) (*model.ResCreateQRIS, error) {
	callbackURL := cfg.Config.Xendit.CallbackQrisUrl
	reqData := fmt.Sprintf("external_id=%s&type=%s&callback_url=%s&amount=%d",
		util.GenerateTrxID(), "DYNAMIC", callbackURL, amount)
	postReq, err := http.NewRequest("POST", xenditHost+"/qr_codes", bytes.NewBufferString(reqData))
	if err != nil {
		log.Printf("ERROR: %s", err.Error())
		return nil, err
	}
	header := http.Header{}
	for k, v := range xenditHeader {
		header[k] = v
	}
	fmt.Printf("Header: %v", header)
	postReq.Header = header
	postReq.Header.Add("for-user-id", xenditID)

	var body []byte
	if body, err = util.DoRequest(postReq); err != nil {
		return nil, err
	}
	respModelXnd := new(model.ResCreateQRISXnd)
	respModel := new(model.ResCreateQRIS)
	if err := json.Unmarshal(body, &respModelXnd); err != nil {
		log.Printf("ERROR: %s", err.Error())
		return nil, err
	}
	copier.Copy(&respModel, respModelXnd)
	// save trx
	trx := new(entity.Trx)
	trx.Status = util.TrxStatusPending
	trx.ID = respModel.ExternalID
	trx.Types = types
	trx.UsersID = usersID
	trx.Address = address
	trx.QrCode = respModelXnd.QrString
	if db, err := s.TrxRepo.Create(trx); err == nil {
		// save log xendit
		logXendit := new(entity.LogXenditRequest)
		logXendit.TrxID = respModel.ExternalID
		logXendit.ResponseData = string(body)
		bodyReqRead, _ := postReq.GetBody()
		bodyReq, _ := ioutil.ReadAll(bodyReqRead)
		logXendit.RequestData = string(bodyReq)
		if err := s.LogXenditReqRepo.Create(db, logXendit); err == nil {
			db.Commit()
			return respModel, nil
		}
		db.Rollback()
	}
	return nil, errors.New("Failed")
}

// GetQrCodeByTrxID get qr code xendit
func (s *PaymentService) GetQrCodeByTrxID(trxID string) (*model.ResCreateQRIS, error) {
	getReq, err := http.NewRequest("GET", xenditHost+"/qr_codes/"+trxID, nil)
	if err != nil {
		log.Printf("ERROR: %s", err.Error())
		return nil, err
	}
	getReq.Header = xenditHeader

	var body []byte
	if body, err = util.DoRequest(getReq); err != nil {
		return nil, err
	}
	respModelXnd := new(model.ResCreateQRISXnd)
	respModel := new(model.ResCreateQRIS)
	if err := json.Unmarshal(body, &respModelXnd); err != nil {
		log.Printf("ERROR: %s", err.Error())
		return nil, err
	}
	copier.Copy(&respModel, respModelXnd)
	return respModel, nil
}

// SimulateQR ...
func (s *PaymentService) SimulateQR(trxID string) (*model.ResSimulateQR, *entity.Trx, error) {
	if data := s.TrxRepo.FindByID(trxID); data.Status == util.TrxStatusSuccess {
		return nil, nil, errors.New("Already Success")
	}
	postReq, err := http.NewRequest("POST",
		xenditHost+"/qr_codes/"+trxID+"/payments/simulate", nil)
	if err != nil {
		log.Printf("ERROR: %s", err.Error())
		return nil, nil, err
	}
	postReq.Header = xenditHeader

	var body []byte
	if body, err = util.DoRequest(postReq); err != nil {
		return nil, nil, err
	}
	respModelXnd := new(model.ResSimulateQRXnd)
	respModel := new(model.ResSimulateQR)
	if err := json.Unmarshal(body, &respModelXnd); err != nil {
		log.Printf("ERROR: %s", err.Error())
		return nil, nil, err
	}
	copier.Copy(&respModel, respModelXnd)
	// update status to success
	trx := s.TrxRepo.FindByID(trxID)
	trx.Status = util.TrxStatusSuccess
	// runes := []rune(respModel.ID)
	// trx.ExternalID = string(runes[5:])
	s.TrxRepo.Update(trx, trxID)
	data := s.TrxRepo.FindByID(trxID)
	return respModel, data, nil
}

// XenditCallback ..
func (s *PaymentService) XenditCallback(req *model.ReqXndQRISCallback) (*entity.Trx, error) {
	status := req.Status
	data := s.TrxRepo.FindByID(req.QRCode.ExternalID)
	// save log xendit
	logXendit := new(entity.LogXenditRequest)
	logXendit.TrxID = req.QRCode.ExternalID
	logXendit.RequestData = "CALLBACK"
	resp, _ := json.Marshal(req)
	logXendit.ResponseData = string(resp)
	if err := s.LogXenditReqRepo.Create(nil, logXendit); err != nil {
		log.Printf("ERROR: %s", err.Error())
	}
	// get status from xendit
	if status == "COMPLETED" {
		trxOrder := s.TrxOrderRepo.FindByTrxID(data.ID)
		trxOrderBill := s.TrxOrderBillRepo.FindByTrxOrderID(trxOrder.ID)
		if len(*trxOrderBill) > 0 {
			for _, billData := range *trxOrderBill {
				qryProductSales := entity.TrxOrderProductSales{TrxOrderBillID: billData.ID}
				trxProductSales := s.TrxOrderProductSalesRepo.FindAll(&qryProductSales)
				if len(*trxProductSales) > 0 {
					for _, trxOrderProductSales := range *trxProductSales {
						s.MerchantMenusRepo.UpdateStock(trxOrderProductSales.Qty, trxOrderProductSales.MerchantMenuID)
					}
				}
			}
		}
		data.Status = util.TrxStatusSuccess
		runes := []rune(req.ID)
		data.ExternalID = string(runes[5:])
		s.TrxRepo.Update(data, req.QRCode.ExternalID)
		return data, nil
	}
	data.Status = util.TrxStatusFailed
	s.TrxRepo.Update(data, req.QRCode.ExternalID)
	log.Printf("INFO: Transaction %s Failed", req.QRCode.ExternalID)
	return data, errors.New("Failed")
}
