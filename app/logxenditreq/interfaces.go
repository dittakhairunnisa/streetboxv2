package logxenditreq

import (
	"github.com/jinzhu/gorm"
	"streetbox.id/entity"
)

// RepoInterface ..
type RepoInterface interface {
	Create(*gorm.DB, *entity.LogXenditRequest) error
}
