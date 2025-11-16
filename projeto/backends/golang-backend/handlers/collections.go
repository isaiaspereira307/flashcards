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

type CollectionsHandlers struct {
	db  *gorm.DB
	cfg *config.Config
}

func NewCollectionsHandlers(db *gorm.DB, cfg *config.Config) *CollectionsHandlers {
	return &CollectionsHandlers{db: db, cfg: cfg}
}

func (h *CollectionsHandlers) Create(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	var req struct {
		Name     string `json:"name" binding:"required,min=1,max=255"`
		IsPublic bool   `json:"is_public"`
		MaxCards int    `json:"max_cards"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos: " + err.Error()})
		return
	}

	collection := models.Collection{
		ID:       uuid.New(),
		UserID:   userID.(uuid.UUID),
		Name:     req.Name,
		IsPublic: req.IsPublic,
		MaxCards: req.MaxCards,
	}

	if err := h.db.Create(&collection).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar coleção: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, collection)
}

func (h *CollectionsHandlers) List(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	var collections []models.Collection
	if err := h.db.Where("user_id = ?", userID.(uuid.UUID)).
		Preload("Flashcards").
		Find(&collections).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao listar coleções: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, collections)
}

func (h *CollectionsHandlers) GetByID(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	collectionID := c.Param("id")
	var collection models.Collection

	if err := h.db.Where("id = ? AND user_id = ?", collectionID, userID.(uuid.UUID)).
		Preload("Flashcards").
		First(&collection).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Coleção não encontrada"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar coleção: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, collection)
}

func (h *CollectionsHandlers) Update(c *gin.Context) {
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

	var req struct {
		Name     string `json:"name"`
		IsPublic *bool  `json:"is_public"`
		MaxCards *int   `json:"max_cards"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos: " + err.Error()})
		return
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.IsPublic != nil {
		updates["is_public"] = *req.IsPublic
	}
	if req.MaxCards != nil {
		updates["max_cards"] = *req.MaxCards
	}

	if err := h.db.Model(&collection).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar coleção: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, collection)
}

// Delete - DELETE /collections/:id
func (h *CollectionsHandlers) Delete(c *gin.Context) {
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

	if err := h.db.Delete(&collection).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao deletar coleção: " + err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
