package main

import (
	"TaskManagement/database"
	"TaskManagement/routes"
)

func main() {
	database.DataBaseConnection()

	// Chama a função HandleRequest
	routes.HandleRequest()
}
