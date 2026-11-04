package main

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

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
