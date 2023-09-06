# Go PDF Generator
A simple PDF generator written in Go.
wkhtmltopdf is used to generate the PDF from HTML.

## Features
- Generate PDF from HTML
- Password protection

## Setup

1. Copy `fonts/` folder to your project root
2. (Optional) If you use Elastic Beanstalk, copy `install-wkhtmltopdf.config` in the `.ebextensions/` folder to your `.ebextensions/` folder
3. (Optional) If you use other OS, download wkhtmltopdf from [https://wkhtmltopdf.org/downloads.html](https://wkhtmltopdf.org/downloads.html) and set the `PATH` environment variable to the wkhtmltopdf executable

## Usage

**Import package**
```go
import "github.com/golfz/gopdfgen"
```

**Example : Generate PDF from URL**

```go
import "github.com/golfz/gopdfgen"

func main() {
    pdfg, err := gopdfgen.NewPDFGenerator()
    if err != nil {
        log.Panic(err)
    }
    defer pdfg.Cleanup()
    
    // not necessary because there's default, 
    // but you can set the specific path
    pdfg.SetTempDir("_custom_gopdfgen_temp")
    
    // set body url
    pdfg.SetBodyURL("https://iamgolfz.com")
    
    // set header and footer url
    pdfg.SetHeaderURL("https://example.com/header.html")
    pdfg.SetFooterURL("https://example.com/footer.html")
    
    // set password
    pdfg.SetPassword("12345678")
    
    // generate pdf to internal buffer
    err = pdfg.Generate()
    if err != nil {
        log.Panic(err)
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
```

**Example : Generate PDF from HTML**

```go
import "github.com/golfz/gopdfgen"

func main() {
    pdfg, err := gopdfgen.NewPDFGenerator()
    if err != nil {
        log.Panic(err)
    }
    defer pdfg.Cleanup()
    
    // not necessary because there's default, 
    // but you can set the specific path
    pdfg.SetTempDir("_custom_gopdfgen_temp")
    
    // set body url
    pdfg.SetBodyHTML("<h1>Hello World</h1>")
    
    // set header and footer url
    pdfg.SetHeaderHTML("<h1>Header</h1><hr>")
    pdfg.SetFooterHTML("<hr><h1>Footer</h1>")
    
    // set password
    pdfg.SetPassword("12345678")
    
    // generate pdf to internal buffer
    err = pdfg.Generate()
    if err != nil {
        log.Panic(err)
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
```

**Example : Page number in header and footer**

```html
<!doctype html>
<html>
<head>
    <meta charset="utf-8">
    <script>
        function substitutePdfVariables() {

            function getParameterByName(name) {
                var match = RegExp('[?&]' + name + '=([^&]*)').exec(window.location.search);
                return match && decodeURIComponent(match[1].replace(/\+/g, ' '));
            }

            function substitute(name) {
                var value = getParameterByName(name);
                var elements = document.getElementsByClassName(name);

                for (var i = 0; elements && i < elements.length; i++) {
                    elements[i].textContent = value;
                }
            }

            ['frompage', 'topage', 'page', 'webpage', 'section', 'subsection', 'subsubsection']
                .forEach(function (param) {
                    substitute(param);
                });


            function showPage1Header(name) {
                if (name === 'page') {
                    var value = getParameterByName(name);
                    var elements = document.getElementById('special-header');

                    if (value === '1') {
                        elements.style.display = 'block';
                    }
                }
            }

            ['frompage', 'topage', 'page', 'webpage', 'section', 'subsection', 'subsubsection']
                .forEach(function (param) {
                    showPage1Header(param);
                });
        }
    </script>
</head>
<body onload="substitutePdfVariables()">
    <h1 style="text-align: center">Header 1</h1>
    <p style="text-align: right">
        Page <span class="page"></span> of <span class="topage"></span>
    </p>
    <hr>
</body>
</html>
```

## Git
add `_gopdfgen_temp/` to your .gitignore 




