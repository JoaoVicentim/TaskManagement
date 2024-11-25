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
	"TaskManagement/app/repository"

	"github.com/gin-gonic/gin"
	"github.com/zeebo/assert"
)

func SetupRouter(repo database.TaskRepository) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	controller := &controller.TaskController{Repo: repo}

	r.GET("/task", controller.GetTask)
	r.POST("/task", controller.CreateTask)
	r.GET("/task/:id", controller.SearchTask)
	r.DELETE("/task/:id", controller.DeleteTask)
	r.PATCH("/task/:id", controller.EditTask)
	r.PUT("/task/:id/complete", controller.MarkTaskAsCompleted)
	r.GET("/task/pending", controller.GetPendingTasks)

	return r
}

func TestGetTask(t *testing.T) {
	// Cria uma instância do repositório mock
	mockRepo := &repository.MockTaskRepository{Tasks: make(map[uint]*models.Task)}

	// Cria uma tarefa de teste
	testTask := &models.Task{Title: "Test Task", Description: "Test Description"}
	mockRepo.Create(testTask) // Adiciona a tarefa ao repositório mock

	// Configura o roteador com o repositório mock
	router := SetupRouter(mockRepo)

	// Cria uma requisição GET
	req, _ := http.NewRequest("GET", "/task", nil)

	// Armazena a resposta
	w := httptest.NewRecorder()

	// Envia a requisição
	router.ServeHTTP(w, req)

	// Verifica se o código de status da resposta é 200 OK
	assert.Equal(t, http.StatusOK, w.Code)

	// Verifica se a resposta contém a tarefa criada
	var tasks []models.Task
	err := json.Unmarshal(w.Body.Bytes(), &tasks)
	assert.NoError(t, err) // Verifica se não houve erro ao deserializar

	assert.Equal(t, 1, len(tasks))                            // Verifica se uma tarefa foi retornada
	assert.Equal(t, "Test Task", tasks[0].Title)              // Verifica o título da tarefa
	assert.Equal(t, "Test Description", tasks[0].Description) // Verifica a descrição da tarefa
}

