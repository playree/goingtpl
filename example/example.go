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
	http.HandleFunc("/example3", handleExample3)
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

func handleExample3(w http.ResponseWriter, r *http.Request) {
	start := time.Now().UnixNano()

	funcMap := template.FuncMap{
		"repeat": func(s string, i int) string {
			return strings.Repeat(s, i)
		}}
	tpl := template.Must(goingtpl.ParseFileFuncs("page2.html", funcMap))

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
