package main

import (
	"fmt"
	"net"
	"strings"
	"time"
)

func main() {
	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Сервер запущен на порту 9000")

	conn, err := listener.Accept()
	if err != nil {
		fmt.Println("Ошибка подключения клиента:", err)
		return
	}
	defer conn.Close()

	buffer := make([]byte, 1024)

	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Ошибка чтения сообщения:", err)
		return
	}

	message := strings.TrimSpace(string(buffer[:n]))

	fmt.Println("Клиент:", conn.RemoteAddr())
	fmt.Println("Получено:", message)

	var response string

	switch message {
	case "/help":
		response = "Commands: /help, /time"
	case "/time":
		response = time.Now().Format("15:04:05")
	default:
		response = "Echo: " + message
	}

	_, err = conn.Write([]byte(response))
	if err != nil {
		fmt.Println("Ошибка отправки ответа:", err)
		return
	}
}
