// package main

// import "fmt"

// func hello(){
// 		fmt.Println("Hello")
// }

// func sum(a int, b int) int{
// 	return a + b
// }


// func main(){
// 	fmt.Printf("Hello World")
// 	fmt.Println("Hello World")
	
// 	// var name string = "alex"
// 	// var age int = 20
// 	// fmt.Println(name)
// 	// fmt.Println(age)

// 	name := "Alex"
// 	age := 20
// 	fmt.Println(name)
// 	fmt.Println(age)

// 	fmt.Printf("%T\n", name)

// 	x := 10 
// 	x = 30
// 	fmt.Println(x)

// 	var name2 string

// 	fmt.Print("Введите имя: ")
// 	fmt.Scanln(&name2)

// 	fmt.Println("Привет", name2)

// 	age2 := 18

// 	if age2 >= 18 {
// 		fmt.Println("Доступ разрешен")
// 	}else{
// 		fmt.Println("qwerty")
// 	}

// 	for i := 0; i < 5; i++{
// 		fmt.Println(i)
// 	}

// 	hello()
	
// 	result := sum(10,20)
// 	fmt.Println(result)
	
// }

package main

import "fmt"

func main() {
    var name string
    var age int

    fmt.Println("Введите имя:")
    fmt.Scanln(&name)

    fmt.Println("Введите возраст:")
    fmt.Scanln(&age)

    fmt.Println("Привет,", name, "! Тебе", age, "лет.")
}