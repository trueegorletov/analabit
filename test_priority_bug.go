package main

import (
	"analabit/core/source/rzgmu"
	"fmt"
	"log"
)

func main() {
	// Test with actual format from rzgmu_p_b-text-converted.txt
	actualFormatData := `Направление подготовки: Педиатрия
№ Код Балл ВИ ИД ПП Приоритет Согласие
ПЕДИАТРИЯ (БЮДЖЕТ) (Специалитет)
Мест: 5
Конкурсная группа: ОСНОВНЫЕ МЕСТА
88. 3950875 275     Химия - 93          5         2               Нет`

	fmt.Println("Testing with actual single-line format...")
	programs, err := rzgmu.parseRZGMUTextData(actualFormatData)
	if err != nil {
		log.Fatalf("Failed to parse: %v", err)
	}

	if len(programs) > 0 && len(programs[0].Applications) > 0 {
		app := programs[0].Applications[0]
		fmt.Printf("Student ID: %s\n", app.StudentID)
		fmt.Printf("Priority: %d (expected: 2)\n", app.Priority)
		fmt.Printf("Score: %d\n", app.ScoresSum)
		fmt.Printf("Original Submitted: %t (expected: false)\n", app.OriginalSubmitted)

		if app.Priority != 2 {
			fmt.Printf("BUG CONFIRMED: Got priority %d, expected 2\n", app.Priority)
		} else {
			fmt.Println("Priority parsing works correctly!")
		}
	}

	// Test with the format from rzgmu_l_b-text-converted.txt
	fmt.Println("\nTesting Лечебное дело format...")
	actualFormatData2 := `Направление подготовки: Лечебное дело
№ Код Балл ВИ ИД ПП Приоритет Согласие
ЛЕЧЕБНОЕ ДЕЛО (БЮДЖЕТ) (Специалитет)
Мест: 19
Конкурсная группа: ОСНОВНЫЕ МЕСТА
93. 3950875 278     Химия - 93          8         1               Нет`

	programs2, err := rzgmu.parseRZGMUTextData(actualFormatData2)
	if err != nil {
		log.Fatalf("Failed to parse: %v", err)
	}

	if len(programs2) > 0 && len(programs2[0].Applications) > 0 {
		app := programs2[0].Applications[0]
		fmt.Printf("Student ID: %s\n", app.StudentID)
		fmt.Printf("Priority: %d (expected: 1)\n", app.Priority)
		fmt.Printf("Score: %d\n", app.ScoresSum)
		fmt.Printf("Original Submitted: %t (expected: false)\n", app.OriginalSubmitted)

		if app.Priority != 1 {
			fmt.Printf("BUG CONFIRMED: Got priority %d, expected 1\n", app.Priority)
		} else {
			fmt.Println("Priority parsing works correctly!")
		}
	}
}
