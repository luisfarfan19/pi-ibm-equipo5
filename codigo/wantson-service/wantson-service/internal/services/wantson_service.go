package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"io"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
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
	log.Println("Iniciando validacion de estante", r.Body)

	var req response.SectionImageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON invalido", http.StatusBadRequest)
		return
	}

	image64 := fmt.Sprintf("data:%s;base64,%s", req.ImageType, req.Image64)

	watsonPrompt := buildShelvePrompt(req.SectionJson)
	responseApi, err := callIbmWatsonApi(image64, watsonPrompt, utils.WatsonGraniteModelId)

	if err != nil {
		log.Println("Error during watson api call")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Respuesta del prompt2 [estante]: ", responseApi)
	resp, _ := buildWatsonAPIResponse(responseApi)
	log.Println("Finalizando validacion de estante con respuesta: ", resp)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func buildShelvePrompt(sectionJson response.WatsonSection) string {
	shelveCount := sectionJson.WatsonPromptResponse.ShelveCount

	// Ordenamos los estantes por clave
	var keys []string
	for k := range sectionJson.WatsonPromptResponse.ShelveDetail {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var descParts []string
	for _, key := range keys {
		product := sectionJson.WatsonPromptResponse.ShelveDetail[key]
		desc := fmt.Sprintf("En el %s deberan estar los productos %s", key, product)
		descParts = append(descParts, desc)
	}

	descShelves := strings.Join(descParts, ", ")
	result := fmt.Sprintf("%d niveles de estantes de los cuales: %s", shelveCount, descShelves)
	result = fmt.Sprintf(strings.Replace(utils.WatsonShelvePrompt, "{parsedResponse}", result, 1))
	return result
}

func ValidatePlanogramHandler(w http.ResponseWriter, r *http.Request) {
	var req response.SectionImageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	log.Println("section: ", req.SectionId)
	log.Println("store: ", req.StoreName)

	log.Println("Iniciando validacion del planograma seccion y tienda: ", req.SectionId, req.StoreName)

	image64 := fmt.Sprintf("data:%s;base64,%s", req.ImageType, req.Image64)

	responseApi, err := callIbmWatsonApi(image64, utils.WatsonPlanogramPrompt, utils.WatsonLlamaModelId)
	if err != nil {
		log.Println("Error during watson api call")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("Respuesta del prompt1 [planograma]: ", responseApi)

	resp := saveWatsonResponse(responseApi, req.StoreName, req.SectionId)
	log.Println("Finalizando validacion del planograma con respuesta: ", resp)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func GetIbmAccessTokenHandler(w http.ResponseWriter, r *http.Request) {
	token, err := getIbmAccessToken()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"access_token": token,
	})
}

// function created with Watson code assistant
func getIbmAccessToken() (string, error) {
	// Prepare form data
	form := url.Values{}
	form.Add("grant_type", utils.TokenGrantType)
	form.Add("apikey", utils.TokenApiKey)

	// Prepare request
	req, err := http.NewRequest("POST", utils.TokenApiUrl, bytes.NewBufferString(form.Encode()))
	if err != nil {
		return "", fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("sending request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("non-200 response: %d, body: %s", resp.StatusCode, body)
	}

	// Parse JSON
	var result struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("parsing JSON: %w", err)
	}

	return result.AccessToken, nil
}

func callIbmWatsonApi(imageBase64 string, watsonPrompt string, watsonModel string) (map[string]interface{}, error) {
	log.Println("Llamando al api de watson con el siguiente modelo y prompt: ", watsonModel, watsonPrompt)
	payload := response.RequestPayload{
		Messages: []response.Message{
			{
				Role: "user",
				Content: []response.ContentItem{
					{
						Type: "text",
						Text: watsonPrompt,
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
		ModelID:          watsonModel,
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

	bearerToken, err := getIbmAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get IbmAccessToken: %w", err)
	}

	req, err := http.NewRequest("POST", utils.WatsonUrl, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+bearerToken)

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

	// Parse message.content (which is a JSON string) into Shelve struct
	var parsedContent response.Shelve
	if err := json.Unmarshal([]byte(contentStr), &parsedContent); err != nil {
		return nil, fmt.Errorf("failed to parse message.content: %w", err)
	}

	return &response.WatsonAPIResponse{
		WatsonResponse: responseApi,
		ShelveMessage:  parsedContent,
	}, nil
}

func saveWatsonResponse(apiResponse map[string]interface{}, storeName, idSection string) bson.M {

	choices, ok := apiResponse["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		log.Println("formato inválido de la respuesta de watson")
		return nil
	}

	choice, ok := choices[0].(map[string]interface{})
	if !ok {
		log.Println("formato inválido de la respuesta de watson")
		return nil
	}

	message, ok := choice["message"].(map[string]interface{})
	if !ok {
		log.Println("formato inválido de la respuesta de watson")
		return nil
	}

	contentStr, ok := message["content"].(string)
	if !ok {
		log.Println("el campo 'content' no es string")
		return nil
	}
	var parsed response.PlanogramWatsonApiResponse
	log.Println(contentStr)
	err := json.Unmarshal([]byte(contentStr), &parsed)
	if err != nil {
		log.Println("ignorando respuesta de watson - el campo 'content' no es un JSON válido: %v", err)
		return nil
	}

	return SaveWatsonResponseToMongo(parsed, storeName, idSection)
}
