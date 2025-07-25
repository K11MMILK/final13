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
		writeError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if task.ID == "" {
		writeError(w, http.StatusBadRequest, "Task ID not specified")
		return
	}

	if task.Title == "" {
		writeError(w, http.StatusBadRequest, "Task title not specified")
		return
	}

	if err := checkDate(&task); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = db.UpdateTask(&task)
	if err != nil {
		writeError(w, http.StatusNotFound, "Task not found")
		return
	}

	writeJSON(w, struct{}{})
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "Task ID not specified")
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeError(w, http.StatusNotFound, "Task not found")
		return
	}

	writeJSON(w, task)
}
