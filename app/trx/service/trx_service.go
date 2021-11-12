package service

import (
	"errors"
	"log"
	"time"

	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
	"streetbox.id/app/homevisitsales"
	"streetbox.id/app/merchantmenu"
	"streetbox.id/app/merchanttax"
	"streetbox.id/app/merchantusers"
	"streetbox.id/app/trx"
	"streetbox.id/app/trxorder"
	"streetbox.id/app/trxorderbill"
	"streetbox.id/app/trxorderpaymentsales"
	"streetbox.id/app/trxorderproductsales"
	"streetbox.id/app/trxordertaxsales"
	"streetbox.id/app/trxvisit"
	"streetbox.id/app/trxvisitmenusales"
	"streetbox.id/app/trxvisitsales"
	"streetbox.id/app/user"
	"streetbox.id/entity"
	"streetbox.id/model"
	"streetbox.id/util"
)

// TrxService ..
type TrxService struct {
	TrxRepo                  trx.RepoInterface
	TrxOrderBillRepo         trxorderbill.RepoInterface
	TrxVisitSalesRepo        trxvisitsales.RepoInterface
	TrxOrderRepo             trxorder.RepoInterface
	TrxOrderProductSalesRepo trxorderproductsales.RepoInterface
	TrxOrderTaxSalesRepo     trxordertaxsales.RepoInterface
	TrxOrderPaymentRepo      trxorderpaymentsales.RepoInterface
	TrxVisitRepo             trxvisit.RepoInterface
	TrxVisitMenuSalesRepo    trxvisitmenusales.Repository
	VisitSalesRepo           homevisitsales.RepoInterface
	MerchantTaxRepo          merchanttax.RepoInterface
	MerchantMenusRepo        merchantmenu.RepoInterface
	UsersRepo                user.RepoInterface
	MerchantUsersRepo        merchantusers.RepoInterface
}

// New ..
func New(
	repo trx.RepoInterface,
	trxVisitSales trxvisitsales.RepoInterface,
	trxOrderBill trxorderbill.RepoInterface,
	trxOrder trxorder.RepoInterface,
	trxOrderProductSales trxorderproductsales.RepoInterface,
	trxOrderTaxSales trxordertaxsales.RepoInterface,
	trxOrderPaymentRepo trxorderpaymentsales.RepoInterface,
	trxVisitRepo trxvisit.RepoInterface,
	trxVisitMenuSalesRepo trxvisitmenusales.Repository,
	visitSales homevisitsales.RepoInterface,
	merchantTaxRepo merchanttax.RepoInterface,
	merchantMenusRepo merchantmenu.RepoInterface,
	usersRepo user.RepoInterface,
	merchantUsersRepo merchantusers.RepoInterface,
) trx.ServiceInterface {
	return &TrxService{repo, trxOrderBill,
		trxVisitSales, trxOrder, trxOrderProductSales,
		trxOrderTaxSales, trxOrderPaymentRepo, trxVisitRepo, trxVisitMenuSalesRepo,
		visitSales, merchantTaxRepo, merchantMenusRepo, usersRepo, merchantUsersRepo}
}

// CountTrx ..
func (r *TrxService) CountTrx() int64 {
	return r.TrxRepo.CountTrx()
}

// GetOrderTrx ..
func (r *TrxService) GetOrderTrx(merchantID int64, startDates, endDates, keyword string) *model.ResTrxOrderList {
	return r.TrxRepo.GetOrderTrx(merchantID, startDates, endDates, keyword)
}

// CreateSyncTrx ..
func (r *TrxService) CreateSyncTrx(data *model.ReqCreateSyncTrx, merchantID int64) *entity.TrxSync {
	trxSync := new(entity.TrxSync)
	copier.Copy(&trxSync, data)
	trxSync.MerchantID = merchantID
	trxSync.SyncDate = util.MillisToTime(data.SyncDate)
	trxSync.BusinessDate = util.MillisToTime(data.BusinessDate)
	trxSync.Status = 0
	r.TrxRepo.CreateSyncTrx(trxSync)
	if trxSync.ID == 0 {
		log.Printf("ERROR: Create Sync Transaction Failed")
		return nil
	}
	return trxSync
}

