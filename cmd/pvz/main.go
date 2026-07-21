package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DanilaSemenovvv/pvz/internal/mygrpc"
	"github.com/DanilaSemenovvv/pvz/internal/service"
	"github.com/DanilaSemenovvv/pvz/internal/storage"
	desc "github.com/DanilaSemenovvv/pvz/pkg/api/pvz/v1"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	gwMux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err = desc.RegisterPVZServiceHandlerFromEndpoint(ctx, gwMux, "localhost:50051", opts)
	if err != nil {
		log.Fatalf("ошибка регистрации эндпоинтов: %v", err)
	}

	httpServer := &http.Server{Addr: ":8080", Handler: gwMux}

	go func() {
		fmt.Println("HTTP сервер запущен на порту 8080...")
		err := httpServer.ListenAndServe()
		if err == http.ErrServerClosed {
			fmt.Println("Сервер закрыт")
		} else {
			log.Fatalf("Ошибка работы сервера: %v", err)
		}
	}()

	stopCh := make(chan os.Signal, 1)

	signal.Notify(stopCh, os.Interrupt, syscall.SIGTERM)

	<-stopCh

	log.Println("Получен сигнал ОС, начинаем плавную остановку сервера...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	err = httpServer.Shutdown(shutdownCtx)
	if err != nil {
		log.Fatalf("Ошибка закрытия сервера: %v", err)
	}
	grpcServer.GracefulStop()
	db.Close()

}
