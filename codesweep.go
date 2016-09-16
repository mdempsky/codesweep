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
			if !interesting(sym, obj) {
				continue
			}
			used[obj.Id()] = true
		}
	}
	for _, pkg := range prog.InitialPackages() {
		for sym, obj := range pkg.Defs {
			if !interesting(sym, obj) {
				continue
			}
			if !used[obj.Id()] {
				fmt.Println(conf.Fset.Position(sym.Pos()), obj)
			}
		}
	}
}

func interesting(sym *ast.Ident, obj types.Object) bool {
	// Skip blank identifiers.
	if sym.Name == "_" {
		return false
	}

	// Skip identifiers that don't actually declare objects (e.g., package names).
	if obj == nil {
		return false
	}

	// Skip universal objects.
	if obj.Pkg() == nil {
		return false
	}

	// Fields and methods are interesting.
	if obj, ok := obj.(*types.Var); ok && obj.IsField() {
		return true
	}
	if sig, ok := obj.Type().(*types.Signature); ok && sig.Recv() != nil {
		return true
	}

	// Objects declared at packages scope are interesting,
	// except for init and main.
	if obj.Parent() == obj.Pkg().Scope() && !(sym.Name == "init" || sym.Name == "main" && obj.Pkg().Name() == "main") {
		return true
	}

	// Otherwise, ignore for now.
	return false
}
