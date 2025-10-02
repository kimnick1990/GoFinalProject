package api

import (
	"net/http"
)

// Init регистрирует обработчики API
func Init() {

	http.HandleFunc("GET /api/nextdate", NextDayHandler) // Регистрируем обработчик /api/nextdate

	http.HandleFunc("POST /api/task", AddTaskHandler) // Регистрируем обработчик /api/task

	http.HandleFunc("GET /api/tasks", TasksHandler) // Регистрируем обработчик для /api/tasks

	http.HandleFunc("PUT /api/task", UpdateTaskHandler) // Для PUT запроса /api/task

	http.HandleFunc("GET /api/task", GetTaskHandler) // Для GET запроса /api/task

	http.HandleFunc("POST /api/task/done", DoneTaskHandler) // Для POST запроса /api/task/done

	http.HandleFunc("DELETE /api/task", DeleteTaskHandler) // Для DELETE запроса /api/task
}
