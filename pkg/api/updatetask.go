package api

import (
	"encoding/json"
	"go_final/pkg/db"
	"net/http"
	"time"
)

// updateTaskHandler обрабатывает PUT-запрос /api/task
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}

	// Проверка ID
	if task.ID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "task id is required"})
		return
	}

	// Проверка заголовка
	if task.Title == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "title is required"})
		return
	}

	// Проверка даты
	now := time.Now()
	today := now.Format(DateFormat)

	if task.Date == "today" || task.Date == "" {
		task.Date = today
	}

	t, err := time.Parse(DateFormat, task.Date)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid date format"})
		return
	}

	todayStart, _ := time.Parse(DateFormat, today)

	if t.Before(todayStart) {
		if task.Repeat == "" {
			task.Date = today
		} else {
			next, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
				return
			}
			task.Date = next
		}
	} else {
		if task.Repeat != "" {
			if err := validateRepeat(task.Repeat); err != nil {
				writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
				return
			}
		}
	}

	// Обновление в БД
	if err := db.UpdateTask(&task); err != nil {
		if err.Error() == "task not found" {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	// Успешный ответ пустой JSON
	writeJSON(w, http.StatusOK, map[string]string{})
}
