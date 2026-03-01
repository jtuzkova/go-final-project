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
		task.Date = now.Format("20060102")
	}

	t, err := time.Parse("20060102", task.Date)
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
            task.Date = now.Format("20060102")
        } else {
            task.Date = next
        }
    } 
	return nil
}

func writeError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	resp := map[string]string{"error": err.Error()}
	jsonResponse, _ := json.Marshal(resp)

	w.WriteHeader(http.StatusInternalServerError)
	w.Write(jsonResponse)
}

func writeJson(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	resp, err := json.Marshal(data)
	if err != nil {
		writeError(w, err)
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
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
        writeError(w, err)
		return
    }

	if task.Title == "" {
		myErr := "title is empty"
		writeError(w, fmt.Errorf("%s", myErr))
		return
	}

	if err := checkDate(&task); err != nil {
		writeError(w, err)
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		myErr := "database error"
		writeError(w, fmt.Errorf("%s", myErr))
		return
	}

	writeJson(w, map[string]int64{
		"id": id,
	})
}