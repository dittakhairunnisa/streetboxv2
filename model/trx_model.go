package model

import (
	"time"

	"streetbox.id/entity"
)

// ReqCreateTrxSales ...
type ReqCreateTrxSales struct {
	ParkingSpaceSalesID int64 `json:"parkingSpaceSalesId" binding:"required"`
	MerchantID          int64 `json:"merchantId" binding:"required"`
	TotalSlot           int   `json:"totalSlot" binding:"required"`
}

// ReqRefundParkingSpaceSales ...
type ReqRefundParkingSpaceSales struct {
	Amount                 int64 `json:"amount" binding:"required"`
	TrxParkingSpaceSalesID int64 `json:"trxParkingSpaceSalesId" binding:"required"`
}

// ReqRefundHomeVisit ...
type ReqRefundHomeVisit struct {
	Amount int64 `json:"amount" binding:"required"`
	ID     int64 `json:"id" binding:"required"`
}

// ReqCreateSyncTrx ..
type ReqCreateSyncTrx struct {
	UniqueID     string `json:"uniqueId" binding:"required"`
	BusinessDate int64  `json:"businessDate" binding:"required"`
	SyncDate     int64  `json:"syncDate" binding:"required"`
	Data         string `json:"data" binding:"required"`
}

// ReqTrxOrderList ..
type ReqTrxOrderList struct {
	Order        TrxOrder               `json:"order"`
	ProductSales []TrxOrderProductSales `json:"productSales"`
	PaymentSales []TrxOrderPaymentSales `json:"paymentSales"`
	OrderBills   []TrxOrderBill         `json:"orderBills"`
	TaxSales     []TrxOrderTaxSales     `json:"taxSales"`
}

// ResTrxList ...
type ResTrxList struct {
	ID           int64     `json:"id"`
	MerchantName string    `json:"merchantName"`
	SpaceName    string    `json:"spaceName"`
	StartDate    time.Time `json:"startDate"`
	EndDate      time.Time `json:"endDate"`
	TotalSlot    int       `json:"totalSlot"`
	Point        int64     `json:"point"`
	Status       string    `json:"status"`
}

