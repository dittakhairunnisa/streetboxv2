package model

import "streetbox.id/entity"

type RespCallsByFoodtruckID struct {
	entity.CanvassingCall
	LongitudeUser      float64 `json:"longitude_user"`
	LatitudeUser       float64 `json:"latitude_user"`
	LongitudeFoodtruck      float64 `json:"longitude_foodtruck"`
	LatitudeFoodtruck       float64 `json:"latitude_foodtruck"`
	Name           string  `json:"name"`
	ProfilePicture string  `json:"profile_picture"`
	Phone          string  `json:"phone"`
}

type RespCallsByUserID struct {
	entity.CanvassingCall
	LongitudeUser      float64 `json:"longitude_user"`
	LatitudeUser       float64 `json:"latitude_user"`
	LongitudeFoodtruck      float64 `json:"longitude_foodtruck"`
	LatitudeFoodtruck       float64 `json:"latitude_foodtruck"`
	Name      string  `json:"name"`
	Logo      string  `json:"logo"`
	IgAccount string  `json:"ig_account"`
	PlatNo    string  `json:"plat_no"`
	Phone     string  `json:"phone"`
}

type RespNotifByUserID struct {
	entity.CanvassingNotif
	LongitudeFoodtruck float64 `json:"longitude_foodtruck"`
	LatitudeFoodtruck  float64 `json:"latitude_foodtruck"`
	LongitudeUser float64 `json:"longitude_user"`
	LatitudeUser  float64 `json:"latitude_user"`
	Name      string  `json:"name"`
	Logo      string  `json:"logo"`
	IgAccount string  `json:"ig_account"`
	PlatNo    string  `json:"plat_no"`
}

type RespCanvassing struct {
	entity.Canvassing
	LastBlast   string `json:"last_blast"`
	IsAutoBlast bool   `json:"is_auto_blast"`
}

type AutoCanvassing struct {
	ID          int64
	FoodtruckID int64
	MerchantLogo string
	IsAutoBlast bool
}
