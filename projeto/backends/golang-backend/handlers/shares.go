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

type SharesHandlers struct {
	db  *gorm.DB
	cfg *config.Config
}

func NewSharesHandlers(db *gorm.DB, cfg *config.Config) *SharesHandlers {
	return &SharesHandlers{db: db, cfg: cfg}
}

// Create - POST /collections/:id/shares
func (h *SharesHandlers) Create(c *gin.Context) {
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar coleção"})
		return
	}

	var req struct {
		SharedWithID string `json:"shared_with_id" binding:"required"`
		Permissions  string `json:"permissions" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	sharedWithID, err := uuid.Parse(req.SharedWithID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de usuário inválido"})
		return
	}

	share := models.Share{
		ID:           uuid.New(),
		CollectionID: uuid.MustParse(collectionID),
		SharedWithID: sharedWithID,
		Permissions:  req.Permissions,
		ShareID:      uuid.New().String()[:8],
	}

	if err := h.db.Create(&share).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao compartilhar"})
		return
	}

	c.JSON(http.StatusCreated, share)
}

// ListCollectionShares - GET /collections/:id/shares
func (h *SharesHandlers) ListCollectionShares(c *gin.Context) {
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar coleção"})
		return
	}

	var shares []models.Share
	if err := h.db.Where("collection_id = ?", collectionID).Find(&shares).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao listar"})
		return
	}

	c.JSON(http.StatusOK, shares)
}

// ListShared - GET /shares
func (h *SharesHandlers) ListShared(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	var shares []models.Share
	if err := h.db.Where("shared_with_id = ?", userID.(uuid.UUID)).
		Preload("Collection").
		Find(&shares).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao listar"})
		return
	}

	c.JSON(http.StatusOK, shares)
}

// Delete - DELETE /collections/:collectionID/shares/:shareID
func (h *SharesHandlers) Delete(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	collectionID := c.Param("collectionID")
	shareID := c.Param("shareID")

	var share models.Share
	if err := h.db.Joins("JOIN collections ON shares.collection_id = collections.id").
		Where("shares.id = ? AND collections.user_id = ? AND shares.collection_id = ?", shareID, userID.(uuid.UUID), collectionID).
		First(&share).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Compartilhamento não encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar"})
		return
	}

	if err := h.db.Delete(&share).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao deletar"})
		return
	}

	c.Status(http.StatusNoContent)
}

// UpdatePermissions - PUT /collections/:collectionID/shares/:shareID
func (h *SharesHandlers) UpdatePermissions(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	collectionID := c.Param("collectionID")
	shareID := c.Param("shareID")

	var share models.Share
	if err := h.db.Joins("JOIN collections ON shares.collection_id = collections.id").
		Where("shares.id = ? AND collections.user_id = ? AND shares.collection_id = ?", shareID, userID.(uuid.UUID), collectionID).
		First(&share).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Compartilhamento não encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar"})
		return
	}

	var req struct {
		Permissions string `json:"permissions" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	if err := h.db.Model(&share).Update("permissions", req.Permissions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar"})
		return
	}

	c.JSON(http.StatusOK, share)
}
