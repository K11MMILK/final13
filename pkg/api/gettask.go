package api

import (
	"final13/pkg/db"
	"net/http"
	"time"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func getTasksHandler(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")

	var tasks []*db.Task
	var err error

	if search != "" {
		if parsed, errDate := time.Parse("02.01.2006", search); errDate == nil {
			search = parsed.Format("20060102")
		}
		tasks, err = db.TasksSearch(search, 50)
	} else {
		tasks, err = db.Tasks(50)
	}

	if err != nil {
		writeError(w, http.StatusInternalServerError, "error reading tasks")
		return
	}

	writeJSON(w, TasksResp{Tasks: tasks})
}
