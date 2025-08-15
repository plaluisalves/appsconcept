package fake

import (
	"github.com/brianvoe/gofakeit/v7"
)

var (
	// Pre-generate up to 20 random first names for str1 and str2
	str1Values = generateRandomNames(2)
	str2Values = generateRandomNames(2)

	// Pre-generate integer ranges for int1, int2, limit
	int1Values  = generateIntRange(1, 2)
	int2Values  = generateIntRange(1, 2)
	limitValues = generateIntRange(1, 2)
)

// Generates up to `count` random first names
func generateRandomNames(count int) []string {
	names := make([]string, count)
	for i := 0; i < count; i++ {
		names[i] = gofakeit.FirstName()
	}
	return names
}

// Generates a slice of integers between min and max (inclusive)
func generateIntRange(min, max int) []int {
	values := make([]int, 0, max-min+1)
	for i := min; i <= max; i++ {
		values = append(values, i)
	}
	return values
}

// Random pickers
func Int1() int {
	return gofakeit.RandomInt(int1Values)
}

func Int2() int {
	return gofakeit.RandomInt(int2Values)
}

func Limit() int {
	return gofakeit.RandomInt(limitValues)
}

func Str1() string {
	return gofakeit.RandomString(str1Values)
}

func Str2() string {
	return gofakeit.RandomString(str2Values)
}
