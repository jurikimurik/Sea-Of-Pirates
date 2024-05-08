package util

import (
	"errors"
	"strconv"
	"strings"
)

// Substring function to get a substring from the string
func substring(str string, start, end int) string {
	return strings.TrimSpace(str[start:end])
}

// Function that translates coords (like 'B10') to two separate numbers (like 2 and 10)
//
//	Arguments:
//
// coords - String coordinate (eg. "B10")
//
//	Returns:
//
// int - Translated letter to number as integer
//
// int - Translated string number to integer
//
// error - If error occurs, it will return -1, -1 and occured error
func CoordToIntegers(coords string) (int, int, error) {
	//Checking size of coords
	if len(coords) < 2 {
		return -1, -1, errors.New("coords can't be less than size of 2")
	}

	//Taking out letter from coords and make it small
	letter := substring(coords, 0, 1)
	letter = strings.ToLower(letter)

	// Taking out numbers from coords
	numbers := substring(coords, 1, len(coords))

	//String numbers to integer numbers
	newNumbers, err := strconv.Atoi(numbers)
	if err != nil {
		return -1, -1, err
	}
	return int(rune(letter[0]) - 96), newNumbers, nil
}
