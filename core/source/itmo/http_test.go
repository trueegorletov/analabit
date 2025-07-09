package itmo

import (
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/trueegorletov/analabit/core"
	"golang.org/x/net/html"
)

// TestCapacityExtractionOnRealURLs tests that capacity extraction works on actual ITMO URLs
func TestCapacityExtractionOnRealURLs(t *testing.T) {
	// Test with a few representative URLs
	testCases := []struct {
		url           string
		expectedTotal int
		name          string
		shouldExtract bool // whether we expect real extraction vs fallback
	}{
		{
			url:           "https://abit.itmo.ru/rating/bachelor/budget/2190",
			expectedTotal: 170,
			name:          "Прикладная математика и информатика",
			shouldExtract: true,
		},
		{
			url:           "https://abit.itmo.ru/rating/bachelor/budget/2191",
			expectedTotal: 30,
			name:          "Математическое обеспечение и администрирование информационных систем",
			shouldExtract: true,
		},
		{
			url:           "https://abit.itmo.ru/rating/bachelor/budget/2192",
			expectedTotal: 25,
			name:          "Физика",
			shouldExtract: true,
		},
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Fetch the actual page
			resp, err := client.Get(tc.url)
			if err != nil {
				t.Skipf("Failed to fetch URL %s: %v", tc.url, err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				t.Skipf("Got non-200 status %d for URL %s", resp.StatusCode, tc.url)
				return
			}

			// Read response body
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Failed to read response body: %v", err)
			}

			// Parse HTML
			doc, err := html.Parse(strings.NewReader(string(body)))
			if err != nil {
				t.Fatalf("Failed to parse HTML: %v", err)
			}

			// Test capacity extraction
			extractedCapacities := extractCapacities(doc)

			// Verify total matches expected
			totalExtracted := extractedCapacities.Regular + extractedCapacities.TargetQuota +
				extractedCapacities.SpecialQuota + extractedCapacities.DedicatedQuota

			if totalExtracted != tc.expectedTotal {
				t.Errorf("Expected total capacity %d, got %d", tc.expectedTotal, totalExtracted)
			}

			// If we expect real extraction (not fallback), verify the quotas are properly distributed
			if tc.shouldExtract {
				if extractedCapacities.TargetQuota == 0 && extractedCapacities.SpecialQuota == 0 && extractedCapacities.DedicatedQuota == 0 {
					t.Errorf("Expected real capacity extraction but got fallback distribution for URL %s", tc.url)
					t.Logf("Body snippet: %s", body[:min(500, len(body))])
				} else {
					t.Logf("Successfully extracted capacities: %+v", extractedCapacities)
				}
			}

		})
	}
}

// Helper function for min
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// TestExtractCapacitiesWithRealData tests the extraction function with real fetched data
func TestExtractCapacitiesWithRealData(t *testing.T) {
	// Real data examples from actual ITMO pages
	testCases := []struct {
		name     string
		htmlText string
		expected core.Capacities
	}{
		{
			name:     "Applied Math page",
			htmlText: "Количество мест: 170 (17 ЦК, 17 ОcК, 17 ОтК)",
			expected: core.Capacities{
				Regular:        119, // 170 - 17 - 17 - 17
				TargetQuota:    17,
				SpecialQuota:   17,
				DedicatedQuota: 17,
			},
		},
		{
			name:     "Math Support page",
			htmlText: "Количество мест: 30 (5 ЦК, 3 ОcК, 3 ОтК)",
			expected: core.Capacities{
				Regular:        19, // 30 - 5 - 3 - 3
				TargetQuota:    5,
				SpecialQuota:   3,
				DedicatedQuota: 3,
			},
		},
		{
			name:     "No capacity info",
			htmlText: "Some random text without capacity information",
			expected: core.Capacities{
				Regular:        0, // extraction failed, returns zeros
				TargetQuota:    0,
				SpecialQuota:   0,
				DedicatedQuota: 0,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Parse HTML
			doc, err := html.Parse(strings.NewReader(tc.htmlText))
			if err != nil {
				t.Fatalf("Failed to parse HTML: %v", err)
			}

			result := extractCapacities(doc)

			if result.Regular != tc.expected.Regular {
				t.Errorf("Regular: expected %d, got %d", tc.expected.Regular, result.Regular)
			}
			if result.TargetQuota != tc.expected.TargetQuota {
				t.Errorf("TargetQuota: expected %d, got %d", tc.expected.TargetQuota, result.TargetQuota)
			}
			if result.SpecialQuota != tc.expected.SpecialQuota {
				t.Errorf("SpecialQuota: expected %d, got %d", tc.expected.SpecialQuota, result.SpecialQuota)
			}
			if result.DedicatedQuota != tc.expected.DedicatedQuota {
				t.Errorf("DedicatedQuota: expected %d, got %d", tc.expected.DedicatedQuota, result.DedicatedQuota)
			}

			total := result.Regular + result.TargetQuota + result.SpecialQuota + result.DedicatedQuota
			expectedTotal := tc.expected.Regular + tc.expected.TargetQuota + tc.expected.SpecialQuota + tc.expected.DedicatedQuota

			if total != expectedTotal {
				t.Errorf("Total mismatch: expected %d, got %d", expectedTotal, total)
			}
		})
	}
}

