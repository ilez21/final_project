package api

import (
	"final_sprint/pkg/db"
	"net/http"
	"time"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	search := r.FormValue("search")

	// Проверяем, является ли search датой в формате 02.01.2006
	if search != "" {
		if t, err := time.Parse("02.01.2006", search); err == nil {
			search = t.Format(DateFormat)
		}
	}

	tasks, err := db.Tasks(50, search)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, TasksResp{Tasks: tasks})
}
