package common

import (
	"time"

	"gorm.io/gorm"
)

type MARKET_MODEL struct {
	ID        uint           `json:"-" gorm:"primarykey"`
	Status    int            `json:"-" gorm:"comment:status"` // 1:ACTIVE 2:INACTIVE
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"created_time"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
