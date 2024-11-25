package controller

import (
	"TaskManagement/app/database"
	"TaskManagement/app/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
	if err := c.ShouldBindJSON(&task); err != nil || task.Title == "" || task.Description == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Título e descrição são obrigatórios"})
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

	// Tenta encontrar a tarefa pelo ID
	if err := database.DB.First(&task, id).Error; err != nil {
		// Se a tarefa não for encontrada, retorna 404
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Tarefa não encontrada"})
			return
		}
		// Se houver outro erro ao acessar o banco de dados, retorna 500
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao acessar o banco de dados"})
		return
	}

	// Se a tarefa for encontrada, tenta deletá-la
	if err := database.DB.Delete(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao deletar a tarefa"})
		return
	}

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
	if err := database.DB.First(&task, id).Error; err != nil {
		// Se a tarefa não for encontrada, retorna 404
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Tarefa não encontrada"})
			return
		}

		// Caso ele entre no if mas não seja um erro de tarefa não encontrada, é um erro ao acessar o db
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao acessar o banco de dados"})
		return
	}

	// Variavel com o novo valor da tarefa
	var updatedTask models.Task
	if err := c.ShouldBindJSON(&updatedTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	if updatedTask.Title == "" || updatedTask.Description == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Título e descrição são obrigatórios"})
		return
	}

	// Atualiza a tarefa com os novos dados
	database.DB.Model(&task).UpdateColumns(updatedTask)

	// Retorna a tarefa atualizada
	c.JSON(http.StatusOK, task)

}

// Função que marca uma tarefa como concluída
func MarkTaskAsCompleted(c *gin.Context) {
	var task models.Task

	// Pegar o ID passado na URL
	id := c.Params.ByName("id")

	// Procurar a tarefa com o ID passado
	if err := database.DB.First(&task, id).Error; err != nil {
		// Se a tarefa não for encontrada, retorna 404
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Tarefa não encontrada"})
			return
		}
		// Caso ele entre no if mas não seja um erro de tarefa não encontrada, é um erro ao acessar o db
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao acessar o banco de dados"})
		return
	}

	if task.Status == "completed" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tarefa já foi concluída"})
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
	if err := database.DB.Where("status = ?", "not completed").Find(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao acessar o banco de dados"})
		return
	}

	// Retorna as tarefas encontradas
	c.JSON(http.StatusOK, task)
}
