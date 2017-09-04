package main

import (
	"bytes"
	"testing"

	"github.com/elpinal/gec/parser"
)

func TestProgram(t *testing.T) {
	input := []byte(` n  =  3+1
m = n * (3 +1)
m / n`)
	wd, err := parser.Parse(input)
	if err != nil {
		t.Errorf("parsing input: %v", err)
	}
	var buf bytes.Buffer
	f := newFormatter(&buf)
	err = f.program(wd)
	if err != nil {
		t.Errorf("formatting input: %v", err)
	}
	if want := `n = 3 + 1
m = n * (3 + 1)
m / n
`; buf.String() != want {
		t.Errorf("f.program = %q; want %q", buf.String(), want)
	}
}
