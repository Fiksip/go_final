package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"go_final/pkg/api"
	"go_final/pkg/db"
)

func main() {
	// Инициализация БД
	dbFile := "scheduler.db"
	if envFile := os.Getenv("TODO_DBFILE"); envFile != "" {
		dbFile = envFile
	}

	if err := db.Init(dbFile); err != nil {
		log.Fatalf("Ошибка инициализации БД: %v", err)
	}
	defer db.DB.Close()

	// Инициализация API-обработчиков
	api.Init()

	// Определяем порт
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	// Файловый сервер для фронтенда
	fs := http.FileServer(http.Dir("./web"))
	http.Handle("/", fs)

	fmt.Printf("Сервер запущен на http://localhost:%s\n", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Ошибка сервера: %v", err)
	}
}
