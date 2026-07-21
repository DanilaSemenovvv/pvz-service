package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/DanilaSemenovvv/pvz/internal/mygrpc"
	"github.com/DanilaSemenovvv/pvz/internal/service"
	"github.com/DanilaSemenovvv/pvz/internal/storage"
	desc "github.com/DanilaSemenovvv/pvz/pkg/api/pvz/v1"
	"google.golang.org/grpc"
)

func main() {
	connString := "postgres://postgres:qwerty@localhost:5432/pvz_db"

	db, err := storage.NewPostgresStorage(connString)
	if err != nil {
		log.Fatalf("Ошибка запуска БД: %v", err)
	}

	orderSrv := service.NewOrderService(db)

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Ошибка TCP: %v", err)
	}

	grpcServer := grpc.NewServer()

	pvzGrpcHandler := mygrpc.NewServer(orderSrv)

	desc.RegisterPVZServiceServer(grpcServer, pvzGrpcHandler)
	go func() {
		fmt.Println("gRPC сервер запущен на порту 50051...")
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("Ошибка работы сервера: %v", err)
		}
	}()

	stopCh := make(chan os.Signal, 1)

	signal.Notify(stopCh, os.Interrupt, syscall.SIGTERM)

	<-stopCh

	log.Println("Получен сигнал ОС, начинаем плавную остановку сервера...")

	grpcServer.GracefulStop()

	db.Close()

}
