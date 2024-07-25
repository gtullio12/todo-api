package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Setup MongoDB client
var collection *mongo.Collection
var ctx = context.TODO()
var coll = connectToDatabase()

type Todo struct {
	Id      primitive.ObjectID `bson:"_id"`
	Title   string             `bson:"title"`
	Content string             `bson:"content"`
	IsDone  bool               `bson:"isDone"`
}

func main() {
	router := gin.Default()
	router.GET("/getTodos", getTodos)
	router.PUT("/editTodo", editTodo)
	router.DELETE("/deleteTodo", deleteTodo)
	router.POST("/createTodo", createTodo)

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", os.Getenv("FRONT_END_ENDPOINT")},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Length"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080" // Default port for development
	}

	// Start the router
	router.Run(":" + port)
}

func connectToDatabase() (coll *mongo.Collection) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_URL"))
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	coll = client.Database("todo").Collection("todo_collection")
	return
}

func createTodo(c *gin.Context) {
	var todo Todo
	c.Bind(&todo)
	todo.Id = primitive.NewObjectID()

	res, err := coll.InsertOne(context.TODO(), todo)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Inserted document with _id: %v\n", res.InsertedID)
}

func getTodos(c *gin.Context) {
	var todos []Todo
	findOptions := options.Find()
	cursor, err := coll.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		log.Fatal(err)
	}
	if err = cursor.All(context.TODO(), &todos); err != nil {
		log.Fatal(err)
	}

	c.IndentedJSON(http.StatusOK, todos)
}

func editTodo(c *gin.Context) {
	var todoToEdit Todo
	c.BindJSON(&todoToEdit)

	update := bson.D{{"$set", bson.D{{"title", todoToEdit.Title}}}, {"$set", bson.D{{"content", todoToEdit.Content}}}, {"$set", bson.D{{"isDone", todoToEdit.IsDone}}}}
	_, err := coll.UpdateOne(context.TODO(), bson.D{{"_id", todoToEdit.Id}}, update)
	if err != nil {
		log.Fatal(err)
	}
}

func deleteTodo(c *gin.Context) {
	var todoToDelete Todo
	c.BindJSON(&todoToDelete)

	_, err := coll.DeleteOne(context.TODO(), bson.D{{"_id", todoToDelete.Id}})
	if err != nil {
		log.Fatal(err)
	}
}
