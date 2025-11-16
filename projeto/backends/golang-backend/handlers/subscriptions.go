package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/isaiaspereira307/flashcards-golang/config"
	"github.com/isaiaspereira307/flashcards-golang/models"
	"gorm.io/gorm"
)

type SubscriptionsHandlers struct {
	db  *gorm.DB
	cfg *config.Config
}

func NewSubscriptionsHandlers(db *gorm.DB, cfg *config.Config) *SubscriptionsHandlers {
	return &SubscriptionsHandlers{db: db, cfg: cfg}
}

// GetCurrent - GET /subscriptions/current
func (h *SubscriptionsHandlers) GetCurrent(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	var user models.User
	if err := h.db.Where("id = ?", userID.(uuid.UUID)).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuário não encontrado"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"plan": user.Plan,
	})
}

// Upgrade - POST /subscriptions/upgrade
func (h *SubscriptionsHandlers) Upgrade(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	var req struct {
		Plan      string `json:"plan" binding:"required"`
		Months    int    `json:"months" binding:"required"`
		PaymentID string `json:"payment_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	var user models.User
	if err := h.db.Where("id = ?", userID.(uuid.UUID)).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuário não encontrado"})
		return
	}

	startDate := time.Now()
	endDate := startDate.AddDate(0, req.Months, 0)

	payment := models.Payment{
		ID:             uuid.New(),
		UserID:         userID.(uuid.UUID),
		SubscriptionID: req.PaymentID,
		Status:         "active",
		StartDate:      &startDate,
		EndDate:        &endDate,
	}

	if err := h.db.Create(&payment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar subscription"})
		return
	}

	if err := h.db.Model(&user).Update("plan", req.Plan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar plano"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"payment": payment,
		"user":    user,
	})
}

// Cancel - POST /subscriptions/cancel
func (h *SubscriptionsHandlers) Cancel(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	var payment models.Payment
	if err := h.db.Where("user_id = ? AND status = ?", userID.(uuid.UUID), "active").
		Order("end_date DESC").
		First(&payment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Nenhuma subscription ativa"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar subscription"})
		return
	}

	if err := h.db.Model(&payment).Update("status", "canceled").Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao cancelar"})
		return
	}

	var user models.User
	if err := h.db.Where("id = ?", userID.(uuid.UUID)).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar usuário"})
		return
	}

	if err := h.db.Model(&user).Update("plan", "free").Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar plano"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Subscription cancelada",
		"payment": payment,
	})
}

// GetHistory - GET /subscriptions/history
func (h *SubscriptionsHandlers) GetHistory(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	var payments []models.Payment
	if err := h.db.Where("user_id = ?", userID.(uuid.UUID)).
		Order("created_at DESC").
		Find(&payments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao listar"})
		return
	}

	c.JSON(http.StatusOK, payments)
}
