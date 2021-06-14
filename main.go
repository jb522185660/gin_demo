package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type (
	todoModel struct {
		gorm.Model
		Title      string `json:"title"`
		Completed  int    `json:"completed"`
	}

	transformedTodo struct {
		ID			uint 	`json:"id"`
		Title		string	`json:"title"`
		Completed   bool  	`json:"completed"`
	}
)

var db *gorm.DB

func init() {
	fmt.Println("初始化")
	fmt.Println("开始连接数据库")
	var err error
	dsn := "root:123456@tcp(127.0.0.1:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err == nil {
		fmt.Println("连接数据库成功")
	} else {
		fmt.Println("连接数据库失败")
	}
	err = db.AutoMigrate(&todoModel{})
	if err != nil {
		fmt.Println("迁移数据表失败")
		return
	} else {
		fmt.Println("迁移数据表成功")
	}
}

func main() {
	engine := gin.Default()
	v1 := engine.Group("/api/v1/todos")
	{
		v1.POST("/", createTodo)
		v1.GET("/", fetchAllTodos)
		v1.GET("/:id", fetchSingleTodo)
	}
	engine.Run(":8090")
}

func createTodo(c *gin.Context)  {
	completed, _ := strconv.Atoi(c.PostForm("completed"))
	todo := todoModel{Title: c.PostForm("title"), Completed: completed}
	db.Save(&todo)
	c.JSON(http.StatusCreated, gin.H{
		"status": http.StatusCreated,
		"message": "Todo item created successfully!",
		"id": todo.ID,
	})
}

func fetchAllTodos(c *gin.Context)  {
	var todos []todoModel
	var _todos	[]transformedTodo

	db.Find(&todos)
	if len(todos) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"message": "No todo found!",
		})
		return
	}

	for _, item := range todos {
		completed := false
		if item.Completed == 1 {
			completed = true
		} else {
			completed = false
		}
		_todos = append(_todos, transformedTodo{ID: item.ID, Title: item.Title, Completed: completed})
	}
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data": _todos,
	})
}

func fetchSingleTodo(c *gin.Context)  {
	var todo todoModel
	todoID := c.Param("id")
	db.First(&todo, todoID)
	if todo.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"message": "No todo found!",
		})
	}
	completed := false
	if todo.Completed == 1 {
		completed = true
	} else {
		completed = false
	}
	_todo := transformedTodo{ID: todo.ID, Title: todo.Title, Completed: completed}
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data": _todo,
	})
}

