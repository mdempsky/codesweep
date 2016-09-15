package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/types"
	"log"

	"github.com/kisielk/gotool"
	"golang.org/x/tools/go/loader"
)

func main() {
	flag.Parse()

	importPaths := gotool.ImportPaths(flag.Args())
	if len(importPaths) == 0 {
		return
	}

	var conf loader.Config
	for _, importPath := range importPaths {
		conf.Import(importPath)
	}
	prog, err := conf.Load()
	if err != nil {
		log.Fatal(err)
	}

	used := make(map[string]bool)
	for _, pkg := range prog.InitialPackages() {
		for sym, obj := range pkg.Uses {
			if skip(sym, obj) {
				continue
			}
			used[obj.Id()] = true
		}
	}
	for _, pkg := range prog.InitialPackages() {
		for sym, obj := range pkg.Defs {
			if skip(sym, obj) {
				continue
			}
			if !used[obj.Id()] {
				fmt.Println(conf.Fset.Position(sym.Pos()), obj)
			}
		}
	}
}

func skip(sym *ast.Ident, obj types.Object) bool {
	if sym.Name == "_" || obj == nil {
		return true
	}
	return obj.Parent() == nil || obj.Parent().Parent() == nil
}
