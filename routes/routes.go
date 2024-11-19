package routes

import (
	controller "TaskManagement/controllers"

	"github.com/gin-gonic/gin"
)

// Função que lida com as requisições
func HandleRequest() {
	// Cria um novo roteador padrão
	r := gin.Default()
	r.GET("/task", controller.GetTask)
	r.POST("/task", controller.CreateTask)
	r.GET("/task/:id", controller.SearchTask)
	r.DELETE("/task/:id", controller.DeleteTask)
	r.PATCH("/task/:id", controller.EditTask)
	r.PUT("/task/:id/complete", controller.MarkTaskAsCompleted)
	r.GET("task/pending", controller.GetPendingTasks)
	r.Run()
}