// UpdateStatusSyncTrx ...
func (r *TrxService) UpdateStatusSyncTrx(uniqueID string, merchantID int64, status int, db *gorm.DB) (*gorm.DB, error) {
	return r.TrxRepo.UpdateStatusSyncTrx(uniqueID, merchantID, status, db)
}

// CreateTrxVisit ...
func (r *TrxService) CreateTrxVisit(
	req *model.ReqCreateVisitTrx) (*entity.TrxVisit, error) {
	trxVisit := new(entity.TrxVisit)
	copier.Copy(&trxVisit, req)
	if db, err := r.TrxVisitRepo.Create(trxVisit); err == nil {
		for _, salesData := range req.VisitSales {
			trxVisitSales := new(entity.TrxHomevisitSales)
			copier.Copy(&trxVisitSales, salesData)
			trxVisitSales.TrxVisitID = trxVisit.ID
			trxVisitSales.Status = util.TrxVisitStatusOpen
			if err := r.TrxVisitSalesRepo.Create(db, trxVisitSales); err != nil {
				db.Rollback()
				return trxVisit, err
			}
			for _, menu := range salesData.Menus {
				trxVisitMenuSales := new(entity.TrxHomevisitMenuSales)
				copier.Copy(&trxVisitMenuSales, menu)
				trxVisitMenuSales.TrxHomevisitSalesID = trxVisitSales.ID
				if err := r.TrxVisitMenuSalesRepo.Create(db, trxVisitMenuSales); err == nil {
					continue
				} else {
					db.Rollback()
					return trxVisit, err
				}
			}
		}
		db.Commit()
		trxVisit := r.TrxVisitRepo.FindOne(trxVisit.ID)
		return trxVisit, nil
	}
	return nil, errors.New("Failed")
}

// CreateTrxOrderOnline ...
func (r *TrxService) CreateTrxOrderOnline(req *model.ReqTrxOrderOnline) (*entity.Trx, error) {
	if order := r.TrxOrderRepo.FindByTrxID(req.TrxID); order.ID == 0 {
		if db, trxOrderID, err := r.TrxOrderRepo.CreateOnline(req); err == nil {
			for _, orderBill := range req.OrderBills {
				if trxOrderBillID := r.TrxOrderBillRepo.
					CreateOnline(&orderBill, trxOrderID, db); trxOrderBillID > 0 {
					if len(req.ProductSales) > 0 {
						for _, productSales := range req.ProductSales {
							// if err := r.MerchantMenusRepo.UpdateStock(productSales.Qty, productSales.MerchantMenuID); err != nil {
							// 	return nil, err
							// }
							if err := r.TrxOrderProductSalesRepo.CreateOnline(
								&productSales, db, trxOrderID, trxOrderBillID); err != nil {
								return nil, err
							}
						}
					}
					if len(req.TaxSales) > 0 {
						for _, taxSales := range req.TaxSales {
							if err := r.TrxOrderTaxSalesRepo.CreateOnline(
								&taxSales, trxOrderID, trxOrderBillID, db); err != nil {
								return nil, err
							}
						}
					}
					if len(req.PaymentSales) > 0 {
						for _, trxOrderPayment := range req.PaymentSales {
							if err := r.TrxOrderPaymentRepo.CreateOnline(
								&trxOrderPayment, trxOrderBillID, db); err != nil {
								return nil, err
							}
						}
					}
				}
			}

			db.Commit()
			trx := r.TrxRepo.FindByID(req.TrxID)
			return trx, nil
		}
	}

	return nil, errors.New("Failed")
}

