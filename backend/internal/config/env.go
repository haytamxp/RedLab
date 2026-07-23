package config

import (
    "fmt"
    "github.com/joho/godotenv"
)

func LoadEnv() error {
    err := godotenv.Load(".env")
    if err != nil {
        fmt.Println("❌ Failed to load .env:", err)
        return err
    }

    fmt.Println("✅ .env loaded successfully")
    return nil
}