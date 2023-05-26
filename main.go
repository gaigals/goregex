package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

const (
	paramRegexPattern = "regexPattern"
	paramTextValue    = "textValue"
)

type Result struct {
	HTML  string
	Error string
}

func main() {
	http.HandleFunc("/", handleHTMLContent)
	http.HandleFunc("/static/", serveStaticFiles)
	http.HandleFunc("/regex", handleRegexPost)

	fmt.Println("Server started on http://localhost:80")

	log.Fatalln(http.ListenAndServe(":80", nil))
}

func serveStaticFiles(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Path[len("/static/"):]
	serveFile := filepath.Join("static", filePath)

	http.ServeFile(w, r, serveFile)
}

func handleHTMLContent(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	tmpl, err := template.ParseFiles("templ/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleRegexPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the form values
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	input := r.FormValue(paramTextValue)
	regexPattern := r.FormValue(paramRegexPattern)

	jsonBytes, err := NewBuilder(input).
		CompileRegex(regexPattern).
		MatchString().
		BuildResponse()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set the JSON content type
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(jsonBytes)
}

// func matchString(regexPattern, value string) Result {
// 	regex, err := regexp.Compile(regexPattern)
// 	if err != nil {
// 		return Result{Error: err.Error()}
// 	}
//
// 	return Result{Matched: regex.FindAllStringSubmatchIndex(value, -1)}
// }
