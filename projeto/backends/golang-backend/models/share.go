package models

import (
	"time"

	"github.com/google/uuid"
)

type Share struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	CollectionID uuid.UUID `gorm:"type:uuid;not null;column:collection_id;index" json:"collection_id"`
	SharedWithID uuid.UUID `gorm:"type:uuid;not null;column:shared_with_id;index" json:"shared_with_id"`
	Permissions  string    `gorm:"type:varchar(20);default:read" json:"permissions"` // read, write
	ShareID      string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"share_id"`
	CreatedAt    time.Time `json:"created_at"`

	// Relacionamentos
	Collection Collection `gorm:"foreignKey:CollectionID;constraint:OnDelete:CASCADE" json:"-"`
	SharedWith User       `gorm:"foreignKey:SharedWithID;constraint:OnDelete:CASCADE" json:"-"`
}

type ShareRequest struct {
	Permissions string `json:"permissions" binding:"required,oneof=read write"`
}

type ShareResponse struct {
	ID           uuid.UUID `json:"id"`
	CollectionID uuid.UUID `json:"collection_id"`
	Permissions  string    `json:"permissions"`
	ShareID      string    `json:"share_id"`
	CreatedAt    time.Time `json:"created_at"`
}

func (s *Share) TableName() string {
	return "shares"
}
