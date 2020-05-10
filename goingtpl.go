// Package goingtpl is template parser that supports file inclusion
package goingtpl

import (
	"html/template"
	"io/ioutil"
	"strings"
)

const (
	leftDelim         = "{{"
	rightDelim        = "}}"
	includeFuncString = "include"
	extendsFuncString = "extends"
)

var baseDir = ""
var cacheON = false
var tplCache = map[string]*template.Template{}
var funcMap = template.FuncMap{
	includeFuncString: func(templatefile string) string { return "" },
	extendsFuncString: func(templatefile string) string { return "" },
}

// SetBaseDir Specify the directory where the template is placed.
// By specifying it, you can shorten the template file specification.
func SetBaseDir(dir string) {
	if dir == "" || strings.HasSuffix(dir, "/") {
		baseDir = dir
	} else {
		baseDir = dir + "/"
	}
}

// EnableCache Enable template caching.
// When the cache function is enabled, the template once read is cached in memory.
// You can reduce disk access.
func EnableCache(flg bool) {
	cacheON = flg
	if !flg {
		ClearCache()
	}
}

// ClearCache Clear the cached template.
func ClearCache() {
	tplCache = map[string]*template.Template{}
}

// AddFixedFunc Add a fixed function to use.
func AddFixedFunc(name string, fnc interface{}) {
	funcMap[name] = fnc
}

// ParseFile is template parser that supports file inclusion.
// If BaseDIr is specified, specify a file path based on BaseDir.
// e.g. {{include "xxx.html"}}
func ParseFile(filename string) (*template.Template, error) {
	return ParseFileFuncs(filename, template.FuncMap{})
}

// ParseFileFuncs is template parser that supports file inclusion.
// If BaseDIr is specified, specify a file path based on BaseDir.
// e.g. {{include "xxx.html"}}
func ParseFileFuncs(filename string, funcs template.FuncMap) (*template.Template, error) {
	if cacheON {
		// search template cache
		if tpl, ok := tplCache[filename]; ok {
			return tpl, nil
		}
	}

	tpl, err := nextParse(
		template.New(filename).Funcs(funcMap).Funcs(funcs),
		filename,
		map[string]bool{},
	)
	if err != nil {
		return nil, err
	}

	if cacheON {
		// add template cache
		tplCache[filename] = tpl
	}

	return tpl, nil
}

func nextParse(tpl *template.Template, filename string, comp map[string]bool) (*template.Template, error) {
	buf, err := ioutil.ReadFile(baseDir + filename)
	if err != nil {
		return nil, err
	}
	contents := string(buf)

	contents, tpl, err = nextExt(tpl, contents, comp)
	if err != nil {
		return nil, err
	}

	tpl, err = tpl.Parse(contents)
	if err != nil {
		return nil, err
	}
	comp[filename] = true

	incList := nextInc(contents, []string{})
	for _, inc := range incList {
		if _, ok := comp[inc]; ok {
			continue
		}

		tpl, err = nextParse(tpl, inc, comp)
		if err != nil {
			return nil, err
		}
	}
	return tpl, nil
}

func nextInc(contents string, list []string) []string {
	if start := strings.Index(contents, leftDelim); start >= 0 {
		start += 2
		if end := strings.Index(contents[start:], rightDelim); end >= 0 {
			end += start
			param := strings.Fields(contents[start:end])
			if len(param) == 2 && param[0] == includeFuncString {
				list = append(list, strings.Trim(param[1], `"`))
			}
			return nextInc(contents[end:], list)
		}
	}
	return list
}

func nextExt(tpl *template.Template, contents string, comp map[string]bool) (string, *template.Template, error) {
	if start := strings.Index(contents, leftDelim); start >= 0 {
		start += 2
		if end := strings.Index(contents[start:], rightDelim); end >= 0 {
			end += start
			param := strings.Fields(contents[start:end])
			if len(param) == 2 && param[0] == extendsFuncString {
				tpl, err := nextParse(tpl, strings.Trim(param[1], `"`), comp)
				return contents[end+2:], tpl, err
			}
		}
	}
	return contents, tpl, nil
}
