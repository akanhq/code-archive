package hashtable

import (
	"fmt"
	"testing"
)

var roman = map[string]int{
	"I": 1,
	"V": 5,
	"X": 10,
	"L": 50,
	"C": 100,
	"D": 500,
	"M": 1000,
}

func romanToInt(s string) int {
	var toInt int
	runes := []rune(s)

	length := len(runes)

	for i := 0; i < len(runes); i++ {
		val := roman[string(runes[i])]
		if i < length-1 && roman[string(runes[i+1])] > val {
			toInt -= val
		} else {
			toInt += val
		}
	}
	return toInt
}
func TestRomanToInt(t *testing.T) {
	fmt.Println("MCMXCIV = ", romanToInt("MCMXCIV"))
}
