package main

import (
	"strings"
	"testing"
)

func TestFunctions(t *testing.T) {
	if isUrgent(1) || !isUrgent(2) || !isUrgent(3) {
		t.Fatal("проверь границу срочности 2")
	}
	if !strings.Contains(describe("чай", 2), "чай") {
		t.Fatal("пропало название")
	}
}
