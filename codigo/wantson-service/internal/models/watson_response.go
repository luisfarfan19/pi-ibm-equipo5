package models

type Health struct {
	Status string `json:"status"`
}

type WatsonAPIResponse struct {
	WatsonResponse map[string]interface{} `json:"watson_response"`
	Message        MessageContent         `json:"message"`
}

type MessageContent struct {
	Analysis string         `json:"analysis"`
	Scores   map[string]int `json:"scores"`
	Alerts   []string       `json:"alerts"`
}

// Used for watson API
type ImageURL struct {
	URL string `json:"url"`
}

type ContentItem struct {
	Type     string    `json:"type"`
	Text     string    `json:"text,omitempty"`
	ImageURL *ImageURL `json:"image_url,omitempty"`
}

type Message struct {
	Role    string        `json:"role"`
	Content []ContentItem `json:"content"`
}

type RequestPayload struct {
	Messages         []Message `json:"messages"`
	ProjectID        string    `json:"project_id"`
	ModelID          string    `json:"model_id"`
	FrequencyPenalty int       `json:"frequency_penalty"`
	MaxTokens        int       `json:"max_tokens"`
	PresencePenalty  int       `json:"presence_penalty"`
	Temperature      int       `json:"temperature"`
	TopP             int       `json:"top_p"`
}
