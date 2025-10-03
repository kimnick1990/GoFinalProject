package db

import (
	"database/sql"
	"fmt"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// Добавление функции AddTask
func AddTask(task *Task) (int64, error) {
	var id int64
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES ($1, $2, $3, $4)`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err == nil {
		id, err = res.LastInsertId()
	}
	return id, err
}

// Добавление функции Tasks
func Tasks(limit int) ([]*Task, error) {
	rows, err := db.Query(fmt.Sprintf("SELECT ID, Date, Title, Comment, Repeat FROM scheduler ORDER BY date ASC LIMIT %d", limit))
	if err != nil { // Обработка ошибки при выполнении запроса
		return nil, err
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat) // Попытка сканирования данных
		if err != nil {                                                                  // Обработка ошибки при сканировании данных
			return nil, err
		}
		tasks = append(tasks, &task)
	}

	if len(tasks) == 0 {
		return []*Task{}, nil // Возвращаем пустой слайс
	}

	return tasks, nil
}

// Добавление функции GetTask
func GetTask(id string) (*Task, error) {
	var task Task
	err := db.QueryRow("SELECT ID, Date, Title, Comment, Repeat FROM scheduler WHERE id = $1", id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("задача не найдена")
	} else if err != nil {
		return nil, err
	}
	return &task, nil
}

// Добавление функции UpdateTask
func UpdateTask(task *Task) error {
	query := `UPDATE scheduler SET date = $1, title = $2, comment = $3, repeat = $4 WHERE id = $5`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("неверный идентификатор для обновления задачи")
	}
	return nil
}

// Добавление функции DeleteTask
func DeleteTask(id string) error {
	res, err := db.Exec("DELETE FROM scheduler WHERE id = $1", id)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("задача с указанным идентификатором не найдена")
	}
	return nil
}

// Добавление функции UpdateDate
func UpdateDate(nextDate string, id string) error {
	query := "UPDATE scheduler SET date = $1 WHERE id = $2"
	res, err := db.Exec(query, nextDate, id)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("неверный идентификатор для обновления задачи")
	}
	return nil
}
