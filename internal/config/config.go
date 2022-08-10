package config

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

func LogLevel() string {
	return os.Getenv("LOG_LEVEL")
}

func Token() string {
	cfg := os.Getenv("TOKEN")
	if cfg == "" {
		logrus.Fatal("token value is empty")
	}
	return cfg
}

func PostgresDSN() string {
	host := os.Getenv("POSTGRES_HOST")
	db := os.Getenv("POSTGRES_DATABASE")
	user := os.Getenv("POSTGRES_USER")
	pw := os.Getenv("POSTGRES_PASSWORD")
	port := os.Getenv("POSTGRES_PORT")

	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, pw, db, port)
}

// PrivateKey get private key from env
func PrivateKey() string {
	key := os.Getenv("PRIVATE_KEY")
	if key == "" {
		logrus.Error("PRIVATE_KEY is unset. May cause danger in encryption method")
	}

	return key
}

// IvKey get private key from env
func IvKey() string {
	key := os.Getenv("IV_KEY")
	if key == "" {
		logrus.Error("IV_KEY is unset. May cause danger in encryption method")
	}

	return key
}

func HimatroAPIHost() string {
	return os.Getenv("HIMATRO_API_HOST")
}
