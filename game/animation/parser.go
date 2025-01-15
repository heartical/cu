package animation

import (
	"log"
	"strconv"
	"strings"
)

// parseInterval парсит интервал из различных типов данных (int, float64, string).
// Возвращает минимальное значение, максимальное значение и шаг (1 или -1 в зависимости от порядка).
// В случае ошибки парсинга программа завершается с фатальной ошибкой.
func parseInterval(value interface{}) (min, max, step int) {
	switch v := value.(type) {
	case int:
		return v, v, 1
	case float64:
		return int(v), int(v), 1
	case string:
		return parseIntervalFromString(v)
	default:
		log.Fatalf("unsupported type for interval parsing: %T", value)
	}
	return 0, 0, 0
}

// parseIntervalFromString парсит интервал из строки.
// Ожидаемый формат строки: "min-max" или просто число.
// Возвращает минимальное значение, максимальное значение и шаг.
// В случае ошибки парсинга программа завершается с фатальной ошибкой.
func parseIntervalFromString(s string) (min, max, step int) {
	s = strings.TrimSpace(s)

	// Попытка парсинга как одиночного числа.
	if num, err := strconv.Atoi(s); err == nil {
		return num, num, 1
	}

	// Попытка парсинга как интервала в формате "min-max".
	parts := strings.Split(s, "-")
	if len(parts) != 2 {
		log.Fatalf("invalid interval format: %s. Expected 'min-max' or a single number", s)
	}

	min, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		log.Fatalf("failed to parse min value from interval: %s", s)
	}

	max, err = strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		log.Fatalf("failed to parse max value from interval: %s", s)
	}

	if min > max {
		return min, max, -1
	}

	return min, max, 1
}
