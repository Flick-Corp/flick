/*
** FLICK PROJECT, 2026
** flick/internal/api/code/code.go
** File description:
** code.go
 */

package code

import (
	"bufio"
	_ "embed"
	"fmt"
	"math/rand/v2"
	"strings"
)

//go:embed words.txt
var file string

// PrintCode: Prints a code to the format word-word-number.
//
// Params:
// - code (string): the code to display.
func PrintCode(code string) {
	fmt.Printf("share code: %s\n", code)
}

// CodeGen: Generates a random code to the format word-word-number.
//
// Returns:
// - code: The random code.
func CodeGen() string {
	var word1 string
	var word2 string
	var max int

	number := rand.IntN(1000)
	index1 := rand.IntN(999)
	index2 := rand.IntN(999)
	scanner := bufio.NewScanner(strings.NewReader(file))

	if index1 > index2 {
		max = index1
	} else {
		max = index2
	}

	for i := 0; i <= max; i++ {
		if i == index1 {
			if scanner.Scan() {
				word1 = scanner.Text()
				continue
			}
		}
		if i == index2 {
			if scanner.Scan() {
				word2 = scanner.Text()
				continue
			}
		}
		scanner.Scan()
	}

	code := fmt.Sprintf("%s-%s-%d", word1, word2, number)
	return code
}
