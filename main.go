package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)



type Todo struct {
    ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
    Completed bool `json:"completed"`
    Body string `json:"body"`
}

var collection *mongo.Collection



func main() {
	fmt.Println("Hello, World!")

    err := godotenv.Load(".env")
    if err != nil {
        log.Fatalf("Error loading .env file: %v", err)
    }

    mongoURI := os.Getenv("MONGO_URI")
    clientOptions := options.Client().ApplyURI(mongoURI)
    client,err := mongo.Connect(context.Background(), clientOptions)
    if err != nil {
        log.Fatal(err)
    }
    defer client.Disconnect(context.Background())

    err = client.Ping(context.Background(), nil)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Connected to MongoDB")

    collection = client.Database("golang_db").Collection("todos")

    app := fiber.New()

    app.Get("/api/todos", getTodos)
    app.Post("/api/todos", createTodo)
    // app.Patch("/api/todos/:id", updateTodo)
    // app.Delete("/api/todos/:id", deleteTodo)

    port := os.Getenv("PORT")
    if port == "" {
        port = "4000"
    }

    log.Fatal(app.Listen(":" + port))
}

 func getTodos(c *fiber.Ctx) error {
    todos := []Todo{}
    cursor, err := collection.Find(context.Background(), bson.M{})
    if err != nil {
        return err
    }

    defer cursor.Close(context.Background())

    for cursor.Next(context.Background()) {
        var todo Todo
        if err := cursor.Decode(&todo); err != nil {
            return err
        }
        todos = append(todos, todo)
    }
    return c.JSON(todos)
}

func createTodo(c *fiber.Ctx) error {
    todo := new(Todo)
    if err := c.BodyParser(todo); err != nil {
        return err
    }

    if todo.Body == "" {
        return c.Status(400).JSON(fiber.Map{
            "message": "Body is required",
        })
    }

    result, err := collection.InsertOne(context.Background(), todo)
    if err != nil {
        return err
    }

    todo.ID = result.InsertedID.(primitive.ObjectID)
    
    return c.Status(201).JSON(todo)
}