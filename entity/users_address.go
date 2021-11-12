package entity

type UsersAddress struct {
	ID          int64   `json:"id" gorm:"primary_key"`
	UserID      int64   `json:"user_id"`
	Person      string  `json:"person"`
	Address     string  `json:"address"`
	Phone       string  `json:"phone"`
	Primary     bool    `json:"primary"`
	AddressName string  `json:"address_name"`
	Longitude   float64 `json:"longitude"`
	Latitude    float64 `json:"latitude"`
}
