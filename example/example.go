package main

import (
	"bytes"
	"embed"
	"fmt"
	"github.com/golfz/gopdfgen"
	"html/template"
	"log"
	"strconv"
	"time"
)

const outputFilePath = "example.pdf"

const testPassword = "1q2w3e4r"

//go:embed webserver/html/*
var htmlFS embed.FS

//go:embed webserver/templates/*
var templateFS embed.FS

const bodyTemplateFilePath = "webserver/templates/body.go.html"
const headerFilePath = "webserver/html/header.html"
const footerFilePath = "webserver/html/footer.html"

const bodyURL = "http://localhost:18080/body"
const headerURL = "http://localhost:18080/header"
const footerURL = "http://localhost:18080/footer"

func main() {
	//testWithURL()
	testWithHTML()
}

// testWithURL generates pdf with url
// please run webserver/webserver.go first
func testWithURL() {
	pdfg, err := gopdfgen.NewPDFGenerator()
	if err != nil {
		log.Println(err)
	}
	defer pdfg.Cleanup()

	// not necessary because it's default, but you can set the specific path
	//pdfg.SetTempDir("_new_gopdfgen_temp")

	// set body url
	pdfg.SetBodyURL(bodyURL)

	// set header and footer url
	pdfg.SetHeaderURL(headerURL)
	pdfg.SetFooterURL(footerURL)

	// set password
	pdfg.SetPassword(testPassword)

	// generate pdf to internal buffer
	err = pdfg.Generate()
	if err != nil {
		log.Println(err)
	}

	// write pdf to file
	err = pdfg.WriteFile(outputFilePath)
	if err != nil {
		log.Println(err)
	}

	// get pdf as bytes
	b := pdfg.Bytes()

	fmt.Printf("Done, %d bytes", len(b))
}

func testWithHTML() {
	pdfg, err := gopdfgen.NewPDFGenerator()
	if err != nil {
		log.Println(err)
	}
	defer pdfg.Cleanup()

	// not necessary because it's default, but you can set the specific path
	//pdfg.SetTempDir("_new_gopdfgen_temp")

	// set body html
	bodyHTML := getBodyHTML()
	pdfg.SetBodyHTML(bodyHTML)

	headerHTML, err := htmlFS.ReadFile(headerFilePath)
	if err != nil {
		log.Println(err)
	}

	footerHTML, err := htmlFS.ReadFile(footerFilePath)
	if err != nil {
		log.Println(err)
	}

	// set header and footer html
	pdfg.SetHeaderHTML(headerHTML)
	pdfg.SetFooterHTML(footerHTML)

	// set password
	//pdfg.SetPassword(testPassword)

	// generate pdf to internal buffer
	err = pdfg.Generate()
	if err != nil {
		log.Println(err)
	}

	// write pdf to file
	err = pdfg.WriteFile(outputFilePath)
	if err != nil {
		log.Println(err)
	}

	// get pdf as bytes
	b := pdfg.Bytes()

	fmt.Printf("Done, %d bytes", len(b))
}

// ============================================================================================

type BodyData struct {
	ID   int
	Name string
	Time string
}

func createBodyData() []BodyData {
	bodyData := make([]BodyData, 0)

	for i := 1; i <= 60; i++ {
		bodyData = append(bodyData, BodyData{
			ID:   i,
			Name: "Name " + strconv.Itoa(i),
			Time: time.Now().String(),
		})
	}

	return bodyData
}

func getBodyHTML() []byte {
	templateStr, err := templateFS.ReadFile(bodyTemplateFilePath)
	if err != nil {
		log.Println(err)
	}

	t, err := template.New("body").Parse(string(templateStr))
	if err != nil {
		log.Println(err)
	}

	bodyData := createBodyData()

	var outbuf bytes.Buffer

	err = t.Execute(&outbuf, bodyData)
	if err != nil {
		log.Println(err)
	}

	return outbuf.Bytes()
}
