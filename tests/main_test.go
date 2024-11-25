package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	controller "TaskManagement/app/controllers"
	"TaskManagement/app/models"
	"TaskManagement/app/repository" // Importa o repositório

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// Mock do repositório de tarefas
type MockTaskRepository struct {
	tasks  map[uint]models.Task
	nextID uint
}

func NewMockTaskRepository() *MockTaskRepository {
	return &MockTaskRepository{
		tasks:  make(map[uint]models.Task),
		nextID: 1,
	}
}

func (m *MockTaskRepository) Create(task *models.Task) error {
	task.ID = m.nextID
	m.tasks[m.nextID] = *task
	m.nextID++
	return nil
}

func (m *MockTaskRepository) FindByID(id uint) (*models.Task, error) {
	task, exists := m.tasks[id]
	if !exists {
		return nil, gorm.ErrRecordNotFound
	}
	return &task, nil
}

func (m *MockTaskRepository) Delete(id uint) error {
	delete(m.tasks, id)
	return nil
}

func (m *MockTaskRepository) Update(task *models.Task) error {
	m.tasks[task.ID] = *task
	return nil
}

func (m *MockTaskRepository) FindAll() ([]models.Task, error) {
	var tasks []models.Task
	for _, task := range m.tasks {
		tasks = append(tasks, task)
	}
	return tasks, nil
}

// Configura o roteador com o repositório mock
func SetupRouter(repo repository.TaskRepository) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/task", func(c *gin.Context) { controller.GetTask(c, repo) })
	r.POST("/task", func(c *gin.Context) { controller.CreateTask(c, repo) })
	r.GET("/task/:id", func(c *gin.Context) { controller.SearchTask(c, repo) })
	r.DELETE("/task/:id", func(c *gin.Context) { controller.DeleteTask(c, repo) })
	r.PATCH("/task/:id", func(c *gin.Context) { controller.EditTask(c, repo) })
	r.PUT("/task/:id/complete", func(c *gin.Context) { controller.MarkTaskAsCompleted(c, repo) })
	r.GET("/task/pending", func(c *gin.Context) { controller.GetPendingTasks(c, repo) })
	return r
}

// Testes
func TestCreateTask(t *testing.T) {
	repo := NewMockTaskRepository() // Usando o repositório mock
	router := SetupRouter(repo)

	task := models.Task{Title: "Test Task", Description: "Test Description"}
	jsonValue, _ := json.Marshal(task)
	req, _ := http.NewRequest("POST", "/task", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetTask(t *testing.T) {
	repo := NewMockTaskRepository()
	router := SetupRouter(repo)

	// Cria uma tarefa para testar
	task := models.Task{Title: "Test Task", Description: "Test Description"}
	repo.Create(&task)

	// Cria uma requisição GET
	req, _ := http.NewRequest("GET", "/task", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSearchTask(t *testing.T) {
	repo := NewMockTaskRepository()
	router := SetupRouter(repo)

	// Cria uma tarefa para testar
	task := models.Task{Title: "Test Task", Description: "Test Description"}
	repo.Create(&task)

	path := "/task/" + strconv.Itoa(int(task.ID))
	req, _ := http.NewRequest("GET", path, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var taskMock models.Task
	json.Unmarshal(w.Body.Bytes(), &taskMock)

	assert.Equal(t, "Test Task", taskMock.Title)
	assert.Equal(t, "Test Description", taskMock.Description)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestEditTask(t *testing.T) {
	repo := NewMockTaskRepository()
	router := SetupRouter(repo)

	// Cria uma tarefa para testar
	task := models.Task{Title: "Test Task", Description: "Test Description"}
	repo.Create(&task)

	updatedTask := models.Task{Title: "Updated Task", Description: "Updated Description"}
	jsonValue, _ := json.Marshal(updatedTask)

	path := "/task/" + strconv.Itoa(int(task.ID))
	req, _ := http.NewRequest("PATCH", path, bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Verifica se a tarefa foi atualizada
	updatedTaskFromRepo, _ := repo.FindByID(task.ID)
	assert.Equal(t, "Updated Task", updatedTaskFromRepo.Title)
	assert.Equal(t, "Updated Description", updatedTaskFromRepo.Description)
}

func TestDeleteTask(t *testing.T) {
	repo := NewMockTaskRepository()
	router := SetupRouter(repo)

	// Cria uma tarefa para testar
	task := models.Task{Title: "Test Task", Description: "Test Description"}
	repo.Create(&task)

	path := "/task/" + strconv.Itoa(int(task.ID))
	req, _ := http.NewRequest("DELETE", path, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestMarkTaskAsCompleted(t *testing.T) {
	repo := NewMockTaskRepository()
	router := SetupRouter(repo)

	// Cria uma tarefa para testar
	task := models.Task{Title: "Test Task", Description: "Test Description"}
	repo.Create(&task)

	path := "/task/" + strconv.Itoa(int(task.ID)) + "/complete"
	req, _ := http.NewRequest("PUT", path, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetPendingTasks(t *testing.T) {
	repo := NewMockTaskRepository()
	router := SetupRouter(repo)

	// Cria uma tarefa pendente
	task := models.Task{Title: "Test Task", Description: "Test Description"}
	repo.Create(&task)

	// Cria uma requisição GET para tarefas pendentes
	req, _ := http.NewRequest("GET", "/task/pending", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
