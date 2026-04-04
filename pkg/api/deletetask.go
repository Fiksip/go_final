package api

import (
	"go_final/pkg/db"
	"net/http"
)

// deleteTaskHandler обрабатывает DELETE-запрос /api/task?id=<id>
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "task id is required"})
		return
	}

	if err := db.DeleteTask(id); err != nil {
		if err.Error() == "task not found" {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{})
}
