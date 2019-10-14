package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"strings"
)

func walk(dir string, handle func(string)) {
	fs, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	for _, file := range fs {
		if file.IsDir() {
			walk(dir+"/"+file.Name(), handle)
			continue
		}

		handle(dir + "/" + file.Name())
	}
}

func fixPrefixUpper(s string) string {
	if len(s) == 0 {
		return s
	}
	var p = make([]byte, len(s))
	copy(p, s)
	for i := 0; i < len(s); i++ {
		if s[i] > 'A' && s[i] <= 'Z' {
			p[i] = uint8(s[i] + 32)
		}
	}

	return string(p)
}

func processFile(filename string) {
	if !strings.HasSuffix(filename, ".go") {
		return
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, nil, 0)
	if err != nil {
		panic(err)
	}

	for _, d := range f.Decls {
		if fn, isFn := d.(*ast.FuncDecl); isFn && len(fn.Type.Params.List) > 0 {
			for _, param := range fn.Type.Params.List {
				for _, paramName := range param.Names {
					if len(paramName.Name) > 0 && paramName.Name[0] >= 'A' && paramName.Name[0] <= 'Z' {
						fixed := fixPrefixUpper(paramName.Name)
						if fixed != paramName.Name {
							fmt.Printf("filename: %s, line: %v, func: %s, param: %s should startwith lowercase.\n", filename, d.Pos(), fn.Name, paramName.Name)
						}
					}
				}
			}
		}
	}
}

func main() {
	var dir = "."
	if len(os.Args) > 1 {
		dir = os.Args[1]
	}
	walk(dir, processFile)
}
