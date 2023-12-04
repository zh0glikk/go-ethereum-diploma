package configger

import (
	"errors"
	"github.com/ethereum/go-ethereum/log"
	"os"

	"github.com/joho/godotenv"
)

var envPath = ""

func Init(path string) {
	envPath = path
}

func Get(key string) string {
	if _, err := os.Stat(envPath); !errors.Is(err, os.ErrNotExist) {
		err := godotenv.Load(envPath)
		if err != nil {
			log.Info("failed to load sando.env")
		}
	}

	return os.Getenv(key)
}
