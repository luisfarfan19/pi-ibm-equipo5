package models

type Health struct {
	Status string `json:"status"`
}

type WatsonAPIResponse struct {
	WatsonResponse map[string]interface{} `bson:"watsonResponse" json:"watsonResponse"`
	ShelveMessage  Shelve                 `bson:"shelveMessage" json:"shelveMessage"`
}

type Shelve map[string]ShelveDetail

type ShelveDetail struct {
	Obs        string `bson:"obs" json:"obs"`
	Percentage string `bson:"porcentaje" json:"porcentaje"`
}

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

type PlanogramWatsonApiResponse struct {
	ShelveCount  int               `json:"estantes"`
	ShelveDetail map[string]string `json:"detalle"`
}

type SectionImageRequest struct {
	Image64     string        `json:"image64"`
	ImageType   string        `json:"imageType"`
	StoreName   string        `json:"storeName"`
	SectionId   string        `json:"sectionId"`
	SectionJson WatsonSection `json:"sectionJson"`
}

type CredentialsRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type StoreName struct {
	Name string `bson:"nombre" json:"nombre"`
}
