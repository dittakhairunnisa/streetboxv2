package repository

import (
	"log"
	"strings"
	"time"

	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
	"streetbox.id/app/trxvisitsales"
	"streetbox.id/entity"
	"streetbox.id/model"
	"streetbox.id/util"
)

// TrxVisitSalesRepo ..
type TrxVisitSalesRepo struct {
	DB *gorm.DB
}

// New ..
func New(db *gorm.DB) trxvisitsales.RepoInterface {
	return &TrxVisitSalesRepo{db}
}

// Create ..
func (r *TrxVisitSalesRepo) Create(db *gorm.DB, data *entity.TrxHomevisitSales) error {
	if err := db.Create(&data).Error; err != nil {
		db.Rollback()
		log.Printf("ERROR: %s", err.Error())
		return err
	}
	log.Printf("INFO: Create TrxVisitSales: %+v", data)
	return nil
}

// FindByTrxID ..
func (r *TrxVisitSalesRepo) FindByTrxID(id string) *[]model.TrxHomevisitSales {
	data := new([]model.TrxHomevisitSales)
	r.DB.Select("tv.*, thv.id as trx_homevisit_sales_id, m.name as merchant_name ,m.logo as merchant_logo, "+
		"hv.id as homevisit_sales_id, hv.deposit, hv.merchant_id,hv.available, p.name as payment_name").Joins("JOIN "+
		"trx_homevisit_sales thv on tv.id = thv.trx_visit_id").Joins("JOIN "+
		"homevisit_sales hv on thv.homevisit_sales_id = hv.id").Joins("JOIN "+
		"merchant m on hv.merchant_id = m.id").Joins("JOIN "+
		"payment_method p on tv.payment_method_id = p.id").
		Table("trx_visit tv").Where("tv.deleted_at is null and "+
		"hv.deleted_at is null and tv.trx_id = ?", id).Scan(&data)
	return data
}

// FindByMerchantID ..
func (r *TrxVisitSalesRepo) FindByMerchantID(id int64) *[]model.HomeVisitSales {
	data := new([]model.HomeVisitSales)
	r.DB.Select("tv.*,m.name as merchant_name ,m.logo as merchant_logo, "+
		"hv.deposit, hv.start_date, hv.end_date, u.profile_picture, "+
		"thv.id as trx_visit_sales_id").Joins("JOIN "+
		"trx_visit tv on t.id = tv.trx_id").Joins("JOIN "+
		"users u on t.users_id = u.id").Joins("JOIN "+
		"trx_homevisit_sales thv on tv.id = thv.trx_visit_id").Joins("JOIN "+
		"homevisit_sales hv on thv.homevisit_sales_id = hv.id").Joins("JOIN "+
		"merchant m on hv.merchant_id = m.id").
		Table("trx t").Where("tv.deleted_at is null and thv.deleted_at is null and "+
		"hv.deleted_at is null and hv.merchant_id = ? and thv.status = ? "+
		"and hv.end_date > ? and t.status = ?",
		id, util.TrxVisitStatusOpen, time.Now(), util.TrxStatusSuccess).Scan(&data)
	return data
}