// CreateTrxOrder ..
func (r *TrxService) CreateTrxOrder(
	trxSync *model.ReqTrxOrderList, userID int64, merchantID int64, uniqueID string) error {
	var trxOrderID int64

	if trxSync.Order.TrxID == "" {
		trx := new(entity.Trx)
		//businessDates := util.MillisToTime(trxSync.Order.BusinessDate).Local()
		trx.ID = util.GenerateTrxID()
		trx.Types = util.TrxOrder
		trx.Status = util.TrxStatusSuccess
		trx.UsersID = trxSync.Order.UserID
		trx.CreatedAt = util.MillisToTime(trxSync.Order.CreatedAt).Local()
		updatedAt := util.MillisToTime(trxSync.Order.UpdatedAt).Local()
		trx.UpdatedAt = &updatedAt
		db, err := r.TrxRepo.Create(trx)
		if err != nil {
			return err
		}

		trxOrderID, db, err = r.TrxOrderRepo.Create(trx.ID, &trxSync.Order, db)
		if err != nil {
			return err
		}
		trxOrderBill := new(entity.TrxOrderBill)
		var billIDs []int64
		var billID int64

		for _, value := range trxSync.OrderBills {
			copier.Copy(&trxOrderBill, value)
			trxOrderBill.ID = 0
			trxOrderBill.TrxOrderID = trxOrderID
			trxOrderBill.BusinessDate = util.MillisToTime(value.BusinessDate).Local()
			trxOrderBill.CreatedAt = util.MillisToTime(value.CreatedAt).Local()
			updatedAt = util.MillisToTime(value.UpdatedAt).Local()
			trxOrderBill.UpdatedAt = &updatedAt
			billID, db, err = r.TrxOrderBillRepo.Create(trxOrderBill, db)
			if err != nil {
				return err
			}
			billIDs = append(billIDs, billID)
		}

		trxOrderProductSales := new(entity.TrxOrderProductSales)
		var trxOrderProductSalesIDs []int64
		var trxOrderProductSalesID int64
		trxOrderTaxSales := new(entity.TrxOrderTaxSales)
		trxOrderPaymentSales := new(entity.TrxOrderPaymentSales)
		for _, value := range billIDs {
			for _, value2 := range trxSync.ProductSales {
				copier.Copy(&trxOrderProductSales, value2)
				trxOrderProductSales.ID = 0
				trxOrderProductSales.TrxOrderBillID = value
				trxOrderProductSales.CreatedAt = util.MillisToTime(value2.CreatedAt).Local()
				trxOrderProductSales.BusinessDate = util.MillisToTime(value2.BusinessDate).Local()
				updatedAt = util.MillisToTime(value2.UpdatedAt).Local()
				trxOrderProductSales.UpdatedAt = &updatedAt
				if err := r.MerchantMenusRepo.UpdateStock(trxOrderProductSales.Qty, trxOrderProductSales.MerchantMenuID); err != nil {
					return err
				}
				trxOrderProductSalesID, db, err = r.TrxOrderProductSalesRepo.Create(trxOrderProductSales, db)
				if err != nil {
					return err
				}
				trxOrderProductSalesIDs = append(trxOrderProductSalesIDs, trxOrderProductSalesID)
			}
			for _, value2 := range trxSync.TaxSales {
				copier.Copy(&trxOrderTaxSales, value2)
				trxOrderTaxSales.ID = 0
				trxOrderTaxSales.TrxOrderBillID = value
				trxOrderTaxSales.CreatedAt = util.MillisToTime(value2.CreatedAt).Local()
				updatedAt = util.MillisToTime(value2.UpdatedAt).Local()
				trxOrderTaxSales.UpdatedAt = &updatedAt
				merchantTax := r.MerchantTaxRepo.GetTax(merchantID)
				trxOrderTaxSales.MerchantTaxID = merchantTax.ID
				db, err = r.TrxOrderTaxSalesRepo.CreateOffline(trxOrderTaxSales, db)
				if err != nil {
					return err
				}
			}
			for _, value2 := range trxSync.PaymentSales {
				copier.Copy(&trxOrderPaymentSales, value2)
				trxOrderPaymentSales.ID = 0
				trxOrderPaymentSales.TrxOrderBillID = value
				trxOrderPaymentSales.CreatedAt = util.MillisToTime(value2.CreatedAt).Local()
				updatedAt = util.MillisToTime(value2.UpdatedAt).Local()
				trxOrderPaymentSales.UpdatedAt = &updatedAt
				db, err = r.TrxOrderPaymentRepo.CreateOffline(trxOrderPaymentSales, db)
				if err != nil {
					return err
				}
			}
		}

		db, err = r.UpdateStatusSyncTrx(uniqueID, merchantID, 1, db)
		if err != nil {
			return err
		}
		db.Commit()
	} else {
		trx := r.TrxRepo.FindByID(trxSync.Order.TrxID)
		updatedAt := util.MillisToTime(trxSync.Order.UpdatedAt).Local()

		trx.UsersID = trxSync.Order.UserID
		db, err := r.TrxRepo.UpdateByTrxID(trx, trx.ID)
		if err != nil {
			return err
		}

		trxOrderID, db, err = r.TrxOrderRepo.Create(trx.ID, &trxSync.Order, db)
		if err != nil {
			return err
		}
		trxOrderBill := new(entity.TrxOrderBill)
		var billIDs []int64
		var billID int64

		for _, value := range trxSync.OrderBills {
			copier.Copy(&trxOrderBill, value)
			trxOrderBill.ID = 0
			trxOrderBill.TrxOrderID = trxOrderID
			trxOrderBill.BusinessDate = util.MillisToTime(value.BusinessDate).Local()
			trxOrderBill.CreatedAt = util.MillisToTime(value.CreatedAt).Local()
			updatedAt = util.MillisToTime(value.UpdatedAt).Local()
			trxOrderBill.UpdatedAt = &updatedAt
			billID, db, err = r.TrxOrderBillRepo.Create(trxOrderBill, db)
			if err != nil {
				return err
			}
			billIDs = append(billIDs, billID)
		}

		trxOrderProductSales := new(entity.TrxOrderProductSales)
		var trxOrderProductSalesIDs []int64
		var trxOrderProductSalesID int64
		trxOrderTaxSales := new(entity.TrxOrderTaxSales)
		trxOrderPaymentSales := new(entity.TrxOrderPaymentSales)
		for _, value := range billIDs {
			for _, value2 := range trxSync.ProductSales {
				copier.Copy(&trxOrderProductSales, value2)
				trxOrderProductSales.ID = 0
				trxOrderProductSales.TrxOrderBillID = value
				trxOrderProductSales.CreatedAt = util.MillisToTime(value2.CreatedAt).Local()
				trxOrderProductSales.BusinessDate = util.MillisToTime(value2.BusinessDate).Local()
				updatedAt = util.MillisToTime(value2.UpdatedAt).Local()
				trxOrderProductSales.UpdatedAt = &updatedAt
				trxOrderProductSalesID, db, err = r.TrxOrderProductSalesRepo.Create(trxOrderProductSales, db)
				if err != nil {
					return err
				}
				trxOrderProductSalesIDs = append(trxOrderProductSalesIDs, trxOrderProductSalesID)
			}
			for _, value2 := range trxSync.TaxSales {
				copier.Copy(&trxOrderTaxSales, value2)
				trxOrderTaxSales.ID = 0
				trxOrderTaxSales.TrxOrderBillID = value
				trxOrderTaxSales.CreatedAt = util.MillisToTime(value2.CreatedAt).Local()
				updatedAt = util.MillisToTime(value2.UpdatedAt).Local()
				trxOrderTaxSales.UpdatedAt = &updatedAt
				db, err = r.TrxOrderTaxSalesRepo.CreateOffline(trxOrderTaxSales, db)
				if err != nil {
					return err
				}
			}
			for _, value2 := range trxSync.PaymentSales {
				copier.Copy(&trxOrderPaymentSales, value2)
				trxOrderPaymentSales.ID = 0
				trxOrderPaymentSales.TrxOrderBillID = value
				trxOrderPaymentSales.CreatedAt = util.MillisToTime(value2.CreatedAt).Local()
				updatedAt = util.MillisToTime(value2.UpdatedAt).Local()
				trxOrderPaymentSales.UpdatedAt = &updatedAt
				db, err = r.TrxOrderPaymentRepo.CreateOffline(trxOrderPaymentSales, db)
				if err != nil {
					return err
				}
			}
		}

		db, err = r.UpdateStatusSyncTrx(uniqueID, merchantID, 1, db)
		if err != nil {
			return err
		}
		db.Commit()
	}

	return nil
}

