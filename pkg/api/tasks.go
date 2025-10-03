package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/kimnick1990/GoFinalProject/pkg/db"
)

const MaxTaskLimit = 50 // Определение константы для максимального количества записей

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

// Хэндлер GET запроса /api/tasks
func TasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := db.Tasks(MaxTaskLimit) // в параметре максимальное количество записей MaxTaskLimit
	if err != nil {
		writeError(w, "error", err.Error(), http.StatusInternalServerError)
		return
	}
	WriteJson(w, TasksResp{
		Tasks: tasks,
	})
}

// Хэндлер GET запроса /api/task?id=<идентификатор>
func GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeError(w, "error", "Не указан идентификатор", http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeError(w, "error", err.Error(), http.StatusInternalServerError)
		return
	}

	WriteJson(w, task)
}

// Хэндлер PUT запроса /api/task
func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		writeError(w, "error", "Ошибка при декодировании данных", http.StatusBadRequest)
		return
	}

	// Проверка корректности task.Title
	if task.Title == "" {
		writeError(w, "error", "Не указан заголовок задачи", http.StatusBadRequest)
		return
	}

	// Проверка даты и правила повторения
	err = checkDate(&task)
	if err != nil {
		writeError(w, "error", err.Error(), http.StatusInternalServerError)
		return
	}

	err = db.UpdateTask(&task)
	if err != nil {
		writeError(w, "error", err.Error(), http.StatusInternalServerError)
		return
	}

	WriteJson(w, map[string]string{})
}

// Хэндлер POST запроса /api/task/done
func DoneTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeError(w, "error", "Идентификатор задачи не указан", http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeError(w, "error", fmt.Sprintf("Ошибка при получении задачи: %v", err), http.StatusInternalServerError)
		return
	}

	if task.Repeat == "" { // Одноразовая задача
		err := db.DeleteTask(id)
		if err != nil {
			writeError(w, "error", fmt.Sprintf("Ошибка при удалении задачи: %v", err), http.StatusInternalServerError)
			return
		}
	} else { // Периодическая задача
		nextDate, err := NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			writeError(w, "error", fmt.Sprintf("Ошибка при расчёте следующей даты: %v", err), http.StatusInternalServerError)
			return
		}
		err = db.UpdateDate(nextDate, id)
		if err != nil {
			writeError(w, "error", fmt.Sprintf("Ошибка при обновлении даты задачи: %v", err), http.StatusInternalServerError)
			return
		}
	}

	WriteJson(w, map[string]string{})
}

// Хэндлер DELETE запроса /api/task
func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeError(w, "error", "Идентификатор задачи не указан", http.StatusBadRequest)
		return
	}

	err := db.DeleteTask(id)
	if err != nil {
		writeError(w, "error", fmt.Sprintf("Ошибка при удалении задачи: %v", err), http.StatusInternalServerError)
		return
	}

	WriteJson(w, map[string]string{})
}

func writeError(w http.ResponseWriter, errorKey, message string, status int) {
	errorResponse := map[string]string{errorKey: message}
	jsonError, err := json.Marshal(errorResponse)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(jsonError)
}
