package errstack

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ChainedErrorWithDefaults(t *testing.T) {
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
					expect: fmt.Sprintf("Error 1=>github.com/nnishant776/errstack.Test_ChainedErrorWithDefaults.func1.3@%s/chain_test.go:28;github.com/nnishant776/errstack.Test_ChainedErrorWithDefaults.func1@%s/chain_test.go:31;testing.tRunner@/usr/local/go/src/testing/testing.go:1690;runtime.goexit@/usr/local/go/src/runtime/asm_amd64.s:1700, Error 2=>github.com/nnishant776/errstack.Test_ChainedErrorWithDefaults.func1.1@%s/chain_test.go:22;github.com/nnishant776/errstack.Test_ChainedErrorWithDefaults.func1.2@%s/chain_test.go:25;github.com/nnishant776/errstack.Test_ChainedErrorWithDefaults.func1.3@%s/chain_test.go:28;github.com/nnishant776/errstack.Test_ChainedErrorWithDefaults.func1@%s/chain_test.go:31;testing.tRunner@/usr/local/go/src/testing/testing.go:1690;runtime.goexit@/usr/local/go/src/runtime/asm_amd64.s:1700", currDir, currDir, currDir, currDir, currDir, currDir),
				},
				{
					format: "+v",
					expect: fmt.Sprintf("Error 1\ngithub.com/nnishant776/errstack.Test_ChainedErrorWithDefaults.func1.3@%s/chain_test.go:28\ngithub.com/nnishant776/errstack.Test_ChainedErrorWithDefaults.func1@%s/chain_test.go:31\ntesting.tRunner@/usr/local/go/src/testing/testing.go:1690\nruntime.goexit@/usr/local/go/src/runtime/asm_amd64.s:1700\nError 2\ngithub.com/nnishant776/errstack.Test_ChainedErrorWithDefaults.func1.1@%s/chain_test.go:22\ngithub.com/nnishant776/errstack.Test_ChainedErrorWithDefaults.func1.2@%s/chain_test.go:25\ngithub.com/nnishant776/errstack.Test_ChainedErrorWithDefaults.func1.3@%s/chain_test.go:28\ngithub.com/nnishant776/errstack.Test_ChainedErrorWithDefaults.func1@%s/chain_test.go:31\ntesting.tRunner@/usr/local/go/src/testing/testing.go:1690\nruntime.goexit@/usr/local/go/src/runtime/asm_amd64.s:1700", currDir, currDir, currDir, currDir, currDir, currDir),
				},
				{
					format: "#v",
					expect: fmt.Sprintf("Error 1\n#3: github.com/nnishant776/errstack.Test_ChainedErrorWithDefaults.func1.3@%s/chain_test.go:28\n#2: github.com/nnishant776/errstack.Test_ChainedErrorWithDefaults.func1@%s/chain_test.go:31\n#1: testing.tRunner@/usr/local/go/src/testing/testing.go:1690\n#0: runtime.goexit@/usr/local/go/src/runtime/asm_amd64.s:1700\nError 2\n#5: github.com/nnishant776/errstack.Test_ChainedErrorWithDefaults.func1.1@%s/chain_test.go:22\n#4: github.com/nnishant776/errstack.Test_ChainedErrorWithDefaults.func1.2@%s/chain_test.go:25\n#3: github.com/nnishant776/errstack.Test_ChainedErrorWithDefaults.func1.3@%s/chain_test.go:28\n#2: github.com/nnishant776/errstack.Test_ChainedErrorWithDefaults.func1@%s/chain_test.go:31\n#1: testing.tRunner@/usr/local/go/src/testing/testing.go:1690\n#0: runtime.goexit@/usr/local/go/src/runtime/asm_amd64.s:1700", currDir, currDir, currDir, currDir, currDir, currDir),
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
					expect: fmt.Sprintf("Error 1=>github.com/nnishant776/errstack.Test_ChainedErrorWithDefaults.func2.3@%s/chain_test.go:128;github.com/nnishant776/errstack.Test_ChainedErrorWithDefaults.func2@%s/chain_test.go:131, Error 2=>github.com/nnishant776/errstack.Test_ChainedErrorWithDefaults.func2.1@%s/chain_test.go:122;github.com/nnishant776/errstack.Test_ChainedErrorWithDefaults.func2.2@%s/chain_test.go:125", currDir, currDir, currDir, currDir),
				},
				{
					format: "+v",
					expect: fmt.Sprintf("Error 1\ngithub.com/nnishant776/errstack.Test_ChainedErrorWithDefaults.func2.3@%s/chain_test.go:128\ngithub.com/nnishant776/errstack.Test_ChainedErrorWithDefaults.func2@%s/chain_test.go:131\nError 2\ngithub.com/nnishant776/errstack.Test_ChainedErrorWithDefaults.func2.1@%s/chain_test.go:122\ngithub.com/nnishant776/errstack.Test_ChainedErrorWithDefaults.func2.2@%s/chain_test.go:125", currDir, currDir, currDir, currDir),
				},
				{
					format: "#v",
					expect: fmt.Sprintf("Error 1\n#1: github.com/nnishant776/errstack.Test_ChainedErrorWithDefaults.func2.3@%s/chain_test.go:128\n#0: github.com/nnishant776/errstack.Test_ChainedErrorWithDefaults.func2@%s/chain_test.go:131\nError 2\n#1: github.com/nnishant776/errstack.Test_ChainedErrorWithDefaults.func2.1@%s/chain_test.go:122\n#0: github.com/nnishant776/errstack.Test_ChainedErrorWithDefaults.func2.2@%s/chain_test.go:125", currDir, currDir, currDir, currDir),
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
