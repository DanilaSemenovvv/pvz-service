package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/DanilaSemenovvv/pvz/service"
	"github.com/DanilaSemenovvv/pvz/storage"
)

func main() {
	db, err := storage.NewJSONStorage("orders.json")
	if err != nil {
		log.Fatalf("Ошибка запуска БД: %v", err)
	}

	srv := service.NewOrderService(db)

	fmt.Println("Система ПВЗ успешно запущена! Ожидание команд...")

	var command string

	for {
		fmt.Print("\nВведите команду > ")
		fmt.Scan(&command)

		switch command {
		case "exit":
			fmt.Println("Завершение работы...")
			return
		case "accept":
			var orderID, clientID int
			var dateStr string

			fmt.Print("Введите ID заказа: ")
			fmt.Scan(&orderID)

			fmt.Print("Введите ID клиента: ")
			fmt.Scan(&clientID)

			fmt.Print("Введите срок хранения (в формате YYYY-MM-DD): ")
			fmt.Scan(&dateStr)

			deadline, err := time.Parse(time.DateOnly, dateStr)
			if err != nil {
				fmt.Println("Ошибка: неверный формат даты. Попробуйте еще раз.")
				continue // Прерываем текущий кейс и возвращаемся в начало цикла
			}

			err = srv.AcceptOrder(orderID, clientID, deadline)
			if err != nil {
				// Если бизнес-логика вернула ошибку, вежливо показываем ее оператору
				fmt.Printf("Ошибка при принятии заказа: %v\n", err)
			} else {
				fmt.Println("Заказ успешно принят от курьера!")
			}
		case "deliver":
			var clientID int
			var orderIDs []int
			var ordersInput string

			fmt.Print("Введите ID клиента: ")
			fmt.Scan(&clientID)

			fmt.Print("Введите ID заказов через запятую (например: 1,2,3): ")
			fmt.Scan(&ordersInput)

			strIDs := strings.Split(ordersInput, ",")

			isValid := true

			for _, str := range strIDs {
				id, err := strconv.Atoi(str)
				if err != nil {
					fmt.Printf("Ошибка: '%s' не является корректным числом. Повторите ввод команды.\n", str)
					isValid = false
					break
				}
				orderIDs = append(orderIDs, id)
			}

			if !isValid {
				continue
			}

			err := srv.ProcessClient(clientID, orderIDs, service.ActionType(service.ActionDeliver))
			if err != nil {
				fmt.Printf("Ошибка выдачи заказа: %v\n", err)
			}
		case "return":
			var clientID int
			var orderIDs []int
			var ordersInput string

			fmt.Print("Введите ID клиента: ")
			fmt.Scan(&clientID)

			fmt.Print("Введите ID заказов через запятую (например: 1,2,3): ")
			fmt.Scan(&ordersInput)

			strIDs := strings.Split(ordersInput, ",")

			isValid := true

			for _, str := range strIDs {
				id, err := strconv.Atoi(str)
				if err != nil {
					fmt.Printf("Ошибка: '%s' не является корректным числом. Повторите ввод команды.\n", str)
					isValid = false
					break
				}
				orderIDs = append(orderIDs, id)
			}

			if !isValid {
				continue
			}

			err := srv.ProcessClient(clientID, orderIDs, service.ActionType(service.ActionReturn))
			if err != nil {
				fmt.Printf("Ошибка возврата заказа: %v\n", err)
			}
		case "return_to_courier":
			var orderID int

			fmt.Print("Введите ID заказа: ")
			fmt.Scan(&orderID)

			err := srv.ReturnToCourier(orderID)
			if err != nil {
				fmt.Printf("Ошибка возврата заказа курьеру: %v\n", err)
			}
		case "history":
			var limit, page int
			fmt.Print("Размер странницы: ")
			fmt.Scan(&limit)

			fmt.Print("Введите номер страницы: ")
			fmt.Scan(&page)

			if page < 1 {
				page = 1
			}

			offset := (page - 1) * limit

			data, err := srv.GetOrderHistory(limit, offset)
			if err != nil {
				fmt.Printf("Ошибка получения заказов: %v\n", err)
			}

			if len(data) == 0 {
				fmt.Println("На этой странице нет заказов.")
				continue
			}

			fmt.Println("\n=== ИСТОРИЯ ЗАКАЗОВ ===")
			for _, order := range data {
				fmt.Printf("Заказ #%d | Клиент: %d | Статус: %s | Обновлен: %s\n",
					order.OrderID, order.ClientID, order.Status, order.UpdatedAt.Format("2006-01-02 15:04"))
			}
			fmt.Println("=======================")
		default:
			fmt.Println("Неизвестная команда. Доступные команды: accept, return_to_courier, deliver, return, history, exit")
		}
	}
}
