package models

import (
	"time"

	"github.com/google/uuid"
)

type Collection struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	Name      string    `gorm:"not null" json:"name"`
	IsPublic  bool      `gorm:"default:false" json:"is_public"`
	MaxCards  int       `gorm:"default:10" json:"max_cards"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relacionamentos
	User       User        `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
	Flashcards []Flashcard `gorm:"foreignKey:CollectionID;constraint:OnDelete:CASCADE" json:"-"`
	Shares     []Share     `gorm:"foreignKey:CollectionID;constraint:OnDelete:CASCADE" json:"-"`
}

type CollectionRequest struct {
	Name     string `json:"name" binding:"required,min=3,max=255"`
	IsPublic bool   `json:"is_public"`
	MaxCards int    `json:"max_cards" binding:"min=1,max=1000"`
}

type CollectionResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	IsPublic  bool      `json:"is_public"`
	MaxCards  int       `json:"max_cards"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (c *Collection) TableName() string {
	return "collections"
}