// GetOrderHistoryByUsersID ..
func (r *TrxService) GetOrderHistoryByUsersID(limit, page int, usersID int64, filter string) model.Pagination {
	resp := make([]model.ResTrxHistory, 0)
	var status = ""
	if filter == "ongoing" {
		status = "= 'PENDING'"
	} else if filter == "history" {
		status = "<> 'PENDING'"
	}
	trxParent, count, offset := r.TrxRepo.FindByUsersID(limit, page, usersID, status)

	for _, v := range *trxParent {
		data := new(model.ResTrxHistory)
		copier.Copy(&data, v)

		if v.Types == util.TrxHomeVisit {
			visit := r.TrxVisitRepo.FindByTrxID(v.ID)
			data.Amount = visit.Deposit
			data.MerchantName = visit.Name
			data.Logo = visit.Logo
			data.Address = visit.Address
			data.Notes = visit.Notes
			data.Phone = visit.Phone

			if ft := r.UsersRepo.GetFoodtruckByTrxVisitID(visit.ID); ft != nil {
				data.Phone = ft.Phone
			}
			// init detail
			data.Detail = *r.getOrderHistoryDetail(v.ID, util.TrxHomeVisit, nil)
			if filter == "ongoing" && data.Status == "PENDING" {
				resp = append(resp, *data)
			} else if filter == "history" && data.Status != "PENDING" {
				resp = append(resp, *data)
			} else if filter == "all" {
				resp = append(resp, *data)
			}
			continue
		}

		order := r.TrxOrderRepo.FindByTrxID(v.ID)
		data.Amount = int64(int(order.GrandTotal))
		data.MerchantName = order.MerchantName
		data.Logo = order.MerchantLogo
		data.CreatedAt = order.CreatedAt
		data.Address = order.Address
		data.Status = order.Status
		merchantUsers := r.UsersRepo.GetByMerchantUsersID(order.MerchantUsersID)
		if merchantUsers != nil {
			data.Phone = merchantUsers.Phone
		}
		data.Notes = order.Note
		// init detail
		data.Detail = *r.getOrderHistoryDetail(v.ID, util.TrxOrder, order)

		// data.Amount = int64(int(data.Detail.PaymentDetails.Total))
		// if filter == "ongoing" && data.Status == "PENDING" {
		// 	resp = append(resp, *data)
		// } else if filter == "history" && data.Status != "PENDING" {
		// 	resp = append(resp, *data)
		// } else if filter == "all" {
		// 	resp = append(resp, *data)
		// }

		resp = append(resp, *data)
	}
	totalPages := util.TotalPages(count, limit)
	model := model.Pagination{
		Data:         resp,
		Limit:        limit,
		NextPage:     util.NextPage(page, totalPages),
		Offset:       offset,
		Page:         page,
		PrevPage:     util.PrevPage(page),
		TotalPages:   totalPages,
		TotalRecords: count,
	}
	return model
}

