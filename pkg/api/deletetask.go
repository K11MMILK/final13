package api

import (
	"final13/pkg/db"
	"net/http"
)

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, map[string]string{"error": "id not specified"})
		return
	}
	err := db.DeleteTask(id)
	if err != nil {
		writeJSON(w, map[string]string{"error": "failed to delete task"})
		return
	}
	writeJSON(w, map[string]string{})
}
