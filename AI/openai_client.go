package ai

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	openai "github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

var openaiClient *openai.Client

func init() {
	err := godotenv.Load("secrets.env")
	if err != nil {
		fmt.Println("❌ Error loading .env file in init openai client:", err)
		return
	}

	OpenaiAPIKey := os.Getenv("OpenaiAPIKey")
	client := openai.NewClient(option.WithAPIKey(OpenaiAPIKey))
	openaiClient = &client
}
