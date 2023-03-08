package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestInsertArticle(t *testing.T) {
	// Set up a test MongoDB database and collection
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		t.Fatalf("failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())

	collection := client.Database("mydatabase").Collection("articles")

	// Insert an article into the test collection
	article := Article{
		Title:   "Test Article1234",
		Content: "This is a test article.",
		Author:  "younus",
	}
	_, err = collection.InsertOne(context.Background(), article)
	if err != nil {
		t.Fatalf("failed to insert article into collection: %v", err)
	}

	// Verify that the article has been inserted successfully
	result := Article{}
	err = collection.FindOne(context.Background(), Article{Title: "Test Article1234"}).Decode(&result)
	if err != nil {
		t.Fatalf("failed to find article in collection: %v", err)
	}

	expected := Article{
		Title:   "Test Article1234",
		Content: "This is a test article.",
		Author:  "younus",
	}
	/*if result != expected {
		t.Errorf("unexpected result: got %v, want %v", result, expected)
	}*/
	if result.Title != expected.Title ||
		result.Content != expected.Content ||
		result.Author != expected.Author {
		t.Errorf("unexpected result: got %+v, want %+v", result, expected)
	}
}

// list all articles

/*func TestFindAllArticles(t *testing.T) {
	// Seed the database with some test data
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	articles := []*Article{
		{
			Title:     "Article 1",
			Content:   "This is article 1.",
			Author:    "John Doe",
			CreatedAt: time.Now(),
		},
		{
			Title:     "Article 2",
			Content:   "This is article 2.",
			Author:    "Jane Smith",
			CreatedAt: time.Now(),
		},
		{
			Title:     "Article 3",
			Content:   "This is article 3.",
			Author:    "Bob Johnson",
			CreatedAt: time.Now(),
		},
	}

	for _, article := range articles {
		err := CreateArticleEndpoint(ctx, article)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Call the function to retrieve all the articles
	result, err := findAllArticles(ctx)
	if err != nil {
		t.Fatal(err)
	}

	// Verify that the returned result matches the expected result
	if len(result) != len(articles) {
		t.Errorf("unexpected number of articles: got %d, want %d", len(result), len(articles))
	}

	for i, article := range result {
		if article.Title != articles[i].Title ||
			article.Content != articles[i].Content ||
			article.Author != articles[i].Author ||
			article.CreatedAt.Unix() != articles[i].CreatedAt.Unix() {
			t.Errorf("unexpected result for article %d: got %+v, want %+v", i+1, article, articles[i])
		}
	}
}*/

func TestListArticlesEndpoint(t *testing.T) {
	// create a mock database and collection
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	defer client.Disconnect(context.Background())

	db := client.Database("testdb")
	collection := db.Collection("testcollection")

	// insert some test data
	_, err = collection.InsertMany(context.Background(), []interface{}{
		bson.M{
			"title":      "Test Article 1",
			"content":    "This is a test article.",
			"author":     "John Doe",
			"created_at": time.Now(),
		},
		bson.M{
			"title":      "Test Article 2",
			"content":    "This is another test article.",
			"author":     "Jane Doe",
			"created_at": time.Now(),
		},
	})
	if err != nil {
		t.Fatalf("failed to insert test data: %v", err)
	}

	// create a request
	req := httptest.NewRequest(http.MethodGet, "/articles", nil)
	w := httptest.NewRecorder()

	// call the ListArticlesEndpoint function
	ListArticlesEndpoint(w, req.WithContext(context.Background()))

	// check the response
	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code %d, but got %d", http.StatusOK, resp.StatusCode)
	}

	var articles []Article
	err = json.NewDecoder(resp.Body).Decode(&articles)
	if err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if len(articles) != 2 {
		t.Errorf("expected %d articles, but got %d", 2, len(articles))
	}
}
