package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:9000")
	if err != nil {
		fmt.Println("Ошибка подключения к серверу:", err)
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Введите сообщение: ")

	message, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Ошибка ввода:", err)
		return
	}

	message = strings.TrimSpace(message)

	_, err = conn.Write([]byte(message))
	if err != nil {
		fmt.Println("Ошибка отправки сообщения:", err)
		return
	}

	buffer := make([]byte, 1024)

	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Ошибка получения ответа:", err)
		return
	}

	fmt.Println("Server:", string(buffer[:n]))
}