// GetSyncTrx ..
func (r *TrxService) GetSyncTrx() *[]entity.TrxSync {
	return r.TrxRepo.GetSyncTrx()
}

// GetOnlineOrder ..
func (r *TrxService) GetOnlineOrder(merchantUsers entity.MerchantUsers) *model.ResTrxOrderList {
	orderList := new(model.ResTrxOrderList)
	orders := make([]model.TrxOrder, 0)
	productSales := make([]model.TrxOrderProductSales, 0)
	paymentSales := make([]model.TrxOrderPaymentSales, 0)
	orderBills := make([]model.TrxOrderBill, 0)
	orderTaxSales := make([]model.TrxOrderTaxSales, 0)

	trxOrder := r.TrxOrderRepo.FindOpenByMerchantUsersID(merchantUsers.ID)
	if len(*trxOrder) > 0 {
		for _, orderData := range *trxOrder {
			// init order
			data := new(model.TrxOrder)

			trx := r.TrxRepo.FindByID(orderData.TrxID)
			user := r.UsersRepo.FindEndUserByID(trx.UsersID)
			if user != nil {
				data.Phone = user.Phone
			}

			merchantTax := r.MerchantTaxRepo.GetTax(merchantUsers.MerchantID)

			copier.Copy(&data, orderData)
			data.BusinessDate = util.DateTimeToMilliSeconds(orderData.BusinessDate)
			data.CreatedAt = util.DateTimeToMilliSeconds(orderData.CreatedAt)
			data.UpdatedAt = util.DateTimeToMilliSeconds(*orderData.UpdatedAt)
			data.DateCreated = util.DateTimeToMilliSeconds(orderData.CreatedAt.Add(time.Duration(-7) * time.Hour))
			orders = append(orders, *data)

			trxOrderBill := r.TrxOrderBillRepo.FindByTrxOrderID(orderData.ID)
			if len(*trxOrderBill) > 0 {
				for _, billData := range *trxOrderBill {
					// init order bill
					data := new(model.TrxOrderBill)
					copier.Copy(&data, billData)
					data.BusinessDate = util.DateTimeToMilliSeconds(billData.BusinessDate)
					data.CreatedAt = util.DateTimeToMilliSeconds(billData.CreatedAt)
					data.UpdatedAt = util.DateTimeToMilliSeconds(*billData.UpdatedAt)
					orderBills = append(orderBills, *data)

					qryProductSales := entity.TrxOrderProductSales{TrxOrderBillID: billData.ID}
					trxProductSales := r.TrxOrderProductSalesRepo.FindAll(&qryProductSales)
					if len(*trxProductSales) > 0 {
						for _, productData := range *trxProductSales {
							// init order product sales
							data := new(model.TrxOrderProductSales)
							copier.Copy(&data, productData)
							data.BusinessDate = util.DateTimeToMilliSeconds(productData.BusinessDate)
							data.CreatedAt = util.DateTimeToMilliSeconds(productData.CreatedAt)
							data.UpdatedAt = util.DateTimeToMilliSeconds(*productData.UpdatedAt)
							productSales = append(productSales, *data)
						}
					}

					qryPaymentSales := entity.TrxOrderPaymentSales{TrxOrderBillID: billData.ID}
					trxPaymentSales := r.TrxOrderPaymentRepo.FindAll(&qryPaymentSales)
					if len(*trxPaymentSales) > 0 {
						for _, paymentData := range *trxPaymentSales {
							// init order payment sales
							data := new(model.TrxOrderPaymentSales)
							copier.Copy(&data, paymentData)
							data.CreatedAt = util.DateTimeToMilliSeconds(paymentData.CreatedAt)
							data.UpdatedAt = util.DateTimeToMilliSeconds(*paymentData.UpdatedAt)
							paymentSales = append(paymentSales, *data)
						}
					}

					qryTaxSales := entity.TrxOrderTaxSales{TrxOrderBillID: billData.ID}
					trxTaxSales := r.TrxOrderTaxSalesRepo.FindAll(&qryTaxSales)
					if len(*trxTaxSales) > 0 {
						for _, taxData := range *trxTaxSales {
							// init order tax sales
							data := new(model.TrxOrderTaxSales)
							copier.Copy(&data, taxData)
							data.CreatedAt = util.DateTimeToMilliSeconds(taxData.CreatedAt)
							data.UpdatedAt = util.DateTimeToMilliSeconds(*taxData.UpdatedAt)
							data.IsActive = *merchantTax.IsActive
							orderTaxSales = append(orderTaxSales, *data)
						}
					}
				}
			}
		}
		orderList.Order = orders
		orderList.OrderBills = orderBills
		orderList.PaymentSales = paymentSales
		orderList.ProductSales = productSales
		orderList.TaxSales = orderTaxSales

	}
	return orderList
}

