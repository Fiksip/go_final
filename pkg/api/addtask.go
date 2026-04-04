package api

import (
	"encoding/json"
	"fmt"
	"go_final/pkg/db"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func writeJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// addTaskHandler обрабатывает добавление задачи
func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}

	// Проверка заголовка
	if task.Title == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "title is required"})
		return
	}

	// Проверка даты
	if err := checkDate(&task); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// Добавление в БД
	id, err := db.AddTask(&task)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	// Возвращаем ID
	writeJSON(w, http.StatusOK, map[string]string{"id": fmt.Sprintf("%d", id)})
}

// checkDate проверяет и корректирует дату задачи
func checkDate(task *db.Task) error {
	now := time.Now()
	today := now.Format(DateFormat)

	// Специальное значение "today" или пустая строка - всегда сегодняшняя дата
	if task.Date == "today" || task.Date == "" {
		task.Date = today
		// Проверяем валидность правила, но не меняем дату
		if err := validateRepeat(task.Repeat); err != nil {
			return err
		}
		return nil
	}

	// Парсим указанную дату
	t, err := time.Parse(DateFormat, task.Date)
	if err != nil {
		return fmt.Errorf("invalid date format")
	}

	// Начало сегодняшнего дня
	todayStart, _ := time.Parse(DateFormat, today)

	// Если дата меньше сегодняшней
	if t.Before(todayStart) {
		if task.Repeat == "" {
			task.Date = today
		} else {
			next, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				return err
			}
			task.Date = next
		}
	} else {
		// Дата в будущем или сегодня
		if err := validateRepeat(task.Repeat); err != nil {
			return err
		}
	}

	return nil
}

// validateRepeat проверяет формат правила без вычисления даты
func validateRepeat(repeat string) error {
	if repeat == "" {
		return nil
	}

	parts := strings.Split(repeat, " ")
	if len(parts) == 0 {
		return fmt.Errorf("invalid repeat format")
	}

	switch parts[0] {
	case "d":
		if len(parts) != 2 {
			return fmt.Errorf("invalid d format")
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil || days < 1 || days > 400 {
			return fmt.Errorf("invalid days count")
		}
	case "y":
		if len(parts) != 1 {
			return fmt.Errorf("invalid y format")
		}
	case "w":
		if len(parts) != 2 {
			return fmt.Errorf("invalid w format")
		}
		days := strings.Split(parts[1], ",")
		for _, d := range days {
			n, err := strconv.Atoi(strings.TrimSpace(d))
			if err != nil || n < 1 || n > 7 {
				return fmt.Errorf("invalid weekday")
			}
		}
	case "m":
		if len(parts) < 2 {
			return fmt.Errorf("invalid m format")
		}
		daysStr := strings.Split(parts[1], ",")
		// Проверяет, что список дней не пуст
		hasValidDay := false
		for _, d := range daysStr {
			trimmed := strings.TrimSpace(d)
			if trimmed == "" {
				continue
			}
			n, err := strconv.Atoi(trimmed)
			if err != nil {
				return fmt.Errorf("invalid day")
			}
			if n != -1 && n != -2 && (n < 1 || n > 31) {
				return fmt.Errorf("invalid day")
			}
			hasValidDay = true
		}
		if !hasValidDay {
			return fmt.Errorf("invalid m format: no valid days")
		}
		if len(parts) > 2 {
			monthsStr := strings.Split(parts[2], ",")
			for _, m := range monthsStr {
				n, err := strconv.Atoi(strings.TrimSpace(m))
				if err != nil || n < 1 || n > 12 {
					return fmt.Errorf("invalid month")
				}
			}
		}
	default:
		return fmt.Errorf("invalid repeat format")
	}

	return nil
}
