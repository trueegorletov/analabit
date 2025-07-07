package main

import (
	"analabit/core/source/rzgmu"
	"fmt"
	"log"
)

func main() {
	// Test priority parsing specifically for student 3950875
	testData := `Направление подготовки: Педиатрия
№ Код Балл ВИ ИД ПП Приоритет Согласие
ПЕДИАТРИЯ (БЮДЖЕТ) (Специалитет)
Мест: 5
Конкурсная группа: ОСНОВНЫЕ МЕСТА
88. 3950875 275     Химия - 93          5         2               Нет`

	// Test individual function parsing
	app := rzgmu.parseApplicationStart("88. 3950875 275     Химия - 93          5         2               Нет")

	if app != nil {
		fmt.Printf("parseApplicationStart results:\n")
		fmt.Printf("  Student ID: %s\n", app.StudentID)
		fmt.Printf("  Priority: %d\n", app.Priority)
		fmt.Printf("  Score: %d\n", app.ScoresSum)
		fmt.Printf("  Original Submitted: %t\n", app.OriginalSubmitted)
	} else {
		fmt.Println("parseApplicationStart returned nil")
	}

	// Test remaining fields parsing
	remaining := rzgmu.parseRemainingFields("Химия - 93          5         2               Нет")
	if remaining != nil {
		fmt.Printf("\nparseRemainingFields results:\n")
		fmt.Printf("  Priority: %d\n", remaining.Priority)
		fmt.Printf("  Original Submitted: %t\n", remaining.OriginalSubmitted)
		fmt.Printf("  Bonus Points: %d\n", remaining.BonusPoints)
	} else {
		fmt.Println("parseRemainingFields returned nil")
	}
}
