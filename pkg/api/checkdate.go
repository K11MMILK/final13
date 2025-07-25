package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"final13/pkg/db"
	scheduler "final13/pkg/scheduler"
)

func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("writeJSON error: %v", err)
	}
}

func writeError(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)
	writeJSON(w, map[string]string{
		"error": msg,
	})
}

func checkDate(task *db.Task) error {
	now := time.Now()

	if task.Date == "" {
		task.Date = now.Format("20060102")
	}

	t, err := time.Parse("20060102", task.Date)
	if err != nil {
		return errors.New("invalid date")
	}

	if strings.HasPrefix(task.Repeat, "m ") {
		parts := strings.Fields(task.Repeat)
		if len(parts) == 2 {
			task.Repeat += " 1,2,3,4,5,6,7,8,9,10,11,12"
		}
	}

	if task.Repeat != "" {
		next, err := scheduler.NextDate(now, task.Date, task.Repeat, false)
		if err != nil {
			return errors.New("error in repeat rule")
		}

		if !t.After(now) {
			task.Date = next
		}
	} else {
		if !t.After(now) {
			task.Date = now.Format("20060102")
		}
	}

	return nil
}
