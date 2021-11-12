package entity

import "time"

type TrxHomevisitMenuSales struct {
	ID                  int64      `json:"id"`
	TrxHomevisitSalesID int64      `json:"trx_homevisit_sales_id"`
	MenuID              int64      `json:"menu_id"`
	Quantity            int64      `json:"quantity"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           *time.Time `json:"updated_at"`
	DeletedAt           *time.Time `json:"-"`
}
