package repository

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
	"streetbox.id/app/trx"
	"streetbox.id/entity"
	"streetbox.id/model"
	"streetbox.id/util"
)

// TrxRepo ..
type TrxRepo struct {
	DB *gorm.DB
}

// New ..
func New(db *gorm.DB) trx.RepoInterface {
	return &TrxRepo{db}
}

// CountTrx ..
func (r *TrxRepo) CountTrx() int64 {
	var count int64
	r.DB.Find(new(entity.Trx)).Count(&count)
	return count
}

// Create ..
func (r *TrxRepo) Create(data *entity.Trx) (*gorm.DB, error) {
	trx := r.DB.Begin()
	if err := trx.Create(&data).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
		trx.Rollback()
		return nil, err
	}
	return trx, nil
}

// UpdateStatusSyncTrx ...
func (r *TrxRepo) UpdateStatusSyncTrx(UniqueID string, merchantID int64, status int, db *gorm.DB) (*gorm.DB, error) {
	trxSync := new(entity.TrxSync)
	if err := db.Where("unique_id = ? AND merchant_id = ?", UniqueID, merchantID).Find(&trxSync).Update("status", status).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
		db.Rollback()
		return nil, err
	}
	return db, nil
}

func orderPos(data *[]entity.TrxOrder) ([]string, []int64, []model.TrxOrder) {
	var orderPosIDs []int64
	var orderTrxIDs []string
	var orderPosNewSlice []model.TrxOrder
	var orderPosNew model.TrxOrder

	for _, value := range *data {
		copier.Copy(&orderPosNew, value)
		orderPosNew.BusinessDate = util.DateTimeToMilliSeconds(value.BusinessDate)
		orderPosNew.CreatedAt = util.DateTimeToMilliSeconds(value.CreatedAt)
		if value.UpdatedAt != nil {
			orderPosNew.UpdatedAt = util.DateTimeToMilliSeconds(*value.UpdatedAt)
		} else {
			orderPosNew.UpdatedAt = 0
		}
		orderPosNew.DateCreated = util.DateTimeToMilliSeconds(value.CreatedAt.Add(time.Duration(-7) * time.Hour))

		orderPosNewSlice = append(orderPosNewSlice, orderPosNew)
		orderPosIDs = append(orderPosIDs, value.ID)
		orderTrxIDs = append(orderTrxIDs, value.TrxID)
	}

	return orderTrxIDs, orderPosIDs, orderPosNewSlice
}

func orderBill(data *[]entity.TrxOrderBill) ([]int64, []model.TrxOrderBill) {
	var orderBillNewSlice []model.TrxOrderBill
	var orderBillNew model.TrxOrderBill
	var orderBillIDs []int64

	for _, value := range *data {
		copier.Copy(&orderBillNew, value)
		orderBillNew.BusinessDate = util.DateTimeToMilliSeconds(value.BusinessDate)
		orderBillNew.CreatedAt = util.DateTimeToMilliSeconds(value.CreatedAt)
		if value.UpdatedAt != nil {
			orderBillNew.UpdatedAt = util.DateTimeToMilliSeconds(*value.UpdatedAt)
		} else {
			orderBillNew.UpdatedAt = 0
		}

		orderBillNewSlice = append(orderBillNewSlice, orderBillNew)
		orderBillIDs = append(orderBillIDs, value.ID)
	}
	return orderBillIDs, orderBillNewSlice
}