// GetOneTrxOrderByTrxID get trx_order for fcm
func (r *TrxService) GetOneTrxOrderByTrxID(id string) *model.TrxOrderMerchant {
	return r.TrxOrderRepo.FindByTrxID(id)
}

// GetMerchantIDByTrxID get trx_visit for fcm
func (r *TrxService) GetMerchantIDByTrxID(id string) int64 {
	return r.TrxVisitRepo.GetMerchantIDByTrxID(id)
}

// ListBookingTrxVisitSale ..
func (r *TrxService) ListBookingTrxVisitSale(merchantID int64, limit, page int, sort []string, filter string) model.Pagination {
	return r.TrxVisitSalesRepo.ListBookingTrxVisitSale(merchantID, limit, page, sort, filter)
}

// ListBookingTrxVisitSalesByID ..
func (r *TrxService) ListBookingTrxVisitSalesByID(ID int64, merchantID int64) *model.ResHomeVisitBookingDetailTimeNew {
	return r.TrxVisitSalesRepo.ListBookingTrxVisitSalesByID(ID, merchantID)
}

// ClosedOnlineOrderByTrxID by POS
func (r *TrxService) ClosedOnlineOrderByTrxID(trxID string) *entity.TrxOrder {
	return r.TrxOrderRepo.UpdateByTrxID(&entity.TrxOrder{IsClose: true}, trxID)
}

