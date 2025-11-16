package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Flashcard struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	CollectionID uuid.UUID `gorm:"type:uuid;not null;index" json:"collection_id"`
	Front        string    `gorm:"type:text;not null" json:"front"`
	Back         string    `gorm:"type:text;not null" json:"back"`
	Extra        JSONB     `gorm:"type:jsonb" json:"extra,omitempty"`
	CreatedByIA  bool      `gorm:"default:true" json:"created_by_ia"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// Relacionamentos
	Collection Collection `gorm:"foreignKey:CollectionID;constraint:OnDelete:CASCADE" json:"-"`
}

type JSONB map[string]interface{}

func (j JSONB) Value() (driver.Value, error) {
	return json.Marshal(j)
}

func (j *JSONB) Scan(value interface{}) error {
	bytes := value.([]byte)
	return json.Unmarshal(bytes, &j)
}

type FlashcardRequest struct {
	Front string `json:"front" binding:"required,max=1000"`
	Back  string `json:"back" binding:"required,max=5000"`
	Extra JSONB  `json:"extra,omitempty"`
}

type FlashcardResponse struct {
	ID          uuid.UUID `json:"id"`
	Front       string    `json:"front"`
	Back        string    `json:"back"`
	Extra       JSONB     `json:"extra,omitempty"`
	CreatedByIA bool      `json:"created_by_ia"`
	CreatedAt   time.Time `json:"created_at"`
}

func (f *Flashcard) TableName() string {
	return "flashcards"
}
