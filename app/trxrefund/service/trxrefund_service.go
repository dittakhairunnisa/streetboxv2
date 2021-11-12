package service

import (
	"time"

	"github.com/jinzhu/copier"
	"streetbox.id/app/trxrefund"
	"streetbox.id/app/trxvisitsales"
	"streetbox.id/entity"
	"streetbox.id/model"
	"streetbox.id/util"
)

// TrxRefundService ...
type TrxRefundService struct {
	Repo         trxrefund.RepoInterface
	TrxHomeVisit trxvisitsales.RepoInterface
}

// New ..
func New(repo trxrefund.RepoInterface, trxHomeVisit trxvisitsales.RepoInterface) trxrefund.ServiceInterface {
	return &TrxRefundService{repo, trxHomeVisit}
}

// CreateRefundParkingSpace ...
func (r *TrxRefundService) CreateRefundParkingSpace(req *model.ReqRefundParkingSpaceSales) error {
	//Create Refund
	refund := new(entity.TrxRefund)
	refund.Types = util.TrxRefundParkingSpace
	refund.CreatedAt = time.Now()
	refundID, db, err := r.Repo.CreateRefund(refund)
	if err != nil {
		return err
	}

	//Create Refund Parking Space Sales
	refundParkingSpace := new(entity.TrxRefundSpace)
	copier.Copy(&refundParkingSpace, req)
	refundParkingSpace.TrxRefundID = refundID
	refundParkingSpace.CreatedAt = time.Now()
	db, err = r.Repo.CreateRefundSpace(refundParkingSpace, db)
	if err != nil {
		return err
	}
	db.Commit()
	return nil
}

// CreateRefundHomeVisit ...
func (r *TrxRefundService) CreateRefundHomeVisit(req *model.ReqRefundHomeVisit, merchantID int64) error {
	//Create Refund
	refund := new(entity.TrxRefund)
	refund.Types = util.TrxRefundHomeVisit
	refund.CreatedAt = time.Now()
	refundID, db, err := r.Repo.CreateRefund(refund)
	if err != nil {
		return err
	}

	refundHomeVisit := new(entity.TrxRefundVisit)
	refundHomeVisit.Amount = req.Amount
	refundHomeVisit.TrxHomevisitSalesID = req.ID
	refundHomeVisit.TrxRefundID = refundID
	refundHomeVisit.CreatedAt = time.Now()
	db, err = r.Repo.CreateRefundVisit(refundHomeVisit, db)
	if err != nil {
		return err
	}

	db.Commit()
	return nil
}
