package api

import (
	"encoding/json"
	"net/http"

	"final13/pkg/db"
)

func editTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		writeJSON(w, map[string]string{"error": "Invalid JSON"})
		return
	}

	if task.ID == "" {
		writeJSON(w, map[string]string{"error": "Task ID not specified"})
		return
	}

	if task.Title == "" {
		writeJSON(w, map[string]string{"error": "Task title not specified"})
		return
	}

	if err := checkDate(&task); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	err = db.UpdateTask(&task)
	if err != nil {
		writeJSON(w, map[string]string{"error": "Task not found"})
		return
	}

	writeJSON(w, struct{}{})
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, map[string]string{"error": "Task ID not specified"})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, map[string]string{"error": "Task not found"})
		return
	}

	writeJSON(w, task)
}