// GetOrderTrx ..
func (r *TrxRepo) GetOrderTrx(merchantUsersID int64, startDate, endDate, keyword string) *model.ResTrxOrderList {
	results := new(model.ResTrxOrderList)
	orderPosOld := make([]entity.TrxOrder, 1)
	merchantUsers := new(entity.MerchantUsers)
	merchantTax := new(entity.MerchantTax)
	merchantType := 0
	merchantTaxIsActive := false
	if startDate == "" {
		now := time.Now()
		filterStartDates := now.AddDate(0, -1, 0)
		startDate = filterStartDates.Format("2006-01-02")
	}
	fmt.Println("endDate : " + endDate)
	if keyword == "" {
		if endDate == "" {
			r.DB.Where("merchant_users_id = ? AND business_date >= ?", merchantUsersID, startDate).Find(&orderPosOld)
		} else {
			r.DB.Where("merchant_users_id = ? AND business_date >= ? AND business_date <= ?", merchantUsersID, startDate, endDate).Find(&orderPosOld)
		}
	} else {
		if endDate == "" {
			r.DB.Where("merchant_users_id = ? AND business_date >= ? AND (trx_id like ? OR order_no like ?)", merchantUsersID, startDate, "%"+keyword+"%", "%"+keyword+"%").Find(&orderPosOld)
		} else {
			r.DB.Where("merchant_users_id = ? AND business_date >= ? AND business_date <= ? AND (trx_id like ? OR order_no like ?)",
				merchantUsersID, startDate, endDate, "%"+keyword+"%", "%"+keyword+"%").Find(&orderPosOld)
		}
	}
	r.DB.Where("id = ?", merchantUsersID).Select("*").Table("merchant_users").Find(&merchantUsers)
	if err := r.DB.Find(&merchantTax, "merchant_id = ? AND is_active = ?", merchantUsers.MerchantID, true).Error; err != nil {
		if err := r.DB.Find(&merchantTax, "merchant_id = ?", merchantUsers.MerchantID).Order(merchantTax.UpdatedAt).Limit(1).Error; err != nil {
			merchantType = 0
			merchantTaxIsActive = false
		} else {
			merchantType = *merchantTax.Type
			merchantTaxIsActive = *merchantTax.IsActive
		}
	} else {
		merchantType = *merchantTax.Type
		merchantTaxIsActive = *merchantTax.IsActive
	}

	if len(orderPosOld) > 0 {
		orderTrxIDs, orderPosIDs, orderPosNewSlice := orderPos(&orderPosOld)
		results.Order = orderPosNewSlice
		orderBillOld := make([]entity.TrxOrderBill, 1)
		trxOld := make([]model.Trx, 1)
		paymentSalesNews := make([]model.TrxOrderPaymentSales, 0)
		paymentSalesNew := new(model.TrxOrderPaymentSales)
		paymentSalesOld := make([]model.ResTrxOrderPaymentSales, 1)
		taxSalesOld := make([]model.ResTrxOrderTaxSales, 1)
		taxSalesNews := make([]model.TrxOrderTaxSales, 0)
		taxSalesNew := new(model.TrxOrderTaxSales)
		var (
			orderBillIDs      []int64
			orderBillNewSlice []model.TrxOrderBill
		)
		print("orderTrxIDs")
		print(orderTrxIDs)
		if len(orderTrxIDs) > 0 {
			r.DB.Where("id IN (?)", orderTrxIDs).Select("*, id AS trx_id").Table("trx").Find(&trxOld)
			results.Trx = trxOld
		}
		if len(orderPosIDs) > 0 {
			r.DB.Where("trx_order_id IN (?)", orderPosIDs).Find(&orderBillOld)
			orderBillIDs, orderBillNewSlice = orderBill(&orderBillOld)
			results.OrderBills = orderBillNewSlice
			r.DB.Where("trx_order_bill_id IN (?)", orderBillIDs).
				Joins("LEFT JOIN trx_order_bill tob ON tops.trx_order_bill_id = tob.id").
				Select("tops.id, tops.unique_id, tops.order_bill_unique_id, tops.name, tops.amount, tops.created_at, " +
					"tops.updated_at, tops.payment_method_id, tob.order_unique_id").Table("trx_order_payment_sales tops").Scan(&paymentSalesOld)
			for _, value := range paymentSalesOld {
				copier.Copy(paymentSalesNew, &value)
				paymentSalesNew.CreatedAt = util.DateTimeToMilliSeconds(value.CreatedAt)
				if value.UpdatedAt != nil {
					paymentSalesNew.UpdatedAt = util.DateTimeToMilliSeconds(*value.UpdatedAt)
				} else {
					paymentSalesNew.UpdatedAt = 0
				}
				paymentSalesNews = append(paymentSalesNews, *paymentSalesNew)
			}
			results.PaymentSales = paymentSalesNews

			r.DB.Where("trx_order_bill_id IN (?)", orderBillIDs).
				Joins("LEFT JOIN trx_order_bill tob ON tots.trx_order_bill_id = tob.id").
				Select("tots.id, tots.unique_id, tots.name, tots.amount, tots.types as type, " +
					"tots.created_at, tots.updated_at, tots.merchant_tax_id, tots.order_bill_unique_id" +
					", tob.order_unique_id").Table("trx_order_tax_sales tots").Scan(&taxSalesOld)
			for _, value := range taxSalesOld {
				copier.Copy(taxSalesNew, &value)
				taxSalesNew.CreatedAt = util.DateTimeToMilliSeconds(value.CreatedAt)
				if value.UpdatedAt != nil {
					taxSalesNew.UpdatedAt = util.DateTimeToMilliSeconds(*value.UpdatedAt)
				} else {
					taxSalesNew.UpdatedAt = 0
				}
				taxSalesNew.Type = merchantType
				taxSalesNew.IsActive = merchantTaxIsActive
				taxSalesNews = append(taxSalesNews, *taxSalesNew)
			}
			results.TaxSales = taxSalesNews
		}
		productSalesOld := make([]model.ResTrxOrderProductSales, 1)
		productSalesNew := new(model.TrxOrderProductSales)
		productSalesNews := make([]model.TrxOrderProductSales, 0)
		if len(orderBillIDs) > 0 {
			r.DB.Select("tops.order_bill_unique_id, tops.unique_id, tops.merchant_menu_id as product_id, tops.name "+
				", tops.price, tops.qty, tops.notes, "+
				"tob.order_unique_id, t.qr_code").Table("trx_order_product_sales tops").Group("tops.order_bill_unique_id, tops.unique_id, tops.merchant_menu_id , tops.name "+
				", tops.price, tops.qty, tops.notes, "+
				"tob.order_unique_id, t.qr_code	").
				Where("trx_order_bill_id IN (?)", orderBillIDs).
				Joins("JOIN trx_order_bill tob ON tops.trx_order_bill_id = tob.id").
				Joins("LEFT JOIN trx_order tor ON tor.bill_no = tob.bill_no").
				Joins("LEFT JOIN trx t ON t.id = tor.trx_id").
				Scan(&productSalesOld)
			for _, value := range productSalesOld {
				copier.Copy(productSalesNew, &value)
				// productSalesNew.CreatedAt = util.DateTimeToMilliSeconds(value.CreatedAt)
				// productSalesNew.BusinessDate = util.DateTimeToMilliSeconds(value.BusinessDate)
				// if value.UpdatedAt != nil {
				// 	productSalesNew.UpdatedAt = util.DateTimeToMilliSeconds(*value.UpdatedAt)
				// } else {
				// 	productSalesNew.UpdatedAt = 0
				// }

				productSalesNews = append(productSalesNews, *productSalesNew)
			}
			results.ProductSales = productSalesNews
		}
	}

	return results
}

