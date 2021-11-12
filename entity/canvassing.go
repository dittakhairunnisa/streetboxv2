package entity

type Canvassing struct {
	ID            int64   `json:"merchant_id"`
	Radius        float64 `json:"radius"`
	Interval      int     `json:"interval"`
	Expire        int     `json:"expire"`
	Cooldown      int     `json:"cooldown"`
	IsActive      bool    `json:"is_active"`
	LastAutoBlast string  `json:"last_auto_blast"`
}

type CanvassingNotif struct {
	ID            int64  `json:"id"`
	FoodtruckID   int64  `json:"foodtruck_id"`
	CustomerID    int64  `json:"customer_id"`
	CustomerToken string `json:"-"`
	Status        string `json:"status"`
	CreatedAt     string `json:"created_at"`
}

type CanvassingCall struct {
	ID          int64  `json:"id"`
	NotifID     int64  `json:"notif_id"`
	FoodtruckID int64  `json:"foodtruck_id"`
	CustomerID  int64  `json:"customer_id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	DeletedAt   string `json:"deleted_at"`
	Status      string `json:"status"`
	QueueNo     int    `json:"queue_no"`
}
type UserTokens struct {
	ID                int64
	RegistrationToken string
}

type FoodtruckLocation struct {
	ID        int64   `json:"foodtruck_id"`
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

type UsersLocation struct {
	ID        int64   `json:"user_id"`
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}
