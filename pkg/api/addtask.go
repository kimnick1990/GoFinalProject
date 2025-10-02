package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/kimnick1990/GoFinalProject/pkg/db"
)

// Хэндлер GETPOST запроса к /api/nextdate
func AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var task db.Task

	// Десериализация JSON в task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		WriteJson(w, map[string]string{"error": "Ошибка десериализации JSON"})
		return
	}

	// Проверка корректности task.Title
	if task.Title == "" {
		WriteJson(w, map[string]string{"error": "Не указан заголовок задачи"})
		return
	}

	// Проверка даты и правила повторения
	err = checkDate(&task)
	if err != nil {
		WriteJson(w, map[string]string{"error": err.Error()})
		return
	}

	// Вызов db.AddTask для добавления задачи в БД
	id, err := db.AddTask(&task)
	if err != nil {
		WriteJson(w, map[string]string{"error": err.Error()})
		return
	}

	// Сериализация результата в JSON и отправка ответа
	WriteJson(w, map[string]interface{}{"id": id})
}

func WriteJson(w http.ResponseWriter, data interface{}) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func checkDate(task *db.Task) error {
	now := time.Now().Truncate(24 * time.Hour) // Устанавливаем время на начало дня
	// Проверка формата даты
	if task.Date == "" {
		task.Date = now.Format("20060102") // Если дата пустая, ставим сегодняшнее число
	} else if !isValidDateFormat(task.Date) {
		return fmt.Errorf("некорректный формат даты: %s", task.Date)
	}

	t, err := time.Parse("20060102", task.Date)
	if err != nil {
		return err
	}
	t = t.Truncate(24 * time.Hour) // Устанавливаем время даты задачи на начало дня

	if t.Before(now) { // Если дата в прошлом
		if len(task.Repeat) > 0 {
			next, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				return err
			}
			task.Date = next // Обновляем дату задачи
		} else {
			task.Date = now.Format("20060102") // Ставим сегодняшнюю дату, если нет правила повторения
		}
	}
	// Если дата сегодняшняя или будущая, ничего не делаем

	return nil
}

// Функция для проверки формата даты
func isValidDateFormat(date string) bool {
	_, err := time.Parse("20060102", date)
	return err == nil
}
