package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Article represents an article object.
type Article struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title     string             `json:"title,omitempty" bson:"title,omitempty"`
	Content   string             `json:"content,omitempty" bson:"content,omitempty"`
	Author    string             `json:"author,omitempty" bson:"author,omitempty"`
	CreatedAt time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
}

// Response represents a generic JSON response.
type Response struct {
	Status string `json:"status"`
}

// ErrorResponse represents an error JSON response.
type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// MongoDB connection string.
//const connectionString = "mongodb://localhost:27017"

//for docker use this  connection string
const connectionString = "mongodb://mongo:27017"

// MongoDB database name.
const dbName = "mydatabase"

// MongoDB collection name.
const collectionName = "articles"

// MongoDB client instance.
var client *mongo.Client

// MongoDB collection instance.
var collection *mongo.Collection

// CreateArticleEndpoint creates a new article.
func CreateArticleEndpoint(w http.ResponseWriter, req *http.Request) {
	log.Print("Insert")
	w.Header().Set("Content-Type", "application/json")
	var article Article
	err := json.NewDecoder(req.Body).Decode(&article)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Status:  "error",
			Message: "invalid request body",
		})
		return
	}
	if article.Title == "" || article.Content == "" || article.Author == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Status:  "error",
			Message: "title, content and author are required",
		})
		return
	}
	article.CreatedAt = time.Now()
	result, err := collection.InsertOne(req.Context(), article)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			Status:  "error",
			Message: "failed to create article",
		})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{
		Status: "success",
	})
	json.NewEncoder(w).Encode(bson.M{"id": result.InsertedID})
}

// GetArticleEndpoint retrieves an existing article.
func GetArticleEndpoint(w http.ResponseWriter, req *http.Request) {
	log.Print("GetArticleEndpoint")
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(req)
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Status:  "error",
			Message: "invalid article ID",
		})
		return
	}
	var article Article
	err = collection.FindOne(req.Context(), bson.M{"_id": id}).Decode(&article)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{
			Status:  "error",
			Message: "article not found",
		})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(article)
}

// ListArticlesEndpoint retrieves all existing articles.
func ListArticlesEndpoint(w http.ResponseWriter, req *http.Request) {
	log.Print("ListArticlesEndpoint")
	w.Header().Set("Content-Type", "application/json")
	cursor, err := collection.Find(req.Context(), bson.M{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			Status:  "error",
			Message: "failed to list articles",
		})
		return
	}
	var articles []Article
	for cursor.Next(req.Context()) {
		var article Article
		err := cursor.Decode(&article)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{
				Status:  "error",
				Message: "failed to decode article",
			})
			return
		}
		articles = append(articles, article)
	}
	if err := cursor.Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			Status:  "error",
			Message: "failed to list articles",
		})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(articles)
}

func main() {
	// Set up MongoDB connection.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	clientOptions := options.Client().ApplyURI(connectionString)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	collection = client.Database(dbName).Collection(collectionName)

	// Set up HTTP server.
	router := mux.NewRouter()
	router.HandleFunc("/articles", CreateArticleEndpoint).Methods("POST")
	router.HandleFunc("/articles/{id}", GetArticleEndpoint).Methods("GET")
	router.HandleFunc("/articles", ListArticlesEndpoint).Methods("GET")
	log.Fatal(http.ListenAndServe(":8123", router))
}

//curl -X POST -H "Content-Type: application/json" -d '{"title":"My first article", "content":"This is some content", "author":"John Doe"}' http://localhost:8123/articles

//curl http://localhost:8123/articles/6404e49c65c0a3503bad5c8e

//curl http://localhost:8123/articles

// docker pull mongo

//mongo docker cmd
// docker run -d -p 27017:27017 --name mymongo mongo

//go app docker cmd
//docker run -p 8123:8123 --name myblogs blogs

// docker stop mymongo
// docker rm mymongo
