package main

import (
	"fmt"
	"github.com/golfz/gopdfgen"
	"log"
)

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
	pdfg.SetBodyURL("https://iamgolfz.com")

	// set header and footer url
	//pdfg.HeaderHTML.Set("./asdfas/asdfasdf")
	//pdfg.FooterHTML.Set("./asdfas/asdfasdf")

	pdfg.SetPassword("1q2w3e4r")

	// generate pdf to internal buffer
	err = pdfg.Generate()
	if err != nil {
		log.Fatal(err)
	}

	// write pdf to file
	err = pdfg.WriteFile("./test-url.pdf")
	if err != nil {
		log.Fatal(err)
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
