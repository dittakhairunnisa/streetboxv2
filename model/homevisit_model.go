package model

import "time"

// ReqCreateHomevisitSales ..

type ReqBatchCreateHomevisitSales struct {
	Request []ReqCreateHomevisitSales `json:"request"`
}
type ReqCreateHomevisitSales struct {
	Date    string                          `json:"date"    binding:"required"`
	Deposit int64                           `json:"deposit" binding:"required"`
	Summary []ReqCreateHomeDetailVisitSales `json:"summary"`
}

// ReqCreateHomeDetailVisitSales ..
type ReqCreateHomeDetailVisitSales struct {
	StartTime         string `json:"startTime" binding:"required"`
	EndTime           string `json:"endTime" binding:"required"`
	NumberOfFoodtruck int    `json:"numberOfFoodtruck" binding:"required"`
}

// ReqUpdateHomevisitSales ..
type ReqUpdateHomevisitSales struct {
	Date    string                          `json:"date"    binding:"required"`
	Deposit int64                           `json:"deposit" binding:"required"`
	Summary []ReqUpdateHomeDetailVisitSales `json:"summary" binding:"required"`
}

// ReqUpdateHomeDetailVisitSales ..
type ReqUpdateHomeDetailVisitSales struct {
	ID                int64  `json:"id" binding:"required"`
	StartTime         string `json:"startTime" binding:"required"`
	EndTime           string `json:"endTime" binding:"required"`
	NumberOfFoodtruck int    `json:"numberOfFoodtruck" binding:"required"`
}

// ResHomeVisitSales ..
type ResHomeVisitSales struct {
	ID         int64      `json:"id"`
	MerchantID int64      `json:"merchantId"`
	StartDate  time.Time  `json:"startDate"`
	EndDate    time.Time  `json:"endDate"`
	Deposit    int64      `json:"deposit"`
	Total      int        `json:"total"`
	Available  int        `json:"available"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  *time.Time `json:"updatedAt"`
}

// ResHomeVisitGetInfo ..
type ResHomeVisitGetInfo struct {
	Date    string                 `json:"date"`
	Deposit int64                  `json:"deposit"`
	Summary *[]ResHomeVisitDetails `json:"summary"`
}

// ResHomeVisitDetails ..
type ResHomeVisitDetails struct {
	ID                int64  `json:"id"`
	StartTime         string `json:"startTime"`
	EndTime           string `json:"endTime"`
	NumberOfFoodtruck int    `json:"numberOfFoodtruck"`
}

// ResCountFoodtruck ..
type ResCountFoodtruck struct {
	Foodtrucks int `json:"foodtrucks"`
}
