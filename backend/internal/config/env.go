package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() error {
	if err := godotenv.Load(".env"); err != nil {
		if !os.IsNotExist(err) {
			fmt.Println("❌ Failed to load .env:", err)
			return err
		}

		fmt.Println("⚠️ .env file not found; using process environment")
		return nil
	}

	fmt.Println("✅ .env loaded successfully")

	return nil
}
