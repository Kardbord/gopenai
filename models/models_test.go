package models_test

import (
	"os"
	"testing"

	"github.com/Kardbord/gopenai/authentication"
	"github.com/Kardbord/gopenai/models"
	_ "github.com/joho/godotenv/autoload"
)

const OpenAITokenEnv = "OPENAI_API_KEY"

func init() {
	key := os.Getenv(OpenAITokenEnv)
	authentication.SetAPIKey(key)
}

func TestModels(t *testing.T) {
	listresp, err := models.MakeListModelsRequest(nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(listresp.Data) < 1 {
		t.Fatal("No model information retrieved")
	}

	resp, err := models.MakeRetrieveModelRequest(listresp.Data[0].ID, nil)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("Nil response returned")
	}
}
