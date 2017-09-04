package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/elpinal/gec/ast"
	"github.com/elpinal/gec/parser"

	"github.com/pkg/errors"
)

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stdout, "eltfmt: no Elacht source file given")
		os.Exit(1)
	}
	b, err := ioutil.ReadFile(flag.Arg(0))
	if err != nil {
		fmt.Fprintf(os.Stdout, "eltfmt: %v\n", err)
		os.Exit(1)
	}
	err = run(b, flag.Arg(0))
	if err != nil {
		fmt.Fprintln(os.Stdout, err)
		os.Exit(1)
	}
}

func run(src []byte, filename string) error {
	wd, err := parser.Parse(src)
	if err != nil {
		return errors.Wrapf(err, "parsing %s", filename)
	}
	w := bufio.NewWriter(os.Stdout)
	f := newFormatter(w)
	err = f.program(wd)
	if err != nil {
		return errors.Wrapf(err, "formatting %s", filename)
	}
	return w.Flush()
}

type writer interface {
	io.Writer
	WriteString(string) (int, error)
	WriteRune(rune) (int, error)
}

type formatter struct {
	w writer
}

func newFormatter(w writer) formatter {
	return formatter{w: w}
}

func (f formatter) program(wd *ast.WithDecls) error {
	err := f.decls(wd.Decls)
	if err != nil {
		return err
	}
	err = f.expr(wd.Expr)
	if err != nil {
		return err
	}
	f.w.WriteRune('\n')
	return nil
}

func (f formatter) decls(decls []*ast.Decl) error {
	for _, d := range decls {
		err := f.decl(d)
		if err != nil {
			return err
		}
	}
	return nil
}

func (f formatter) decl(decl *ast.Decl) error {
	f.w.WriteString(decl.LHS.Lit)
	f.w.WriteString(" = ")
	err := f.expr(decl.RHS)
	if err != nil {
		return err
	}
	f.w.WriteRune('\n')
	return nil
}

func (f formatter) expr(e ast.Expr) error {
	switch x := e.(type) {
	case *ast.Int:
		f.w.WriteString(x.X.Lit)
	case *ast.Bool:
		f.w.WriteString(x.X.Lit)
	case *ast.Ident:
		f.w.WriteString(x.Name.Lit)
	case *ast.Add:
		err := f.expr(x.X)
		if err != nil {
			return err
		}
		f.w.WriteString(" + ")
		err = f.expr(x.Y)
		if err != nil {
			return err
		}
	case *ast.Sub:
		err := f.expr(x.X)
		if err != nil {
			return err
		}
		f.w.WriteString(" - ")
		err = f.expr(x.Y)
		if err != nil {
			return err
		}
	case *ast.Mul:
		err := f.expr(x.X)
		if err != nil {
			return err
		}
		f.w.WriteString(" * ")
		err = f.expr(x.Y)
		if err != nil {
			return err
		}
	case *ast.Div:
		err := f.expr(x.X)
		if err != nil {
			return err
		}
		f.w.WriteString(" / ")
		err = f.expr(x.Y)
		if err != nil {
			return err
		}
	case *ast.App:
		err := f.expr(x.Fn)
		if err != nil {
			return err
		}
		f.w.WriteRune(' ')
		err = f.expr(x.Arg)
		if err != nil {
			return err
		}
	case *ast.Abs:
		f.w.WriteRune('\\')
		f.w.WriteString(x.Param.Lit)
		f.w.WriteString(" -> ")
		err := f.expr(x.Body)
		if err != nil {
			return err
		}
	case *ast.If:
		f.w.WriteString("if ")
		err := f.expr(x.Cond)
		if err != nil {
			return err
		}
		f.w.WriteString(" then ")
		err = f.expr(x.E1)
		if err != nil {
			return err
		}
		f.w.WriteString(" else ")
		err = f.expr(x.E2)
		if err != nil {
			return err
		}
	case *ast.ParenExpr:
		f.w.WriteRune('(')
		f.expr(x.X)
		f.w.WriteRune(')')
	case *ast.Cmp:
		err := f.expr(x.LHS)
		if err != nil {
			return err
		}
		switch x.Op {
		case ast.Eq:
			f.w.WriteString(" == ")
		case ast.NE:
			f.w.WriteString(" /= ")
		case ast.LT:
			f.w.WriteString(" < ")
		case ast.GT:
			f.w.WriteString(" > ")
		case ast.LE:
			f.w.WriteString(" <= ")
		case ast.GE:
			f.w.WriteString(" >= ")
		default:
			return fmt.Errorf("unknown operator: %v (in %v)", x.Op, e)
		}
		err = f.expr(x.RHS)
		if err != nil {
			return err
		}
	case *ast.NilList:
		f.w.WriteString("[]")
	default:
		return fmt.Errorf("unknown expression: %[1]v (type: %[1]T)", e)
	}
	return nil
}