// Trx ...
type Trx struct {
	ID         string    `json:"id"`
	TrxID      string    `json:"trxId"`
	Types      string    `json:"types"`
	UserID     int64     `json:"userId"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
	DeletedAt  time.Time `json:"deletedAt"`
	Status     string    `json:"status"`
	ExternalID string    `json:"externalId"`
	Address    string    `json:"address"`
	QrCode     string    `json:"qrCode"`
}

// TrxOrder ...
type TrxOrder struct {
	ID                int64   `json:"id"`
	UniqueID          string  `json:"uniqueId"`
	OrderNo           string  `json:"orderNo"`
	UserID            int64   `json:"userId"`
	BillNo            string  `json:"billNo"`
	IsClose           bool    `json:"isClose"`
	Note              string  `json:"note"`
	Types             int     `json:"types"`
	BusinessDate      int64   `json:"businessDate"`
	TotalDiscount     float64 `json:"totalDiscount"`
	GrandTotal        float64 `json:"grandTotal"`
	CreatedAt         int64   `json:"createdAt"`
	UpdatedAt         int64   `json:"updatedAt"`
	MerchantUsersID   int64   `json:"merchantUsersId"`
	TrxID             string  `json:"trxId"`
	PaymentMethodId   string  `json:"paymentMethodId"`
	PaymentMethodName string  `json:"paymentMethodName"`
	TypeOrder         string  `json:"typeOrder"`
	TypePayment       string  `json:"typePayment"`
	Phone             string  `json:"phone"`
	DateCreated       int64   `json:"dateCreated"`
	MerchantID        int64   `json:"merchantId"`
}

// TrxOrderBill ..
type TrxOrderBill struct {
	ID            int64   `json:"id"`
	UniqueID      string  `json:"uniqueId"`
	OrderUniqueID string  `json:"orderUniqueId"`
	BillNo        string  `json:"billNo"`
	IsClose       bool    `json:"isClose"`
	TotalDiscount float32 `json:"totalDiscount"`
	SubTotal      float64 `json:"subTotal"`
	TotalTax      float64 `json:"totalTax"`
	GrandTotal    float64 `json:"grandTotal"`
	BusinessDate  int64   `json:"businessDate"`
	CreatedAt     int64   `json:"createdAt"`
	UpdatedAt     int64   `json:"updatedAt"`
}

// TrxOrderPaymentSales ..
type TrxOrderPaymentSales struct {
	ID                int64   `json:"id"`
	UniqueID          string  `json:"uniqueId"`
	OrderUniqueID     string  `json:"orderUniqueId"`
	OrderBillUniqueID string  `json:"orderBillUniqueId"`
	Name              string  `json:"name"`
	Amount            float64 `json:"amount"`
	CreatedAt         int64   `json:"createdAt"`
	UpdatedAt         int64   `json:"updatedAt"`
	PaymentMethodID   int64   `json:"paymentMethodId"`
}

// ResTrxOrderPaymentSales ..
type ResTrxOrderPaymentSales struct {
	ID                int64      `json:"id"`
	UniqueID          string     `json:"uniqueId"`
	OrderUniqueID     string     `json:"orderUniqueId"`
	OrderBillUniqueID string     `json:"orderBillUniqueId"`
	Name              string     `json:"name"`
	Amount            float64    `json:"amount"`
	CreatedAt         time.Time  `json:"createdAt"`
	UpdatedAt         *time.Time `json:"updatedAt"`
	PaymentMethodID   int64      `json:"paymentMethodId"`
}

// TrxOrderProductSales ..
type TrxOrderProductSales struct {
	ID                int64   `json:"id"`
	OrderUniqueID     string  `json:"orderUniqueId"`
	OrderBillUniqueID string  `json:"orderBillUniqueId"`
	UniqueID          string  `json:"uniqueId"`
	MerchantMenuID    int64   `json:"productId"`
	Name              string  `json:"name"`
	Price             float64 `json:"price"`
	Qty               int     `json:"qty"`
	Notes             string  `json:"notes"`
	BusinessDate      int64   `json:"businessDate"`
	CreatedAt         int64   `json:"createdAt"`
	UpdatedAt         int64   `json:"updatedAt"`
	QrCode            string  `json:"qrCode"`
}

// ResTrxOrderProductSales ..
type ResTrxOrderProductSales struct {
	ID                int64      `json:"id"`
	OrderUniqueID     string     `json:"orderUniqueId"`
	OrderBillUniqueID string     `json:"orderBillUniqueId"`
	UniqueID          string     `json:"uniqueId"`
	MerchantMenuID    int64      `json:"productId"`
	Name              string     `json:"name"`
	Price             float64    `json:"price"`
	Qty               int        `json:"qty"`
	Notes             string     `json:"notes"`
	BusinessDate      time.Time  `json:"businessDate"`
	CreatedAt         time.Time  `json:"createdAt"`
	UpdatedAt         *time.Time `json:"updatedAt"`
	QrCode            string     `json:"qrCode"`
}

// ResTrxOrderTaxSales ..
type ResTrxOrderTaxSales struct {
	ID                int64      `json:"id"`
	UniqueID          string     `json:"uniqueId"`
	OrderUniqueID     string     `json:"orderUniqueId"`
	OrderBillUniqueID string     `json:"orderBillUniqueId"`
	Name              string     `json:"name"`
	Amount            float32    `json:"amount"`
	Type              int        `json:"type"`
	CreatedAt         time.Time  `json:"createdAt"`
	UpdatedAt         *time.Time `json:"updatedAt"`
	MerchantTaxID     int64      `json:"merchantTaxId"`
}

// TrxOrderTaxSales ..
type TrxOrderTaxSales struct {
	ID                int64   `json:"id"`
	UniqueID          string  `json:"uniqueId"`
	OrderUniqueID     string  `json:"orderUniqueId"`
	OrderBillUniqueID string  `json:"orderBillUniqueId"`
	Name              string  `json:"name"`
	Amount            float32 `json:"amount"`
	Type              int     `json:"type"`
	CreatedAt         int64   `json:"createdAt"`
	UpdatedAt         int64   `json:"updatedAt"`
	MerchantTaxID     int64   `json:"merchantTaxId"`
	IsActive          bool    `json:"isActive"`
}

// ResTrxOrderList ..
type ResTrxOrderList struct {
	Trx          []Trx                  `json:"trx"`
	Order        []TrxOrder             `json:"order"`
	ProductSales []TrxOrderProductSales `json:"productSales"`
	PaymentSales []TrxOrderPaymentSales `json:"paymentSales"`
	OrderBills   []TrxOrderBill         `json:"orderBills"`
	TaxSales     []TrxOrderTaxSales     `json:"taxSales"`
}

// ReqTrxOrderOnline ..
type ReqTrxOrderOnline struct {
	Order        TrxOrder               `json:"order"`
	ProductSales []TrxOrderProductSales `json:"productSales"`
	PaymentSales []TrxOrderPaymentSales `json:"paymentSales"`
	OrderBills   []TrxOrderBill         `json:"orderBills"`
	TaxSales     []TrxOrderTaxSales     `json:"taxSales"`
	TrxID        string                 `json:"trxId"`
}

// ResTrx response transaction by consumer
type ResTrx struct {
	ID string `json:"id"`
}

// ReqCreateVisitTrx request create homevisit trx by consumer
type ReqCreateVisitTrx struct {
	VisitSales      []Visit `json:"visitSales" binding:"required"`
	Notes           string  `json:"notes" binding:"required"`
	GrandTotal      float64 `json:"grandTotal" binding:"required"`
	Address         string  `json:"address" binding:"required"`
	Latitude        float64 `json:"latitude" binding:"required"`
	Longitude       float64 `json:"longitude" binding:"required"`
	CustomerName    string  `json:"customerName" binding:"required"`
	TrxID           string  `json:"trxId" binding:"required"`
	Phone           string  `json:"phone" binding:"required"`
	PaymentMethodID int64   `json:"paymentMethodId" binding:"required"`
}

// Visit ..
type Visit struct {
	HomevisitSalesID int64  `json:"salesId"  binding:"required"`
	Total            int64  `json:"deposit" binding:"required"`
	Menus            []Menu `json:"menus" binding:"required"`
}

type Menu struct {
	MenuID   int64 `json:"menu_id"`
	Quantity int   `json:"quantity"`
}

type ResMenu struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}

// ResTrxHistory response for order history end user apps
type ResTrxHistory struct {
	ID           string         `json:"trxId"`
	Logo         string         `json:"logo"`
	MerchantName string         `json:"merchantName"`
	CreatedAt    time.Time      `json:"date"`
	Status       string         `json:"status"`
	Amount       int64          `json:"amount"`
	Types        string         `json:"types"`
	Address      string         `json:"address"`
	QrCode       string         `json:"qrCode"`
	Phone        string         `json:"phone"`
	Notes        string         `json:"notes"`
	Detail       DetailOrderHis `json:"detail"`
}

// TrxHomevisitSales response for order history
type TrxHomevisitSales struct {
	ID                  int64      `json:"id" `
	TrxHomevisitSalesID int64      `json:"trxHomeVisitSalesId"`
	HomevisitSalesID    int64      `json:"homeVisitSalesId"`
	Notes               string     `json:"notes"`
	GrandTotal          float64    `json:"total"`
	Address             string     `json:"address"`
	Latitude            float64    `json:"latitude"`
	Longitude           float64    `json:"longitude"`
	TrxID               string     `json:"trxId"`
	CustomerName        string     `json:"customerName"`
	CreatedAt           time.Time  `json:"createdAt"`
	UpdatedAt           *time.Time `json:"updatedAt"`
	DeletedAt           *time.Time `json:"-"`
	Phone               string     `json:"phone"`
	MerchantName        string     `json:"merchantName"`
	MerchantLogo        string     `json:"merchantLogo"`
	Deposit             int64      `json:"deposit"`
	MerchantID          int64      `json:"merchantId"`
	Available           int        `json:"-"`
	PaymentName         string     `json:"-"`
}

// HomeVisitSales for myparking manager app
type HomeVisitSales struct {
	ID               int64      `json:"id" `
	HomevisitSalesID int64      `json:"homeVisitSalesId"`
	Notes            string     `json:"notes"`
	Total            int64      `json:"total"`
	Address          string     `json:"address"`
	Latitude         float64    `json:"latitude"`
	Longitude        float64    `json:"longitude"`
	TrxID            string     `json:"trxId"`
	CustomerName     string     `json:"customerName"`
	CreatedAt        time.Time  `json:"createdAt"`
	UpdatedAt        *time.Time `json:"updatedAt"`
	DeletedAt        *time.Time `json:"-"`
	MerchantName     string     `json:"merchantName"`
	MerchantLogo     string     `json:"merchantLogo"`
	Deposit          int64      `json:"deposit"`
	StartDate        time.Time  `json:"startDate"`
	EndDate          time.Time  `json:"endDate"`
	ProfilePicture   string     `json:"profilePicture"`
	TrxVisitSalesID  int64      `json:"trxVisitSalesId"`
}

// TrxOrderMerchant response for fcm
type TrxOrderMerchant struct {
	ID              int64      `json:"id"`
	UniqueID        string     `json:"uniqueId"`
	OrderNo         string     `json:"orderNo"`
	BillNo          string     `json:"billNo"`
	IsClose         bool       `json:"isClose"`
	Note            string     `json:"note"`
	Types           int        `json:"types"`
	BusinessDate    time.Time  `json:"businessDate"`
	TotalDiscount   float64    `json:"totalDiscount"`
	GrandTotal      float64    `json:"grandTotal"`
	TrxID           string     `json:"trxId"`
	MerchantUsersID int64      `json:"merchantUsersId"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       *time.Time `json:"updatedAt,omitempty"`
	DeletedAt       *time.Time `json:"-,omitempty"`
	MerchantName    string     `json:"merchantName"`
	MerchantLogo    string     `json:"merchantLogo"`
	Address         string     `json:"address"`
	Status          string     `json:"status"`
}

