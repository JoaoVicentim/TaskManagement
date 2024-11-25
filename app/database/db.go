package database

import (
	"TaskManagement/app/models"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type TaskRepository interface {
	Create(task *models.Task) error
	Delete(id uint) error
	FindAll() ([]models.Task, error)
	FindByID(id uint) (*models.Task, error)
	Update(task *models.Task) error
}

var (
	DB  *gorm.DB // Armazena a conexão com o banco de dados
	err error    // Armazena o erro retornado pela função Open
)

// Função que estabelece a conexão com o banco de dados
func DataBaseConnection() {
	stringconnection := "host=postgres user=root password=root dbname=root port=5432 sslmode=disable" // String de conexão com banco de dados

	// Tentar conectar ao banco de dados 10 vezes
	for i := 0; i < 10; i++ {
		DB, err = gorm.Open(postgres.Open(stringconnection)) // Tentativa de conexão com o banco de dados
		if err == nil {
			log.Println("Conexão com o banco de dados estabelecida com sucesso!")
			break
		}
		log.Printf("Falha ao conectar ao banco de dados, tentativa %d: %v", i+1, err)
		time.Sleep(2 * time.Second) // Espera 2 segundos antes da próxima tentativa
	}
	if err != nil {
		log.Panic("Não foi possível conectar ao banco de dados após 10 tentativas.")
	}

	// Criar uma tabela no banco de dados com base no modelo de Tarefa
	DB.AutoMigrate(&models.Task{})
}
