package api

import (
	"final13/pkg/db"
	"final13/pkg/scheduler"
	"net/http"
	"time"
)

func taskDoneHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "id not specified")
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeError(w, http.StatusNotFound, "task not found")
		return
	}

	if task.Repeat == "" {
		if err := db.DeleteTask(id); err != nil {
			writeError(w, http.StatusInternalServerError, "failed to delete task")
			return
		}
		writeJSON(w, map[string]string{})
		return
	}

	now, err := parseDateString(task.Date)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid task date")
		return
	}

	now = truncateToDate(now)

	next, err := scheduler.NextDate(now, task.Date, task.Repeat, true)
	if err != nil {
		writeError(w, http.StatusBadRequest, "failed to calculate next date")
		return
	}

	if err := db.UpdateDate(next, id); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update date")
		return
	}

	writeJSON(w, map[string]string{})
}

func truncateToDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}

func parseDateString(s string) (time.Time, error) {
	return time.Parse("20060102", s)
}
