package main

import (
	"TaskManagement/app/database"
	"TaskManagement/app/routes"
)

func main() {
	database.DataBaseConnection()

	// Chama a função HandleRequest
	routes.HandleRequest()
}
