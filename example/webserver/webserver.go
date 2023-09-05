package main

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
)

//go:embed html/*
var htmlFS embed.FS

//go:embed templates/*
var templateFS embed.FS

const (
	bodyTemplateFilePath = "templates/body.go.html"
	headerFilePath       = "html/header.html"
	footerFilePath       = "html/footer.html"
)

func main() {
	http.HandleFunc("/body", handlerBody)
	http.HandleFunc("/header", handlerHeader)
	http.HandleFunc("/footer", handlerFooter)

	log.Panic(http.ListenAndServe(":18080", nil))
}

func handlerBody(w http.ResponseWriter, r *http.Request) {
	templateStr, err := templateFS.ReadFile(bodyTemplateFilePath)
	if err != nil {
		log.Println(err)
	}

	t, err := template.New("body").Parse(string(templateStr))
	if err != nil {
		log.Println(err)
	}

	bodyData := createBodyData()

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	err = t.Execute(w, bodyData)
	if err != nil {
		log.Println(err)
	}
}

func handlerHeader(w http.ResponseWriter, r *http.Request) {
	headerHTML, err := htmlFS.ReadFile(headerFilePath)
	if err != nil {
		log.Println(err)
	}
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write(headerHTML)
}

func handlerFooter(w http.ResponseWriter, r *http.Request) {
	headerHTML, err := htmlFS.ReadFile(footerFilePath)
	if err != nil {
		log.Println(err)
	}
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write(headerHTML)
}

// ============================================================================================

type BodyData struct {
	ID   int
	Name string
	Time string
}

func createBodyData() []BodyData {
	bodyData := make([]BodyData, 0)

	for i := 1; i <= 1000; i++ {
		bodyData = append(bodyData, BodyData{
			ID:   i,
			Name: "Name " + strconv.Itoa(i),
			Time: time.Now().String(),
		})
	}

	return bodyData
}
