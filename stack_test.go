package errstack

import (
	"errors"
	"fmt"
	"io"
	"testing"

	pkgerrs "github.com/pkg/errors"
)

func Test_stackerr(t *testing.T) {
	err := stackabc1()
	t.Logf("%-v\n", err)
	t.Logf("%+v\n", pkgerrs.WithStack(errors.New("pkgerrs")))
}

func stackabc1() Error {
	return stackabc2().Throw()
}

func stackabc2() Error {
	return stackabc3().Throw()
}

func stackabc3() Error {
	return New(errors.New("Hello Errors!")).Throw()
}

func Benchmark_pkgerrors(b *testing.B) {
	type ifce interface {
		StackTrace() pkgerrs.StackTrace
	}

	b.Run("pkgerrrs", func(b *testing.B) {
		b.Run("print error string only", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				err := pkgerrs.WithStack(pkgerrs.New("pkgerrs"))
				fmt.Fprintf(io.Discard, "%s", err)
			}
		})

		b.Run("print error and stacktrace", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				err := pkgerrs.WithStack(pkgerrs.New("pkgerrs"))
				fmt.Fprintf(io.Discard, "%+v", err)
			}
		})
	})

	b.Run("errstack", func(b *testing.B) {
		b.Run("print error string only", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				err := NewString("errstk", WithStack())
				fmt.Fprintf(io.Discard, "%s", err)
			}
		})

		b.Run("print error and stacktrace", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				err := NewString("errstk", WithStack())
				fmt.Fprintf(io.Discard, "%-j", err)
			}
		})
	})
}
