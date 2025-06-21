package errstack

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ChainedErrorWithDefaults(t *testing.T) {
	goRoot := runtime.GOROOT()
	currDir, err := os.Getwd()
	if err != nil {
		t.Errorf("Failed to get current directory: error: %s", err)
	}

	t.Run("automatic stacktrace collection", func(t *testing.T) {
		chainabc3 := func() ChainedError {
			return NewChain(NewString("Error 2", WithStack()))
		}
		chainabc2 := func() ChainedError {
			return chainabc3()
		}
		chainabc1 := func() ChainedError {
			return Chain(NewString("Error 1", WithStack()), chainabc2())
		}

		err := chainabc1()

		t.Run("printing formats", func(t *testing.T) {
			type testCase struct {
				format string
				output strings.Builder
				expect string
			}

			cases := []testCase{
				{
					format: "s",
					expect: "Error 1, Error 2",
				},
				{
					format: "+s",
					expect: "Error 1: Error 2",
				},
				{
					format: "v",
					expect: "Error 1=>github.com/nnishant776/errstack.Test_ChainedErrorWithDefaults.func1.3;github.com/nnishant776/errstack.Test_ChainedErrorWithDefaults.func1;testing.tRunner;runtime.goexit, Error 2=>github.com/nnishant776/errstack.Test_ChainedErrorWithDefaults.func1.1;github.com/nnishant776/errstack.Test_ChainedErrorWithDefaults.func1.2;github.com/nnishant776/errstack.Test_ChainedErrorWithDefaults.func1.3;github.com/nnishant776/errstack.Test_ChainedErrorWithDefaults.func1;testing.tRunner;runtime.goexit",
				},
				{
					format: " v",
					expect: "Error 1\ngithub.com/nnishant776/errstack.Test_ChainedErrorWithDefaults.func1.3\ngithub.com/nnishant776/errstack.Test_ChainedErrorWithDefaults.func1\ntesting.tRunner\nruntime.goexit\nError 2\ngithub.com/nnishant776/errstack.Test_ChainedErrorWithDefaults.func1.1\ngithub.com/nnishant776/errstack.Test_ChainedErrorWithDefaults.func1.2\ngithub.com/nnishant776/errstack.Test_ChainedErrorWithDefaults.func1.3\ngithub.com/nnishant776/errstack.Test_ChainedErrorWithDefaults.func1\ntesting.tRunner\nruntime.goexit",
				},
				{
					format: "-v",
					expect: fmt.Sprintf("Error 1=>github.com/nnishant776/errstack.Test_ChainedErrorWithDefaults.func1.3@%s/chain_test.go:30;github.com/nnishant776/errstack.Test_ChainedErrorWithDefaults.func1@%s/chain_test.go:33;testing.tRunner@%s/src/testing/testing.go:1690;runtime.goexit@%s/src/runtime/asm_amd64.s:1700, Error 2=>github.com/nnishant776/errstack.Test_ChainedErrorWithDefaults.func1.1@%s/chain_test.go:24;github.com/nnishant776/errstack.Test_ChainedErrorWithDefaults.func1.2@%s/chain_test.go:27;github.com/nnishant776/errstack.Test_ChainedErrorWithDefaults.func1.3@%s/chain_test.go:30;github.com/nnishant776/errstack.Test_ChainedErrorWithDefaults.func1@%s/chain_test.go:33;testing.tRunner@%s/src/testing/testing.go:1690;runtime.goexit@%s/src/runtime/asm_amd64.s:1700", currDir, currDir, goRoot, goRoot, currDir, currDir, currDir, currDir, goRoot, goRoot),
				},
				{
					format: "+v",
					expect: fmt.Sprintf("Error 1\ngithub.com/nnishant776/errstack.Test_ChainedErrorWithDefaults.func1.3@%s/chain_test.go:30\ngithub.com/nnishant776/errstack.Test_ChainedErrorWithDefaults.func1@%s/chain_test.go:33\ntesting.tRunner@%s/src/testing/testing.go:1690\nruntime.goexit@%s/src/runtime/asm_amd64.s:1700\nError 2\ngithub.com/nnishant776/errstack.Test_ChainedErrorWithDefaults.func1.1@%s/chain_test.go:24\ngithub.com/nnishant776/errstack.Test_ChainedErrorWithDefaults.func1.2@%s/chain_test.go:27\ngithub.com/nnishant776/errstack.Test_ChainedErrorWithDefaults.func1.3@%s/chain_test.go:30\ngithub.com/nnishant776/errstack.Test_ChainedErrorWithDefaults.func1@%s/chain_test.go:33\ntesting.tRunner@%s/src/testing/testing.go:1690\nruntime.goexit@%s/src/runtime/asm_amd64.s:1700", currDir, currDir, goRoot, goRoot, currDir, currDir, currDir, currDir, goRoot, goRoot),
				},
				{
					format: "#v",
					expect: fmt.Sprintf("Error 1\n#3: github.com/nnishant776/errstack.Test_ChainedErrorWithDefaults.func1.3@%s/chain_test.go:30\n#2: github.com/nnishant776/errstack.Test_ChainedErrorWithDefaults.func1@%s/chain_test.go:33\n#1: testing.tRunner@%s/src/testing/testing.go:1690\n#0: runtime.goexit@%s/src/runtime/asm_amd64.s:1700\nError 2\n#5: github.com/nnishant776/errstack.Test_ChainedErrorWithDefaults.func1.1@%s/chain_test.go:24\n#4: github.com/nnishant776/errstack.Test_ChainedErrorWithDefaults.func1.2@%s/chain_test.go:27\n#3: github.com/nnishant776/errstack.Test_ChainedErrorWithDefaults.func1.3@%s/chain_test.go:30\n#2: github.com/nnishant776/errstack.Test_ChainedErrorWithDefaults.func1@%s/chain_test.go:33\n#1: testing.tRunner@%s/src/testing/testing.go:1690\n#0: runtime.goexit@%s/src/runtime/asm_amd64.s:1700", currDir, currDir, goRoot, goRoot, currDir, currDir, currDir, currDir, goRoot, goRoot),
				},
				{
					format: "j",
					expect: func() string {
						errList := []Error{}
						for chErr := err; chErr != nil; chErr = chErr.Next() {
							errList = append(errList, chErr.Inner())
						}
						b, _ := json.Marshal(errList)
						return string(b) + "\n"
					}(),
				},
				{
					format: "+j",
					expect: func() string {
						errList := []Error{}
						for chErr := err; chErr != nil; chErr = chErr.Next() {
							errList = append(errList, chErr.Inner())
						}
						b, _ := json.MarshalIndent(errList, "", "  ")
						return string(b) + "\n"
					}(),
				},
				{
					format: "+4j",
					expect: func() string {
						errList := []Error{}
						for chErr := err; chErr != nil; chErr = chErr.Next() {
							errList = append(errList, chErr.Inner())
						}
						b, _ := json.MarshalIndent(errList, "", "    ")
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
		chainabc3 := func() ChainedError {
			return NewChain(errors.New("Error 2")).Throw()
		}
		chainabc2 := func() ChainedError {
			return chainabc3().Throw()
		}
		chainabc1 := func() ChainedError {
			return Chain(NewString("Error 1"), chainabc2()).Throw()
		}

		err := chainabc1().Throw()

		t.Run("printing formats", func(t *testing.T) {
			type testCase struct {
				format string
				output strings.Builder
				expect string
			}

			cases := []testCase{
				{
					format: "s",
					expect: "Error 1, Error 2",
				},
				{
					format: "+s",
					expect: "Error 1: Error 2",
				},
				{
					format: "v",
					expect: "Error 1=>github.com/nnishant776/errstack.Test_ChainedErrorWithDefaults.func2.3;github.com/nnishant776/errstack.Test_ChainedErrorWithDefaults.func2, Error 2=>github.com/nnishant776/errstack.Test_ChainedErrorWithDefaults.func2.1;github.com/nnishant776/errstack.Test_ChainedErrorWithDefaults.func2.2",
				},
				{
					format: " v",
					expect: "Error 1\ngithub.com/nnishant776/errstack.Test_ChainedErrorWithDefaults.func2.3\ngithub.com/nnishant776/errstack.Test_ChainedErrorWithDefaults.func2\nError 2\ngithub.com/nnishant776/errstack.Test_ChainedErrorWithDefaults.func2.1\ngithub.com/nnishant776/errstack.Test_ChainedErrorWithDefaults.func2.2",
				},
				{
					format: "-v",
					expect: fmt.Sprintf("Error 1=>github.com/nnishant776/errstack.Test_ChainedErrorWithDefaults.func2.3@%s/chain_test.go:130;github.com/nnishant776/errstack.Test_ChainedErrorWithDefaults.func2@%s/chain_test.go:133, Error 2=>github.com/nnishant776/errstack.Test_ChainedErrorWithDefaults.func2.1@%s/chain_test.go:124;github.com/nnishant776/errstack.Test_ChainedErrorWithDefaults.func2.2@%s/chain_test.go:127", currDir, currDir, currDir, currDir),
				},
				{
					format: "+v",
					expect: fmt.Sprintf("Error 1\ngithub.com/nnishant776/errstack.Test_ChainedErrorWithDefaults.func2.3@%s/chain_test.go:130\ngithub.com/nnishant776/errstack.Test_ChainedErrorWithDefaults.func2@%s/chain_test.go:133\nError 2\ngithub.com/nnishant776/errstack.Test_ChainedErrorWithDefaults.func2.1@%s/chain_test.go:124\ngithub.com/nnishant776/errstack.Test_ChainedErrorWithDefaults.func2.2@%s/chain_test.go:127", currDir, currDir, currDir, currDir),
				},
				{
					format: "#v",
					expect: fmt.Sprintf("Error 1\n#1: github.com/nnishant776/errstack.Test_ChainedErrorWithDefaults.func2.3@%s/chain_test.go:130\n#0: github.com/nnishant776/errstack.Test_ChainedErrorWithDefaults.func2@%s/chain_test.go:133\nError 2\n#1: github.com/nnishant776/errstack.Test_ChainedErrorWithDefaults.func2.1@%s/chain_test.go:124\n#0: github.com/nnishant776/errstack.Test_ChainedErrorWithDefaults.func2.2@%s/chain_test.go:127", currDir, currDir, currDir, currDir),
				},
				{
					format: "j",
					expect: func() string {
						errList := []Error{}
						for chErr := err; chErr != nil; chErr = chErr.Next() {
							errList = append(errList, chErr.Inner())
						}
						b, _ := json.Marshal(errList)
						return string(b) + "\n"
					}(),
				},
				{
					format: "+j",
					expect: func() string {
						errList := []Error{}
						for chErr := err; chErr != nil; chErr = chErr.Next() {
							errList = append(errList, chErr.Inner())
						}
						b, _ := json.MarshalIndent(errList, "", "  ")
						return string(b) + "\n"
					}(),
				},
				{
					format: "+4j",
					expect: func() string {
						errList := []Error{}
						for chErr := err; chErr != nil; chErr = chErr.Next() {
							errList = append(errList, chErr.Inner())
						}
						b, _ := json.MarshalIndent(errList, "", "    ")
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
