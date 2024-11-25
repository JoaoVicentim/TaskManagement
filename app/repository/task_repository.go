package repository

import (
	"TaskManagement/app/models"

	"gorm.io/gorm"
)

type GormTaskRepository struct {
	DB *gorm.DB
}

func (r *GormTaskRepository) Create(task *models.Task) error {
	return r.DB.Create(task).Error
}

func (r *GormTaskRepository) Delete(id uint) error {
	return r.DB.Delete(&models.Task{}, id).Error
}

func (r *GormTaskRepository) FindAll() ([]models.Task, error) {
	var tasks []models.Task
	err := r.DB.Find(&tasks).Error
	return tasks, err
}

func (r *GormTaskRepository) FindByID(id uint) (*models.Task, error) {
	var task models.Task
	err := r.DB.First(&task, id).Error
	return &task, err
}

func (r *GormTaskRepository) Update(task *models.Task) error {
	return r.DB.Save(task).Error
}
