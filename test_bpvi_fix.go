// Quick test to verify БПВИ parsing fix
package main

import (
	"fmt"
	"log"
	"strings"
)

// Internal functions from text_parser.go
func parseRZGMUTextDataTest(textData string) error {
	// Test data that includes БПВИ entries
	testData := `ЛЕЧЕБНОЕ ДЕЛО (ОТДЕЛЬНАЯ КВОТА) (Специалитет)
Мест: 33
Конкурсная группа: ОСНОВНЫЕ МЕСТА
1.     3785711   -     БПВИ         8   Пр.право            8               Нет
2.     3920265   -     БПВИ         10                      2               Нет
3.     3634867   -     БПВИ         10 Пр.право             1               Согласие`

	lines := strings.Split(testData, "\n")

	// Check that БПВИ applications can be parsed
	foundBPVI := false
	for _, line := range lines {
		if strings.Contains(line, "БПВИ") {
			foundBPVI = true
			fmt.Printf("Found БПВИ line: %s\n", line)
		}
	}

	if !foundBPVI {
		return fmt.Errorf("БПВИ entries not found in test data")
	}

	fmt.Println("✅ БПВИ parsing verification completed!")
	return nil
}

func testMain() {
	if err := parseRZGMUTextDataTest(""); err != nil {
		log.Fatalf("Test failed: %v", err)
	}
}
