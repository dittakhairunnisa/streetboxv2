package repository

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	"streetbox.id/entity"
	"streetbox.id/model"
)

type CanvassingRepo struct {
	DB *gorm.DB
}

func (c *CanvassingRepo) CreateCanvas(canv *entity.Canvassing) (err error) {
	if err = c.DB.Exec("INSERT INTO canvassing (id, radius, interval, is_active, expire, cooldown) VALUES (?, ?, ?, ?, ?, ?)", canv.ID, canv.Radius, canv.Interval, canv.IsActive, canv.Expire, canv.Cooldown).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
	}
	return
}

func (c *CanvassingRepo) GetCanvas(id int64) (canv entity.Canvassing, err error) {
	if err = c.DB.Where("id = ?", id).First(&canv).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
	}
	return
}

func (c *CanvassingRepo) GetFoodtruckCanvas(id, merchID int64) (canv model.RespCanvassing, err error) {
	if err = c.DB.Raw("SELECT id, radius, interval, cooldown, expire, is_active, to_char(last_auto_blast, 'yyyy-MM-dd HH24:mi:ss.MS') last_auto_blast FROM canvassing WHERE id = ?", merchID).Scan(&canv).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
	}

	if err = c.DB.Raw("SELECT to_char(last_blast, 'yyyy-MM-dd HH24:mi:ss.MS') last_blast, is_auto_blast FROM foodtruck_blast WHERE id = ?", id).Scan(&canv).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
		err = nil
	}
	return
}

func (c *CanvassingRepo) GetTriggeredCanvas() (canv []model.AutoCanvassing, err error) {
	if err = c.DB.Raw("select c.id, mu.id foodtruck_id, f.is_auto_blast, m.logo merchant_logo from canvassing c, merchant_users mu left join foodtruck_blast f on f.id = mu.id left join merchant m on m.id = mu.merchant_id where extract(minute from (current_timestamp - c.last_auto_blast)) >= c.interval and c.id = mu.merchant_id and is_active").Scan(&canv).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
	}
	return
}
func (c *CanvassingRepo) UpdateCanvas(canv entity.Canvassing) (err error) {
	if err = c.DB.Exec("UPDATE canvassing SET radius = ?, interval = ?, expire = ?, cooldown = ?, is_active = ? WHERE id = ?", canv.Radius, canv.Interval, canv.Expire, canv.Cooldown, canv.IsActive, canv.ID).Error; err != nil {
		log.Printf("ERROR: %s\n", err.Error())
	}
	return
}

func (c *CanvassingRepo) GetNearby(id int64, long, lat, rad float64) (users []entity.UserTokens, err error) {
	result := c.DB.Table("users u").Select("DISTINCT u.registration_token, u.id").Joins("LEFT JOIN users_location l ON l.id = u.id").Where("distance(?, ?, l.latitude, l.longitude) <= ? AND u.registration_token IS NOT NULL AND u.id NOT IN (SELECT customer_id FROM canvassing_call WHERE foodtruck_id = ? AND status = 'REQUEST')", lat, long, rad, id).Scan(&users)

	if result.Error != nil {
		log.Printf("ERROR: %s", err.Error())
	}
	return
}

func (c *CanvassingRepo) GetFoodtruckLocation(id int64) (loc entity.FoodtruckLocation, err error) {
	if err = c.DB.Where("id = ?", id).Find(&loc).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
	}
	return
}

func (c *CanvassingRepo) GetUserLocation(id int64) (loc entity.UsersLocation, err error) {
	if err = c.DB.Where("id = ?", id).Find(&loc).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
	}
	return
}

func (c *CanvassingRepo) UpdateFoodtruckLocation(loc entity.FoodtruckLocation) (err error) {
	sql := "INSERT INTO foodtruck_location as fl (id, longitude, latitude) VALUES(?, ?, ?) ON CONFLICT (id) DO UPDATE SET longitude = ?, latitude = ? WHERE fl.id = ?"
	if err = c.DB.Exec(sql, loc.ID, loc.Longitude, loc.Latitude, loc.Longitude, loc.Latitude, loc.ID).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
	}
	return
}

