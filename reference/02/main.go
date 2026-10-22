package main

import "fmt"

func main() {
	firstTitle := "Записаться к врачу"
	secondTitle := "Купить чай"
	thirdTitle := "Разобрать почту"
	firstMinutes := 25
	secondMinutes := 15
	thirdMinutes := 35
	const breakMinutes = 5
	availableMinutes := 90

	workMinutes := firstMinutes + secondMinutes + thirdMinutes
	breaksMinutes := 2 * breakMinutes
	totalMinutes := workMinutes + breaksMinutes
	fitsEvening := totalMinutes <= availableMinutes

	fmt.Printf("%s: %d мин\n", firstTitle, firstMinutes)
	fmt.Printf("%s: %d мин\n", secondTitle, secondMinutes)
	fmt.Printf("%s: %d мин\n", thirdTitle, thirdMinutes)
	fmt.Println("Работа, мин:", workMinutes)
	fmt.Println("Перерывы, мин:", breaksMinutes)
	fmt.Printf("Всего: %d ч %d мин\n", totalMinutes/60, totalMinutes%60)
	fmt.Println("План помещается:", fitsEvening)
}
