package routes

import (
	"TaskManagement/app/database"

	controller "TaskManagement/app/controllers"

	"github.com/gin-gonic/gin"
)

// Função que lida com as requisições
func HandleRequest(repo database.TaskRepository) {
	r := gin.Default()
	controller := &controller.TaskController{Repo: repo} // Passa o repositório para o controlador

	r.GET("/task", controller.GetTask)
	r.POST("/task", controller.CreateTask)
	r.GET("/task/:id", controller.SearchTask)
	r.DELETE("/task/:id", controller.DeleteTask)
	r.PATCH("/task/:id", controller.EditTask)
	r.PUT("/task/:id/complete", controller.MarkTaskAsCompleted)
	r.GET("/task/pending", controller.GetPendingTasks)
	r.Run()
}
