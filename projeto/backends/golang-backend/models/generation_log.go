package models

import (
	"time"

	"github.com/google/uuid"
)

type GenerationLog struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	Date      time.Time `gorm:"type:date;not null" json:"date"`
	Count     int       `gorm:"default:0" json:"count"`
	CreatedAt time.Time `json:"created_at"`

	// Relacionamentos
	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
}

func (g *GenerationLog) TableName() string {
	return "generation_logs"
}