// ResHomeVisitBookingList ..
type ResHomeVisitBookingList struct {
	ID              int64     `json:"id"`
	CustomerName    string    `json:"customerName"`
	StartDate       time.Time `json:"startDate"`
	EndDate         time.Time `json:"endDate"`
	Deposit         int64     `json:"deposit"`
	Status          string    `json:"status"`
	TrxID           string    `json:"trxID"`
	TransactionDate time.Time `json:"transactionDate"`
	GrandTotal      float64   `json:"grandTotal"`
	PaymentMethod   string    `json:"paymentMethod"`
}

// ResHomeVisitBookingDetail ..
type ResHomeVisitBookingDetail struct {
	ID           int64  `json:"id"`
	CustomerName string `json:"customerName"`
	Date         string `json:"date"`
	Address      string `json:"address"`
	Deposit      int64  `json:"deposit"`
	Status       string `json:"status"`
	Phone1       string `json:"phone1"`
	Phone2       string `json:"phone2"`
	Notes        string `json:"notes"`
}

// ResScheduleTime ..
type ResScheduleTime struct {
	ScheduleTimeStart string `json:"scheduleTimeStart"`
	ScheduleTimeEnd   string `json:"scheduleTimeEnd"`
}

// ResHomeVisitBookingListTime ..
type ResHomeVisitBookingListTime struct {
	ID              int64     `json:"id"`
	CustomerName    string    `json:"customerName"`
	StartDate       string    `json:"startDate"`
	EndDate         string    `json:"endDate"`
	Deposit         int64     `json:"deposit"`
	Status          string    `json:"status"`
	TrxID           string    `json:"trxID"`
	TransactionDate time.Time `json:"transactionDate"`
	GrandTotal      float64   `json:"grandTotal"`
	PaymentMethod   string    `json:"paymentMethod"`
}

