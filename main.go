package main

import (
	"net/http"
	"strconv"
	"todo-item-app/handlers"
	"todo-item-app/models"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
	var err error
	db, err = gorm.Open(sqlite.Open("todos.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&models.Todo{}, &models.Item{})

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")

	// Todo Routes
	r.GET("/", getTodos)
	r.GET("/todos", getTodosPartial)
	r.POST("/todos", createTodo)
	r.PUT("/todos/:id", updateTodo)
	r.DELETE("/todos/:id", deleteTodo)

	itemHandler := &handlers.NewItemHandler(db)
	r.GET("/items", itemHandler.GetItems)
	r.GET("/items/list", itemHandler.GetItemsPartial)
	r.POST("/items", itemHandler.CreateItem)
	r.DELETE("/items/:id", itemHandler.DeleteItem)

	r.Run(":8080")
}

// Todo Handlers
func getTodos(c *gin.Context) {
	var todos []models.Todo
	db.Find(&todos)
	c.HTML(http.StatusOK, "index.html", gin.H{
		"todos":  todos,
		"active": "todos",
	})
}

func getTodosPartial(c *gin.Context) {
	var todos []models.Todo
	db.Find(&todos)
	c.HTML(http.StatusOK, "todos.html", gin.H{
		"todos": todos,
	})
}

func createTodo(c *gin.Context) {
	var todo models.Todo
	if err := c.Bind(&todo); err != nil {
		c.String(http.StatusBadRequest, "Bad request")
		return
	}
	db.Create(&todo)
	c.HTML(http.StatusCreated, "todo.html", gin.H{
		"todo": todo,
	})
}

func updateTodo(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var todo models.Todo
	if err := db.First(&todo, id).Error; err != nil {
		c.String(http.StatusNotFound, "Not found")
		return
	}
	todo.Completed = !todo.Completed
	db.Save(&todo)
	c.HTML(http.StatusOK, "todo.html", gin.H{
		"todo": todo,
	})
}

func deleteTodo(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var todo models.Todo
	if err := db.First(&todo, id).Error; err != nil {
		c.String(http.StatusNotFound, "Not found")
		return
	}
	db.Delete(&todo)
	c.String(http.StatusOK, "")
}
