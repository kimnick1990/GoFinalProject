package api

import (
	"net/http"
)

// Init регистрирует обработчики API
func Init() {
	http.HandleFunc("/api/nextdate", NextDayHandler) // Регистрируем обработчик api/nextdate

	http.HandleFunc("/api/task", AddTaskHandler) // Регистрируем обработчик /api/task

	http.HandleFunc("/api/tasks", TasksHandler) // Регистрируем обработчик для /api/tasks

	http.HandleFunc("/api/task/update", UpdateTaskHandler) // Для PUT запроса

	http.HandleFunc("/api/task/", GetTaskHandler) // Для GET запроса
}
