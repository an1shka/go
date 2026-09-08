package main

import (
	"fmt"
)

func lessonMain() {
	// names := []string{"Amir", "Danya", "Tagir", "Misha"}

	// fmt.Println("Names:", names[0])

	// for i := 0; i < len(names); i++ {
	// 	fmt.Println("Names:", names[i])
	// }

	// for i := 0; i < len(names); i++ {
	// 	fmt.Println(i,":", names[i])
	// }

	// names := []string{
	// 	{name := "Amir", age := 20, city := "Moscow"},
	// 	{name := "Danya", age := 21, city := "Moscow"},
	// }

	// for index, value := range names {
	// 	fmt.Println(index, value)
	// }

	// users := []string{"Amirka", "Danya", "Tagir", "Misha"}
	// fmt.Println(users[0])
	// users[0] = "Amir"
	// fmt.Println(users)

	// numbers := []int{1, 2, 3, 4, 5}
	// part := numbers[1:5]
	// fmt.Println("Numbers:", part)

	ages := map[string]int{
		"Amir":  17,
		"Anya":  16,
		"Manya": 18,
	}
	// delete(ages, "Anya")
	// fmt.Println("Ages:", ages)

	age, exists := ages["Amir"]
	if exists {
		fmt.Println("Age:", age)
	} else {
		fmt.Println("User not found")
	}
}

// func printNames() {
// 	for index, value := range names {
//  		fmt.Println(index, value)
// 	}
// }

// var index int

// fmt.Print("Введите индекс: ")
// fmt.Scan(&index)

// fmt.Println(names[index])
