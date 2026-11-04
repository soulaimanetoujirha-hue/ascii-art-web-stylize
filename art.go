package main

import "strings"

func startLine(char rune) int {
	return (int(char)-32)*9 + 1
}

func drawing(text string, lines []string) string {
	words := strings.Split(text, "\n")
	if strings.ReplaceAll(text, "\n", "") == "" {
		return strings.Repeat("\n", len(words)-1)
	}

	var art strings.Builder
	for _, word := range words {
		if word == "" {
			art.WriteString("\n")
			continue
		}
		for row := 0; row < 8; row++ {
			for _, char := range word {
				art.WriteString(lines[startLine(char)+row])
			}
			art.WriteString("\n")
		}
	}
	return art.String()
}
