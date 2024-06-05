package core

import (
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

var _ = godotenv.Load()

var DBConfig map[string]string = map[string]string{
	"HOST":     os.Getenv("DB_HOST"),
	"PORT":     os.Getenv("DB_PORT"),
	"USER":     os.Getenv("DB_USER"),
	"PASSWORD": os.Getenv("DB_PASSWORD"),
	"NAME":     os.Getenv("DATABASE"),
	"SSL_MODE": "disable",
}

var JWTConfig map[string]string = map[string]string{
	"ACCESS_KEY":  os.Getenv("JWT_ACCESS_KEY"),
	"REFRESH_KEY": os.Getenv("JWT_REFRESH_KEY"),
}

var CORSConfig map[string]any = map[string]any{
	"ALLOWED_ORIGIN": []string{
		"*",
	},
	"ALLOW_CREDENTIALS": "true",
	"ALLOWED_HEADERS": []string{
		"Content-Type", "Content-Length", "Accept-Encoding",
		"X-CSRF-Token", "Authorization", "accept", "origin", "Cache-Control",
		"X-Requested-With",
	},
	"ALLOWED_METHODS": []string{
		"POST", "OPTIONS", "GET", "PUT", "PATCH",
	},
}

var serverPort, _ = strconv.Atoi(os.Getenv("SERVER_PORT"))
var ServerConfig map[string]any = map[string]any{
	"PORT": serverPort,
}
