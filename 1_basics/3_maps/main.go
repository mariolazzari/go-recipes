package main

import (
	"fmt"
	"strings"
)

func main() {
	msg := "To be or not to be"
	fmt.Println("Frequency:", frequency(msg))
}

func frequency(str string) map[string]int {
	freq := make(map[string]int)

	words := strings.SplitSeq(strings.ToLower(str), " ")
	for word := range words {
		freq[word] += 1
	}

	return freq
}