// TestFallbackCapacityCalculation tests the fallback capacity calculation logic
func TestFallbackCapacityCalculation(t *testing.T) {
	testCases := []struct {
		name     string
		totalKCP int
		expected core.Capacities
	}{
		{
			name:     "Total 100",
			totalKCP: 100,
			expected: core.Capacities{
				Regular:        70, // 100 - 10 - 10 - 10
				TargetQuota:    10,
				SpecialQuota:   10,
				DedicatedQuota: 10,
			},
		},
		{
			name:     "Total 170",
			totalKCP: 170,
			expected: core.Capacities{
				Regular:        119, // 170 - 17 - 17 - 17
				TargetQuota:    17,
				SpecialQuota:   17,
				DedicatedQuota: 17,
			},
		},
		{
			name:     "Total 25",
			totalKCP: 25,
			expected: core.Capacities{
				Regular:        16, // 25 - 3 - 3 - 3
				TargetQuota:    3,
				SpecialQuota:   3,
				DedicatedQuota: 3,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := CalculateFallbackCapacities(tc.totalKCP)

			if result.Regular != tc.expected.Regular {
				t.Errorf("Regular: expected %d, got %d", tc.expected.Regular, result.Regular)
			}
			if result.TargetQuota != tc.expected.TargetQuota {
				t.Errorf("TargetQuota: expected %d, got %d", tc.expected.TargetQuota, result.TargetQuota)
			}
			if result.SpecialQuota != tc.expected.SpecialQuota {
				t.Errorf("SpecialQuota: expected %d, got %d", tc.expected.SpecialQuota, result.SpecialQuota)
			}
			if result.DedicatedQuota != tc.expected.DedicatedQuota {
				t.Errorf("DedicatedQuota: expected %d, got %d", tc.expected.DedicatedQuota, result.DedicatedQuota)
			}

			total := result.Regular + result.TargetQuota + result.SpecialQuota + result.DedicatedQuota
			if total != tc.totalKCP {
				t.Errorf("Total mismatch: expected %d, got %d", tc.totalKCP, total)
			}
		})
	}
}

