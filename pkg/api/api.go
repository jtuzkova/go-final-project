package api

import (
	"fmt"
	"net/http"
)

const DateFormat = "20060102"

func taskHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        GetTaskHandler(w, r)
    case http.MethodPost:
        AddTaskHandler(w, r)
    case http.MethodPut:
        UpdateTaskHandler(w, r)
    case http.MethodDelete:
        DeleteTaskHandler(w, r)
    default:
        writeError(w, fmt.Errorf("failed "), http.StatusMethodNotAllowed)
    }
}

func Init() {
    http.HandleFunc("/api/nextdate", NextDayHandler)
	http.HandleFunc("/api/task", Auth(taskHandler))
    http.HandleFunc("/api/tasks", Auth(tasksHandler))
    http.HandleFunc("/api/task/done", Auth(DoneTaskHandler))
    http.HandleFunc("/api/signin", SignInHandler)
}