package repository

import "TaskManagement/app/models"

type MockTaskRepository struct {
	Tasks map[uint]*models.Task
}

func (m *MockTaskRepository) Create(task *models.Task) error {
	m.Tasks[task.ID] = task
	return nil
}

func (m *MockTaskRepository) Delete(id uint) error {
	delete(m.Tasks, id)
	return nil
}

func (m *MockTaskRepository) FindAll() ([]models.Task, error) {
	var tasks []models.Task
	for _, task := range m.Tasks {
		tasks = append(tasks, *task)
	}
	return tasks, nil
}

func (m *MockTaskRepository) FindByID(id uint) (*models.Task, error) {
	task, exists := m.Tasks[id]
	if !exists {
		return nil, nil // ou um erro, dependendo da sua lógica
	}
	return task, nil
}

func (m *MockTaskRepository) Update(task *models.Task) error {
	m.Tasks[task.ID] = task
	return nil
}
