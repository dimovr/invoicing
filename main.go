package main

import (
	"fmt"
	"net/http"
	"strconv"
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

	// models.Item Routes
	r.GET("/items", getItems)
	r.GET("/items/list", getItemsPartial)
	r.POST("/items", createItem)
	r.DELETE("/items/:id", deleteItem)

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

func getItems(c *gin.Context) {
	var items []models.Item
	db.Find(&items)
	fmt.Println(items)
	c.HTML(http.StatusOK, "index.html", gin.H{
		"items":  items,
		"active": "items",
	})
}

func getItemsPartial(c *gin.Context) {
	var items []models.Item
	db.Find(&items)
	fmt.Println(items)
	c.HTML(http.StatusOK, "items.html", gin.H{
		"items": items,
	})
}

func createItem(c *gin.Context) {
	var item models.Item
	if err := c.Bind(&item); err != nil {
		c.String(http.StatusBadRequest, "Bad request")
		return
	}
	fmt.Println("creating item: ", item.ID)
	db.Create(&item)
	fmt.Println("item created: ", item.ID)
	c.HTML(http.StatusCreated, "item.html", item)
}

func deleteItem(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var item models.Item
	if err := db.First(&item, id).Error; err != nil {
		c.String(http.StatusNotFound, "Not found")
		return
	}
	db.Delete(&item)
	c.String(http.StatusOK, "")
}
