package configs

import "fmt"

type Config struct {
	ENV string `envconfig:"ENV" default:"dev"`

	DBHost  string `envconfig:"DB_HOST"`
	DBPort  string `envconfig:"DB_PORT"`
	DBUser  string `envconfig:"DB_USER"`
	DBPass  string `envconfig:"DB_PASS"`
	DBName  string `envconfig:"DB_NAME"`
	SSLMode string `envconfig:"DB_SSL_MODE" default:"disable"`

	TelegramBotToken string `envconfig:"TELEGRAM_BOT_TOKEN"`

	AllowedChatIDs []int64 `envconfig:"ALLOWED_CHAT_IDS"`
}

func (c Config) IsProd() bool {
	return c.ENV == "prod"
}

func (c Config) IsDev() bool {
	return !c.IsProd()
}

func (c Config) DBDSNString() string {
	if c.SSLMode != "disable" {
		c.SSLMode = "enable"
	}

	return fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		c.DBHost,
		c.DBPort,
		c.DBUser,
		c.DBName,
		c.DBPass,
		c.SSLMode,
	)
}