// ResHomeVisitBookingDetailTime ..
type ResHomeVisitBookingDetailTime struct {
	ID              int64     `json:"id"`
	CustomerName    string    `json:"customerName"`
	StartDate       time.Time `json:"startDate"`
	EndDate         time.Time `json:"endDate"`
	Address         string    `json:"address"`
	Deposit         int64     `json:"deposit"`
	Status          string    `json:"status"`
	Phone1          string    `json:"phone1"`
	Phone2          string    `json:"phone2"`
	Notes           string    `json:"notes"`
	Menus           []ResMenu `json:"menus"`
	TrxID           string    `json:"trxID"`
	GrandTotal      float64   `json:"grandTotal"`
	PaymentMethod   string    `json:"paymentMethod"`
	TransactionDate time.Time `json:"transactionDate"`
}

// ResHomeVisitBookingDetailTimeNew ..
type ResHomeVisitBookingDetailTimeNew struct {
	ID              int64     `json:"id"`
	CustomerName    string    `json:"customerName"`
	StartDate       string    `json:"startDate"`
	EndDate         string    `json:"endDate"`
	Address         string    `json:"address"`
	Deposit         int64     `json:"deposit"`
	Status          string    `json:"status"`
	Phone1          string    `json:"phone1"`
	Phone2          string    `json:"phone2"`
	Notes           string    `json:"notes"`
	Menus           []ResMenu `json:"menus"`
	TrxID           string    `json:"trxID"`
	GrandTotal      float64   `json:"grandTotal"`
	PaymentMethod   string    `json:"paymentMethod"`
	TransactionDate time.Time `json:"transactionDate"`
}

