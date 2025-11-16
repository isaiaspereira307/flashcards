package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Email        string    `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string    `gorm:"-:all" json:"-"`           // Não serializar
	Plan         string    `gorm:"default:free" json:"plan"` // free, pro, admin
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// Relacionamentos
	Collections       []Collection    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
	SharedCollections []Share         `gorm:"foreignKey:SharedWithID;constraint:OnDelete:CASCADE" json:"-"`
	Payments          []Payment       `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
	GenerationLogs    []GenerationLog `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
}

type UserRegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type UserLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token        string       `json:"token"`
	RefreshToken string       `json:"refresh_token,omitempty"`
	User         UserResponse `json:"user"`
}

type UserResponse struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
	Plan  string    `json:"plan"`
}

func (u *User) BeforeSave(tx *gorm.DB) error {
	// Adicionar lógica de validação pré-save se necessário
	return nil
}

func (u *User) TableName() string {
	return "users"
}
