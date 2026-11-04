package main

import (
	"bytes"
	"errors"
	"html/template"
	"log"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"
)

// The fonts users can pick. Each one lives in <name>.txt next to main.go.
var bannerNames = []string{"standard", "shadow", "thinkertoy"}

// Banner is one font card on the page: its name and a small sample.
type Banner struct {
	Name    string
	Preview string
}

// PageData feeds templates/index.html.
type PageData struct {
	Text    string
	Banner  string
	Result  string
	Banners []Banner
}

// ErrorData feeds templates/error.html.
type ErrorData struct {
	Code    int
	Title   string
	Message string
	Art     string
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		renderError(w, http.StatusNotFound, "This page doesn't exist. Head back and make some art instead.")
		return
	}
	if r.Method != http.MethodGet {
		renderError(w, http.StatusMethodNotAllowed, "This page can only be opened, not posted to.")
		return
	}
	renderPage(w, PageData{Banner: "standard"})
}

func asciiHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		renderError(w, http.StatusMethodNotAllowed, "Use the form on the home page to generate art.")
		return
	}

	// Textareas send Windows line breaks, and typing \n by hand still works.
	text := strings.ReplaceAll(r.FormValue("text"), "\r\n", "\n")
	text = strings.ReplaceAll(text, `\n`, "\n")
	banner := r.FormValue("banner")

	if text == "" {
		renderError(w, http.StatusBadRequest, "Type some text first, then hit generate.")
		return
	}
	for _, char := range text {
		if char != '\n' && (char < 32 || char > 126) {
			renderError(w, http.StatusBadRequest, "Only printable ASCII can be drawn: letters, numbers, spaces and common symbols.")
			return
		}
	}
	if !slices.Contains(bannerNames, banner) {
		renderError(w, http.StatusBadRequest, "Pick one of the three fonts: standard, shadow or thinkertoy.")
		return
	}

	lines, err := loadBanner(banner)
	if os.IsNotExist(err) {
		renderError(w, http.StatusNotFound, "The font file "+banner+".txt is missing.")
		return
	}
	if err != nil {
		log.Println("banner:", err)
		renderError(w, http.StatusInternalServerError, "The font file "+banner+".txt couldn't be read.")
		return
	}

	renderPage(w, PageData{Text: text, Banner: banner, Result: drawing(text, lines)})
}

// renderPage draws the generator page. The template is read on every request,
// so a missing file gives a 404 instead of stopping the server.
func renderPage(w http.ResponseWriter, data PageData) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		renderError(w, http.StatusNotFound, "The page template is missing.")
		return
	}

	data.Banners = previews()

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		log.Println("template:", err)
		renderError(w, http.StatusInternalServerError, "Something went wrong while drawing the page.")
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	buf.WriteTo(w)
}

// renderError shows the styled error page, with the status code drawn in ASCII.
// If the error template itself is missing, it falls back to plain text.
func renderError(w http.ResponseWriter, code int, message string) {
	data := ErrorData{Code: code, Title: http.StatusText(code), Message: message}
	if lines, err := loadBanner("standard"); err == nil {
		data.Art = drawing(strconv.Itoa(code), lines)
	}

	var buf bytes.Buffer
	tmpl, err := template.ParseFiles("templates/error.html")
	if err == nil {
		err = tmpl.Execute(&buf, data)
	}
	if err != nil {
		http.Error(w, data.Title, code)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(code)
	buf.WriteTo(w)
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

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/ascii-art", asciiHandler)

	log.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
