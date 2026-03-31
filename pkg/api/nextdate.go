package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const DateFormat = "20060102"

// NextDate вычисляет следующую дату выполнения задачи
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("repeat is empty")
	}

	date, err := time.Parse(DateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("invalid date format: %s", dstart)
	}

	nowDate, _ := time.Parse(DateFormat, now.Format(DateFormat))
	parts := strings.Split(repeat, " ")

	switch parts[0] {
	case "d":
		if len(parts) != 2 {
			return "", fmt.Errorf("invalid d format")
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil || days < 1 || days > 400 {
			return "", fmt.Errorf("invalid days count")
		}
		for {
			date = date.AddDate(0, 0, days)
			if date.After(nowDate) {
				break
			}
		}
		return date.Format(DateFormat), nil

	case "y":
		if len(parts) != 1 {
			return "", fmt.Errorf("invalid y format")
		}
		for {
			date = date.AddDate(1, 0, 0)
			if date.Month() == time.February && date.Day() == 29 {
				date = time.Date(date.Year(), time.March, 1, 0, 0, 0, 0, time.UTC)
			}
			if date.After(nowDate) {
				break
			}
		}
		return date.Format(DateFormat), nil

	case "w":
		// w дни недели 1-7 через запятую
		if len(parts) != 2 {
			return "", fmt.Errorf("invalid w format")
		}
		days := strings.Split(parts[1], ",")
		allowedDays := make(map[int]bool)
		for _, d := range days {
			n, err := strconv.Atoi(strings.TrimSpace(d))
			if err != nil || n < 1 || n > 7 {
				return "", fmt.Errorf("invalid weekday")
			}
			allowedDays[n] = true
		}

		for {
			date = date.AddDate(0, 0, 1)
			weekday := int(date.Weekday())
			if weekday == 0 {
				weekday = 7 // Воскресенье = 7
			}
			if allowedDays[weekday] && date.After(nowDate) {
				break
			}
		}
		return date.Format(DateFormat), nil

	case "m":
		// m дни месяца 1-31, -1, -2 [месяцы 1-12]
		if len(parts) < 2 {
			return "", fmt.Errorf("invalid m format")
		}
		daysStr := strings.Split(parts[1], ",")
		var days []int
		var lastDay, preLastDay bool
		for _, d := range daysStr {
			n, err := strconv.Atoi(strings.TrimSpace(d))
			if err != nil {
				return "", fmt.Errorf("invalid day")
			}
			if n == -1 {
				lastDay = true
			} else if n == -2 {
				preLastDay = true
			} else if n >= 1 && n <= 31 {
				days = append(days, n)
			} else {
				return "", fmt.Errorf("invalid day")
			}
		}

		allowedMonths := make(map[int]bool)
		if len(parts) > 2 {
			monthsStr := strings.Split(parts[2], ",")
			for _, m := range monthsStr {
				n, err := strconv.Atoi(strings.TrimSpace(m))
				if err != nil || n < 1 || n > 12 {
					return "", fmt.Errorf("invalid month")
				}
				allowedMonths[n] = true
			}
		} else {
			// Все месяцы разрешены
			for i := 1; i <= 12; i++ {
				allowedMonths[i] = true
			}
		}

		for {
			date = date.AddDate(0, 0, 1)
			year, month, day := date.Date()

			if !allowedMonths[int(month)] {
				continue
			}

			lastDayOfMonth := time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()

			valid := false
			for _, d := range days {
				if d == day {
					valid = true
					break
				}
			}
			if lastDay && day == lastDayOfMonth {
				valid = true
			}
			if preLastDay && day == lastDayOfMonth-1 {
				valid = true
			}

			if valid && date.After(nowDate) {
				break
			}
		}
		return date.Format(DateFormat), nil

	default:
		return "", fmt.Errorf("invalid repeat format")
	}
}

// nextDateHandler обрабатывает GET /api/nextdate
func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	var now time.Time
	var err error
	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(DateFormat, nowStr)
		if err != nil {
			http.Error(w, "invalid now format", http.StatusBadRequest)
			return
		}
	}

	nextDate, err := NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte(nextDate))
}