// CreateSyncTrx ..
func (r *TrxRepo) CreateSyncTrx(data *entity.TrxSync) {
	if err := r.DB.Create(&data).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
	}
}

// FindByID ...
func (r *TrxRepo) FindByID(id string) *entity.Trx {
	data := new(entity.Trx)
	r.DB.Find(&data, "id = ?", id)
	return data
}

// FindByUsersID ...
func (r *TrxRepo) FindByUsersID(limit, page int, usersID int64, filter string) (*[]entity.Trx, int, int) {
	if filter == "" {
		data := new([]entity.Trx)
		count := 0
		offset := util.Offset(page, limit)
		fmt.Println(usersID)
		r.DB.Raw("SELECT t.* FROM trx t LEFT JOIN trx_order tro ON tro.trx_id = t.id WHERE t.users_id = ? AND t.types = 'ORDER' GROUP BY t.id HAVING COUNT(tro.id) > 0 union SELECT t.* FROM trx t WHERE t.users_id = ? AND t.types = 'VISIT' ORDER BY created_at DESC, status DESC LIMIT ? OFFSET ?", usersID, usersID, limit, offset).Scan(&data)
		r.DB.Raw("SELECT t.* FROM trx t LEFT JOIN trx_order tro ON tro.trx_id = t.id WHERE t.users_id = ? AND t.types = 'ORDER' GROUP BY t.id HAVING COUNT(tro.id) > 0 union SELECT t.* FROM trx t WHERE t.users_id = ? AND t.types = 'VISIT' ORDER BY created_at DESC, status DESC", usersID, usersID).Scan(&count)
		return data, count, offset
	} else {
		return r.FindByfilter(limit, page, usersID, filter)
	}
}

