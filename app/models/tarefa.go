package models

import "gorm.io/gorm"

type Task struct {
	gorm.Model
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

// hook do GORM que é executado antes de criar uma nova tarefa no banco de dados
func (task *Task) BeforeCreate(tx *gorm.DB) (err error) {
	// Define o status da tarefa como "not completed" se não for especificado
	if task.Status == "" {
		task.Status = "not completed"
	}
	return
}