// ListBookingTrxVisitSale ..
func (r *TrxVisitSalesRepo) ListBookingTrxVisitSale(merchantID int64, limit, page int, sort []string, filter string) model.Pagination {
	data := make([]model.ResHomeVisitBookingListTime, 0)
	dataOld := make([]model.ResHomeVisitBookingList, 1)
	dataNew := new(model.ResHomeVisitBookingListTime)
	count := 0
	offset := util.Offset(page, limit)

	sortString := ""
	if len(sort) > 0 {
		sortString = " order by "
		for _, o := range sort {
			splitString := strings.Split(o, ",")
			sortString += splitString[0]
		}
	}
	filterQuery := ""
	if filter != "" {
		filterQuery = "WHERE status = ? "
	}

	if filter == "" {
		//r.DB.Raw("select * from get_all_trx_visit(?) "+sortString+" limit ? offset ?", merchantID, limit, offset).Scan(&dataOld)
		r.DB.Raw("select ths.id, tv.customer_name as customer_name, hs.start_date, hs.end_date, hs.deposit, tv.trx_id, tv.created_at as transaction_date, tv.grand_total, pm.name as payment_method, case when ((select count(*) from trx_refund_visit trv where ths.id = trv.trx_homevisit_sales_id) > 0) then 'REFUNDED' else t.status end as status from trx_visit tv JOIN trx_homevisit_sales ths on tv.id = ths.trx_visit_id JOIN homevisit_sales hs ON ths.homevisit_sales_id = hs.id JOIN trx t on tv.trx_id = t.id  JOIN payment_method pm  on tv.payment_method_id = pm.id  where hs.merchant_id = ?"+sortString+" ORDER BY tv.created_at DESC limit ? offset ?", merchantID, limit, offset).Scan(&dataOld)

	} else {
		//r.DB.Raw("select * from get_all_trx_visit(?) "+filterQuery+sortString+" limit ? offset ?", merchantID, filter, limit, offset).Scan(&dataOld)
		r.DB.Raw("select ths.id, tv.customer_name as customer_name, hs.start_date, hs.end_date, hs.deposit, tv.trx_id, tv.created_at as transaction_date, tv.grand_total, pm.name as payment_method, case when ((select count(*) from trx_refund_visit trv where ths.id = trv.trx_homevisit_sales_id) > 0) then 'REFUNDED' else t.status end as status from trx_visit tv JOIN trx_homevisit_sales ths on tv.id = ths.trx_visit_id JOIN homevisit_sales hs ON ths.homevisit_sales_id = hs.id JOIN trx t on tv.trx_id = t.id  JOIN payment_method pm  on tv.payment_method_id = pm.id  where hs.merchant_id = ?"+filterQuery+sortString+" ORDER BY tv.created_at DESC limit ? offset ?", merchantID, limit, offset).Scan(&dataOld)
	}
	for _, value := range dataOld {
		copier.Copy(&dataNew, value)

		dataNew.StartDate = value.StartDate.Format("2006-01-02 15:04:05")
		dataNew.EndDate = value.EndDate.Format("2006-01-02 15:04:05")
		data = append(data, *dataNew)
	}

	if filter == "" {
		r.DB.Raw("select count(*) from get_all_trx_visit(?)", merchantID).Count(&count)
	} else {
		r.DB.Raw("select count(*) from get_all_trx_visit(?) "+filterQuery, merchantID, filter).Count(&count)
	}

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

// ListBookingTrxVisitSalesByID ..
func (r *TrxVisitSalesRepo) ListBookingTrxVisitSalesByID(ID int64, merchantID int64) *model.ResHomeVisitBookingDetailTimeNew {
	var menus []model.ResMenu
	data := new(model.ResHomeVisitBookingDetailTime)
	dataNew := new(model.ResHomeVisitBookingDetailTimeNew)
	r.DB.Select("tv.id, tv.trx_id, tv.customer_name, start_date, end_date, tv.address, hs.deposit,tv.grand_total, hs.created_at as transaction_date, "+
		"trx.status, notes, u.phone as phone1, tv.phone as phone2, pm.name as payment_method ").
		Joins("JOIN trx_homevisit_sales ths on tv.id = ths.trx_visit_id").
		Joins("JOIN homevisit_sales hs ON ths.homevisit_sales_id = hs.id").
		Joins("JOIN trx ON tv.trx_id = trx.id").
		Joins("JOIN users u ON trx.users_id = u.id").
		Joins("JOIN payment_method pm ON tv.payment_method_id = pm.id ").
		Table("trx_visit tv").
		Where("hs.merchant_id = ? AND ths.id = ?", merchantID, ID).Scan(&data)

	r.DB.Raw("SELECT m.name, t.quantity FROM merchant_menu m, trx_homevisit_menu_sales t WHERE t.trx_homevisit_sales_id = ? AND t.menu_id = m.id", ID).Scan(&menus)
	data.Menus = menus
	copier.Copy(&dataNew, data)
	dataNew.StartDate = data.StartDate.Format("2006-01-02 15:04:05")
	dataNew.EndDate = data.EndDate.Format("2006-01-02 15:04:05")
	return dataNew
}

// GetHomeVisitData ...
func (r *TrxVisitSalesRepo) GetHomeVisitData(date string, merchantID int64) []int64 {
	var data []int64
	r.DB.Select("ths.id").Joins("JOIN homevisit_sales hs ON ths.homevisit_sales_id  = hs.id").Table("trx_homevisit_sales ths").
		Where("TO_CHAR(hs.start_date, 'yyyy-mm-dd') = ? and merchant_id = ?", date, merchantID).Pluck("ths.id", &data)
	return data
}

// UpdateByID ..
func (r *TrxVisitSalesRepo) UpdateByID(data *entity.TrxHomevisitSales, id int64) error {
	qry := entity.TrxHomevisitSales{ID: id}
	if err := r.DB.Model(qry).Updates(data).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
		return err
	}
	return nil
}
