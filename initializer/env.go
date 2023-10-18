package initializer

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// LoadEnvironmentVariables loads environment variables from a .env file.
func LoadEnvironmentVariables() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading the .env file:", err)
		os.Exit(1)
	}
}
