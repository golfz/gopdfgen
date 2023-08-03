# Go PDF Generator
A simple PDF generator written in Go.
wkhtmltopdf is used to generate the PDF from HTML.

## Features
- Generate PDF from HTML
- Password protection

## Setup

1. Copy `fonts/` folder to your project root
2. (Optional) If you use Elastic Beanstalk, copy `install-wkhtmltopdf.config` in the `.ebextensions/` folder to your `.ebextensions/` folder

## Usage

**Import package**
```go
import "github.com/golfz/gopdfgen"
```

**Generate PDF from HTML String**

with password protection:
```go
b, err := gopdfgen.GenerateFromHTMLString("<h1>render from HTML string</h1>", "password")
```

without password protection:
```go
b, err := gopdfgen.GenerateFromHTMLString("<h1>render from HTML string</h1>", "")
```

**Generate PDF from HTML Template**

with password protection:
```go
b, err := gopdfgen.GenerateFromHTMLTemplate(htmlTemplateAsString, data, "password")
```

without password protection:
```go
b, err := gopdfgen.GenerateFromHTMLTemplate(htmlTemplateAsString, data, "")
```

## Git
add `_gopdfgen_temp/` to your .gitignore 




