package payment

import (
	"streetbox.id/entity"
	"streetbox.id/model"
)

// ServiceInterface ...
type ServiceInterface interface {
	CreateQRIS(amount int64, types string, usersID int64, xenditID, address string) (*model.ResCreateQRIS, error)
	GetQrCodeByTrxID(string) (*model.ResCreateQRIS, error)
	SimulateQR(string) (*model.ResSimulateQR, *entity.Trx, error)
	XenditCallback(*model.ReqXndQRISCallback) (*entity.Trx, error)
}
