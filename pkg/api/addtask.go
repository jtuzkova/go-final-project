package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-final-project/pkg/db"
	"net/http"
	"time"
)

func checkDate(task *db.Task) error {
	now := time.Now()
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	
	if task.Date == "" {
		task.Date = now.Format(DateFormat)
	}

	t, err := time.Parse(DateFormat, task.Date)
	if err != nil {
		return fmt.Errorf("invalid data format: %s", task.Date)
	}

	var next string
	if task.Repeat != "" {
		next, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return fmt.Errorf("invalid rule: %w", err)
		}
	}

	if AfterNow(now, t) || t.Equal(now){
        if len(task.Repeat) == 0 {
            task.Date = now.Format(DateFormat)
        } else {
            task.Date = next
        }
    } 
	return nil
}

func writeError(w http.ResponseWriter, err error, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	resp := map[string]string{"error": err.Error()}
	jsonResponse, _ := json.Marshal(resp)

	w.WriteHeader(statusCode)
	w.Write(jsonResponse)
}

func writeJson(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	resp, err := json.Marshal(data)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func AddTaskHandler(w http.ResponseWriter, req *http.Request) {
	var task db.Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(req.Body)
    if err != nil {
        writeError(w, fmt.Errorf("failed to read request body"), http.StatusBadRequest)
        return
    }
    if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
        writeError(w, err, http.StatusBadRequest)
		return
    }

	if task.Title == "" {
		myErr := "title is empty"
		writeError(w, fmt.Errorf("%s", myErr), http.StatusBadRequest)
		return
	}

	if err := checkDate(&task); err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		myErr := "database error"
		writeError(w, fmt.Errorf("%s", myErr), http.StatusInternalServerError)
		return
	}

	writeJson(w, map[string]int64{
		"id": id,
	})
}