package controller

import (
	"TaskManagement/app/database"
	"TaskManagement/app/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Função que retorna todas as tarefas
func GetTask(c *gin.Context) {
	var task []models.Task
	database.DB.Find(&task) // Busca todas as tarefas no banco de dados
	c.JSON(200, task)       // Retorna todas as tarefas em formato JSON
}

// Função que cria uma nova tarefa
func CreateTask(c *gin.Context) {
	var task models.Task
	// Verifica se o JSON é válido e o armazena em task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.DB.Create(&task)   // Cria uma nova tarefa no banco de dados
	c.JSON(http.StatusOK, task) // Retorna a tarefa criada em formato JSON
}

// Função que busca uma tarefa pelo ID
func SearchTask(c *gin.Context) {
	var task models.Task

	// pegar o ID passado na URL
	id := c.Params.ByName("id")

	// Procurar a tarefa com o ID passado
	database.DB.First(&task, id)

	// Verificação se a tarefa foi encontrada
	if task.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"Not Found": "Tarefa não encontrada"})
		return
	}

	// Retorna a tarefa encontrada
	c.JSON(http.StatusOK, task)
}

// Função que deleta uma tarefa
func DeleteTask(c *gin.Context) {
	var task models.Task

	// extrai o ID da recebido
	id := c.Params.ByName("id")

	// Deleta a tarefa com o ID passado
	database.DB.Delete(&task, id)

	// Retorna uma mensagem de sucesso
	c.JSON(http.StatusOK, gin.H{"data": "Tarefa deletada"})

}

// Função que edita uma tarefa
func EditTask(c *gin.Context) {
	var task models.Task

	// Pegar o ID passado na URL
	id := c.Params.ByName("id")

	// Procurar a tarefa com o ID passado
	database.DB.First(&task, id)

	// Verificação se a tarefa foi encontrada
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Atualiza a tarefa com os novos dados
	database.DB.Model(&task).UpdateColumns(task)

	// Retorna a tarefa atualizada
	c.JSON(http.StatusOK, task)

}

// Função que marca uma tarefa como concluída
func MarkTaskAsCompleted(c *gin.Context) {
	var task models.Task

	// Pegar o ID passado na URL
	id := c.Params.ByName("id")

	// Procurar a tarefa com o ID passado
	database.DB.First(&task, id)

	// Verificação se a tarefa foi encontrada
	if task.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	// Atualiza o status da tarefa para "completed"
	task.Status = "completed"

	// Salva a alteração no banco de dados
	database.DB.Save(&task)

	// Retorna a tarefa atualizada
	c.JSON(http.StatusOK, task)
}

// Função que retorna todas as tarefas pendentes
func GetPendingTasks(c *gin.Context) {
	var task models.Task

	// Procurar todas as tarefas que não foram completadas
	database.DB.Where("status = ?", "not completed").Find(&task)

	// Retorna as tarefas encontradas
	c.JSON(http.StatusOK, task)
}
