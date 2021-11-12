package model

import "time"

// ResLiveTracking live tracking map on consumer apps
type ResLiveTracking struct {
	Typess       string    `json:"types"`
	Lat          float64   `json:"latitude"`
	Lon          float64   `json:"longitude"`
	MerchantID   int64     `json:"merchantId"`
	MerchantName string    `json:"merchantName"`
	Logo         string    `json:"logo"`
	TasksID      int64     `json:"tasksId"`
	Nearby       float64   `json:"distance"`
	LogTime      time.Time `json:"logTime"`
	Banner       string    `json:"banner"`
	IGAccount    string    `json:"merchantIG"`
	Status       int       `json:"status"`
}

// ResParkingSpace show parking space on consumer map
type ResParkingSpace struct {
	ID        int64   `json:"id"`
	Lat       float64 `json:"latitude"`
	Lon       float64 `json:"longitude"`
	Name      string  `json:"name"`
	Address   string  `json:"address"`
	Nearby    string  `json:"distance"`
	IsCheckin bool    `json:"isCheckin"`
}

// ResMerchantNearby already check in
type ResMerchantNearby struct {
	Logo             string    `json:"logo"`
	Banner           string    `json:"banner"`
	MerchantID       int64     `json:"merchantId"`
	MerchantUsersID  int64     `json:"merchantUsersId"`
	MerchantName     string    `json:"name"`
	Address          string    `json:"address"`
	Latitude         float64   `json:"latitude"`
	Longitude        float64   `json:"longitude"`
	Nearby           float64   `json:"distance"`
	Typess           string    `json:"types"`
	IGAccount        string    `json:"merchantIG"`
	MerchantCategory string    `json:"merchantCategory"`
	CategoryColor    string    `json:"categoryColor"`
	City             string    `json:"city"`
	UpdatedAt        time.Time `json:"-"`
	TypesID          int64     `json:"typesId"`
	Status           int       `json:"status"`
}

// NearbySorted ..
type NearbySorted []ResMerchantNearby

// Len ..
func (n NearbySorted) Len() int {
	return len(n)
}

// Less ..
func (n NearbySorted) Less(i, j int) bool {
	return n[i].Nearby < n[j].Nearby
}

// Swap ..
func (n NearbySorted) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}

// ResParkingSpaceDetail after click Parking Space icon
type ResParkingSpaceDetail struct {
	Logo         string      `json:"logo"`
	Banner       string      `json:"banner"`
	MerchantName string      `json:"merchantName"`
	MerchantID   int64       `json:"merchantId"`
	IGAccount    string      `json:"merchantIG"`
	Address      string      `json:"address"`
	IsCheckin    bool        `json:"isCheckin"`
	Schedules    []Schedules `json:"schedules"`
}

// Schedules schedules merchant on parking space
type Schedules struct {
	ID        int64     `json:"salesId"`
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
}

// ResVisitSales show homevisit sales available
type ResVisitSales struct {
	ID             int64  `json:"id"`
	Banner         string `json:"banner"`
	Name           string `json:"merchantName"`
	Address        string `json:"address"`
	City           string `json:"city"`
	Logo           string `json:"logo"`
	Category       string `json:"category"`
	Category_color string `json:"categoryColor"`
	IgAccount      string `json:"merchantIG"`
	Terms          string `json:"terms"`
}

// ResVisitSalesDetail ...
type ResVisitSalesDetail struct {
	ID        int64     `json:"id"`
	StartDate time.Time `json:"startDate"` // timestamp
	EndDate   time.Time `json:"endDate"`   // timestamp
	Deposit   int64     `json:"deposit"`
}

// ResPaymentMethod payment method active
type ResPaymentMethod struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Types        int    `json:"type"`
	IsActive     bool   `json:"isActive"`
	ProviderName string `json:"providerName"`
}

// MerchantSpace ..
type MerchantSpace struct {
	MerchantID   int64  `json:"id"`
	Logo         string `json:"logo"`
	MerchantName string `json:"name"`
	Banner       string `json:"banner"`
	IGAccount    string `json:"igAccount"`
	Address      string `json:"address"`
}
