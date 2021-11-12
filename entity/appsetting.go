package entity

// AppSetting ..
type AppSetting struct {
	ID         	int64   	`json:"id" gorm:"primary_key"`
	Key 		string  	`json:"key" gorm:"not null"`
	Value  		string  	`json:"value"`
}