func (c *CanvassingRepo) UpdateUsersLocation(loc entity.UsersLocation) (err error) {
	sql := "INSERT INTO users_location as ul (id, longitude, latitude) VALUES(?, ?, ?) ON CONFLICT (id) DO UPDATE SET longitude = ?, latitude = ? WHERE ul.id = ?"
	if err = c.DB.Exec(sql, loc.ID, loc.Longitude, loc.Latitude, loc.Longitude, loc.Latitude, loc.ID).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
	}
	return
}

func (c *CanvassingRepo) SaveNotification(notif entity.CanvassingNotif) (err error) {
	if err = c.DB.Exec("INSERT INTO canvassing_notif (foodtruck_id, customer_id, customer_token) VALUES (?, ?, ?)", notif.FoodtruckID, notif.CustomerID, notif.CustomerToken).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
	}
	return
}

func (c *CanvassingRepo) UpdateNotification(id int64, status string) (err error) {
	result := c.DB.Exec("UPDATE canvassing_notif SET status = ? WHERE id = ?", status, id)
	err = result.Error
	if result.RowsAffected < 1 && err == nil {
		err = fmt.Errorf("Notification has been expired")
	}
	if err != nil {
		log.Printf("ERROR: %s", err.Error())
	}
	return
}

func (c *CanvassingRepo) GetNotificationByID(id int64) (notif entity.CanvassingNotif, err error) {
	if err = c.DB.Where("id = ?", id).First(&notif).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
	}
	return
}

func (c *CanvassingRepo) GetNotificationsByUserID(id int64) (notifs []model.RespNotifByUserID, err error) {
	if err = c.DB.Raw("select distinct on(t.status, t.foodtruck_id, t.customer_id, DATE(t.created_at)) * from ( select distinct on (n.status, n.foodtruck_id, n.customer_id, n.created_at) n.*, n.customer_token, f.longitude longitude_foodtruck, f.latitude latitude_foodtruck, ul.longitude longitude_user, ul.latitude latitude_user, m.logo, m.name, m.ig_account, u.plat_no from canvassing_notif n left join foodtruck_location f on n.foodtruck_id = f.id left join merchant_users mu on mu.id = n.foodtruck_id left join merchant m on m.id = mu.merchant_id left join users u on u.id = mu.users_id left join users_location ul on n.customer_id = ul.id where n.customer_id = ? order by n.created_at, n.status desc) t order by DATE(t.created_at) asc", id).Scan(&notifs).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
	}
	return
}

func (c *CanvassingRepo) CreateCall(call *entity.CanvassingCall) (foodtruck entity.UserTokens, err error) {
	if err = c.DB.Raw("INSERT INTO canvassing_call (notif_id, foodtruck_id, customer_id) VALUES (?, ?, ?) RETURNING id", call.NotifID, call.FoodtruckID, call.CustomerID).Scan(call).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
	}
	if err = c.DB.Raw("SELECT users_id as id, registration_token FROM merchant_users WHERE id = ?", call.FoodtruckID).Scan(&foodtruck).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
	}
	return
}

func (c *CanvassingRepo) UpdateCall(id int64, status string, queue int) (call entity.CanvassingCall, err error) {
	if err = c.DB.Raw("UPDATE canvassing_call SET status = ?, queue_no = ?, updated_at = current_timestamp WHERE id = ? RETURNING *", status, queue, id).Scan(&call).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
	}
	return
}

func (c *CanvassingRepo) UpdateStatusCall(id int64, status string) (call entity.CanvassingCall, err error) {
	if err = c.DB.Raw("UPDATE canvassing_call SET status = ?, updated_at = current_timestamp WHERE id = ? RETURNING *", status, id).Scan(&call).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
	}
	return
}

func (c *CanvassingRepo) FinishCall(id int64) (err error) {
	if err = c.DB.Exec("UPDATE canvassing_call SET status = 'FINISH', updated_at = current_timestamp WHERE id = ?", id).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
	}
	return
}

