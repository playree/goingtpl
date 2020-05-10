# goingtpl
[Japanese here](https://github.com/playree/goingtpl/blob/master/README.JP.md)

goingtpl is template parser that supports file include and extends.  
There is also a function to cache templates.

## Table of content
- [Install](#install)
- [Usage](#usage)
- [API](#api)
- [Examples](#examples)
- [Licence](#licence)

## Install
```bash
go get github.com/playree/goingtpl
```

## Usage
### Template include function
Write `{{include "xxx.html"}}` in the template file.
```html
[parent.html]

<!DOCTYPE html>
<html><body>
    <h1>Test code</h1>
{{template "footer"}}{{include "footer.html"}}
</body></html>
```

```html
[footer.html]

{{define "footer"}}
	<p>Footer</p>
{{end}}
```
All you have to do is parse the parent file with the prepared method.
```go
tpl := template.Must(goingtpl.ParseFile("parent.html"))
```

After parsing.
```html
<!DOCTYPE html>
<html><body>
    <h1>Test code</h1>
	<p>Footer</p>
</body></html>
```

### Template extends function
Write `{{extends "xxx.html"}}` in the template file.
```html
[base.html]

<!DOCTYPE html>
<html><body>
	<h1>{{template "title" .}}</h1>
    <div style="background-color: #ddf;">
		{{template "content" .}}
	</div>
	<p>Footer</p>
</body></html>
```

```html
[page1.html]

{{extends "base.html"}}
{{define "title"}}Page1{{end}}
{{define "content"}}
	<p>This is Page1.</p>
{{end}}
```

Be sure to write `{{extends "xxx.html"}}` at the beginning.

All you have to do is parse the file with the prepared method.

```go
tpl := template.Must(goingtpl.ParseFile("page1.html"))
```

After parsing.
```html
<!DOCTYPE html>
<html><body>
	<h1>Page1</h1>
    <div style="background-color: #ddf;">
		<p>This is Page1.</p>
	</div>
	<p>Footer</p>
</body></html>
```

### Template caching function
We call the following method beforehand.
```go
goingtpl.EnableCache(true)
```

## API
- ParseFile(filename string) (*template.Template, error)
- ParseFileFuncs(filename string, funcs template.FuncMap) (*template.Template, error)
- SetBaseDir(dir string)
- EnableCache(flg bool)
- ClearCache()
- AddFixedFunc(name string, fnc interface{})

## Examples

https://github.com/playree/goingtpl/tree/master/example

example.go
```go
package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/playree/goingtpl"
)

func main() {
	// goingtpl setting
	goingtpl.SetBaseDir("./templates")
	goingtpl.EnableCache(true)
	goingtpl.AddFixedFunc(
		"now",
		func() string {
			return time.Now().Format("2006-01-02 15:04:05")
		},
	)

	http.HandleFunc("/example1", handleExample1)
	http.HandleFunc("/example2", handleExample2)
	http.HandleFunc("/clear", handleClear)
	log.Fatal(http.ListenAndServe(":8088", nil))
}

func handleExample1(w http.ResponseWriter, r *http.Request) {
	start := time.Now().UnixNano()

	funcMap := template.FuncMap{
		"repeat": func(s string, i int) string {
			return strings.Repeat(s, i)
		}}
	tpl := template.Must(goingtpl.ParseFileFuncs("parent.html", funcMap))
	// If you do not add a function
	// e.g. goingtpl.ParseFile("xxx.html")

	m := map[string]string{
		"Date": time.Now().Format("2006-01-02"),
		"Time": time.Now().Format("15:04:05"),
	}
	err := tpl.Execute(w, m)
	if err != nil {
		panic(err)
	}

	fmt.Printf("ExecTime=%d MicroSec\n",
		(time.Now().UnixNano()-start)/int64(time.Microsecond))
}

func handleExample2(w http.ResponseWriter, r *http.Request) {
	start := time.Now().UnixNano()

	tpl := template.Must(goingtpl.ParseFileFuncs("page1.html", nil))
	// If you do not add a function
	// e.g. goingtpl.ParseFile("xxx.html")

	m := map[string]string{
		"Date": time.Now().Format("2006-01-02"),
		"Time": time.Now().Format("15:04:05"),
	}
	err := tpl.Execute(w, m)
	if err != nil {
		panic(err)
	}

	fmt.Printf("ExecTime=%d MicroSec\n",
		(time.Now().UnixNano()-start)/int64(time.Microsecond))
}

func handleClear(w http.ResponseWriter, r *http.Request) {
	goingtpl.ClearCache()
	fmt.Println("ClearCache")
}
```

templates/parent.html
```html
<!DOCTYPE html>
<html><body>
    <div>parent.html</div>
    <div style="padding-left: 2rem;">
        {{template "child01" .}}{{include "parts/child01.html"}}
        {{template "child02-1" .}}{{include "parts/child02.html"}}
        {{template "child02-2" .}}
    </div>
</body></html>
```

templates/parts/child01.html
```html
{{define "child01"}}
<div>
    child01.html {{.Date}} {{.Time}}
    <div style="padding-left: 2rem;">
        {{template "child03-1"}}{{include "parts/child03.html"}}
    </div>
</div>
{{end}}
```

templates/parts/child02.html
```html
{{define "child02-1"}}
<div>child02.html - 1</div>
{{end}}

{{define "child02-2"}}
<div>
    child02.html - 2
    <div style="padding-left: 2rem;">
        {{template "child03-2" .}}{{include "parts/child03.html"}}
    </div>
</div>
{{end}}
```

templates/parts/child03.html
```html
{{define "child03-1"}}
<i>child03.html - 1</i>
{{end}}

{{define "child03-2"}}
<i>child03.html - 2 </i><br>
Arguments = {{.Date}} {{.Time}}<br>
Func now = {{now}}<br>
Func repeat = {{repeat "A" 5}}
{{end}}
```

templates/parts/base.html
```html
<!DOCTYPE html>
<html><body>
	<h1>{{template "title" .}}</h1>
	<p>This is a sample of extends.</p>
	<div style="background-color: #ddf;">
		{{template "content" .}}
	</div>
</body></html>
```

templates/page1.html
```html
{{extends "parts/base.html"}}
{{define "title"}}Page1{{end}}
{{define "content"}}
	This is Page1.
	<p>
		Loaded {{.Date}} {{.Time}}
	</p>
{{end}}
```

Result of operation

![Demo](https://user-images.githubusercontent.com/41541796/81497209-f5678100-92f7-11ea-91ac-d77e6add3ea6.png)

![Demo](https://user-images.githubusercontent.com/41541796/81497214-fa2c3500-92f7-11ea-87a8-cc5ac48d1ce0.png)

## Licence
[MIT](https://github.com/playree/goingtpl/blob/master/LICENSE)