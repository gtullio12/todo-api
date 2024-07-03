package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

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
	Content string             `bson:"content"`
	IsDone  bool               `bson:"isDone"`
}

func main() {
	router := gin.Default()
	router.GET("/getTodos", getTodos)
	router.PUT("/editTodo", editTodo)
	router.DELETE("/deleteTodo", deleteTodo)
	router.POST("/createTodo", createTodo)

	router.Run("todo-app-production-8c19.up.railway.app")
}

func connectToDatabase() (coll *mongo.Collection) {
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_PRIVATE_URL"))
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
	c.Bind(&todoToEdit)

	update := bson.D{{"$set", bson.D{{"content", todoToEdit.Content}}}, {"$set", bson.D{{"isDone", todoToEdit.IsDone}}}}
	_, err := coll.UpdateOne(context.TODO(), bson.D{{"_id", todoToEdit.Id}}, update)
	if err != nil {
		log.Fatal(err)
	}
}

func deleteTodo(c *gin.Context) {
	var todoToDelete Todo
	c.Bind(&todoToDelete)

	_, err := coll.DeleteOne(context.TODO(), bson.D{{"_id", todoToDelete.Id}})
	if err != nil {
		log.Fatal(err)
	}
}
