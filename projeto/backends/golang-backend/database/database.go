package database

import (
	"fmt"
	"log"

	"github.com/isaiaspereira307/flashcards-golang/config"
	"github.com/isaiaspereira307/flashcards-golang/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Init(cfg *config.Config) error {
	dsn := cfg.Database.GetDSN()

	var logLevel logger.LogLevel
	if cfg.Server.IsDevelopment() {
		logLevel = logger.Info
	} else {
		logLevel = logger.Error
	}

	db, err := gorm.Open(
		postgres.Open(dsn),
		&gorm.Config{
			Logger: logger.Default.LogMode(logLevel),
		},
	)
	if err != nil {
		return fmt.Errorf("erro ao conectar ao banco de dados: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("erro ao obter instância SQL: %w", err)
	}

	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)

	DB = db

	log.Println("Conectado ao PostgreSQL com sucesso")
	return nil
}

func RunMigrations() error {
	if DB == nil {
		return fmt.Errorf("banco de dados não inicializado")
	}

	if err := DB.Migrator().AutoMigrate(
		&models.User{},
		&models.Collection{},
		&models.Flashcard{},
		&models.Share{},
		&models.GenerationLog{},
		&models.Payment{},
	); err != nil {
		log.Printf("Aviso ao executar migrações: %v", err)
	}

	log.Println("Migrações verificadas com sucesso")
	return nil
}

func Close() error {
	if DB == nil {
		return nil
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}
