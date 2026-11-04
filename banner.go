package main

import (
	"errors"
	"os"
	"strings"
)

// The fonts users can pick. Each one lives in <name>.txt next to main.go.
var bannerNames = []string{"standard", "shadow", "thinkertoy"}

// Banner is one font card on the page: its name and a small sample.
type Banner struct {
	Name    string
	Preview string
}

// loadBanner reads a font file and splits it into lines.
// Each of the 95 printable characters takes 9 lines: a blank one, then 8 rows.
func loadBanner(name string) ([]string, error) {
	content, err := os.ReadFile(name + ".txt")
	if err != nil {
		return nil, err
	}
	lines := strings.Split(strings.ReplaceAll(string(content), "\r\n", "\n"), "\n")
	if len(lines) < 95*9 {
		return nil, errors.New(name + ".txt is incomplete")
	}
	return lines, nil
}

// previews draws "Aa" in every font for the font picker cards.
func previews() []Banner {
	var list []Banner
	for _, name := range bannerNames {
		banner := Banner{Name: name}
		if lines, err := loadBanner(name); err == nil {
			banner.Preview = dropBlankRows(drawing("Aa", lines))
		}
		list = append(list, banner)
	}
	return list
}

func dropBlankRows(art string) string { 
	var rows []string
	for _, row := range strings.Split(art, "\n") {
		if strings.TrimSpace(row) != "" {
			rows = append(rows, row)
		}
	}
	return strings.Join(rows, "\n")
}
