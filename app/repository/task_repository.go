package repository

import "TaskManagement/app/models"

type TaskRepository interface {
	Create(task *models.Task) error
	FindByID(id uint) (*models.Task, error)
	Delete(id uint) error
	Update(task *models.Task) error
	FindAll() ([]models.Task, error)
}
