package main

import (
	"log"
	"net/http"
	"os"
	"slices"
	"strings"
)

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
