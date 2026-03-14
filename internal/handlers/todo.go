package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Todo struct {
	Title     string `json:"title" bson:"title"`
	Completed bool   `json:"completed" bson:"completed"`
}

var TodoList []Todo

func GetTodos(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, TodoList)
}

func AddTodo(ctx *gin.Context) {
	var body Todo
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusNotAcceptable, gin.H{"message": err})
		return
	}
	TodoList = append(TodoList, body)
	ctx.JSON(http.StatusOK, gin.H{"message": "Todo added successfully", "todo": TodoList})
}
