package services

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	response "wantson-service/internal/models"
	"wantson-service/pkg/utils"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("App is running...")
	message := response.Health{Status: "UP"}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(message)
}

func ValidateShelterHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Validating shelter...")
	dataUrl := getImage()
	if dataUrl == "" {
		log.Println("Unable to retrieve shelter image.")
		http.Error(w, "Unable to read image", http.StatusInternalServerError)
		return
	}
	responseApi, err := callIBMWatsonVisionAPI(dataUrl)
	if err != nil {
		log.Println("Error during watson api call")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp, _ := buildWatsonAPIResponse(responseApi)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func getImage() string {
	imagePath := utils.Image
	data, err := os.ReadFile(imagePath)
	if err != nil {
		log.Println(fmt.Errorf("Error reading image", err))
		return ""
	}
	base64Encoding := base64.StdEncoding.EncodeToString(data)
	mimeType := "image/jpeg"
	dataURL := fmt.Sprintf("data:%s;base64,%s", mimeType, base64Encoding)
	return dataURL
}

func callIBMWatsonVisionAPI(imageBase64 string) (map[string]interface{}, error) {
	payload := response.RequestPayload{
		Messages: []response.Message{
			{
				Role: "user",
				Content: []response.ContentItem{
					{
						Type: "text",
						Text: utils.WatsonPrompt,
					},
					{
						Type: "image_url",
						ImageURL: &response.ImageURL{
							URL: imageBase64,
						},
					},
				},
			},
		},
		ProjectID:        utils.WatsonProjectId,
		ModelID:          utils.WatsonModelId,
		FrequencyPenalty: 0,
		MaxTokens:        utils.WatsonMaxTokens,
		PresencePenalty:  0,
		Temperature:      0,
		TopP:             1,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("POST", utils.WatsonUrl, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", utils.WatsonBearerToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to call Watson API: %v", err)
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to to read response from Watson API: %v", err)
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if resp.StatusCode != 200 {
		log.Printf("Failed to call Watson API: %d %s", resp.StatusCode, string(respBody))
		return result, fmt.Errorf("response status code is not 200: %d", resp.StatusCode)
	}

	return result, nil
}

func buildWatsonAPIResponse(responseApi map[string]interface{}) (*response.WatsonAPIResponse, error) {
	// Extract the message.content string
	choices, ok := responseApi["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return nil, fmt.Errorf("choices array missing or invalid")
	}

	choice := choices[0].(map[string]interface{})
	message, ok := choice["message"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("message object missing in choice")
	}

	contentStr, ok := message["content"].(string)
	if !ok {
		return nil, fmt.Errorf("message.content is not a string")
	}

	// Parse message.content (which is a JSON string) into MessageContent struct
	var parsedContent response.MessageContent
	if err := json.Unmarshal([]byte(contentStr), &parsedContent); err != nil {
		return nil, fmt.Errorf("failed to parse message.content: %w", err)
	}

	return &response.WatsonAPIResponse{
		WatsonResponse: responseApi,
		Message:        parsedContent,
	}, nil
}
