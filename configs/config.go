package configs

type Config struct {
	ENV    string `envconfig:"ENV"`
	DBHost string `envconfig:"DB_HOST"`
	DBPort string `envconfig:"DB_PORT"`
	DBUser string `envconfig:"DB_USER"`
	DBPass string `envconfig:"DB_PASS"`
	DBName string `envconfig:"DB_NAME"`
}

func (c Config) IsProd() bool {
	return c.ENV == "prod"
}
