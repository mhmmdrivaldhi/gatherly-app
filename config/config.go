package config

import (
	"fmt"
	"os"
)

func (c *Config) readConfig() error {
	c.DBConfig = DBConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Database: os.Getenv("DB_DATABASE"),
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		Driver:   os.Getenv("DB_DRIVER"),
	}

	c.APIConfig = APIConfig{
		ApiPort: os.Getenv("API_PORT"),
	}
	// Token Configuration
	c.TokenConfig = TokenConfig{
		ApplicationName:     os.Getenv("APP_NAME"),
		JwtSignatureKey:     os.Getenv("JWT_SIGNATURE_KEY"),
		JwtSigningMethod:    os.Getenv("JWT_SIGNING_METHOD"),
		AccessTokenLifeTime: 24, // Default 24 jam
	}

	c.LocationIQAPIKey = os.Getenv("LOCATIONIQ_API_KEY")

	c.MidtransServerKey = os.Getenv("MIDTRANS_SERVER_KEY")

	// Validasi config wajib
	
	required := map[string]string{
		"DB_HOST":           c.Host,
		"DB_PORT":           c.Port,
		"DB_USERNAME":       c.Username,
		"DB_PASSWORD":       c.Password,
		"API_PORT":          c.ApiPort,
		"JWT_SIGNATURE_KEY": c.JwtSignatureKey,
	}

	for key, val := range required {
		if val == "" {
			return fmt.Errorf("config %s is required", key)
		}
	}

	fmt.Println(os.Getenv("HOST"), os.Getenv("PORT"), os.Getenv("DATABASE"), os.Getenv("USERNAME"), os.Getenv("PASSWORD"), os.Getenv("DRIVER"))
	fmt.Println(os.Getenv("API_PORT"))

	if c.Host == "" || c.Port == "" || c.Username == "" || c.Password == "" || c.ApiPort == "" || c.LocationIQAPIKey == "" || c.TokenConfig.JwtSignatureKey == "" || c.MidtransServerKey == "" {
		return fmt.Errorf("required config")
	}

	return nil
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := cfg.readConfig(); err != nil {
		return nil, err
	}
	return cfg, nil
}
