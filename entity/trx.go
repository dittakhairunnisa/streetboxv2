package entity

import "time"

// Trx Parent table transaction order and homevisit
type Trx struct {
	ID         string     `json:"id" gorm:"primary_key"`
	Types      string     `json:"types" gorm:"not null"` // ORDER; VISIT
	UsersID    int64      `json:"usersId" gorm:"not null"`
	CreatedAt  time.Time  `json:"createdAt" gorm:"not null"`
	UpdatedAt  *time.Time `json:"updatedAt"`
	DeletedAt  *time.Time `json:"-"`
	Status     string     `json:"status" gorm:"not null"` // UNPAID;PAID;COMPLETED;FAILED;CANCEL;VOID?
	ExternalID string     `json:"-"`
	Address    string     `json:"address" gorm:"not null"`
	QrCode     string     `json:"qrCode"`
}
