package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DanilaSemenovvv/pvz/internal/handler"
	"github.com/DanilaSemenovvv/pvz/internal/service"
	"github.com/DanilaSemenovvv/pvz/internal/storage"
)

func HandlePing(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong! Система ПВЗ работает!"))
}

func main() {
	connString := "postgres://postgres:qwerty@localhost:5432/pvz_db"

	db, err := storage.NewPostgresStorage(connString)
	if err != nil {
		log.Fatalf("Ошибка запуска БД: %v", err)
	}

	orderSrv := service.NewOrderService(db)

	h := handler.NewHandler(orderSrv)

	http.HandleFunc("GET /history", h.GetHistory)
	http.HandleFunc("POST /acceptOrder", h.AcceptOrder)
	http.HandleFunc("POST /returnOrder", h.ReturnOrder)
	http.HandleFunc("POST /deliverOrder", h.DeliverOrder)
	http.HandleFunc("GET /ping", HandlePing)

	fmt.Println("Веб-сервер запущен на порту 8080...")

	server := &http.Server{
		Addr: ":8080",
	}

	go func() {
		err := server.ListenAndServe()
		if err != http.ErrServerClosed {
			log.Fatalf("Ошибка работы сервера: %v", err)
		}
	}()

	stopCh := make(chan os.Signal, 1)

	signal.Notify(stopCh, os.Interrupt, syscall.SIGTERM)

	<-stopCh

	fmt.Println("Начинается плавная остановка сервера...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	if err = server.Shutdown(ctx); err != nil {
		log.Printf("Сервер принудительно остановлен из-за ошибки: %v", err)
	} else {
		fmt.Println("Сервер успешно и плавно остановлен.")
	}

}
