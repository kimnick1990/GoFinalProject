package api

import (
	"encoding/json"
	"net/http"

	"github.com/kimnick1990/GoFinalProject/pkg/db"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func TasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := db.Tasks(50) // в параметре максимальное количество записей
	if err != nil {
		writeError(w, "error", err.Error())
		return
	}
	WriteJson(w, TasksResp{
		Tasks: tasks,
	})
}

// Новый хендлер для GET запроса /api/task?id=<идентификатор>
func GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeError(w, "error", "Не указан идентификатор")
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeError(w, "error", err.Error())
		return
	}

	WriteJson(w, task)
}

// Новый хендлер для PUT запроса /api/task
func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		writeError(w, "error", "Ошибка при декодировании данных")
		return
	}

	err = db.UpdateTask(&task)
	if err != nil {
		writeError(w, "error", err.Error())
		return
	}

	WriteJson(w, map[string]string{})
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

func writeError(w http.ResponseWriter, errorKey, message string) {
	errorResponse := map[string]string{errorKey: message}
	jsonError, err := json.Marshal(errorResponse)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(jsonError)
}