// GetTrxByID ..
func (r *TrxService) GetTrxByID(id string) *entity.Trx {
	return r.TrxRepo.FindByID(id)
}

// GetOrderHistoryDetailByUsersID method to get order detail history end user
func (r *TrxService) getOrderHistoryDetail(id, types string, trxOrder *model.TrxOrderMerchant) *model.DetailOrderHis {
	detail := new(model.DetailOrderHis)
	isActive := false
	if trxOrder != nil {
		merchantUser := r.MerchantUsersRepo.GetOne(trxOrder.MerchantUsersID)
		if merchantUser.ID > 0 {
			merchantTax := r.MerchantTaxRepo.GetTax(merchantUser.MerchantID)
			isActive = *merchantTax.IsActive
		}
	}

	if types == util.TrxHomeVisit {
		trxVisitSales := r.TrxVisitSalesRepo.FindByTrxID(id)
		orderDetail := make([]model.OrderDetail, 0)
		var paymentName string
		var grandTotal float64
		if len(*trxVisitSales) > 0 {
			for _, v := range *trxVisitSales {
				if paymentName == "" {
					paymentName = v.PaymentName
					grandTotal = v.GrandTotal
				}
				merchantTax2 := r.MerchantTaxRepo.GetTax(v.MerchantID)
				isActive = *merchantTax2.IsActive
				sales := r.VisitSalesRepo.GetByID(v.HomevisitSalesID)
				menus := r.VisitSalesRepo.GetMenuByTrxVisitSalesID(v.TrxHomevisitSalesID)
				item := util.VisitItemHistory(sales.StartDate, sales.EndDate)
				data := new(model.OrderDetail)
				data.Name = item
				data.Menus = menus
				orderDetail = append(orderDetail, *data)
			}
			detail.OrderDetails = orderDetail
			detail.PaymentDetails.Tax = 0
			detail.PaymentDetails.SubTotal = grandTotal
			detail.PaymentDetails.Total = grandTotal
			detail.PaymentDetails.IsActive = isActive
			detail.PaymentName = paymentName
		}
		return detail
	}
	orderBill := r.TrxOrderBillRepo.FindByTrxOrderID(trxOrder.ID)
	if len(*orderBill) > 0 {
		orderDetail := make([]model.OrderDetail, 0)
		var totalTax float64
		var paymentName string
		for _, orderBillData := range *orderBill {
			if paymentName == "" {
				paymentName = *r.TrxOrderPaymentRepo.FindPaymentName(orderBillData.ID)
			}
			qry := entity.TrxOrderProductSales{TrxOrderBillID: orderBillData.ID}
			productSales := r.TrxOrderProductSalesRepo.FindAll(&qry)
			qryTaxSales := &entity.TrxOrderTaxSales{TrxOrderBillID: orderBillData.ID}
			taxSales := r.TrxOrderTaxSalesRepo.Find(qryTaxSales)
			if len(*productSales) > 0 {
				for _, productSalesData := range *productSales {
					data := new(model.OrderDetail)
					data.Name = productSalesData.Name
					data.Qty = productSalesData.Qty
					orderDetail = append(orderDetail, *data)
				}
			}
			totalTax = totalTax + orderBillData.TotalTax
			detail.OrderDetails = orderDetail
			detail.PaymentDetails.Tax = totalTax
			detail.PaymentDetails.TaxName = taxSales.Name
			detail.PaymentDetails.TaxType = taxSales.Types
			detail.PaymentDetails.SubTotal = orderBillData.SubTotal
			if isActive && taxSales.Types == 0 {
				detail.PaymentDetails.Total = trxOrder.GrandTotal + totalTax
			} else {
				detail.PaymentDetails.Total = trxOrder.GrandTotal
			}
			detail.PaymentDetails.IsActive = isActive
			detail.PaymentName = paymentName
		}
	}
	return detail
}

