package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/isaiaspereira307/flashcards-golang/config"
	"github.com/isaiaspereira307/flashcards-golang/database"
	"github.com/isaiaspereira307/flashcards-golang/models"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandlers struct {
	config *config.Config
}

func NewAuthHandlers(cfg *config.Config) *AuthHandlers {
	return &AuthHandlers{config: cfg}
}

func (h *AuthHandlers) Register(c *gin.Context) {
	var req models.UserRegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "invalid_request",
			"message": err.Error(),
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "password_hash_failed",
		})
		return
	}

	user := models.User{
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Plan:         "free",
	}

	if result := database.DB.Create(&user); result.Error != nil {
		c.JSON(http.StatusConflict, gin.H{
			"success": false,
			"error":   "email_already_exists",
			"message": "Email já registrado",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Usuário registrado com sucesso",
		"data": gin.H{
			"user_id": user.ID,
			"email":   user.Email,
			"plan":    user.Plan,
		},
	})
}

func (h *AuthHandlers) Login(c *gin.Context) {
	var req models.UserLoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "invalid_request",
		})
		return
	}

	var user models.User
	if result := database.DB.Where("email = ?", req.Email).First(&user); result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "invalid_credentials",
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "invalid_credentials",
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID.String(),
		"email":   user.Email,
		"plan":    user.Plan,
		"exp":     time.Now().Add(time.Second * time.Duration(h.config.JWT.Expiration)).Unix(),
	})

	tokenString, err := token.SignedString([]byte(h.config.JWT.Secret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "token_generation_failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"token": tokenString,
			"user": models.UserResponse{
				ID:    user.ID,
				Email: user.Email,
				Plan:  user.Plan,
			},
		},
	})
}

func (h *AuthHandlers) GetMe(c *gin.Context) {
	userID := c.GetString("user_id")

	var user models.User
	if result := database.DB.Where("id = ?", userID).First(&user); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "user_not_found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": models.UserResponse{
			ID:    user.ID,
			Email: user.Email,
			Plan:  user.Plan,
		},
	})
}
