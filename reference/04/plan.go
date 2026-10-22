package main

import "fmt"

func describe(title string, priority int) string {
	return fmt.Sprintf("%s (приоритет %d)", title, priority)
}

func validPriority(priority int) bool {
	return priority >= 1 && priority <= 3
}

// isUrgent expects a validated priority in the range 1..3.
func isUrgent(priority int) bool {
	return priority >= 2
}

// totalMinutes expects nonnegative durations for three tasks and two breaks.
func totalMinutes(first int, second int, third int, breakMinutes int) int {
	return first + second + third + 2*breakMinutes
}

// printReminders expects a validated count in the range 0..3.
func printReminders(title string, count int) {
	for number := 1; number <= count; number++ {
		fmt.Println("Напоминание", number, title)
	}
}