// FindByfilter ...
func (r *TrxRepo) FindByfilter(limit, page int, usersID int64, filter string) (*[]entity.Trx, int, int) {
	data := new([]entity.Trx)
	count := 0
	offset := util.Offset(page, limit)
	r.DB.Raw("SELECT t.* FROM trx t LEFT JOIN trx_order tro ON tro.trx_id = t.id WHERE t.users_id = ? AND t.types = 'ORDER' AND t.status "+filter+" GROUP BY t.id HAVING COUNT(tro.id) > 0 union SELECT t.* FROM trx t WHERE t.users_id = ? AND t.types = 'VISIT' AND t.status "+filter+" ORDER BY status DESC, created_at DESC LIMIT ? OFFSET ?", usersID, usersID, limit, offset).Scan(&data)
	r.DB.Raw("SELECT t.* FROM trx t LEFT JOIN trx_order tro ON tro.trx_id = t.id WHERE t.users_id = ? AND t.types = 'ORDER' AND t.status "+filter+" GROUP BY t.id HAVING COUNT(tro.id) > 0 union SELECT t.* FROM trx t WHERE t.users_id = ? AND t.types = 'VISIT' AND t.status "+filter+" ORDER BY status DESC, created_at DESC", usersID, usersID).Scan(&count)
	return data, count, offset
}

// GetSyncTrx ..
func (r *TrxRepo) GetSyncTrx() *[]entity.TrxSync {
	data := new([]entity.TrxSync)
	r.DB.Where("status = ?", 0).Find(&data)
	return data
}

// Update ..
func (r *TrxRepo) Update(data *entity.Trx, id string) {
	if err := r.DB.Model(&entity.Trx{ID: id}).Updates(&data).
		Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
	}
	log.Printf("INFO: Transaction Updated: %+v", data)
}

// UpdateByTrxID ...
func (r *TrxRepo) UpdateByTrxID(data *entity.Trx, id string) (*gorm.DB, error) {
	trx := r.DB.Begin()
	if err := trx.Model(&entity.Trx{ID: id}).Updates(&data).Error; err != nil {
		trx.Rollback()
		log.Printf("ERROR: %s", err.Error())
	}
	return trx, nil
}

// TrxReport ...
func (r *TrxRepo) TrxReport(merchantID int64, month, year string) []model.ResTransactionReport {
	data := make([]model.ResTransactionReport, 1)
	r.DB.Raw("select * from report_transaction(?, ?, ?)", merchantID, month, year).Scan(&data)
	// for i := 0; i < len(data); i++ {
	// 	if data[i].ExternalID == "" {
	// 		data[i].ExternalID = "Cash"
	// 	} else {
	// 		data[i].ExternalID = "QR"
	// 	}
	// }
	return data
}

// TrxReportSingle ...
func (r *TrxRepo) TrxReportSingle(merchantID int64, month, year string) []model.ResTransactionReportSingle {
	data := make([]model.ResTransactionReportSingle, 1)
	r.DB.Raw("select * from report_transaction_single(?, ?, ?)", merchantID, month, year).Scan(&data)
	// for i := 0; i < len(data); i++ {
	// 	if data[i].ExternalID == "" {
	// 		data[i].ExternalID = "Cash"
	// 	} else {
	// 		data[i].ExternalID = "QR"
	// 	}
	// }

	return data
}

