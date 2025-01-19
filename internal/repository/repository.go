package repository

import (
	"database/sql"
	"fmt"

	"github.com/alexeybudnikov/go_final_project/internal/models"
)

type TaskRepository interface {
	Create(task models.Task) (int64, error)
	GetAll() ([]models.Task, error)
	GetByID(id int64) (models.Task, error)
	Update(task models.Task) error
	Delete(id int64) error
}

type taskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) Create(task models.Task) (int64, error) {
	query := `INSERT INTO scheduler ( date, title, comment, repeat) VALUES ($1, $2, $3, $4);`

	res, err := r.db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *taskRepository) GetAll() ([]models.Task, error) {
	var tasks []models.Task
	query := `SELECT * FROM scheduler order by date LIMIT 10;`

	res, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	for res.Next() {
		t := models.Task{}

		err := res.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, t)
	}
	err = res.Err()
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *taskRepository) GetByID(id int64) (models.Task, error) {
	var task models.Task
	query := `SELECT * FROM scheduler where id = :id;`

	res := r.db.QueryRow(query, sql.Named("id", id))

	err := res.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return models.Task{}, fmt.Errorf("erorr while getting task: %w", err)
	}

	return task, nil
}

func (r *taskRepository) Update(task models.Task) error {
	_, checkIfExist := r.GetByID(task.ID)
	if checkIfExist != nil {
		return checkIfExist
	}

	query := `UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id;`
	_, err := r.db.Exec(query,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
		sql.Named("id", task.ID))
	if err != nil {
		return err
	}
	return nil
}

func (r *taskRepository) Delete(id int64) error {
	query := `DELETE FROM scheduler WHERE id = :id`
	_, err := r.db.Exec(query, sql.Named("id", id))
	if err != nil {
		return err
	}
	return nil
}
