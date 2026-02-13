package model

import (
	"os"
	"strings"
)

type Env struct {
	MongoURI      string
	Database      string
	JWTSecret     string
	Port          string
	AdminEmail    string
	AdminPassword string
}

func getEnv(key, def string) string {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return def
	}
	return v
}

func LoadEnv() Env {
	return Env{
		MongoURI:      getEnv("MONGO_URI", "mongodb://127.0.0.1:27017/bazarPO"),
		Database:      getEnv("MONGO_DB", "bazarPO"),
		JWTSecret:     getEnv("JWT_SECRET", "change-me"),
		Port:          getEnv("PORT", "8090"),
		AdminEmail:    strings.ToLower(getEnv("ADMIN_EMAIL", "admin@cps.local")),
		AdminPassword: getEnv("ADMIN_PASSWORD", "admin12345"),
	}
}
