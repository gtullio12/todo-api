package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	Id       primitive.ObjectID `bson:"_id",required`
	Name     string             `bson:"name"`
	Email    string             `bson:"email"`
	Password string             `bson:"password"`
}

// Setup MongoDB client
var collection *mongo.Collection
var ctx = context.TODO()
var coll = connectToDatabase()

type Todo struct {
	Id      int    `json:"Id"`
	Content string `json:"Content"`
	IsDone  bool   `json:"IsDone"`
}

// const todos []Todo
var todos = []Todo{
	Todo{Id: 1, Content: "Update Documentation", IsDone: true},
	Todo{Id: 2, Content: "Add RESTful API", IsDone: false},
	Todo{Id: 3, Content: "Modify function parameters", IsDone: false},
	Todo{Id: 4, Content: "Add method recievers", IsDone: true},
}

func main() {
	router := gin.Default()
	router.GET("/getTodos", getTodos)
	router.GET("/getTodo", getSingleTodo)
	router.POST("/editTodo", editTodo)
	router.DELETE("/deleteTodo", deleteTodo)
	router.POST("/createTodo", createTodo)

	router.Run("localhost:8000")
}

func connectToDatabase() (coll *mongo.Collection) {
	clientOptions := options.Client().ApplyURI("mongodb+srv://{username}:{password}@todocluser.yh1icqj.mongodb.net/?retryWrites=true&w=majority&appName=todocluser")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	coll = client.Database("todo").Collection("todo_collection")
	return
}

func createTodo(c *gin.Context) {

}

func getTodos(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, todos)
}

func editTodo(c *gin.Context) {

}

func deleteTodo(c *gin.Context) {

}

func getSingleTodo(c *gin.Context) {

}
