package errstack

import (
	"encoding/json"
	"strings"
)

type ChainedError interface {
	error
	Chain(err error) ChainedError
	Inner() Error
	Next() ChainedError
	String() string
	Throw() ChainedError
	Unwrap() []error
}

var _ ChainedError = (*ChainedStacktraceError)(nil)

func NewChain(err error) *ChainedStacktraceError {
	return newChainedStacktraceError(err)
}

func newChainedStacktraceError(err error) *ChainedStacktraceError {
	return &ChainedStacktraceError{
		currErr: StacktraceError{
			err: err,
		},
	}
}

func Chain(err1, err2 error) ChainedError {
	chErr1, ok1 := err1.(ChainedError)
	if ok1 {
		return chErr1.Chain(err2)
	}

	var currErr = &ChainedStacktraceError{
		currErr: StacktraceError{
			err: err1,
		},
	}

	if btErr, ok := err2.(interface{ Backtrace() Backtrace }); ok {
		currErr.currErr.backtrace = btErr.Backtrace()
	}

	return currErr.Chain(err2)
}

type ChainedStacktraceError struct {
	nextErr ChainedError
	currErr StacktraceError
}

func (self *ChainedStacktraceError) Chain(err error) ChainedError {
	chErr, ok := err.(ChainedError)
	if ok {
		self.nextErr = chErr
	} else {
		var nextErr = &ChainedStacktraceError{
			currErr: StacktraceError{
				err: err,
			},
		}
		if btErr, ok := err.(Backtracer); ok {
			nextErr.currErr.backtrace = btErr.Backtrace()
		}
		self.nextErr = nextErr
	}

	return self
}

func (self *ChainedStacktraceError) Unwrap() []error {
	if self == nil {
		return nil
	}

	var errList []error
	var chainElem ChainedError = self

	for {
		errList = append(errList, chainElem.Inner())
		chainElem = chainElem.Next()
		if chainElem == nil {
			break
		}
	}

	return errList
}

func (self *ChainedStacktraceError) Next() ChainedError {
	if self == nil {
		return nil
	}

	return self.nextErr
}

func (self *ChainedStacktraceError) Throw() ChainedError {
	if self == nil {
		return nil
	}

	var frames = &self.currErr.backtrace.Frames
	*frames = append(*frames, caller(1))

	return self

}

func (self *ChainedStacktraceError) Error() string {
	if self == nil {
		return NilErrorString
	}

	var sb strings.Builder
	var chainElem ChainedError = self

	for {
		sb.WriteString(chainElem.Inner().Error())
		chainElem = chainElem.Next()
		if chainElem == nil {
			break
		}
		sb.WriteString(ErrorChainSeparator)
	}

	return sb.String()
}

func (self *ChainedStacktraceError) Inner() Error {
	return &self.currErr
}

func (self *ChainedStacktraceError) MarshalJSON() ([]byte, error) {
	if self == nil {
		return json.Marshal(nil)
	}

	var chainElem ChainedError = self
	var errList []Error

	for {
		errList = append(errList, chainElem.Inner())
		chainElem = chainElem.Next()
		if chainElem == nil {
			break
		}
	}

	return json.Marshal(errList)
}

func (self *ChainedStacktraceError) String() string {
	if self == nil {
		return NilErrorString
	}

	var sb strings.Builder
	var chainElem ChainedError = self

	for {
		sb.WriteString(chainElem.Inner().String())
		chainElem = chainElem.Next()
		if chainElem == nil {
			break
		}
	}

	return sb.String()
}
