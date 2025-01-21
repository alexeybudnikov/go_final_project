package service

import (
	"errors"
	"time"

	"github.com/alexeybudnikov/go_final_project/internal/models"
	"github.com/alexeybudnikov/go_final_project/internal/repository"
	"github.com/alexeybudnikov/go_final_project/utils"
)

type TaskService interface {
	CreateTask(task models.Task) (int64, error)
	GetAllTasks() ([]models.Task, error)
	GetTaskByID(id int64) (models.Task, error)
	UpdateTaskByID(task models.Task) error
	DoneTask(id int64) error
	DeleteTask(id int64) error
}

type taskService struct {
	repo repository.TaskRepository
}

func NewTaskService(repo repository.TaskRepository) TaskService {
	return &taskService{repo: repo}
}

func (s *taskService) CreateTask(task models.Task) (int64, error) {
	var taskDate time.Time
	// title обязательное
	if task.Title == "" {
		return 0, errors.New("не указан заголовок задачи")
	}

	// если date не указано или пустая строка - берем сегодняшнюю дату
	if task.Date == "" {
		task.Date = time.Now().Format(utils.DateFotmat)
	}
	// проверяем что дата парсится и имеет формат 20060102
	taskDate, err := time.Parse(utils.DateFotmat, task.Date)
	if err != nil {
		return 0, errors.New("дата представлена в формате, отличном от 20060102")
	}
	// если дата меньше сегодня, то
	// 1) если правило не указано или равно пустой строке - указываем сегодняшнюю дату
	// 2) если правило указано то записываем то что больше сегодняшнего числа nextDate()
	today := time.Now().Format(utils.DateFotmat)
	todayFormated, err := time.Parse(utils.DateFotmat, today)
	if err != nil {
		return 0, err
	}
	if taskDate.Before(todayFormated) {
		if task.Repeat == "" {
			task.Date = todayFormated.Format(utils.DateFotmat)
			taskDate = todayFormated
		} else {
			nextDateStr, err := utils.NextDate(todayFormated, task.Date, task.Repeat)
			if err != nil {
				return 0, err
			}
			taskDate, err = time.Parse(utils.DateFotmat, nextDateStr)
			if err != nil {
				return 0, errors.New("дата представлена в формате, отличном от 20060102")
			}
		}
	}

	task.Date = taskDate.Format(utils.DateFotmat)
	return s.repo.Create(task)
}

func (s *taskService) GetAllTasks() ([]models.Task, error) {
	return s.repo.GetAll()
}

func (s *taskService) GetTaskByID(id int64) (models.Task, error) {
	return s.repo.GetByID(id)
}

func (s *taskService) UpdateTaskByID(task models.Task) error {
	if task.Title == "" {
		return errors.New("не указан заголовок задачи")
	}

	// если date не указано или пустая строка - берем сегодняшнюю дату
	if task.Date == "" {
		task.Date = time.Now().Format(utils.DateFotmat)
	}
	// проверяем что дата парсится и имеет формат 20060102
	taskDate, err := time.Parse(utils.DateFotmat, task.Date)
	if err != nil {
		return errors.New("дата представлена в формате, отличном от 20060102")
	}

	// если дата меньше сегодня, то
	// 1) если правило не указано или равно пустой строке - указываем сегодняшнюю дату
	// 2) если правило указано то записываем то что больше сегодняшнего числа nextDate()
	today := time.Now()
	if taskDate.Before(today) {
		if task.Repeat == "" {
			taskDate = today
		} else {
			nextDateStr, err := utils.NextDate(today, task.Date, task.Repeat)
			if err != nil {
				return err
			}
			taskDate, err = time.Parse(utils.DateFotmat, nextDateStr)
			if err != nil {
				return errors.New("дата представлена в формате, отличном от 20060102")
			}
		}
	}

	var taskEdited = models.Task{
		ID:      task.ID,
		Date:    taskDate.Format(utils.DateFotmat),
		Title:   task.Title,
		Comment: task.Comment,
		Repeat:  task.Repeat,
	}
	return s.repo.Update(taskEdited)
}

func (s *taskService) DoneTask(id int64) error {
	// без повторения просто удаляем
	// с повторением ставим следующую задачу и удаляем текущую
	task, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if task.Repeat == "" {
		return s.repo.Delete(id)
	} else {
		today := time.Now()
		newTaskDate, err := utils.NextDate(today, task.Date, task.Repeat)
		if err != nil {
			return err
		}
		task.Date = newTaskDate
		return s.repo.Update(task)
	}
}

func (s *taskService) DeleteTask(id int64) error {
	task, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(task.ID)
}
