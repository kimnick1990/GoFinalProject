package api

import (
	"net/http"
)

// Init регистрирует обработчики API
func Init() {
	http.HandleFunc("/api/nextdate", nextDayHandler)
}
