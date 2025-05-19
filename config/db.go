package config

type DBConfig struct {
	Host     string
	Port     string
	Database string
	Username string
	Password string
	Driver   string
}

type APIConfig struct {
	ApiPort string
}

type TokenConfig struct {
	ApplicationName     string
	JwtSignatureKey     string
	JwtSigningMethod    string
	AccessTokenLifeTime int // dalam jam
}

type Config struct {
	DBConfig
	APIConfig
	TokenConfig
	LocationIQAPIKey  string
	MidtransServerKey string
}