// TrxReportPagination ...
func (r *TrxRepo) TrxReportPagination(merchantID int64, month, year string, limit, page int, sort []string) model.Pagination {
	data := make([]model.ResTransactionReport, 1)
	count := 0
	offset := util.Offset(page, limit)
	sortString := ""
	if len(sort) > 0 {
		sortString = "order by "
		for _, o := range sort {
			splitString := strings.Split(o, ",")
			sortString += splitString[0]
		}
	}

	// month = "0"
	// year = "0"''

	if month == "0" && year == "0" {
		r.DB.Raw("select TO_CHAR(t.created_at, 'yyyy-mm-dd') dates, TO_CHAR(t.created_at, 'HH24:MI:SS') as times, tops.name as product_name, muu.name as user_name, tops.qty, tob.total_tax, tops.price * tops.qty as grand_total, t.status, tor.trx_id, tor.id, tor.order_no, tor.bill_no, tor.type_payment::text as external_id, t.types::text, t.created_at as transaction_date, tv.customer_name, hs.start_date, hs.end_date FROM trx t left join trx_order tor on t.id = tor.trx_id left join trx_order_bill tob on tor.id = tob.trx_order_id left join trx_order_product_sales tops on tob.id = tops.trx_order_bill_id left join (select mu.id, name, merchant_id from users u left join merchant_users mu on u.id = mu.users_id ) muu on tor.merchant_users_id  = muu.id LEFT JOIN trx_visit tv  on t.id  = tv.trx_id LEFT JOIN trx_homevisit_sales ths on tv.id = ths.trx_visit_id LEFT JOIN homevisit_sales hs ON ths.homevisit_sales_id = hs.id WHERE muu.merchant_id = ? order by dates, times asc;"+sortString+" limit ? offset ?", merchantID, limit, offset).Scan(&data)
	} else if month == "0" {
		r.DB.Raw("select TO_CHAR(t.created_at, 'yyyy-mm-dd') dates, TO_CHAR(t.created_at, 'HH24:MI:SS') as times, tops.name as product_name, muu.name as user_name, tops.qty, tob.total_tax, tops.price * tops.qty as grand_total, t.status, tor.trx_id, tor.id, tor.order_no, tor.bill_no, tor.type_payment::text as external_id, t.types::text, t.created_at as transaction_date, tv.customer_name, hs.start_date, hs.end_date FROM trx t left join trx_order tor on t.id = tor.trx_id left join trx_order_bill tob on tor.id = tob.trx_order_id left join trx_order_product_sales tops on tob.id = tops.trx_order_bill_id left join ( select mu.id, name, merchant_id from users u left join merchant_users mu on u.id = mu.users_id ) muu on tor.merchant_users_id  = muu.id LEFT JOIN trx_visit tv  on t.id  = tv.trx_id LEFT JOIN trx_homevisit_sales ths on tv.id = ths.trx_visit_id LEFT JOIN homevisit_sales hs ON ths.homevisit_sales_id = hs.id WHERE muu.merchant_id = ? and TO_CHAR(t.created_at, 'yyyy') = ? order by dates, times asc"+sortString+" limit ? offset ?", merchantID, year, limit, offset).Scan(&data)
	} else if year == "0" {
		r.DB.Raw("select TO_CHAR(t.created_at, 'yyyy-mm-dd') dates, TO_CHAR(t.created_at, 'HH24:MI:SS') as times, tops.name as product_name, muu.name as user_name, tops.qty, tob.total_tax, tops.price * tops.qty as grand_total, t.status, tor.trx_id, tor.id, tor.order_no, tor.bill_no, tor.type_payment::text as external_id, t.types::text, t.created_at as transaction_date, tv.customer_name, hs.start_date, hs.end_date FROM trx t left join trx_order tor on t.id = tor.trx_id left join trx_order_bill tob on tor.id = tob.trx_order_id left join trx_order_product_sales tops on tob.id = tops.trx_order_bill_id left join (select mu.id, name, merchant_id from users u left join merchant_users mu on u.id = mu.users_id) muu on tor.merchant_users_id  = muu.id LEFT JOIN trx_visit tv  on t.id  = tv.trx_id LEFT JOIN trx_homevisit_sales ths on tv.id = ths.trx_visit_id LEFT JOIN homevisit_sales hs ON ths.homevisit_sales_id = hs.id WHERE muu.merchant_id = ? and TO_CHAR(t.created_at, 'mm') = ? order by dates, times asc"+sortString+" limit ? offset ?", merchantID, month, limit, offset).Scan(&data)
	} else {
		r.DB.Raw("select TO_CHAR(t.created_at, 'yyyy-mm-dd') dates, TO_CHAR(t.created_at, 'HH24:MI:SS') as times, tops.name as product_name, muu.name as user_name, tops.qty, tob.total_tax, tops.price * tops.qty as grand_total, t.status, tor.trx_id, tor.id, tor.order_no, tor.bill_no, tor.type_payment::text as external_id, t.types::text, t.created_at as transaction_date, tv.customer_name, hs.start_date, hs.end_date FROM trx t left join trx_order tor on t.id = tor.trx_id left join trx_order_bill tob on tor.id = tob.trx_order_id left join trx_order_product_sales tops on tob.id = tops.trx_order_bill_id left join (select mu.id, name, merchant_id from users u left join merchant_users mu on u.id = mu.users_id) muu on tor.merchant_users_id  = muu.id LEFT JOIN trx_visit tv  on t.id  = tv.trx_id  LEFT JOIN trx_homevisit_sales ths on tv.id = ths.trx_visit_id LEFT JOIN homevisit_sales hs ON ths.homevisit_sales_id = hs.id WHERE muu.merchant_id = ? and TO_CHAR(t.created_at, 'mm') = ? and TO_CHAR(t.created_at, 'yyyy') = ? order by dates, times asc"+sortString+" limit ? offset ?", merchantID, month, year, limit, offset).Scan(&data)
	}

	//r.DB.Raw("select * from r_trx(?, ?, ?) "+sortString+" limit ? offset ?", merchantID, month, year, limit, offset).Scan(&data)
	r.DB.Raw("select count(*) from report_transaction(?, ?, ?)", merchantID, month, year).Count(&count)

	// merchantUsers := new(entity.MerchantUsers)
	// merchantTax := new(entity.MerchantTax)
	// merchantTaxIsActive := false
	// if err := r.DB.Find(&merchantTax, "merchant_id = ? AND is_active = ?", merchantUsers.MerchantID, true).Error; err != nil {
	// 	if err := r.DB.Find(&merchantTax, "merchant_id = ?", merchantUsers.MerchantID).Order(merchantTax.UpdatedAt).Limit(1).Error; err != nil {
	// 		merchantTaxIsActive = false
	// 	} else {
	// 		merchantTaxIsActive = *merchantTax.IsActive
	// 	}
	// } else {
	// 	merchantTaxIsActive = *merchantTax.IsActive
	// }
	// if !merchantTaxIsActive {
	// 	for i := 0; i < len(data); i++ {
	// 		data[i].TotalTax = 0
	// 	}
	// }
	// for i := 0; i < len(data); i++ {
	// 	if data[i].ExternalID == "" {
	// 		data[i].ExternalID = "Cash"
	// 	} else {
	// 		data[i].ExternalID = "QR"
	// 	}
	// }
	totalPages := util.TotalPages(count, limit)

	model := model.Pagination{
		Data:         data,
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

// func unique(arr []model.ResTrxOrderProductSales) []model.ResTrxOrderProductSales {
// 	occurred := map[model.ResTrxOrderProductSales]bool{}
// 	result := []model.ResTrxOrderProductSales{}
// 	for e := range arr {

// 		// check if already the mapped
// 		// variable is set to true or not
// 		if occurred[arr[e]] != true {
// 			occurred[arr[e]] = true

// 			// Append to result slice.
// 			result = append(result, arr[e])
// 		}
// 	}

// 	return result
// }

// TrxReportSinglePagination ...
func (r *TrxRepo) TrxReportSinglePagination(merchantID int64, month, year string, limit, page int, sort []string) model.Pagination {
	data := make([]model.ResTransactionReportSingle, 1)
	count := 0
	offset := util.Offset(page, limit)
	sortString := ""
	if len(sort) > 0 {
		sortString = "order by "
		for _, o := range sort {
			splitString := strings.Split(o, ",")
			sortString += splitString[0]
		}
	}
	r.DB.Raw("select * from report_transaction_single(?, ?, ?) "+sortString+" limit ? offset ?", merchantID, month, year, limit, offset).Scan(&data)
	r.DB.Raw("select count(*) from report_transaction_single(?, ?, ?)", merchantID, month, year).Count(&count)

	// for i := 0; i < len(data); i++ {
	// 	if data[i].ExternalID == "" {
	// 		data[i].ExternalID = "Cash"
	// 	} else {
	// 		data[i].ExternalID = "QR"
	// 	}
	// }

	// merchantUsers := new(entity.MerchantUsers)
	// merchantTax := new(entity.MerchantTax)
	// merchantTaxIsActive := false
	// if err := r.DB.Find(&merchantTax, "merchant_id = ? AND is_active = ?", merchantUsers.MerchantID, true).Error; err != nil {
	// 	if err := r.DB.Find(&merchantTax, "merchant_id = ?", merchantUsers.MerchantID).Order(merchantTax.UpdatedAt).Limit(1).Error; err != nil {
	// 		merchantTaxIsActive = false
	// 	} else {
	// 		merchantTaxIsActive = *merchantTax.IsActive
	// 	}
	// } else {
	// 	merchantTaxIsActive = *merchantTax.IsActive
	// }
	// if !merchantTaxIsActive {
	// 	for i := 0; i < len(data); i++ {
	// 		data[i].TotalTax = 0
	// 	}
	// }
	totalPages := util.TotalPages(count, limit)
	model := model.Pagination{
		Data:         data,
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

// HardDeleteByID method due failed create detail trx
func (r *TrxRepo) HardDeleteByID(id string) {
	if err := r.DB.Exec("delete from trx where id = ?", id).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
	}
	log.Printf("INFO: HardDeleteByID: trxID -> %s", id)
}
