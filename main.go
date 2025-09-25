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

	// Добавляем вызов инициализации базы данных
	if err := db.Init("scheduler.db"); err != nil {
		log.Fatal("Ошибка при инициализации базы данных:", err)
	}

	api.Init() // Инициализируем API

	log.Println("Слушаю на порту:", *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
