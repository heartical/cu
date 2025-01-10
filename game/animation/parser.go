package animation

import (
	"log"
	"strconv"
	"strings"
)

func parseInterval(val interface{}) (int, int, int) {
	switch v := val.(type) {
	case int:
		return v, v, 1
	case float64:
		return int(v), int(v), 1
	case string:
		if n, err := strconv.Atoi(v); err == nil {
			return n, n, 1
		}
		matches := intervalMatcher.FindStringSubmatch(strings.TrimSpace(v))
		if len(matches) != 3 {
			log.Fatalf("could not parse interval from %s", v)
		}
		min, _ := strconv.Atoi(matches[1])
		max, _ := strconv.Atoi(matches[2])
		if min > max {
			return min, max, -1
		} else {
			return min, max, 1
		}
	default:
		log.Fatalf("Could not parse interval from %v", val)
	}
	panic("unreachable")
}
