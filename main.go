package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/kimnick1990/GoFinalProject/pkg/api" // Импортируем API пакет
	"github.com/kimnick1990/GoFinalProject/pkg/db"  // Импортируем пакет db
)

func main() {
	port := flag.String("port", "7540", "Порт для прослушивания")
	flag.Parse()

	if err := db.Init("scheduler.db"); err != nil {
		log.Fatal("Ошибка при инициализации базы данных:", err) // Добавляем вызов инициализации базы данных
	}

	api.Init() // Инициализируем API

	http.Handle("/", http.FileServer(http.Dir("web"))) // Подключаем директорию web

	log.Println("Слушаю на порту:", *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
