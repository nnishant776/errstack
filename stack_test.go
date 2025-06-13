package errstack

import (
	"errors"
	"testing"
)

func Test_stackerr(t *testing.T) {
	err := stackabc1()
	t.Logf("%-v\n", err)
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
