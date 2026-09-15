package main

import (
	"fmt"
	"net"
	"bufio"
)

func handleClient(conn net.Conn) {
	defer conn.Close()

	fmt.Println("Client connected:", conn.RemoteAddr())

	reader := bufio.NewReader(conn)

	for {
		message, err := reader.ReadString('\n')	
		
		if err != nil {
			fmt.Println("Client disconnected:", conn.RemoteAddr())
			return
		}

		fmt.Println("Received:", message)

		_, err = conn.Write([]byte("Server: " + message))

		if err != nil{
			return
		}
	}

}

func main() {
	listener, err := net.Listen("tcp", ":8080")

	if err != nil{
		fmt.Println("Server error: ", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server started on port 8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Accept error: ", err)
			continue
		}
		go handleClient(conn)
	}

}



// import (
// 	"fmt"
// 	"net"
// )

// func handleClient(conn net.Conn) {
// 	defer conn.Close()

// 	fmt.Println("Working with:", conn.RemoteAddr())
// }

// func main() {
// 	listener, err := net.Listen("tcp", ":8080")

// 	if err != nil {
// 		fmt.Println("server error: ", err)
// 	}

// 	defer listener.Close()

// 	fmt.Println("TSP server started on port 8080")
	
// 	for {
// 		conn, err := listener.Accept()

// 		if err != nil {
// 			fmt.Println("Accept error", err)
// 			continue
// 		}
// 		fmt.Println("New client: ", conn.RemoteAddr())
// 		conn.Close()
// 	}
// }
	
// //Последовательное выполнение

