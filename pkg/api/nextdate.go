package api

import (
	scheduler "final13/pkg/scheduler"
	"net/http"
	"time"
)

func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	dstart := query.Get("date")
	repeat := query.Get("repeat")
	nowStr := query.Get("now")

	var now time.Time
	var err error
	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.Parse("20060102", nowStr)
		if err != nil {
			http.Error(w, "invalid now format", http.StatusBadRequest)
			return
		}
	}

	next, err := scheduler.NextDate(now, dstart, repeat, true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(next))
}
