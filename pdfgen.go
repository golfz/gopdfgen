package gopdfgen

import (
	"bytes"
	"context"
	"fmt"
	"github.com/golfz/gopdfgen/files"
	"github.com/google/uuid"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"log"
	"os"
	"os/exec"
	"strings"
)

const (
	cmdWkhtmltopdf = "wkhtmltopdf"    // command to execute
	DefaultTempDir = "_gopdfgen_temp" // default temp directory
	keySessionUUID = "sessionUUID"    // key for context
)

type fileType string

const (
	FileTypeBody   fileType = "body"
	FileTypeHeader fileType = "header"
	FileTypeFooter fileType = "footer"
)

type PDFGenerator struct {
	bodyURL   BodyURLOption
	headerURL StringOption
	footerURL StringOption
	password  string
	ctx       context.Context
	tempDir   string
	pdfBuf    bytes.Buffer
}

// NewPDFGenerator creates a new PDFGenerator session
func NewPDFGenerator() (*PDFGenerator, error) {
	sessionUUID := uuid.New().String()
	ctx := context.Background()
	ctx = context.WithValue(ctx, keySessionUUID, sessionUUID)

	if _, err := os.Stat(DefaultTempDir); os.IsNotExist(err) {
		err := os.Mkdir(DefaultTempDir, os.ModePerm)
		if err != nil {
			log.Println(err)
			return nil, err
		}
	}

	return &PDFGenerator{
		bodyURL: BodyURLOption{
			bodyURL: "",
		},
		headerURL: StringOption{
			option: "--header-html",
		},
		footerURL: StringOption{
			option: "--footer-html",
		},
		password: "",
		ctx:      ctx,
		tempDir:  DefaultTempDir,
	}, nil
}

// SetBodyURL sets the url for the body of the pdf
func (pdfg *PDFGenerator) SetBodyURL(bodyURL string) {
	pdfg.bodyURL.Set(bodyURL)
}

// SetBodyHTML sets the html for the body of the pdf
func (pdfg *PDFGenerator) SetBodyHTML(bodyHTML []byte) {
	bodyFilePath := pdfg.getFilePath(FileTypeBody)
	err := files.WriteBytesToFile(bodyFilePath, bodyHTML)
	if err != nil {
		log.Println(err)
	}
	pdfg.bodyURL.Set(bodyFilePath)
}

// SetHeaderURL sets the url for the header of the pdf
func (pdfg *PDFGenerator) SetHeaderURL(headerURL string) {
	pdfg.headerURL.Set(headerURL)
}

// SetHeaderHTML sets the html for the header of the pdf
func (pdfg *PDFGenerator) SetHeaderHTML(headerHTML []byte) {
	headerFilePath := pdfg.getFilePath(FileTypeHeader)
	err := files.WriteBytesToFile(headerFilePath, headerHTML)
	if err != nil {
		log.Println(err)
	}
	pdfg.headerURL.Set(headerFilePath)
}

// SetFooterURL sets the url for the footer of the pdf
func (pdfg *PDFGenerator) SetFooterURL(footerURL string) {
	pdfg.footerURL.Set(footerURL)
}

// SetFooterHTML sets the html for the footer of the pdf
func (pdfg *PDFGenerator) SetFooterHTML(footerHTML []byte) {
	footerFilePath := pdfg.getFilePath(FileTypeFooter)
	err := files.WriteBytesToFile(footerFilePath, footerHTML)
	if err != nil {
		log.Println(err)
	}
	pdfg.footerURL.Set(footerFilePath)
}

// SetPassword sets the password for the pdf
// if password is empty, the pdf will not be encrypted
func (pdfg *PDFGenerator) SetPassword(password string) {
	pdfg.password = strings.TrimSpace(password)
}

// SetTempDir sets the temp directory for this session
func (pdfg *PDFGenerator) SetTempDir(tempDir string) {
	err := os.Remove(pdfg.tempDir)
	if err != nil {
		log.Println(err)
	}
	pdfg.tempDir = tempDir
	os.Mkdir(pdfg.tempDir, os.ModePerm)
}

// getFilePath returns the path to the file that will be generated for each session depending on the context and fileType
// pattern: {tempDir}/{UUID}-{fileType}.pdf
// e.g. _gopdfgen_temp/290bf5f8-6dfc-4736-9c4c-eccf043b8e49-body.pdf
func (pdfg *PDFGenerator) getFilePath(fileType fileType) string {
	sessionUUID := pdfg.ctx.Value("sessionUUID").(string)
	filePath := fmt.Sprintf("%s/%s-%s.html", pdfg.tempDir, sessionUUID, fileType)
	return filePath
}

// Cleanup removes all files in the temp directory for this session
func (pdfg *PDFGenerator) Cleanup() {
	sessionFiles := make([]string, 0)
	sessionFiles = append(sessionFiles, pdfg.getFilePath(FileTypeBody))
	sessionFiles = append(sessionFiles, pdfg.getFilePath(FileTypeHeader))
	sessionFiles = append(sessionFiles, pdfg.getFilePath(FileTypeFooter))
	for _, f := range sessionFiles {
		if _, err := os.Stat(f); os.IsNotExist(err) {
			continue // file does not exist
		}
		// if file exists, remove it
		err := files.RemoveFile(f)
		if err != nil {
			log.Println(err)
		}
	}
}

func (pdfg *PDFGenerator) Generate() error {
	args := make([]string, 0)
	args = append(args, pdfg.bodyURL.Parse()...)
	args = append(args, pdfg.headerURL.Parse()...)
	args = append(args, pdfg.footerURL.Parse()...)
	args = append(args, "-") // set output to stdout

	wk := exec.CommandContext(context.Background(), cmdWkhtmltopdf, args...)
	wk.Stdout = &pdfg.pdfBuf
	err := wk.Run()
	if err != nil {
		log.Println(err)
		return err
	}

	// Encrypt the pdf, if password is not empty.
	err = pdfg.encryptPDF()
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (pdfg *PDFGenerator) WriteFile(filepath string) error {
	err := files.WriteBytesToFile(filepath, pdfg.Bytes())
	return err
}

func (pdfg *PDFGenerator) Bytes() []byte {
	return pdfg.pdfBuf.Bytes()
}

// encryptPDF encrypts the pdf if password is not empty
func (pdfg *PDFGenerator) encryptPDF() error {
	if pdfg.password = strings.TrimSpace(pdfg.password); pdfg.password != "" {
		conf := model.NewDefaultConfiguration()

		// Set the passwords for encryption.
		conf.OwnerPW = pdfg.password
		conf.UserPW = pdfg.password

		// Set encryption mode. You can choose from several modes: "rc4_40", "rc4_128", "aes_128", "aes_256"
		conf.EncryptUsingAES = true
		conf.EncryptKeyLength = 256
		conf.ValidationMode = model.ValidationNone

		r := bytes.NewReader(pdfg.pdfBuf.Bytes())

		err := api.Encrypt(r, &pdfg.pdfBuf, conf)
		if err != nil {
			return err
		}
	}
	return nil
}
