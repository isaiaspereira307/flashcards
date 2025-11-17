package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/isaiaspereira307/flashcards-golang/config"
	"github.com/isaiaspereira307/flashcards-golang/models"
	"gorm.io/gorm"
)

type FlashcardsHandlers struct {
	db  *gorm.DB
	cfg *config.Config
}

func NewFlashcardsHandlers(db *gorm.DB, cfg *config.Config) *FlashcardsHandlers {
	return &FlashcardsHandlers{db: db, cfg: cfg}
}

// Create - POST /collections/:id/flashcards
func (h *FlashcardsHandlers) Create(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	collectionID := c.Param("id")
	var collection models.Collection

	if err := h.db.Where("id = ? AND user_id = ?", collectionID, userID.(uuid.UUID)).
		First(&collection).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Coleção não encontrada"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar coleção: " + err.Error()})
		return
	}

	var req models.FlashcardRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos: " + err.Error()})
		return
	}

	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Link inválido: " + err.Error()})
		return
	}

	flashcard := models.Flashcard{
		ID:           uuid.New(),
		CollectionID: uuid.MustParse(collectionID),
		Front:        req.Front,
		Back:         req.Back,
		VideoURL:     req.VideoURL,
		CreatedByIA:  false,
	}

	if err := h.db.Create(&flashcard).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar flashcard: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, flashcard)
}

// List - GET /collections/:id/flashcards
func (h *FlashcardsHandlers) List(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	collectionID := c.Param("id")
	var collection models.Collection

	if err := h.db.Where("id = ? AND user_id = ?", collectionID, userID.(uuid.UUID)).
		First(&collection).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Coleção não encontrada"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar coleção: " + err.Error()})
		return
	}

	var flashcards []models.Flashcard
	if err := h.db.Where("collection_id = ?", collectionID).
		Order("created_at DESC").
		Find(&flashcards).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao listar flashcards: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, flashcards)
}

// GetByID - GET /collections/:id/flashcards/:cardID
func (h *FlashcardsHandlers) GetByID(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	collectionID := c.Param("id")
	cardID := c.Param("cardID")
	var flashcard models.Flashcard

	if err := h.db.Joins("JOIN collections ON flashcards.collection_id = collections.id").
		Where("flashcards.id = ? AND collections.user_id = ? AND flashcards.collection_id = ?", cardID, userID.(uuid.UUID), collectionID).
		First(&flashcard).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Flashcard não encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar flashcard: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, flashcard)
}

// Update - PUT /collections/:id/flashcards/:cardID
func (h *FlashcardsHandlers) Update(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	collectionID := c.Param("id")
	cardID := c.Param("cardID")
	var flashcard models.Flashcard

	if err := h.db.Joins("JOIN collections ON flashcards.collection_id = collections.id").
		Where("flashcards.id = ? AND collections.user_id = ? AND flashcards.collection_id = ?", cardID, userID.(uuid.UUID), collectionID).
		First(&flashcard).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Flashcard não encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar flashcard: " + err.Error()})
		return
	}

	var req struct {
		Front string                 `json:"front"`
		Back  string                 `json:"back"`
		Extra map[string]interface{} `json:"extra"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos: " + err.Error()})
		return
	}

	updates := map[string]interface{}{}
	if req.Front != "" {
		updates["front"] = req.Front
	}
	if req.Back != "" {
		updates["back"] = req.Back
	}

	if err := h.db.Model(&flashcard).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar flashcard: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, flashcard)
}

// Delete - DELETE /collections/:id/flashcards/:cardID
func (h *FlashcardsHandlers) Delete(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	collectionID := c.Param("id")
	cardID := c.Param("cardID")
	var flashcard models.Flashcard

	if err := h.db.Joins("JOIN collections ON flashcards.collection_id = collections.id").
		Where("flashcards.id = ? AND collections.user_id = ? AND flashcards.collection_id = ?", cardID, userID.(uuid.UUID), collectionID).
		First(&flashcard).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Flashcard não encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar flashcard: " + err.Error()})
		return
	}

	if err := h.db.Delete(&flashcard).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao deletar flashcard: " + err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