// TestAllRegistryURLsCapacityExtraction tests ALL URLs from the ITMO registry
func TestAllRegistryURLsCapacityExtraction(t *testing.T) {
	// ALL URLs from the ITMO registry with their expected total КЦП
	allURLs := []struct {
		url           string
		expectedTotal int
		name          string
	}{
		{
			url:           "https://abit.itmo.ru/rating/bachelor/budget/2190",
			expectedTotal: 170,
			name:          "Прикладная математика и информатика",
		},
		{
			url:           "https://abit.itmo.ru/rating/bachelor/budget/2191",
			expectedTotal: 30,
			name:          "Математическое обеспечение и администрирование информационных систем",
		},
		{
			url:           "https://abit.itmo.ru/rating/bachelor/budget/2192",
			expectedTotal: 25,
			name:          "Физика",
		},
		{
			url:           "https://abit.itmo.ru/rating/bachelor/budget/2193",
			expectedTotal: 25,
			name:          "Химия",
		},
		{
			url:           "https://abit.itmo.ru/rating/bachelor/budget/2194",
			expectedTotal: 25,
			name:          "Экология и природопользование",
		},
		{
			url:           "https://abit.itmo.ru/rating/bachelor/budget/2195",
			expectedTotal: 30,
			name:          "Информатика и вычислительная техника",
		},
		{
			url:           "https://abit.itmo.ru/rating/bachelor/budget/2196",
			expectedTotal: 151,
			name:          "Информационные системы и технологии",
		},
		{
			url:           "https://abit.itmo.ru/rating/bachelor/budget/2197",
			expectedTotal: 28,
			name:          "Прикладная информатика",
		},
		{
			url:           "https://abit.itmo.ru/rating/bachelor/budget/2198",
			expectedTotal: 164,
			name:          "Программная инженерия",
		},
		{
			url:           "https://abit.itmo.ru/rating/bachelor/budget/2199",
			expectedTotal: 93,
			name:          "Информационная безопасность",
		},
		{
			url:           "https://abit.itmo.ru/rating/bachelor/budget/2200",
			expectedTotal: 81,
			name:          "Инфокоммуникационные технологии и системы связи",
		},
		{
			url:           "https://abit.itmo.ru/rating/bachelor/budget/2201",
			expectedTotal: 14,
			name:          "Конструирование и технология электронных средств",
		},
		{
			url:           "https://abit.itmo.ru/rating/bachelor/budget/2202",
			expectedTotal: 16,
			name:          "Приборостроение",
		},
		{
			url:           "https://abit.itmo.ru/rating/bachelor/budget/2203",
			expectedTotal: 86,
			name:          "Фотоника и оптоинформатика и Лазерная техника и лазерные технологии",
		},
		{
			url:           "https://abit.itmo.ru/rating/bachelor/budget/2204",
			expectedTotal: 28,
			name:          "Биотехнические системы и технологии",
		},
		{
			url:           "https://abit.itmo.ru/rating/bachelor/budget/2205",
			expectedTotal: 16,
			name:          "Электроэнергетика и электротехника",
		},
		{
			url:           "https://abit.itmo.ru/rating/bachelor/budget/2206",
			expectedTotal: 80,
			name:          "Автоматизация технологических процессов и производств и Мехатроника и робототехника",
		},
		{
			url:           "https://abit.itmo.ru/rating/bachelor/budget/2207",
			expectedTotal: 80,
			name:          "Техническая физика",
		},
		{
			url:           "https://abit.itmo.ru/rating/bachelor/budget/2208",
			expectedTotal: 45,
			name:          "Химическая технология",
		},
		{
			url:           "https://abit.itmo.ru/rating/bachelor/budget/2209",
			expectedTotal: 20,
			name:          "Энерго- и ресурсосберегающие процессы в химической технологии, нефтехимии и биотехнологии",
		},
		{
			url:           "https://abit.itmo.ru/rating/bachelor/budget/2210",
			expectedTotal: 65,
			name:          "Биотехнология",
		},
		{
			url:           "https://abit.itmo.ru/rating/bachelor/budget/2211",
			expectedTotal: 17,
			name:          "Системы управления движением и навигация",
		},
		{
			url:           "https://abit.itmo.ru/rating/bachelor/budget/2212",
			expectedTotal: 15,
			name:          "Управление в технических системах",
		},
		{
			url:           "https://abit.itmo.ru/rating/bachelor/budget/2213",
			expectedTotal: 65,
			name:          "Инноватика",
		},
		{
			url:           "https://abit.itmo.ru/rating/bachelor/budget/2298",
			expectedTotal: 59,
			name:          "Бизнес-информатика",
		},
		{
			url:           "https://abit.itmo.ru/rating/bachelor/budget/2214",
			expectedTotal: 42,
			name:          "Интеллектуальные системы в гуманитарной сфере",
		},
		{
			url:           "https://abit.itmo.ru/rating/bachelor/budget/2215",
			expectedTotal: 10,
			name:          "Дизайн",
		},
	}

	client := &http.Client{
		Timeout: 20 * time.Second,
	}

	for _, tc := range allURLs {
		t.Run(tc.name, func(t *testing.T) {
			// Fetch the actual page
			resp, err := client.Get(tc.url)
			if err != nil {
				t.Logf("%s => %s => ERROR: %v", tc.name, tc.url, err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				t.Logf("%s => %s => ERROR: status %d", tc.name, tc.url, resp.StatusCode)
				return
			}

			// Read response body with size limit
			body, err := io.ReadAll(io.LimitReader(resp.Body, 100000))
			if err != nil {
				t.Logf("%s => %s => ERROR: read body: %v", tc.name, tc.url, err)
				return
			}

			// Parse HTML
			doc, err := html.Parse(strings.NewReader(string(body)))
			if err != nil {
				t.Logf("%s => %s => ERROR: parse HTML: %v", tc.name, tc.url, err)
				return
			}

			// Extract capacity
			extractedCapacities := extractCapacities(doc)

			// Log in requested format
			t.Logf("%s => %s => %d", tc.name, tc.url, extractedCapacities)
		})
	}
}

// Helper functions for min/max
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// TestCapacityParsingLogic tests the capacity parsing logic handles edge cases correctly
func TestCapacityParsingLogic(t *testing.T) {
	testCases := []struct {
		name              string
		htmlContent       string
		shouldUseFallback bool
		expectedTotal     int
	}{
		{
			name:              "Normal capacity extraction",
			htmlContent:       "<html><body>Количество мест: 170 (17 ЦК, 17 ОcК, 17 ОтК)</body></html>",
			shouldUseFallback: false,
			expectedTotal:     170,
		},
		{
			name:              "Zero target quota (legitimate)",
			htmlContent:       "<html><body>Количество мест: 10 (0 ЦК, 1 ОcК, 1 ОтК)</body></html>",
			shouldUseFallback: false,
			expectedTotal:     10,
		},
		{
			name:              "All quotas zero (legitimate)",
			htmlContent:       "<html><body>Количество мест: 20 (0 ЦК, 0 ОcК, 0 ОтК)</body></html>",
			shouldUseFallback: false,
			expectedTotal:     20,
		},
		{
			name:              "No capacity information (should use fallback)",
			htmlContent:       "<html><body>Some random content without capacity info</body></html>",
			shouldUseFallback: true,
			expectedTotal:     100, // This will be set as fallback total
		},
		{
			name:              "Malformed capacity string (should use fallback)",
			htmlContent:       "<html><body>Количество мест: broken format</body></html>",
			shouldUseFallback: true,
			expectedTotal:     50, // This will be set as fallback total
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Parse HTML
			doc, err := html.Parse(strings.NewReader(tc.htmlContent))
			if err != nil {
				t.Fatalf("Failed to parse HTML: %v", err)
			}

			// Extract capacities using the same logic as LoadTo
			parsedCapacities := extractCapacities(doc)
			totalParsed := parsedCapacities.Regular + parsedCapacities.TargetQuota +
				parsedCapacities.SpecialQuota + parsedCapacities.DedicatedQuota

			var finalCapacities core.Capacities
			usedFallback := false

			if totalParsed == 0 {
				// Should use fallback
				usedFallback = true
				finalCapacities = CalculateFallbackCapacities(tc.expectedTotal)
			} else {
				// Should use parsed capacities
				finalCapacities = parsedCapacities
			}

			// Verify fallback usage matches expectation
			if usedFallback != tc.shouldUseFallback {
				t.Errorf("Expected fallback usage: %v, got: %v", tc.shouldUseFallback, usedFallback)
			}

			// Verify total capacity
			finalTotal := finalCapacities.Regular + finalCapacities.TargetQuota +
				finalCapacities.SpecialQuota + finalCapacities.DedicatedQuota

			if finalTotal != tc.expectedTotal {
				t.Errorf("Expected total %d, got %d", tc.expectedTotal, finalTotal)
			}

			t.Logf("Capacities: %+v (fallback used: %v)", finalCapacities, usedFallback)
		})
	}
}
