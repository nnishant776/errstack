package errstack

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"testing"

	pkgerrs "github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func Test_StackTracedErrorWithDefaults(t *testing.T) {
	goRoot := runtime.GOROOT()
	currDir, err := os.Getwd()
	if err != nil {
		t.Errorf("Failed to get current directory: error: %s", err)
	}

	t.Run("automatic stacktrace collection", func(t *testing.T) {
		stackabc3 := func() Error {
			return New(errors.New("Hello Errors!"), WithStack())
		}
		stackabc2 := func() Error {
			return stackabc3()
		}
		stackabc1 := func() Error {
			return stackabc2()
		}

		err := stackabc1()

		t.Run("printing formats", func(t *testing.T) {
			type testCase struct {
				format string
				output strings.Builder
				expect string
			}

			cases := []testCase{
				{
					format: "s",
					expect: "Hello Errors!",
				},
				{
					format: "+s",
					expect: "Hello Errors!",
				},
				{
					format: "v",
					expect: "Hello Errors!=>github.com/nnishant776/errstack.Test_StackTracedErrorWithDefaults.func1.1;github.com/nnishant776/errstack.Test_StackTracedErrorWithDefaults.func1.2;github.com/nnishant776/errstack.Test_StackTracedErrorWithDefaults.func1.3;github.com/nnishant776/errstack.Test_StackTracedErrorWithDefaults.func1;testing.tRunner;runtime.goexit",
				},
				{
					format: " v",
					expect: "Hello Errors!\ngithub.com/nnishant776/errstack.Test_StackTracedErrorWithDefaults.func1.1\ngithub.com/nnishant776/errstack.Test_StackTracedErrorWithDefaults.func1.2\ngithub.com/nnishant776/errstack.Test_StackTracedErrorWithDefaults.func1.3\ngithub.com/nnishant776/errstack.Test_StackTracedErrorWithDefaults.func1\ntesting.tRunner\nruntime.goexit",
				},
				{
					format: "-v",
					expect: fmt.Sprintf("Hello Errors!=>github.com/nnishant776/errstack.Test_StackTracedErrorWithDefaults.func1.1@%s/stack_test.go:26;github.com/nnishant776/errstack.Test_StackTracedErrorWithDefaults.func1.2@%s/stack_test.go:29;github.com/nnishant776/errstack.Test_StackTracedErrorWithDefaults.func1.3@%s/stack_test.go:32;github.com/nnishant776/errstack.Test_StackTracedErrorWithDefaults.func1@%s/stack_test.go:35;testing.tRunner@%s/src/testing/testing.go:1690;runtime.goexit@%s/src/runtime/asm_amd64.s:1700", currDir, currDir, currDir, currDir, goRoot, goRoot),
				},
				{
					format: "+v",
					expect: fmt.Sprintf("Hello Errors!\ngithub.com/nnishant776/errstack.Test_StackTracedErrorWithDefaults.func1.1@%s/stack_test.go:26\ngithub.com/nnishant776/errstack.Test_StackTracedErrorWithDefaults.func1.2@%s/stack_test.go:29\ngithub.com/nnishant776/errstack.Test_StackTracedErrorWithDefaults.func1.3@%s/stack_test.go:32\ngithub.com/nnishant776/errstack.Test_StackTracedErrorWithDefaults.func1@%s/stack_test.go:35\ntesting.tRunner@%s/src/testing/testing.go:1690\nruntime.goexit@%s/src/runtime/asm_amd64.s:1700", currDir, currDir, currDir, currDir, goRoot, goRoot),
				},
				{
					format: "#v",
					expect: fmt.Sprintf("Hello Errors!\n#5: github.com/nnishant776/errstack.Test_StackTracedErrorWithDefaults.func1.1@%s/stack_test.go:26\n#4: github.com/nnishant776/errstack.Test_StackTracedErrorWithDefaults.func1.2@%s/stack_test.go:29\n#3: github.com/nnishant776/errstack.Test_StackTracedErrorWithDefaults.func1.3@%s/stack_test.go:32\n#2: github.com/nnishant776/errstack.Test_StackTracedErrorWithDefaults.func1@%s/stack_test.go:35\n#1: testing.tRunner@%s/src/testing/testing.go:1690\n#0: runtime.goexit@%s/src/runtime/asm_amd64.s:1700", currDir, currDir, currDir, currDir, goRoot, goRoot),
				},
				{
					format: "j",
					expect: func() string {
						m := map[string]any{"error": err.Error(), "trace": err.StackTrace()}
						b, _ := json.Marshal(m)
						return string(b) + "\n"
					}(),
				},
				{
					format: "+j",
					expect: func() string {
						m := map[string]any{"error": err.Error(), "trace": err.StackTrace()}
						b, _ := json.MarshalIndent(m, "", "  ")
						return string(b) + "\n"
					}(),
				},
				{
					format: "+4j",
					expect: func() string {
						m := map[string]any{"error": err.Error(), "trace": err.StackTrace()}
						b, _ := json.MarshalIndent(m, "", "    ")
						return string(b) + "\n"
					}(),
				},
			}

			for _, c := range cases {
				t.Run(c.format, func(t *testing.T) {
					fmt.Fprintf(&c.output, "%"+c.format, err)
					if c.output.String() != c.expect {
						assert.Equal(
							t,
							c.expect,
							c.output.String(),
							"Doesn't match the expected output",
						)
					}
				})
			}
		})
	})

	t.Run("manual stacktrace collection", func(t *testing.T) {
		stackabc3 := func() Error {
			return New(errors.New("Hello Errors!")).Throw()
		}
		stackabc2 := func() Error {
			return stackabc3().Throw()
		}
		stackabc1 := func() Error {
			return stackabc2().Throw()
		}

		err := stackabc1().Throw()

		t.Run("printing formats", func(t *testing.T) {
			type testCase struct {
				format string
				output strings.Builder
				expect string
			}

			cases := []testCase{
				{
					format: "s",
					expect: "Hello Errors!",
				},
				{
					format: "+s",
					expect: "Hello Errors!",
				},
				{
					format: "v",
					expect: "Hello Errors!=>github.com/nnishant776/errstack.Test_StackTracedErrorWithDefaults.func2.1;github.com/nnishant776/errstack.Test_StackTracedErrorWithDefaults.func2.2;github.com/nnishant776/errstack.Test_StackTracedErrorWithDefaults.func2.3;github.com/nnishant776/errstack.Test_StackTracedErrorWithDefaults.func2",
				},
				{
					format: " v",
					expect: "Hello Errors!\ngithub.com/nnishant776/errstack.Test_StackTracedErrorWithDefaults.func2.1\ngithub.com/nnishant776/errstack.Test_StackTracedErrorWithDefaults.func2.2\ngithub.com/nnishant776/errstack.Test_StackTracedErrorWithDefaults.func2.3\ngithub.com/nnishant776/errstack.Test_StackTracedErrorWithDefaults.func2",
				},
				{
					format: "-v",
					expect: fmt.Sprintf("Hello Errors!=>github.com/nnishant776/errstack.Test_StackTracedErrorWithDefaults.func2.1@%s/stack_test.go:117;github.com/nnishant776/errstack.Test_StackTracedErrorWithDefaults.func2.2@%s/stack_test.go:120;github.com/nnishant776/errstack.Test_StackTracedErrorWithDefaults.func2.3@%s/stack_test.go:123;github.com/nnishant776/errstack.Test_StackTracedErrorWithDefaults.func2@%s/stack_test.go:126", currDir, currDir, currDir, currDir),
				},
				{
					format: "+v",
					expect: fmt.Sprintf("Hello Errors!\ngithub.com/nnishant776/errstack.Test_StackTracedErrorWithDefaults.func2.1@%s/stack_test.go:117\ngithub.com/nnishant776/errstack.Test_StackTracedErrorWithDefaults.func2.2@%s/stack_test.go:120\ngithub.com/nnishant776/errstack.Test_StackTracedErrorWithDefaults.func2.3@%s/stack_test.go:123\ngithub.com/nnishant776/errstack.Test_StackTracedErrorWithDefaults.func2@%s/stack_test.go:126", currDir, currDir, currDir, currDir),
				},
				{
					format: "#v",
					expect: fmt.Sprintf("Hello Errors!\n#3: github.com/nnishant776/errstack.Test_StackTracedErrorWithDefaults.func2.1@%s/stack_test.go:117\n#2: github.com/nnishant776/errstack.Test_StackTracedErrorWithDefaults.func2.2@%s/stack_test.go:120\n#1: github.com/nnishant776/errstack.Test_StackTracedErrorWithDefaults.func2.3@%s/stack_test.go:123\n#0: github.com/nnishant776/errstack.Test_StackTracedErrorWithDefaults.func2@%s/stack_test.go:126", currDir, currDir, currDir, currDir),
				},
				{
					format: "j",
					expect: func() string {
						m := map[string]any{"error": err.Error(), "trace": err.StackTrace()}
						b, _ := json.Marshal(m)
						return string(b) + "\n"
					}(),
				},
				{
					format: "+j",
					expect: func() string {
						m := map[string]any{"error": err.Error(), "trace": err.StackTrace()}
						b, _ := json.MarshalIndent(m, "", "  ")
						return string(b) + "\n"
					}(),
				},
				{
					format: "+4j",
					expect: func() string {
						m := map[string]any{"error": err.Error(), "trace": err.StackTrace()}
						b, _ := json.MarshalIndent(m, "", "    ")
						return string(b) + "\n"
					}(),
				},
			}

			for _, c := range cases {
				t.Run(c.format, func(t *testing.T) {
					fmt.Fprintf(&c.output, "%"+c.format, err)
					if c.output.String() != c.expect {
						assert.Equal(
							t,
							c.expect,
							c.output.String(),
							"Doesn't match the expected output",
						)
					}
				})
			}
		})
	})
}

func Benchmark_StacktraceError(b *testing.B) {
	type ifce interface {
		StackTrace() pkgerrs.StackTrace
	}

	b.Run("pkgerrrs", func(b *testing.B) {
		b.Run("print error string only", func(b *testing.B) {
			b.Run("no stack capture", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					err := pkgerrs.New("pkgerrs")
					fmt.Fprintf(io.Discard, "%s", err)
				}
			})

			b.Run("with stack capture", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					err := pkgerrs.WithStack(pkgerrs.New("pkgerrs"))
					fmt.Fprintf(io.Discard, "%s", err)
				}
			})
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
			b.Run("no stack capture", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					err := NewString("errstk")
					fmt.Fprintf(io.Discard, "%s", err)
				}
			})

			b.Run("with stack capture", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					err := NewString("errstk", WithStack())
					fmt.Fprintf(io.Discard, "%s", err)
				}
			})
		})

		b.Run("print error and stacktrace", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				err := NewString("errstk", WithStack())
				fmt.Fprintf(io.Discard, "%-v", err)
			}
		})
	})
}
