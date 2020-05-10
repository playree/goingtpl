# goingtpl
`goingtpl`は、テンプレートファイルから別テンプレートファイルのインクルード記述をサポートしたテンプレートパッケージです。  
ついでにテンプレートのキャッシュ機能も備えています。  
`html/template`互換というか、テンプレートファイルをパースした結果は`template.Template`となっています。

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
### テンプレートのインクルード機能
テンプレートファイルに`{{include xxx.tpl}}`を記述します。
```html
[parent.tpl]

<!DOCTYPE html>
<html><body>
    <h1>Test code</h1>
{{template "footer"}}{{include "footer.tpl"}}
</body></html>
```

```html
[footer.tpl]

{{define "footer"}}
	<p>Footer</p>
{{end}}
```
あとは下記のように専用のパースメソッドを使用して、親となるテンプレートファイルをパースするだけで、`include`指定したファイルも一緒にパースされます。
```go
tpl := template.Must(goingtpl.ParseFile("parent.tpl"))
```

### テンプレートのキャッシュ機能
事前に下記のようにして機能をONにしておくだけです。  
これにより、一度パースしたテンプレートファイルはメモリ上にキャッシュされ、次回からはメモリ上からロードされるようになります。  
ディスクアクセスを減らすことができるので、パフォーマンスの向上に期待できます。  
逆にキャッシュ中はテンプレートファイルを変更しても反映されない為、開発中はOFFにしておくほうが効率的かもしれません。
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

下記にサンプルを配置しています。  
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

	http.HandleFunc("/example", handleExample)
	http.HandleFunc("/clear", handleClear)
	log.Fatal(http.ListenAndServe(":8088", nil))
}

func handleExample(w http.ResponseWriter, r *http.Request) {
	start := time.Now().UnixNano()

	funcMap := template.FuncMap{
		"repeat": func(s string, i int) string {
			return strings.Repeat(s, i)
		}}
	tpl := template.Must(goingtpl.ParseFileFuncs("parent.tpl", funcMap))
	// If you do not add a function
	// e.g. goingtpl.ParseFile("xxx.tpl")

	m := map[string]string{
		"Date": time.Now().Format("2006-01-02"),
		"Time": time.Now().Format("15:04:05"),
	}
	tpl.Execute(w, m)

	fmt.Printf("ExecTime=%d MicroSec\n",
		(time.Now().UnixNano()-start)/int64(time.Microsecond))
}

func handleClear(w http.ResponseWriter, r *http.Request) {
	goingtpl.ClearCache()
	fmt.Println("ClearCache")
}
```

templates/parent.tpl
```html
<!DOCTYPE html>
<html><body>
    <div>parent.tpl</div>
    <div style="padding-left: 2rem;">
        {{template "child01"}}{{include "parts/child01.tpl"}}
        {{template "child02-1"}}{{include "parts/child02.tpl"}}
        {{template "child02-2" .}}
    </div>
</body></html>
```

templates/child01.tpl
```html
{{define "child01"}}
<div>
    child01.tpl {{.Date}} {{.Time}}
    <div style="padding-left: 2rem;">
        {{template "child03-1"}}{{include "parts/child03.tpl"}}
    </div>
</div>
{{end}}
```

templates/child02.tpl
```html
{{define "child02-1"}}
<div>child02.tpl - 1</div>
{{end}}

{{define "child02-2"}}
<div>
    child02.tpl - 2
    <div style="padding-left: 2rem;">
        {{template "child03-2" .}}{{include "parts/child03.tpl"}}
    </div>
</div>
{{end}}
```

templates/child03.tpl
```html
{{define "child03-1"}}
<i>child03.tpl - 1</i>
{{end}}

{{define "child03-2"}}
<i>child03.tpl - 2 </i><br>
Arguments = {{.Date}} {{.Time}}<br>
Func now = {{now}}<br>
Func repeat = {{repeat "A" 5}}
{{end}}
```
下記が実行結果です。  
親ファイルをパースしただけで、`include`指定したファイルもパースされているのが分かると思います。

![Demo](https://user-images.githubusercontent.com/41541796/43353103-8f902b5a-926a-11e8-9234-1abb108ed30f.png)

## Licence
[MIT](https://github.com/playree/goingtpl/blob/master/LICENSE)