func (c *CanvassingRepo) GetCallsByCustomerID(id int64) (calls []model.RespCallsByUserID, err error) {
	filter := "WHERE c.customer_id = ? AND c.status IN ('ACCEPT', 'ONPROCESS')"
	if err = c.DB.Raw("SELECT c.*, ul.longitude longitude_user, ul.latitude latitude_user, fl.longitude longitude_foodtruck, fl.latitude latitude_foodtruck, m.name, m.logo, m.ig_account, u.plat_no, m.phone FROM canvassing_call c LEFT JOIN users_location ul ON ul.id = c.customer_id LEFT JOIN foodtruck_location fl ON fl.id = c.foodtruck_id LEFT JOIN merchant_users mu LEFT JOIN users u on u.id = mu.users_id ON mu.id = c.foodtruck_id LEFT JOIN merchant m ON m.id = mu.merchant_id " + filter, id).Scan(&calls).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
	}
	return
}

func (c *CanvassingRepo) GetCallsByFoodtruckID(id int64, status string) (calls []model.RespCallsByFoodtruckID, err error) {
	var filter string
	if status == "HISTORY" {
		filter = "WHERE c.foodtruck_id = ? AND c.status IN ('EXPIRE', 'REJECT', 'FINISH')"
	} else if status == "ACCEPT" {
		filter = fmt.Sprintf("WHERE c.foodtruck_id = ? AND c.status IN ('ACCEPT', 'ONPROCESS')")
	} else if status == "REQUEST" {
		filter = fmt.Sprintf("WHERE c.foodtruck_id = ? AND c.status = 'REQUEST'")
	}
	if err = c.DB.Raw("SELECT c.*, ul.longitude longitude_user, ul.latitude latitude_user, fl.longitude longitude_foodtruck, fl.latitude latitude_foodtruck, u.name, u.profile_picture, u.phone FROM canvassing_call c LEFT JOIN users_location ul ON ul.id = c.customer_id LEFT JOIN foodtruck_location fl ON fl.id = c.foodtruck_id LEFT JOIN users u ON u.id = c.customer_id " + filter + " ORDER BY c.queue_no", id).Scan(&calls).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
	}
	return
}

func (c *CanvassingRepo) UpdateFoodtruckBlast(id int64) (err error) {
	if err = c.DB.Exec("INSERT INTO foodtruck_blast as f VALUES (?, current_timestamp) ON CONFLICT (id) DO UPDATE SET last_blast = current_timestamp WHERE f.id = ?", id, id).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
	}
	return
}

func (c *CanvassingRepo) UpdateAutoBlast(id int64) {
	if err := c.DB.Exec("UPDATE canvassing SET last_auto_blast = current_timestamp WHERE id = ?", id).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
	}
}

func (c *CanvassingRepo) ToggleAutoBlast(id int64) {
	if err := c.DB.Exec("UPDATE foodtruck_blast SET is_auto_blast = NOT is_auto_blast WHERE id = ?", id).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
	}
}

func (c *CanvassingRepo) getFoodtruckToken(id int64) (foodtruck entity.UserTokens, err error) {
	if err = c.DB.Raw("SELECT id, registration_token FROM users WHERE id = ?", id).Scan(&foodtruck).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
	}
	return
}

func (c *CanvassingRepo) ExpireNotification() {
	if err := c.DB.Exec("UPDATE canvassing_notif SET status = 'EXPIRE' WHERE id IN (SELECT cn.id FROM canvassing_notif cn LEFT JOIN merchant_users mu ON cn.foodtruck_id = mu.id LEFT JOIN canvassing c ON mu.merchant_id = c.id WHERE EXTRACT (minute from (current_timestamp - cn.created_at)) >= c.expire AND (cn.status = 'ONGOING' OR cn.status = 'CALLING'))").Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
	}
}

func (c *CanvassingRepo) ExpireCall() {
	if err := c.DB.Exec("UPDATE canvassing_call SET status = 'EXPIRE' WHERE id IN (SELECT cc.id FROM canvassing_call cc LEFT JOIN merchant_users mu ON cc.foodtruck_id = mu.id LEFT JOIN canvassing c ON mu.merchant_id = c.id WHERE EXTRACT (minute from (current_timestamp - cc.created_at)) >= c.expire AND cc.status = 'REQUEST')").Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
	}
}

func (c *CanvassingRepo) ResetQueue() {
	if err := c.DB.Exec("UPDATE canvassing_call SET status = 'EXPIRE' WHERE status NOT IN ('FINISH', 'ONPROCESS', 'REJECT')").Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
	}
}
func New(db *gorm.DB) *CanvassingRepo {
	return &CanvassingRepo{
		DB: db,
	}
}
