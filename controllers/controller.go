package controller

import (
	"TaskManagement/database"
	"TaskManagement/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetTask(c *gin.Context) {
	var task []models.Task
	database.DB.Find(&task)
	c.JSON(200, task)
}

func CreateTask(c *gin.Context) {
	var task models.Task
	// Verifica se o JSON é válido e o armazena em task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if task.Status == "" {
		task.Status = "not completed"
	}

	database.DB.Create(&task)
	c.JSON(http.StatusOK, task)
}

func SearchTask(c *gin.Context) {
	var task models.Task

	// pegar o ID passado
	id := c.Params.ByName("id")

	// Procurar a tarefa com o ID passado
	database.DB.First(&task, id)

	// Verificaç˜ao se a tarefa foi encontrada
	if task.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"Not Found": "Tarefa não encontrada"})
		return
	}

	// Mostrar tarefa encontrada
	c.JSON(http.StatusOK, task)
}

func DeleteTask(c *gin.Context) {
	var task models.Task

	// extrai o ID da recebido
	id := c.Params.ByName("id")

	database.DB.Delete(&task, id)

	c.JSON(http.StatusOK, gin.H{"data": "Tarefa deletada"})

}

func EditTask(c *gin.Context) {
	var task models.Task

	id := c.Params.ByName("id")

	database.DB.First(&task, id)

	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.DB.Model(&task).UpdateColumns(task)

	c.JSON(http.StatusOK, task)

}

func MarkTaskAsCompleted(c *gin.Context) {
	var task models.Task

	id := c.Params.ByName("id")

	database.DB.First(&task, id)

	if task.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	task.Status = "completed"

	database.DB.Save(&task)

	c.JSON(http.StatusOK, task)
}

func GetPendingTasks(c *gin.Context) {
	var task models.Task

	// Procurar todas as tarefas que não foram completadas
	database.DB.Where("status = ?", "not completed").Find(&task)

	c.JSON(http.StatusOK, task)
}

// func SearchTaskByTitle(c *gin.Context) {
// 	var task models.Task
// 	title := c.Param("title")

// 	// Procurar a tarefa com o título passado e armazenar em task
// 	database.DB.Where(&models.Task{Title: title}).First(&task)

// 	if task.ID == 0 {
// 		c.JSON(http.StatusNotFound, gin.H{"Not Found": "Tarefa não encontrada"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, task)
// }
