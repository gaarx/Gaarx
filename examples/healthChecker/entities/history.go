package entities

import (
	"github.com/jinzhu/gorm"
	"time"
)

type History struct {
	gorm.Model
	Url         string    `json:"url"`
	CheckTime   time.Time `json:"checkTime"`
	StatusCode  int       `json:"statusCode"`
	RequestTime int       `json:"requestTime"`
}
