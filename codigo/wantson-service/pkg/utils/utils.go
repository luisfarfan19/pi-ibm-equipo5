package utils

const (
	//DB
	MongoDbUrl  = "mongodb://localhost:27017"
	MongoDbName = "shelter_db"
	//WATSON
	TokenGrantType        = "urn:ibm:params:oauth:grant-type:apikey"
	TokenApiKey           = "sOxpcgVDQB99xoE0JNwdIJohenmye7ycK3HQD_oqF1Ec"
	TokenApiUrl           = "https://iam.cloud.ibm.com/identity/token"
	WatsonUrl             = "https://us-south.ml.cloud.ibm.com/ml/v1/text/chat?version=2023-05-29"
	WatsonShelvePrompt    = "Estas analizando la foto de un anaquel de un supermercado. Donde cada estante esta representado por niveles, donde el 1 es el mas alto en la foto.\nSegún el planograma deben estar presentes los siguientes productos de forma ordenada y en su estante correcto: {parsedResponse} Que porcentaje aproximado de los productos estan presentes? Que porcentaje aproximado de los productos están en su lugar definido?\nRespuesta en formato JSON Ejemplo de respuesta:\n{ \"estante 1\": { \"obs\":\"Todos los productos se encuentran en lugar definido\", \"porcentaje\":\"100\" }, \"estante 2\": { \"obs\":\"Algunos productos no estan en su lugar correcto\", \"porcentaje\":\"85\" } }"
	WatsonGraniteModelId  = "ibm/granite-vision-3-2-2b"
	WatsonPlanogramPrompt = "Estas analizando el planograma de un supermercado, una representación gráfica del orden de productos en un anaquel. Identifica cuantos estantes hay en el anaquel, numeralos de arriba hacia abajo e identifica que productos hay en cada uno. Genera como respuesta unicamente el detalle de lo encontrado en formato JSON. Ejemplo de respuesta: {\n  \"estantes\":3,\n \"detalle\" :{\"estante 1\":\"Leche normal, leche deslactozada\",\n  \"estante 2\": \"Queso amarillo, queso panela\",\n  \"estante 3\": \"Mantequilla\"\n}}."
	WatsonLlamaModelId    = "meta-llama/llama-3-2-90b-vision-instruct"
	WatsonProjectId       = "97091540-19f6-4cbf-b7d4-c03d5651b45e"
	WatsonMaxTokens       = 1000
)
