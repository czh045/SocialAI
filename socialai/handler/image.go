package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	openAIImageEndpoint  = "https://api.openai.com/v1/images/generations"
	maxImagePromptLength = 4000
)

type imageGenerationRequest struct {
	Prompt string `json:"prompt"`
}

type openAIImageRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Size   string `json:"size"`
}

type openAIImageResponse struct {
	Data []struct {
		Base64JSON string `json:"b64_json"`
	} `json:"data"`
}

// generateImageHandler keeps the OpenAI API key on the Go server. The browser
// sends only a prompt and receives a data URL that can later be uploaded as a post.
func generateImageHandler(w http.ResponseWriter, r *http.Request) {
	apiKey := strings.TrimSpace(os.Getenv("OPENAI_API_KEY"))
	if apiKey == "" {
		http.Error(w, "AI image generation is not configured on this server", http.StatusServiceUnavailable)
		return
	}

	var input imageGenerationRequest
	if err := json.NewDecoder(io.LimitReader(r.Body, 64*1024)).Decode(&input); err != nil {
		http.Error(w, "Cannot decode image generation request", http.StatusBadRequest)
		return
	}

	input.Prompt = strings.TrimSpace(input.Prompt)
	if input.Prompt == "" || len([]rune(input.Prompt)) > maxImagePromptLength {
		http.Error(w, "Prompt must contain 1 to 4000 characters", http.StatusBadRequest)
		return
	}

	model := strings.TrimSpace(os.Getenv("OPENAI_IMAGE_MODEL"))
	if model == "" {
		model = "gpt-image-2"
	}

	payload, err := json.Marshal(openAIImageRequest{
		Model:  model,
		Prompt: input.Prompt,
		Size:   "1024x1024",
	})
	if err != nil {
		http.Error(w, "Could not prepare image generation request", http.StatusInternalServerError)
		return
	}

	request, err := http.NewRequestWithContext(
		r.Context(),
		http.MethodPost,
		openAIImageEndpoint,
		bytes.NewReader(payload),
	)
	if err != nil {
		http.Error(w, "Could not prepare image generation request", http.StatusInternalServerError)
		return
	}
	request.Header.Set("Authorization", "Bearer "+apiKey)
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 90 * time.Second}
	response, err := client.Do(request)
	if err != nil {
		http.Error(w, "Image generation service is unavailable", http.StatusBadGateway)
		return
	}
	defer response.Body.Close()

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		body, _ := io.ReadAll(io.LimitReader(response.Body, 4096))
		fmt.Printf("OpenAI image generation failed with status %d: %s\n", response.StatusCode, string(body))
		http.Error(w, "Image generation failed. Check the server API key, model access, and billing.", http.StatusBadGateway)
		return
	}

	var generated openAIImageResponse
	if err := json.NewDecoder(response.Body).Decode(&generated); err != nil {
		http.Error(w, "Could not read the image generation response", http.StatusBadGateway)
		return
	}
	if len(generated.Data) == 0 || generated.Data[0].Base64JSON == "" {
		http.Error(w, "The image service did not return image data", http.StatusBadGateway)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"imageDataUrl": "data:image/png;base64," + generated.Data[0].Base64JSON,
	})
}
