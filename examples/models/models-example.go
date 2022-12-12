package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/TannerKvarfordt/gopenai/authentication"
	"github.com/TannerKvarfordt/gopenai/models"
	_ "github.com/joho/godotenv/autoload"
)

const OpenAITokenEnv = "OPENAI_API_KEY"

func init() {
	key := os.Getenv(OpenAITokenEnv)
	authentication.SetAPIKey(key)
}

func main() {
	listresp, err := models.MakeListModelsRequest(nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	if listresp == nil {
		fmt.Println("Nil response returned.")
		return
	}
	if listresp.Error.Type != "" {
		fmt.Printf("Error: %s -> %s\n", listresp.Error.Type, listresp.Error.Message)
		return
	}
	if len(listresp.Data) < 1 {
		fmt.Println("No model information retrieved")
		return
	}

	fmt.Printf("Retreived information for %d models. Calling the \"Retrieve\" endpoint with the first one.\n", len(listresp.Data))
	resp, err := models.MakeRetrieveModelRequest(listresp.Data[0].ID, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	if resp == nil {
		fmt.Println("Nil response returned")
	}

	jsonResp, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		fmt.Printf(`Failed to marshal response to JSON for printing, err="%s"\n`, err)
		return
	}
	fmt.Printf("Retrieved Model Info: %s\n", string(jsonResp))
}
