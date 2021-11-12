package model

import (
	"time"

	"streetbox.id/entity"
)

// ReqCreateMerchant ...
type ReqCreateMerchant struct {
	Name       string `json:"name" binding:"required"`
	Address    string `json:"address" binding:"required"`
	Phone      string `json:"phone" binding:"required"`
	Email      string `json:"email" binding:"required"`
	City       string `json:"city"`
	IGAccount  string `json:"igAccount"`
	CategoryID int64  `json:"categoryID"`
	Terms      string `json:"terms"`
}

// ReqXenditGenerateSubAccount ...
type ReqXenditGenerateSubAccount struct {
	ID    int64  `json:"id" binding:"required"`
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required"`
}

// ReqCreateMerchantMenu ...
type ReqCreateMerchantMenu struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       int64   `json:"price" binding:"required"`
	Qty         int     `json:"qty" binding:"required"`
	IsActive    bool    `json:"isActive"`
	IsNearby    bool    `json:"isNearby"`
	IsVisit     bool    `json:"isVisit"`
	Discount    float32 `json:"discount"`
}

// ReqUpdateMerchantMenu ...
type ReqUpdateMerchantMenu struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       int64   `json:"price" binding:"required"`
	Discount    float32 `json:"discount"`
	Qty         int     `json:"qty"`
	IsActive    bool    `json:"isActive"`
	IsNearby    bool    `json:"isNearby"`
	IsVisit     bool    `json:"isVisit"`
}

// MerchantTax ...
type MerchantTax struct {
	Name     string   `json:"name" binding:"required"`
	Amount   *float32 `json:"amount" binding:"required"`
	IsActive *bool    `json:"isActive" binding:"required"`
	Type     *int     `json:"type" binding:"required"`
}

// ResGetAllFoodtruck ...
type ResGetAllFoodtruck struct {
	UsersID int64  `json:"id"`
	Name    string `json:"name"`
}

// ResGetFoodtruckTasks ..
type ResGetFoodtruckTasks struct {
	ID        int64      `json:"id"`
	UserName  string     `json:"userName"`
	Name      string     `json:"name"`
	Phone     string     `json:"phone"`
	Address   string     `json:"address"`
	PlatNo    string     `json:"platNo"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
	DeletedAt *time.Time `json:"-"`
	TasksID   int64      `json:"tasksId"`
	Status    int        `json:"status"`
}

// ReqUpdateMerchant ...
type ReqUpdateMerchant struct {
	Name       string `json:"name,omitempty"`
	Address    string `json:"address,omitempty"`
	Phone      string `json:"phone,omitempty"`
	Email      string `json:"email,omitempty"`
	City       string `json:"city"`
	IGAccount  string `json:"igAccount"`
	CategoryID int64  `json:"categoryID"`
	Terms      string `json:"terms"`
	XenditID   string `json:"xendit_id"`
}

// ReqMerchantNearby ..
type ReqMerchantNearby struct {
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
}

// ResMerchantMenuList ..
type ResMerchantMenuList struct {
	ID                 int64      `json:"id"`
	MerchantID         int64      `json:"merchantId"`
	Name               string     `json:"name"`
	Price              float64    `json:"price"`
	PriceAfterDiscount float64    `json:"priceAfterDiscount"`
	Discount           float32    `json:"discount"`
	Qty                int        `json:"qty"`
	IsActive           bool       `json:"isActive"`
	IsNearby           bool       `json:"isNearby"`
	IsVisit            bool       `json:"isVisit"`
	Description        string     `json:"description"`
	Photo              string     `json:"photo"`
	CreatedAt          time.Time  `json:"createdAt"`
	UpdatedAt          *time.Time `json:"updatedAt"`
	DeletedAt          *time.Time `json:"deletedAt"`
}

// Merchant ..
type Merchant struct {
	entity.Merchant
	Category        string `json:"category"`
	Hexcode         string `json:"hexcode"`
	MerchantUsersID int64  `json:"merchantUsersId"`
}
