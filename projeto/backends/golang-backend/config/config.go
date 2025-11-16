package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Database DatabaseConfig `mapstructure:"database"`
	Server   ServerConfig   `mapstructure:"server"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	OpenAI   OpenAIConfig   `mapstructure:"openai"`
	CORS     CORSConfig     `mapstructure:"cors"`
}

type DatabaseConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	Name         string `mapstructure:"name"`
	SSLMode      string `mapstructure:"sslmode"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

type ServerConfig struct {
	Port           string `mapstructure:"port"`
	Environment    string `mapstructure:"environment"`
	RequestTimeout int    `mapstructure:"request_timeout"`
}

type JWTConfig struct {
	Secret            string `mapstructure:"secret"`
	Expiration        int    `mapstructure:"expiration"`
	RefreshExpiration int    `mapstructure:"refresh_expiration"`
}

type OpenAIConfig struct {
	APIKey string `mapstructure:"api_key"`
	Model  string `mapstructure:"model"`
}

type CORSConfig struct {
	AllowedOrigins []string `mapstructure:"allowed_origins"`
	AllowedMethods []string `mapstructure:"allowed_methods"`
	AllowedHeaders []string `mapstructure:"allowed_headers"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/flashcards/")

	viper.AutomaticEnv()
	viper.SetEnvPrefix("FLASHCARDS")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("erro ao ler arquivo de configuração: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("erro ao desserializar configuração: %w", err)
	}

	return &cfg, nil
}

func (d *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.Name, d.SSLMode,
	)
}

func (s *ServerConfig) IsDevelopment() bool {
	return s.Environment == "development"
}
