package main

import (
	"flag"
	"log"
	"net/http"

	"pkg/db" // Импортируйте ваш пакет db
)

func main() {
	port := flag.String("port", "7540", "Порт для прослушивания")
	flag.Parse()

	// Добавляем вызов инициализации базы данных
	if err := db.Init("scheduler.db"); err != nil {
		log.Fatal("Ошибка при инициализации базы данных:", err)
	}

	webDir := "./web"
	log.Println("Слушаю на порту:", *port)
	log.Fatal(http.ListenAndServe(":"+*port, http.FileServer(http.Dir(webDir))))
}
