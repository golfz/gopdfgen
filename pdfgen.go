package gopdfgen

import (
	"bufio"
	"fmt"
	"github.com/google/uuid"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"io"
	"os"
	"os/exec"
)

func Generate(html string, password string) ([]byte, error) {
	tempFolderName := "gopdfgen_temp"
	s := uuid.New().String()

	htmlFilepath := fmt.Sprintf("%s/%s.html", tempFolderName, s)
	defer func() {
		err := os.Remove(htmlFilepath)
		if err != nil {
			fmt.Println(err)
		}
	}()

	pdfFilepath := fmt.Sprintf("%s/%s.pdf", tempFolderName, s)
	defer func() {
		err := os.Remove(pdfFilepath)
		if err != nil {
			fmt.Println(err)
		}
	}()

	h := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Title</title>
    <style>
        * {
            margin: 0;
        }
    </style>
</head>
<body>
	<h1>Hello %s</h1>
</body>
</html>`

	h = fmt.Sprintf(h, htmlFilepath)

	os.Mkdir(tempFolderName, os.ModePerm)

	f, err := os.Create(htmlFilepath)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer f.Close()

	_, err = f.WriteString(h)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	wk := exec.Command("wk/wkhtmltopdf.exe", htmlFilepath, pdfFilepath)
	out, err := wk.Output()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println(string(out))

	conf := model.NewDefaultConfiguration()

	// Set the passwords for encryption.
	conf.UserPW = password
	conf.OwnerPW = password

	// Set encryption mode. You can choose from several modes: "rc4_40", "rc4_128", "aes_128", "aes_256"
	conf.EncryptUsingAES = true
	conf.EncryptKeyLength = 256

	// Encrypt the pdf
	err = api.EncryptFile(pdfFilepath, pdfFilepath, conf)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("PDF encryption successful!")
	}

	pdfFile, err := os.Open(pdfFilepath)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer pdfFile.Close()

	reader := bufio.NewReader(pdfFile)
	content, _ := io.ReadAll(reader)

	return content, nil
}
