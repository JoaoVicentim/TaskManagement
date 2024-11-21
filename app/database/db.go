package database

import (
	"TaskManagement/app/models"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB  *gorm.DB
	err error
)

func DataBaseConnection() {
	stringconnection := "host=postgres user=root password=root dbname=root port=5432 sslmode=disable"
	DB, err = gorm.Open(postgres.Open(stringconnection))
	if err != nil {
		log.Panic("Falha ao conectar ao banco de dados")
	}
	// Criar uma tabela no banco de dados com base no modelo de Tarefa
	DB.AutoMigrate(&models.Task{})
}
