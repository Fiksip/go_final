package api

import (
	"go_final/pkg/db"
	"net/http"
)

// TasksResp структура ответа со списком задач
type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

// tasksHandler обрабатывает GET-запрос /api/tasks
func tasksHandler(w http.ResponseWriter, r *http.Request) {

	tasks, err := db.Tasks(50)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, TasksResp{
		Tasks: tasks,
	})
}
