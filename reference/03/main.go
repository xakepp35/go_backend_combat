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
	title := firstTitle
	priority := 2
	done := false
	reminders := 3
	minutes := firstMinutes + secondMinutes + thirdMinutes + 2*breakMinutes
	fmt.Printf("%s: %d мин\n", firstTitle, firstMinutes)
	fmt.Printf("%s: %d мин\n", secondTitle, secondMinutes)
	fmt.Printf("%s: %d мин\n", thirdTitle, thirdMinutes)
	fmt.Println("Работа, мин:", minutes-2*breakMinutes)
	fmt.Println("Перерывы, мин:", 2*breakMinutes)
	fmt.Printf("Всего: %d ч %d мин\n", minutes/60, minutes%60)
	if minutes <= availableMinutes {
		fmt.Println("План помещается")
	} else {
		fmt.Println("Перегрузка на", minutes-availableMinutes, "минут")
	}
	if title == "" {
		fmt.Println("Ошибка: название не должно быть пустым")
		return
	}
	if priority < 1 || priority > 3 {
		fmt.Println("Ошибка: приоритет должен быть от 1 до 3")
		return
	}
	if reminders < 0 || reminders > 3 {
		fmt.Println("Ошибка: число напоминаний должно быть от 0 до 3")
		return
	}
	if done {
		fmt.Println("Уже выполнено")
		return
	}
	if priority >= 2 {
		fmt.Println("Срочно")
	} else {
		fmt.Println("Можно запланировать")
	}
	for number := 1; number <= reminders; number++ {
		fmt.Println("Напоминание", number, title)
	}
}
