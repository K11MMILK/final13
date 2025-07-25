package api

import (
	"final13/pkg/db"
	"net/http"
)

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "id not specified")
		return
	}

	err := db.DeleteTask(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete task")
		return
	}

	writeJSON(w, map[string]string{})
}
