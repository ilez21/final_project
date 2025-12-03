package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const DateFormat = "20060102"

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("repeat is empty")
	}

	date, err := time.Parse(DateFormat, dstart)
	if err != nil {
		return "", err
	}

	parts := strings.Split(repeat, " ")
	rule := parts[0]

	switch rule {
	case "d":
		if len(parts) < 2 {
			return "", errors.New("invalid format: days not specified")
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil || days < 1 || days > 400 {
			return "", errors.New("invalid days interval")
		}
		for {
			date = date.AddDate(0, 0, days)
			if afterNow(date, now) {
				break
			}
		}

	case "y":
		for {
			date = date.AddDate(1, 0, 0)
			if afterNow(date, now) {
				break
			}
		}

	case "w":
		if len(parts) < 2 {
			return "", errors.New("invalid format: weekdays not specified")
		}
		weekdays := make(map[int]bool)
		for _, d := range strings.Split(parts[1], ",") {
			day, err := strconv.Atoi(d)
			if err != nil || day < 1 || day > 7 {
				return "", errors.New("invalid weekday")
			}
			weekdays[day] = true
		}
		date = date.AddDate(0, 0, 1)
		for {
			wd := int(date.Weekday())
			if wd == 0 {
				wd = 7
			}
			if weekdays[wd] && afterNow(date, now) {
				break
			}
			date = date.AddDate(0, 0, 1)
		}

	case "m":
		if len(parts) < 2 {
			return "", errors.New("invalid format: days not specified")
		}
		daysPart := strings.Split(parts[1], ",")
		days := make(map[int]bool)
		for _, d := range daysPart {
			day, err := strconv.Atoi(d)
			if err != nil || day < -2 || day == 0 || day > 31 {
				return "", errors.New("invalid day of month")
			}
			days[day] = true
		}

		var months map[int]bool
		if len(parts) > 2 {
			months = make(map[int]bool)
			for _, m := range strings.Split(parts[2], ",") {
				month, err := strconv.Atoi(m)
				if err != nil || month < 1 || month > 12 {
					return "", errors.New("invalid month")
				}
				months[month] = true
			}
		}

		date = date.AddDate(0, 0, 1)
		for {
			if months != nil && !months[int(date.Month())] {
				date = date.AddDate(0, 0, 1)
				continue
			}

			lastDay := time.Date(date.Year(), date.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()
			day := date.Day()

			match := false
			if days[day] {
				match = true
			} else if days[-1] && day == lastDay {
				match = true
			} else if days[-2] && day == lastDay-1 {
				match = true
			}

			if match && afterNow(date, now) {
				break
			}
			date = date.AddDate(0, 0, 1)
		}

	default:
		return "", errors.New("unsupported repeat format")
	}

	return date.Format(DateFormat), nil
}

func afterNow(date, now time.Time) bool {
	y1, m1, d1 := date.Date()
	y2, m2, d2 := now.Date()
	if y1 > y2 {
		return true
	}
	if y1 == y2 && m1 > m2 {
		return true
	}
	if y1 == y2 && m1 == m2 && d1 > d2 {
		return true
	}
	return false
}

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	var now time.Time
	if nowStr == "" {
		now = time.Now()
	} else {
		var err error
		now, err = time.Parse(DateFormat, nowStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	result, err := NextDate(now, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprint(w, result)
}