// TrxReport ...
func (r *TrxService) TrxReport(merchantID int64, month string, year string) []model.ResTransactionReport {
	return r.TrxRepo.TrxReport(merchantID, month, year)
}

// TrxReportSingle ...
func (r *TrxService) TrxReportSingle(merchantID int64, month string, year string) []model.ResTransactionReportSingle {
	return r.TrxRepo.TrxReportSingle(merchantID, month, year)
}

// TrxReportPagination ...
func (r *TrxService) TrxReportPagination(merchantID int64, month string, year string, limit, page int, sort []string) model.Pagination {
	return r.TrxRepo.TrxReportPagination(merchantID, month, year, limit, page, sort)
}

// TrxReportPagination ...
func (r *TrxService) TrxReportSinglePagination(merchantID int64, month string, year string, limit, page int, sort []string) model.Pagination {
	return r.TrxRepo.TrxReportSinglePagination(merchantID, month, year, limit, page, sort)
}

// DeleteTrxByID method Hard Delete due failed create detail trx
func (r *TrxService) DeleteTrxByID(id string) {
	trx := &entity.Trx{Status: util.TrxStatusFailed}
	r.TrxRepo.Update(trx, id)
}

// VoidTrxByID method Void order
func (r *TrxService) VoidTrxByID(id string) {
	trxStatus := r.TrxRepo.FindByID(id)
	trx := &entity.Trx{Status: util.TrxStatusVoid}
	r.TrxRepo.Update(trx, id)
	if trxStatus.Status == util.TrxStatusSuccess {
		trxOrder := r.TrxOrderRepo.FindByTrxID(trxStatus.ID)
		trxOrderBill := r.TrxOrderBillRepo.FindByTrxOrderID(trxOrder.ID)
		if len(*trxOrderBill) > 0 {
			for _, billData := range *trxOrderBill {
				qryProductSales := entity.TrxOrderProductSales{TrxOrderBillID: billData.ID}
				trxProductSales := r.TrxOrderProductSalesRepo.FindAll(&qryProductSales)
				if len(*trxProductSales) > 0 {
					for _, trxOrderProductSales := range *trxProductSales {
						r.MerchantMenusRepo.UpdateStock(-trxOrderProductSales.Qty, trxOrderProductSales.MerchantMenuID)
					}
				}
			}
		}
	}
}

// CheckTrxStatusPendingByTrxHomeVisitSalesID to check trx status by trx home visit sales id
func (r *TrxService) CheckTrxStatusPendingByTrxHomeVisitSalesID(id int64) int {
	return r.TrxVisitRepo.CheckTrxStatusPendingByTrxHomeVisitSalesID(id)
}
