package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestDecisionsAndLoop(t *testing.T) {
	f, err := parser.ParseFile(token.NewFileSet(), "main.go", nil, 0)
	if err != nil {
		t.Fatal(err)
	}
	conditions, loops := 0, 0
	ast.Inspect(f, func(n ast.Node) bool {
		switch n.(type) {
		case *ast.IfStmt:
			conditions++
		case *ast.ForStmt:
			loops++
		}
		return true
	})
	if conditions == 0 || loops == 0 {
		t.Fatal("в упражнении нужны условие и цикл")
	}
}
