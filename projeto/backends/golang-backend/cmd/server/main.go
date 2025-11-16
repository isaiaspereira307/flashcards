package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/isaiaspereira307/flashcards-golang/config"
	"github.com/isaiaspereira307/flashcards-golang/database"
	"github.com/isaiaspereira307/flashcards-golang/handlers"
	"github.com/isaiaspereira307/flashcards-golang/middleware"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("❌ Erro ao carregar configurações: %v", err)
	}

	if err := database.Init(cfg); err != nil {
		log.Fatalf("❌ Erro ao inicializar banco de dados: %v", err)
	}
	defer database.Close()

	if err := database.RunMigrations(); err != nil {
		log.Fatalf("❌ Erro ao executar migrações: %v", err)
	}

	router := gin.Default()

	router.Use(middleware.CORSMiddleware(cfg))

	// ============================================
	// ROTAS PÚBLICAS
	// ============================================

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "flashcards-golang",
			"version": "1.0.0",
		})
	})

	authHandlers := handlers.NewAuthHandlers(cfg)
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/register", authHandlers.Register)
		authGroup.POST("/login", authHandlers.Login)
	}

	// ============================================
	// ROTAS PROTEGIDAS
	// ============================================

	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware(cfg))
	{
		protected.GET("/auth/me", authHandlers.GetMe)

		collectionsHandlers := handlers.NewCollectionsHandlers(database.DB, cfg)
		collectionsGroup := protected.Group("/collections")
		{
			collectionsGroup.GET("", collectionsHandlers.List)
			collectionsGroup.POST("", collectionsHandlers.Create)
			collectionsGroup.GET("/:id", collectionsHandlers.GetByID)
			collectionsGroup.PUT("/:id", collectionsHandlers.Update)
			collectionsGroup.DELETE("/:id", collectionsHandlers.Delete)

			sharesHandlers := handlers.NewSharesHandlers(database.DB, cfg)
			sharesGroup := collectionsGroup.Group("/:id/shares")
			{
				sharesGroup.GET("", sharesHandlers.ListCollectionShares)
				sharesGroup.POST("", sharesHandlers.Create)
				sharesGroup.PUT("/:shareID", sharesHandlers.UpdatePermissions)
				sharesGroup.DELETE("/:shareID", sharesHandlers.Delete)
			}

			flashcardsHandlers := handlers.NewFlashcardsHandlers(database.DB, cfg)
			flashcardsGroup := collectionsGroup.Group("/:id/flashcards")
			{
				flashcardsGroup.GET("", flashcardsHandlers.List)
				flashcardsGroup.POST("", flashcardsHandlers.Create)
				flashcardsGroup.GET("/:cardID", flashcardsHandlers.GetByID)
				flashcardsGroup.PUT("/:cardID", flashcardsHandlers.Update)
				flashcardsGroup.DELETE("/:cardID", flashcardsHandlers.Delete)
			}
		}

		sharesHandlers := handlers.NewSharesHandlers(database.DB, cfg)
		protected.GET("/shares", sharesHandlers.ListShared)

		subsHandlers := handlers.NewSubscriptionsHandlers(database.DB, cfg)
		subsGroup := protected.Group("/subscriptions")
		{
			subsGroup.GET("/current", subsHandlers.GetCurrent)
			subsGroup.POST("/upgrade", subsHandlers.Upgrade)
			subsGroup.POST("/cancel", subsHandlers.Cancel)
			subsGroup.GET("/history", subsHandlers.GetHistory)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = cfg.Server.Port
	}

	log.Printf("Servidor iniciado em http://localhost:%s", port)
	log.Printf("Banco de dados: %s@%s:%d/%s", cfg.Database.User, cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)
	log.Printf("Environment: %s", cfg.Server.Environment)

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v", err)
	}
}
