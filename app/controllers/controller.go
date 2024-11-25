package controller

import (
	"TaskManagement/app/database"
	"TaskManagement/app/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Função que cria uma nova tarefa
type TaskController struct {
	Repo database.TaskRepository
}

// Função que retorna todas as tarefas
func (c *TaskController) GetTask(ctx *gin.Context) {
	tasks, err := c.Repo.FindAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao acessar o banco de dados"})
		return
	}
	ctx.JSON(http.StatusOK, tasks)
}

func (c *TaskController) CreateTask(ctx *gin.Context) {
	var task models.Task
	if err := ctx.ShouldBindJSON(&task); err != nil || task.Title == "" || task.Description == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Título e descrição são obrigatórios"})
		return
	}
	c.Repo.Create(&task)
	ctx.JSON(http.StatusOK, task)
}

// Função que busca uma tarefa pelo ID
func (c *TaskController) SearchTask(ctx *gin.Context) {
	idStr := ctx.Params.ByName("id")            // ID como string
	id, err := strconv.ParseUint(idStr, 10, 32) // Converte para uint
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Procurar a tarefa com o ID passado
	task, err := c.Repo.FindByID(uint(id)) // Passa o ID convertido
	if err != nil {
		// Se a tarefa não for encontrada, retorna 404
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Tarefa não encontrada"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao acessar o banco de dados"})
		return
	}

	// Retorna a tarefa encontrada
	ctx.JSON(http.StatusOK, task)
}

// Função que deleta uma tarefa
func (c *TaskController) DeleteTask(ctx *gin.Context) {
	idStr := ctx.Params.ByName("id")            // ID como string
	id, err := strconv.ParseUint(idStr, 10, 32) // Converte para uint

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Primeiro, tenta encontrar a tarefa pelo ID
	if _, err := c.Repo.FindByID(uint(id)); err != nil {

		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Tarefa não encontrada"})
			return

		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao acessar o banco de dados"})
		return
	}

	// Se a tarefa existir, tenta deletá-la
	if err := c.Repo.Delete(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao deletar a tarefa"})
		return
	}

	// Retorna uma mensagem de sucesso
	ctx.JSON(http.StatusOK, gin.H{"data": "Tarefa deletada"})
}

// Função que edita uma tarefa
func (c *TaskController) EditTask(ctx *gin.Context) {
	idStr := ctx.Params.ByName("id")            // ID como string
	id, err := strconv.ParseUint(idStr, 10, 32) // Converte para uint
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Procurar a tarefa com o ID passado
	task, err := c.Repo.FindByID(uint(id)) // Passa o ID convertido
	// Procurar a tarefa com o ID passado
	if err != nil {
		// Se a tarefa não for encontrada, retorna 404
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Tarefa não encontrada"})
			return
		}

		// Caso ele entre no if mas não seja um erro de tarefa não encontrada, é um erro ao acessar o db
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao acessar o banco de dados"})
		return
	}

	if task == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Tarefa não encontrada"})
		return
	}

	// Variável com o novo valor da tarefa
	var updatedTask models.Task
	if err := ctx.ShouldBindJSON(&updatedTask); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	if updatedTask.Title == "" || updatedTask.Description == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Título e descrição são obrigatórios"})
		return
	}

	// Atualiza a tarefa com os novos dados
	task.Title = updatedTask.Title
	task.Description = updatedTask.Description
	if err := c.Repo.Update(task); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar a tarefa"})
		return
	}

	// Retorna a tarefa atualizada
	ctx.JSON(http.StatusOK, task)
}

// Função que marca uma tarefa como concluída
func (c *TaskController) MarkTaskAsCompleted(ctx *gin.Context) {
	idStr := ctx.Params.ByName("id")            // ID como string
	id, err := strconv.ParseUint(idStr, 10, 32) // Converte para uint
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Procurar a tarefa com o ID passado
	task, err := c.Repo.FindByID(uint(id)) // Passa o ID convertido
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Tarefa não encontrada"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao acessar o banco de dados"})
		return
	}

	if task.Status == "completed" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Tarefa já foi concluída"})
		return
	}

	// Atualiza o status da tarefa para "completed"
	task.Status = "completed"
	if err := c.Repo.Update(task); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar a tarefa"})
		return
	}

	// Retorna a tarefa atualizada
	ctx.JSON(http.StatusOK, task)
}

// Função que retorna todas as tarefas pendentes
func (c *TaskController) GetPendingTasks(ctx *gin.Context) {
	tasks, err := c.Repo.FindAll() // Busca todas as tarefas
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao acessar o banco de dados"})
		return
	}

	// Filtra as tarefas pendentes
	var pendingTasks []models.Task
	for _, task := range tasks {
		if task.Status == "not completed" {
			pendingTasks = append(pendingTasks, task)
		}
	}

	// Retorna as tarefas encontradas
	ctx.JSON(http.StatusOK, pendingTasks)
}
