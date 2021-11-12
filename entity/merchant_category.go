package entity

type MerchantCategory struct {
	ID       int64  `json:"id" gorm:"primary_key"`
	Category string `json:"category"`
	Hexcode  string `json:"hexcode"`
}
