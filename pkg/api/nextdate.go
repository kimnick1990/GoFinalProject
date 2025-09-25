package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const dateFormat = "20060102"
const maxInterval = 400 // Максимальный допустимый интервал в днях

// NextDate вычисляет следующую дату на основе правила повторения
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	// Проверка на пустой параметр repeat
	if repeat == "" {
		return "", fmt.Errorf("пустое правило повторения")
	}
	// Парсинг начальной даты
	date, err := time.Parse(dateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("некорректная дата: %v", err)
	}

	// Разбиение правила повторения
	parts := strings.Split(repeat, " ")
	if len(parts) == 0 || parts[0] == "" {
		return "", fmt.Errorf("пустое правило повторения")
	}

	var interval int
	switch parts[0] {
	case "y":
		for {
			date = date.AddDate(1, 0, 0)
			if afterNow(date, now) {
				break
			}
		}
	case "d":
		if len(parts) < 2 {
			return "", fmt.Errorf("отсутствует интервал для правила 'd'")
		}
		interval, err = strconv.Atoi(parts[1])
		if err != nil {
			return "", fmt.Errorf("некорректный интервал: %v", err)
		}
		if interval > maxInterval {
			return "", fmt.Errorf("превышен максимально допустимый интервал")
		}
		for {
			date = date.AddDate(0, 0, interval)
			if afterNow(date, now) {
				break
			}
		}
	default:
		return "", fmt.Errorf("недопустимое правило повторения")
	}

	return date.Format(dateFormat), nil
}

// afterNow проверяет, больше ли первая дата второй без учёта времени
func afterNow(date, now time.Time) bool {
	return date.After(now)
}

// nextDayHandler обрабатывает GET-запросы к /api/nextdate
func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	// Получение параметров запроса
	nowParam := r.FormValue("now")
	dstart := r.FormValue("date")
	repeat := r.FormValue("repeat")

	var now time.Time
	if nowParam == "" {
		now = time.Now()
	} else {
		var err error
		now, err = time.Parse(dateFormat, nowParam)
		if err != nil {
			http.Error(w, "некорректная дата 'now'", http.StatusBadRequest)
			return
		}
	}

	// Вычисление следующей даты
	nextDate, err := NextDate(now, dstart, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, nextDate)
}
