package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"final13/pkg/db"
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		writeJSON(w, map[string]string{"error": "Invalid JSON"})
		return
	}

	if task.Title == "" {
		writeJSON(w, map[string]string{"error": "Task title is required"})
		return
	}

	if err := checkDate(&task); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeJSON(w, map[string]string{"error": "Error while adding task"})
		return
	}

	writeJSON(w, map[string]string{"id": strconv.FormatInt(id, 10)})
}
