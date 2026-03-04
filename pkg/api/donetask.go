package api

import (
	"fmt"
	"go-final-project/pkg/db"
	"net/http"
	"time"
)

func DoneTaskHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		writeError(w, fmt.Errorf("method not allowed"), http.StatusMethodNotAllowed)
		return
	}

	id := req.URL.Query().Get("id")
	if id == "" {
        writeError(w, fmt.Errorf("id is none"), http.StatusBadRequest)
        return
    }

	task, err := db.GetTask(id)
    if err != nil {
        writeError(w, fmt.Errorf("task not found"), http.StatusNotFound)
        return
    }

	if task.Repeat == "" {
        err = db.DeleteTask(id)
        if err != nil {
            writeError(w, fmt.Errorf("failed to delete task"), http.StatusInternalServerError)
            return
        }
    } else {
		nextDate, err := NextDate(time.Now(), task.Date, task.Repeat)
        if err != nil {
            writeError(w, fmt.Errorf("failed to calculate next date: %v", err), http.StatusBadRequest)
            return
        }

		err = db.UpdateDate(nextDate, id)
		if err != nil {
			writeError(w, fmt.Errorf("failed to update next date: %v", err), http.StatusInternalServerError)
		}
	}

	writeJson(w, map[string]int64{})
}
