package gopdfgen

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/google/uuid"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"html/template"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

const (
	TempFolderName = "_gopdfgen_temp"
)

func GenerateFromHTMLTemplate(htmlTemplate string, data interface{}, password string) ([]byte, error) {
	t, err := template.New("template").Parse(htmlTemplate)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var htmlStr bytes.Buffer
	t.Execute(&htmlStr, data)

	return GenerateFromHTMLString(htmlStr.String(), password)
}

func GenerateFromHTMLString(htmlStr string, password string) ([]byte, error) {
	s := uuid.New().String()

	htmlFilepath := fmt.Sprintf("%s/%s.html", TempFolderName, s)
	defer func() {
		err := os.Remove(htmlFilepath)
		if err != nil {
			log.Println(err)
		}
	}()

	os.Mkdir(TempFolderName, os.ModePerm)

	f, err := os.Create(htmlFilepath)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer f.Close()

	_, err = f.WriteString(htmlStr)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return GenerateFromURL(htmlFilepath, password)
}

func GenerateFromURL(url string, password string) ([]byte, error) {
	s := uuid.New().String()

	pdfFilepath := fmt.Sprintf("%s/%s.pdf", TempFolderName, s)
	defer func() {
		err := os.Remove(pdfFilepath)
		if err != nil {
			log.Println(err)
		}
	}()

	wk := exec.Command("wkhtmltopdf", url, pdfFilepath)
	out, err := wk.Output()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println("[gopdfgen] call command wkhtmltopdf success", string(out))

	// Encrypt the pdf, if password is not empty.
	if password = strings.TrimSpace(password); password != "" {
		conf := model.NewDefaultConfiguration()

		// Set the passwords for encryption.
		conf.OwnerPW = password
		conf.UserPW = password

		// Set encryption mode. You can choose from several modes: "rc4_40", "rc4_128", "aes_128", "aes_256"
		conf.EncryptUsingAES = true
		conf.EncryptKeyLength = 256
		conf.ValidationMode = model.ValidationNone

		// Encrypt the pdf
		err = api.EncryptFile(pdfFilepath, pdfFilepath, conf)
		if err != nil {
			log.Println("[gopdfgen]", err)
		} else {
			log.Println("[gopdfgen] pdf encryption successful")
		}
	}

	pdfFile, err := os.Open(pdfFilepath)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer pdfFile.Close()

	reader := bufio.NewReader(pdfFile)
	content, _ := io.ReadAll(reader)

	return content, nil
}