func TestCreateTask(t *testing.T) {
	mockRepo := &repository.MockTaskRepository{Tasks: make(map[uint]*models.Task)}
	router := SetupRouter(mockRepo)

	task := models.Task{Title: "Test Task", Description: "Test Description"}
	jsonValue, _ := json.Marshal(task)
	req, _ := http.NewRequest("POST", "/task", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
func TestSearchTask(t *testing.T) {
	// Cria uma instância do repositório mock
	mockRepo := &repository.MockTaskRepository{Tasks: make(map[uint]*models.Task)}

	// Cria uma tarefa de teste
	testTask := &models.Task{Title: "Test Task", Description: "Test Description"}
	mockRepo.Create(testTask) // Adiciona a tarefa ao repositório mock

	// Configura o roteador com o repositório mock
	router := SetupRouter(mockRepo)

	// Define o caminho da requisição usando o ID da tarefa de teste
	path := "/task/" + strconv.Itoa(int(testTask.ID))

	// Cria uma requisição GET
	req, _ := http.NewRequest("GET", path, nil)

	// Armazena a resposta
	w := httptest.NewRecorder()

	// Envia a requisição
	router.ServeHTTP(w, req)

	// Verifica se a resposta foi deserializada corretamente
	var taskMock models.Task
	err := json.Unmarshal(w.Body.Bytes(), &taskMock)
	assert.NoError(t, err) // Verifica se não houve erro ao deserializar

	// Verifica se os dados da tarefa estão corretos
	assert.Equal(t, "Test Task", taskMock.Title)
	assert.Equal(t, "Test Description", taskMock.Description)
	assert.Equal(t, http.StatusOK, w.Code) // Verifica se o código de status é 200 OK
}

func TestDeleteTask(t *testing.T) {
	// Cria uma instância do repositório mock
	mockRepo := &repository.MockTaskRepository{Tasks: make(map[uint]*models.Task)}

	// Cria uma tarefa de teste
	testTask := &models.Task{Title: "Test Task", Description: "Test Description"}
	mockRepo.Create(testTask) // Adiciona a tarefa ao repositório mock

	// Configura o roteador com o repositório mock
	router := SetupRouter(mockRepo)

	path := "/task/" + strconv.Itoa(int(testTask.ID)) // Usa o ID gerado
	req, _ := http.NewRequest("DELETE", path, nil)

	w := httptest.NewRecorder()

	// Envia a requisição
	router.ServeHTTP(w, req)

	// Verifica se o código de status da resposta é 200 OK
	assert.Equal(t, http.StatusOK, w.Code)

	// Verifica se a tarefa foi realmente deletada
	_, err := mockRepo.FindByID(testTask.ID)   // Tenta encontrar a tarefa após a deleção
	assert.NoError(t, err)                     // Verifica se não houve erro ao buscar
	assert.Nil(t, mockRepo.Tasks[testTask.ID]) // Verifica se a tarefa foi removida
}

func TestEditTask(t *testing.T) {
	// Cria uma instância do repositório mock
	mockRepo := &repository.MockTaskRepository{Tasks: make(map[uint]*models.Task)}

	// Cria uma tarefa de teste
	testTask := &models.Task{Title: "Test Task", Description: "Test Description"}
	mockRepo.Create(testTask) // Adiciona a tarefa ao repositório mock

	// Configura o roteador com o repositório mock
	router := SetupRouter(mockRepo)

	// Tarefa a ser atualizada
	updatedTask := models.Task{Title: "Updated Task", Description: "Updated Description"}
	jsonValue, _ := json.Marshal(updatedTask)

	path := "/task/" + strconv.Itoa(int(testTask.ID)) // Usa o ID gerado
	req, _ := http.NewRequest("PATCH", path, bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Verifica se a tarefa foi atualizada corretamente
	var taskMockUpdate models.Task
	err := json.Unmarshal(w.Body.Bytes(), &taskMockUpdate)
	assert.NoError(t, err) // Verifica se não houve erro ao deserializar

	assert.Equal(t, "Updated Task", taskMockUpdate.Title)              // Verifica o título atualizado
	assert.Equal(t, "Updated Description", taskMockUpdate.Description) // Verifica a descrição atualizada
}
func TestMarkTaskAsCompleted(t *testing.T) {
	// Cria uma instância do repositório mock
	mockRepo := &repository.MockTaskRepository{Tasks: make(map[uint]*models.Task)}

	// Cria uma tarefa de teste
	testTask := &models.Task{Title: "Test Task", Description: "Test Description", Status: "not completed"}
	mockRepo.Create(testTask) // Adiciona a tarefa ao repositório mock

	// Configura o roteador com o repositório mock
	router := SetupRouter(mockRepo)

	path := "/task/" + strconv.Itoa(int(testTask.ID)) + "/complete" // Usa o ID gerado
	req, _ := http.NewRequest("PUT", path, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Verifica se a tarefa foi marcada como concluída
	var updatedTask models.Task
	err := json.Unmarshal(w.Body.Bytes(), &updatedTask)
	assert.NoError(t, err) // Verifica se não houve erro ao deserializar

	assert.Equal(t, "completed", updatedTask.Status) // Verifica se o status foi atualizado
}

func TestGetPendingTasks(t *testing.T) {
	// Cria uma instância do repositório mock
	mockRepo := &repository.MockTaskRepository{Tasks: make(map[uint]*models.Task)}

	// Cria uma tarefa pendente de teste
	testTask := &models.Task{Title: "Test Task", Description: "Test Description", Status: "not completed"}
	mockRepo.Create(testTask) // Adiciona a tarefa ao repositório mock

	// Configura o roteador com o repositório mock
	router := SetupRouter(mockRepo)

	// Cria uma requisição GET para tarefas pendentes
	req, _ := http.NewRequest("GET", "/task/pending", nil)

	// Armazena a resposta
	w := httptest.NewRecorder()

	// Envia a requisição
	router.ServeHTTP(w, req)

	// Verifica se o status da resposta é OK
	assert.Equal(t, http.StatusOK, w.Code)

	// Verifica se a tarefa criada é pendente
	var tasks []models.Task
	err := json.Unmarshal(w.Body.Bytes(), &tasks)
	assert.NoError(t, err) // Verifica se não houve erro ao deserializar

	assert.Equal(t, 1, len(tasks))                    // Verifica se uma tarefa foi retornada
	assert.Equal(t, "not completed", tasks[0].Status) // Verifica o status da tarefa
}
