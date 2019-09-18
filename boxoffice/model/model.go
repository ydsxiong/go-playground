package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

/*
 * This table is for both successfully registered and reservations in progress guests
 */
type Guest struct {
	gorm.Model
	Name      string     `gorm:"unique" json:"name"`
	ExpiredAt *time.Time `json:"expired_at"`
}
