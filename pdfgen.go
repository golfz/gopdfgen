package gopdfgen

import (
	"embed"
	"os/exec"
)

//go:embed wk/*
var wk embed.FS

func Generate() {
	exec.Command("wk/wkhtmltopdf.exe", "index.html", "example.pdf").Output()
}
