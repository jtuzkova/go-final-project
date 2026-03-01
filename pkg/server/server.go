package server

import (
	"fmt"
	"go-final-project/pkg/api"
	"go-final-project/pkg/db"
	"log"
	"net/http"
	"os"
)

func Run() error {
	fmt.Println("Запускаем сервер")

	dbPath := os.Getenv("DB_PATH")
    if dbPath == "" {
        dbPath = "./scheduler.db" 
    }
    
	err := db.Init("scheduler.db")
	if err != nil {
		log.Fatalf("Ошибка инициализации БД: %v", err)
	}

	api.Init()
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}
	webDir := "./web"
	
	http.Handle("/", http.FileServer(http.Dir(webDir))) 

	return http.ListenAndServe(":" + port, nil)
}