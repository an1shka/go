// package main

// import (
// 	"fmt"
// 	"time"
// 	"sync"
// )

// func worker(id int, wg *sync.WaitGroup) {
// 	defer wg.Done()

// 	fmt.Println(
// 		"Worker",
// 		id,
// 		"started")

// 	time.Sleep(time.Second)
// 	fmt.Println(
// 		"worker",
// 		id,
// 		"finished",
// 	)
// }

// func main() {
// 	var wg sync.WaitGroup

// 	for i := 1; i <= 5; i++ {
// 		wg.Add(1)

// 		go worker (i, &wg)
// 	}
// 	wg.Wait()
// 	fmt.Println("All workers finished")
// }

// // func task1 (){
// // 	fmt.Println("Task 1 started")

// // 	time.Sleep(1 * time.Second)

// // 	fmt.Println("Task 1 finished")
// // }

// // func task2 (){
// // 	fmt.Println("Task 2 started")

// // 	time.Sleep(1 * time.Second)

// // 	fmt.Println("Task 2 finished")
// // }

// // func main() {
// // 	task1()
// // 	task2()
// // }

// //Конкуретность -- все работает физически одновременно

// //goroutins

// // func printNumbers() {
// // 	for i := 1; i <= 5; i++{
// // 		fmt.Println("Number:", i)

// // 		time.Sleep(500 * time.Millisecond)
// // 	}
// // }

// // func printLetters() {
// // 	letters := []string{
// // 		"A",
// // 		"B",
// // 		"C",
// // 		"D",
// // 		"E",
// // 	}

// // 	for _, letter := range letters {
// // 		fmt.Println("Letter:", letter)
// // 		time.Sleep(500 * time.Millisecond)
// // 	}
// // }

// // func main() {
// // 	go printNumbers()
// // 	go printLetters()

// // 	time.Sleep(3 * time.Second)
// // }
