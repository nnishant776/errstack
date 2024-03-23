package errstack

import (
	"errors"
	"testing"
)

func Test_chainerr(t *testing.T) {
	err := chainabc1()
	t.Logf("\n%v", err)
}

func chainabc1() ChainedError {
	return Chain(errors.New("Error 1"), chainabc2())
}

func chainabc2() ChainedError {
	return chainabc3().Throw()
}

func chainabc3() ChainedError {
	return NewChain(errors.New("Error 2")).Throw()
}
