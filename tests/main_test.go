package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	controller "TaskManagement/app/controllers"
	"TaskManagement/app/database"
	"TaskManagement/app/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func SetupRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/task", controller.GetTask)
	r.POST("/task", controller.CreateTask)
	r.GET("/task/:id", controller.SearchTask)
	r.DELETE("/task/:id", controller.DeleteTask)
	r.PATCH("/task/:id", controller.EditTask)
	r.PUT("/task/:id/complete", controller.MarkTaskAsCompleted)
	r.GET("task/pending", controller.GetPendingTasks)
	return r
}

var ID int

func CreateTask() {
	task := models.Task{Title: "Test Task", Description: "Test Description"}
	database.DB.Create(&task)

	ID = int(task.ID)

}

func DeleteTask() {
	var task models.Task
	database.DB.Delete(&task, ID)
}

func TestGetTask(t *testing.T) {
	database.DataBaseConnection()

	CreateTask()
	defer DeleteTask()

	router := SetupRouter()

	// Cria uma requisição GET
	req, _ := http.NewRequest("GET", "/task", nil)

	// Armazena a resposta
	w := httptest.NewRecorder()

	// Envia a requisição
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreateTask(t *testing.T) {
	database.DataBaseConnection()
	router := SetupRouter()

	task := models.Task{Title: "Test Task", Description: "Test Description"}
	jsonValue, _ := json.Marshal(task)
	req, _ := http.NewRequest("POST", "/task", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSearchTask(t *testing.T) {
	database.DataBaseConnection()
	CreateTask()
	defer DeleteTask()
	router := SetupRouter()

	path := "/task/" + strconv.Itoa(ID)

	req, _ := http.NewRequest("GET", path, nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var taskMock models.Task
	json.Unmarshal(w.Body.Bytes(), &taskMock)

	assert.Equal(t, "Test Task", taskMock.Title)
	assert.Equal(t, "Test Description", taskMock.Description)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteTask(t *testing.T) {
	database.DataBaseConnection()
	CreateTask()

	router := SetupRouter()

	path := "/task/" + strconv.Itoa(ID)

	req, _ := http.NewRequest("DELETE", path, nil)

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestEditTask(t *testing.T) {
	database.DataBaseConnection()
	CreateTask()
	defer DeleteTask()

	router := SetupRouter()

	task := models.Task{Title: "Updated Task", Description: "Updated Description"}

	jsonValue, _ := json.Marshal(task)

	path := "/task/" + strconv.Itoa(ID)
	req, _ := http.NewRequest("PATCH", path, bytes.NewBuffer(jsonValue))
	// req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var taskMockUpdate models.Task
	json.Unmarshal(w.Body.Bytes(), &taskMockUpdate)

	assert.Equal(t, "Updated Task", task.Title)
	assert.Equal(t, "Updated Description", task.Description)

}

func TestMarkTaskAsCompleted(t *testing.T) {
	database.DataBaseConnection()
	CreateTask()
	defer DeleteTask()

	router := SetupRouter()

	path := "/task/" + strconv.Itoa(ID) + "/complete"

	req, _ := http.NewRequest("PUT", path, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var updatedTask models.Task
	err := json.Unmarshal(w.Body.Bytes(), &updatedTask)
	assert.NoError(t, err) // Check for JSON unmarshal errors

	assert.Equal(t, "completed", updatedTask.Status)
}

func TestGetPendingTasks(t *testing.T) {
	database.DataBaseConnection()
	CreateTask() // Cria uma tarefa pendente
	defer DeleteTask()

	router := SetupRouter()

	// Cria uma requisição GET para tarefas pendentes
	req, _ := http.NewRequest("GET", "/task/pending", nil)

	// Armazena a resposta
	w := httptest.NewRecorder()

	// Envia a requisição
	router.ServeHTTP(w, req)

	var taskMock models.Task
	json.Unmarshal(w.Body.Bytes(), &taskMock)

	// Verifica se o status da resposta é OK
	assert.Equal(t, http.StatusOK, w.Code)

	// Verifica se a tarefa criada é pendente
	assert.Equal(t, "not completed", taskMock.Status)
}
