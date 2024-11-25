package repository

import (
	"TaskManagement/app/models"

	"gorm.io/gorm"
)

type GormTaskRepository struct {
	DB *gorm.DB
}

func NewGormTaskRepository(db *gorm.DB) *GormTaskRepository {
	return &GormTaskRepository{DB: db}
}

func (r *GormTaskRepository) Create(task *models.Task) error {
	return r.DB.Create(task).Error
}

func (r *GormTaskRepository) FindByID(id uint) (*models.Task, error) {
	var task models.Task
	if err := r.DB.First(&task, id).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *GormTaskRepository) Delete(id uint) error {
	return r.DB.Delete(&models.Task{}, id).Error
}

func (r *GormTaskRepository) Update(task *models.Task) error {
	return r.DB.Save(task).Error
}

func (r *GormTaskRepository) FindAll() ([]models.Task, error) {
	var tasks []models.Task
	if err := r.DB.Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}