// OrderTrx for pos get online order
type OrderTrx struct {
	ID              int64
	OrderNo         string
	BillNo          string
	IsClose         bool
	Note            string
	Types           int
	MerchantUsersID int64
	BusinessDate    time.Time
	TotalDiscount   float64
	GrandTotal      float64
	TrxID           string
	CreatedAt       time.Time
	UpdatedAt       *time.Time
	DeletedAt       *time.Time
	UniqueID        string
}

// DetailOrderHis struct Order History Detail by End Users
type DetailOrderHis struct {
	OrderDetails   []OrderDetail `json:"orderDetails"`
	PaymentDetails PaymentDetail `json:"paymentDetails"`
	PaymentName    string        `json:"paymentName"`
}

// OrderDetail struct detail item in DetailOrderHis
type OrderDetail struct {
	Name  string    `json:"productName"`
	Qty   int       `json:"qty,omitempty"`
	Menus []ResMenu `json:"menus,omitempty"`
}

// PaymentDetail struct payment detail in DetailOrderHis
type PaymentDetail struct {
	SubTotal float64 `json:"subtotal"`
	Total    float64 `json:"total"`
	Tax      float64 `json:"tax"`
	TaxName  string  `json:"taxName"`
	TaxType  int     `json:"taxType"`
	IsActive bool    `json:"isActive"`
}

// TrxVisit struct order history visit by end users
type TrxVisit struct {
	ID        int64     `json:"id"`
	TrxID     string    `json:"trxId"`
	Logo      string    `json:"merchantLogo"`
	Name      string    `json:"merchantName"`
	Address   string    `json:"address"`
	Deposit   int64     `json:"amount"`
	CreatedAt time.Time `json:"createdAt"`
	Status    string    `json:"status"`
	Notes     string    `json:"notes"`
	Phone     string    `json:"phone"`
}

// ResTransactionReport ...
type ResTransactionReport struct {
	ID              int64     `json:"id"`
	TransactionDate time.Time `json:"transactionDate"`
	StartDate       time.Time `json:"startDate"`
	EndDate         time.Time `json:"endDate"`
	CustomerName    string    `json:"customerName"`
	TrxID           string    `json:"trxID"`
	BillNo          string    `json:"billNo"`
	OrderNo         string    `json:"orderNo"`
	Dates           string    `json:"dates"`
	Times           string    `json:"times"`
	ProductName     string    `json:"productName"`
	UserName        string    `json:"userName"`
	Qty             int       `json:"qty"`
	TotalTax        float64   `json:"totalTax"`
	GrandTotal      float64   `json:"grandTotal"`
	Status          string    `json:"status"`
	ExternalID      string    `json:"paymentMethod"`
	TypeOrder       string    `json:"typeOrder"`
}

// ResTransactionReportSingle ...
type ResTransactionReportSingle struct {
	ID         string  `json:"ID"`
	TrxID      string  `json:"trxID"`
	BillNo     string  `json:"billNo"`
	OrderNo    string  `json:"orderNo"`
	Dates      string  `json:"dates"`
	Times      string  `json:"times"`
	UserName   string  `json:"userName"`
	TotalTax   float64 `json:"totalTax"`
	GrandTotal float64 `json:"grandTotal"`
	Status     string  `json:"status"`
	ExternalID string  `json:"paymentMethod"`
	TypeOrder  string  `json:"typeOrder"`
}

type ResTrxOnlineClosed struct {
	entity.TrxOrder
	ExternalID    string `json:"externalID"`
	PaymentMethod string `json:"paymentMethod"`
}
