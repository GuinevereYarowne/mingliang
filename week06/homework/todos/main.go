package main

import (
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

type Todo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

var (
	todos  []Todo
	mu     sync.Mutex
	nextID int = 1
)

func indexHandler(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()
	c.HTML(http.StatusOK, "index.html", gin.H{
		"todos": todos,
	})
}

func addTodoHandler(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()

	title := c.PostForm("title")
	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "任务标题不能为空"})
		return
	}

	todos = append(todos, Todo{
		ID:    nextID,
		Title: title,
		Done:  false,
	})
	nextID++

	c.Redirect(http.StatusFound, "/")
}

func toggleTodoHandler(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的任务ID"})
		return
	}

	for i := range todos {
		if todos[i].ID == id {
			todos[i].Done = !todos[i].Done
			break
		}
	}

	c.Redirect(http.StatusFound, "/")
}

func clearTodosHandler(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()

	todos = []Todo{}
	nextID = 1

	c.Redirect(http.StatusFound, "/")
}

func deleteTodoHandler(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的任务ID"})
		return
	}

	for i := range todos {
		if todos[i].ID == id {
			todos = append(todos[:i], todos[i+1:]...)
			break
		}
	}

	c.Redirect(http.StatusFound, "/")
}

func editTodoHandler(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的任务ID"})
		return
	}

	newTitle := c.PostForm("newTitle")
	if newTitle == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "新标题不能为空"})
		return
	}

	for i := range todos {
		if todos[i].ID == id {
			todos[i].Title = newTitle
			break
		}
	}

	c.Redirect(http.StatusFound, "/")
}

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	r.GET("/", indexHandler)
	r.POST("/add", addTodoHandler)
	r.POST("/toggle/:id", toggleTodoHandler)
	r.POST("/clear", clearTodosHandler)

	r.POST("/delete/:id", deleteTodoHandler)
	r.POST("/edit/:id", editTodoHandler)

	r.Run(":8082")
}
