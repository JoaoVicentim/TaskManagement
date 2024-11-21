package models

import "gorm.io/gorm"

type Task struct {
	gorm.Model
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

func (task *Task) BeforeCreate(tx *gorm.DB) (err error) {
	if task.Status == "" {
		task.Status = "not completed"
	}
	return
}
