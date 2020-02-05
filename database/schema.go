package database

import "time"

type MediaItem struct {
	Id                int    `gorm:"primary_key;auto_increment:false"`
	PvrName           string `gorm:"primary_key"`
	WantedType        string `gorm:"primary_key"`
	AirDateUtc        time.Time
	LastSearchDateUtc *time.Time `gorm:"null"`
}
