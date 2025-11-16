package models

import (
	"time"

	"github.com/google/uuid"
)

type Payment struct {
	ID             uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID         uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	SubscriptionID string     `json:"subscription_id"`
	Status         string     `gorm:"default:pending;index" json:"status"` // pending, active, cancelled
	StartDate      *time.Time `json:"start_date"`
	EndDate        *time.Time `json:"end_date"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`

	// Relacionamentos
	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
}

type PaymentResponse struct {
	ID             uuid.UUID  `json:"id"`
	SubscriptionID string     `json:"subscription_id"`
	Status         string     `json:"status"`
	StartDate      *time.Time `json:"start_date"`
	EndDate        *time.Time `json:"end_date"`
	CreatedAt      time.Time  `json:"created_at"`
}

func (p *Payment) TableName() string {
	return "payments"
}
