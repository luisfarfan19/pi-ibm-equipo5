package services

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"time"
	"wantson-service/internal/db"
	"wantson-service/internal/models"
	"wantson-service/pkg/utils"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var loginReq models.CredentialsRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		http.Error(w, "Error al leer JSON", http.StatusBadRequest)
		return
	}

	collection := db.MongoClient.Database(utils.MongoDbName).Collection("usuarios")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Buscar el usuario por username y password
	var usuario models.User
	err := collection.FindOne(ctx, bson.M{
		"username": loginReq.Username,
		"password": loginReq.Password,
	}).Decode(&usuario)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Usuario o contraseña incorrectos", http.StatusUnauthorized)
		} else {
			log.Println("Error al consultar Mongo:", err)
			http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		}
		return
	}

	// Usuario encontrado
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usuario)
}

func GetStoresHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Recuperando todas las tiendas")
	collection := db.MongoClient.Database(utils.MongoDbName).Collection("tienda")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	opts := options.Find().SetProjection(bson.M{
		"nombre": 1,
		"_id":    0,
	})

	cursor, err := collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		log.Println("Error al consultar tiendas:", err)
		http.Error(w, "Error interno", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var stores []models.StoreName
	if err := cursor.All(ctx, &stores); err != nil {
		log.Println("Error al parsear resultados:", err)
		http.Error(w, "Error interno", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stores)
}

func GetStoreDataHandler(w http.ResponseWriter, r *http.Request) {
	storeName := r.URL.Query().Get("name")
	if storeName == "" {
		http.Error(w, "El parámetro 'name' es requerido", http.StatusBadRequest)
		return
	}
	log.Println("Buscando informacion de la tienda : ", storeName)

	collection := db.MongoClient.Database(utils.MongoDbName).Collection("tienda")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var store models.Store
	err := collection.FindOne(ctx, bson.M{"nombre": storeName}).Decode(&store)
	if err != nil {
		http.Error(w, "Store no encontrada", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(store)
}

func GetSectionResponseHandler(w http.ResponseWriter, r *http.Request) {
	storeName := r.URL.Query().Get("storeName")
	sectionId := r.URL.Query().Get("sectionId")

	if storeName == "" || sectionId == "" {
		http.Error(w, "storeName and sectionId are required", http.StatusBadRequest)
		return
	}

	log.Println("Buscando informacion de la tienda y seccion: ", storeName, sectionId)

	collection := db.MongoClient.Database(utils.MongoDbName).Collection("planogram_response")

	filter := bson.M{
		"storeName":           storeName,
		"secciones.idSeccion": sectionId,
	}

	var result models.PlanogramResponse
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "No se encontró información", http.StatusNotFound)
		} else {
			log.Println("Error al consultar resultados:", err)
			http.Error(w, "Error al consultar la base de datos", http.StatusInternalServerError)
		}
		return
	}

	var matchedSection *models.WatsonSection
	for _, sec := range result.Sections {
		if sec.SectionId == sectionId {
			matchedSection = &sec
			break
		}
	}

	if matchedSection == nil {
		http.Error(w, "Sección no encontrada", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(matchedSection)
}

func SaveWatsonResponseToMongo(parsedContent models.PlanogramWatsonApiResponse, storeName, idSection string) bson.M {
	collection := db.MongoClient.Database(utils.MongoDbName).Collection("planogram_response")

	section := bson.M{
		"idSeccion":            idSection,
		"watsonPromptResponse": parsedContent,
	}

	doc := bson.M{
		"storeName": storeName,
		"secciones": []bson.M{section},
	}

	_, err := collection.InsertOne(context.TODO(), doc)
	if err != nil {
		log.Printf("Error insertando en Mongo: %v", err)
		return nil
	}

	log.Println("Se insertó el documento en Mongo correctamente con la respuesta de Watson")
	return section
}
