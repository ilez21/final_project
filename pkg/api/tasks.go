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
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	search := r.FormValue("search")

	if search != "" {
		if t, err := time.Parse("02.01.2006", search); err == nil {
			search = t.Format(DateFormat)
		}
	}

	tasks, err := db.Tasks(50, search)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, TasksResp{Tasks: tasks})
}
