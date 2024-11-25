package main

import (
	"TaskManagement/app/database"
	"TaskManagement/app/repository"
	"TaskManagement/app/routes"
)

func main() {
	database.DataBaseConnection()

	// Cria uma instância do repositório GORM
	repo := &repository.GormTaskRepository{DB: database.DB}
	routes.HandleRequest(repo) // Passa o repositório para as rotas
}
