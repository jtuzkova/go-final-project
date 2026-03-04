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
	fmt.Println("Запускаем сервер на порту")

	dbPath := os.Getenv("DB_PATH")
    if dbPath == "" {
        dbPath = "./scheduler.db" 
    }

	err := db.Init(dbPath)
	if err != nil {
		log.Fatalf("Ошибка инициализации БД: %v", err)
	}

	defer db.Close()

	api.Init()
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	fmt.Printf("Порт: %s\n", port)

	api.ExpectedPass = os.Getenv("TODO_PASSWORD")

	webDir := "./web"
	
	http.Handle("/", http.FileServer(http.Dir(webDir))) 

	return http.ListenAndServe(":" + port, nil)
}