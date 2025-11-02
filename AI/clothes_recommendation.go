package ai

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	openai "github.com/openai/openai-go"
	"github.com/openai/openai-go/responses"
)

func GetClothingRecommendation(input string) (string, error) {
	err := godotenv.Load("secrets.env")
	if err != nil {
		fmt.Println("❌ Error loading .env file in getClothingRecommendation:", err)
		return "", err
	}
	promptID := os.Getenv("ClothRecomendationPromptID")

	ctx := context.Background()
	aiResp, err := openaiClient.Responses.New(ctx, responses.ResponseNewParams{
		Prompt: responses.ResponsePromptParam{
			ID: promptID,
			Variables: map[string]responses.ResponsePromptVariableUnionParam{
				"weather": {
					OfString: openai.String(input),
				},
			},
		},
	})

	if err != nil {
		fmt.Println("❌ Error getting AI response:", err)
		return "", err
	}
	fmt.Println("✅ AI is thinking....")
	fmt.Println("🤖Clothing Recommendation:")
	fmt.Println(aiResp.OutputText())

	return aiResp.OutputText(), nil
}
