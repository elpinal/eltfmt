package main

import (
	"bytes"
	"testing"

	"github.com/elpinal/gec/parser"
)

func TestProgram(t *testing.T) {
	input := []byte(` n  =  3+1
m = n * (3 +1)
f = \ x -> if true  then x-1 else (x * 2)
x = f n>4
lt = 1< 2
gt = 5 > 2
le = m <=n * 2
ge =1>=2
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
f = \x -> if true then x - 1 else (x * 2)
x = f n > 4
lt = 1 < 2
gt = 5 > 2
le = m <= n * 2
ge = 1 >= 2
m / n
`; buf.String() != want {
		t.Errorf("f.program = %q; want %q", buf.String(), want)
	}
}

func BenchmarkProgram(b *testing.B) {
	input := []byte(` n  =  3+1
m = n * (3 +1)
f = \ x -> if true  then x-1 else (x * 2)
x = f n>4
lt = 1< 2
gt = 5 > 2
le = m <=n * 2
ge =1>=2
m / n`)
	wd, err := parser.Parse(input)
	if err != nil {
		b.Errorf("parsing input: %v", err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		f := newFormatter(&buf)
		f.program(wd)
	}
}
