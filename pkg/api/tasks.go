package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-final-project/pkg/db"
	"net/http"
	"time"
)

type TasksResp struct {
    Tasks []*db.Task `json:"tasks"`
}

const tasksLimit = 50

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	var tasks []*db.Task
	var err error

	if search == "" {
		tasks, err = db.Tasks(tasksLimit) 
	} else {
		date, err := time.Parse("02.01.2006", search)
		if err != nil {
			pattern := "%" + search + "%"
			tasks, err = db.SearchTaskString(pattern, tasksLimit)
		} else {
			newFormatDate := date.Format(DateFormat)
			tasks, err = db.SearchTaskDate(newFormatDate, tasksLimit)
		}
	}

	if err != nil {
        writeError(w, err,  http.StatusInternalServerError)
        return
    }

    writeJson(w, TasksResp{
        Tasks: tasks,
    })
}

func GetTaskHandler(w http.ResponseWriter, req *http.Request) {
    id := req.URL.Query().Get("id")
    
    if id == "" {
        writeError(w, fmt.Errorf("id is none"), http.StatusBadRequest)
        return
    }

    task, err := db.GetTask(id)
    if err != nil {
        writeError(w, fmt.Errorf("task not found"), http.StatusBadRequest)
        return
    }

    writeJson(w, task)
}

func UpdateTaskHandler(w http.ResponseWriter, req *http.Request) {
	var task db.Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
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

	err = db.UpdateTask(&task)
	if err != nil {
		myErr := "task not found"
		writeError(w, fmt.Errorf("%s", myErr), http.StatusNotFound)
		return
	}

	writeJson(w, map[string]int64{})
}

func DeleteTaskHandler(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")
    
    if id == "" {
        writeError(w, fmt.Errorf("id is none"), http.StatusBadRequest)
        return
    }

	err := db.DeleteTask(id)
    if err != nil {
        writeError(w, fmt.Errorf("failed to delete task"), http.StatusNotFound)
        return
    }

    writeJson(w, map[string]int64{})
}