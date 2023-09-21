package initiializer

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func Getenv() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("error loading env file")
		os.Exit(1)
	}

}
