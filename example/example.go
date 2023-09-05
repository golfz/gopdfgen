package main

import (
	"embed"
	"fmt"
	"github.com/golfz/gopdfgen"
	"log"
)

var websiteURL = "https://iamgolfz.com"

var outputFilePath = "test-url.pdf"

var headerFilePath = "html/header.html"
var footerFilePath = "html/footer.html"

var testPassword = "1q2w3e4r"

//go:embed html/*
var htmlFS embed.FS

func main() {
	testWithURL()
	//testWithHTML()
}

func testWithURL() {
	pdfg, err := gopdfgen.NewPDFGenerator()
	if err != nil {
		log.Println(err)
	}
	defer pdfg.Cleanup()

	// not necessary because it's default, but you can set the specific path
	//pdfg.SetTempDir("_new_gopdfgen_temp")

	// set body url
	pdfg.SetBodyURL(websiteURL)

	headerHTML, err := htmlFS.ReadFile(headerFilePath)
	if err != nil {
		log.Println(err)
	}

	footerHTML, err := htmlFS.ReadFile(footerFilePath)
	if err != nil {
		log.Println(err)
	}

	pdfg.SetHeaderHTML(headerHTML)
	pdfg.SetFooterHTML(footerHTML)

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

//func testWithHTML() {
//	pdfg, err := NewPDFGenerator()
//	if err != nil {
//		log.Panic(err)
//	}
//	defer pdfg.Cleanup()
//
//	// not necessary because it's default, but you can set the specific path
//	pdfg.TempDir.Set("_gopdfgen_temp")
//
//	// set body url
//	pdfg.BodyHTML.Set("<html><body><h1>Hello world</h1></body></html>")
//
//	// set header and footer url
//	pdfg.HeaderHTML.Set("<html><body><h1>Hello world</h1></body></html>")
//	pdfg.FooterHTML.Set("<html><body><h1>Hello world</h1></body></html>")
//
//	// generate pdf to internal buffer
//	err = pdfg.Generate()
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// write pdf to file
//	err = pdfg.WriteFile("./simplesample.pdf")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// get pdf as bytes
//	b := pdfg.Bytes()
//
//	fmt.Sprintf("Done, %d bytes", len(b))
//}
