package models

type User struct {
	Username string `bson:"username"`
	Password string `bson:"password"`
	Rol      string `bson:"rol"`
}

type Section struct {
	SectionName  string `bson:"nombreSeccion" json:"nombreSeccion"`
	PlanogramImg string `bson:"imgPlanograma" json:"imgPlanograma"`
	IdSection    string `bson:"idSeccion" json:"idSeccion"`
}

type Aisle struct {
	AisleName string    `bson:"nombrePasillo" json:"nombrePasillo"`
	Section   []Section `bson:"secciones" json:"secciones"`
}

type Data struct {
	Address string  `bson:"direccion" json:"direccion"`
	Aisle   []Aisle `bson:"pasillos" json:"pasillos"`
}

type Store struct {
	Name string `bson:"nombre" json:"nombre"`
	Data Data   `bson:"data" json:"data"`
}

type WatsonSection struct {
	SectionId            string         `bson:"idSeccion" json:"idSeccion"`
	WatsonPromptResponse WatsonResponse `bson:"watsonPromptResponse" json:"watsonPromptResponse"`
}

type WatsonResponse struct {
	ShelveCount  int               `bson:"estantes" json:"estantes"`
	ShelveDetail map[string]string `bson:"detalle" json:"detalle"`
}

type PlanogramResponse struct {
	StoreName string          `bson:"storeName" json:"storeName"`
	Sections  []WatsonSection `bson:"secciones" json:"secciones"`
}
