package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/kimnick1990/GoFinalProject/pkg/db"
)

func AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var task db.Task

	// Десериализация JSON в task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		writeJson(w, map[string]string{"error": "Ошибка десериализации JSON"})
		return
	}

	// Проверка корректности task.Title
	if task.Title == "" {
		writeJson(w, map[string]string{"error": "Не указан заголовок задачи"})
		return
	}

	// Проверка даты и правила повторения
	err = checkDate(&task)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	// Вызов db.AddTask для добавления задачи в БД
	id, err := db.AddTask(&task)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	// Сериализация результата в JSON и отправка ответа
	writeJson(w, map[string]interface{}{"id": id})
}

func writeJson(w http.ResponseWriter, data interface{}) {
	jsonData, _ := json.Marshal(data)
	w.Write(jsonData)
}

func checkDate(task *db.Task) error {
	now := time.Now()
	if task.Date == "" {
		task.Date = now.Format("20060102")
	}

	t, err := time.Parse("20060102", task.Date)
	if err != nil {
		return err
	}

	if len(task.Repeat) > 0 {
		next, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
		task.Date = next
	} else if afterNow(now, t) {
		// Если дата меньше текущего времени и нет правила повторения
		task.Date = now.Format("20060102")
	}

	return nil
}
