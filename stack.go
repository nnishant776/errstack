package errstack

import (
	"encoding/json"
	"strings"
)

var _ Error = (*StacktraceError)(nil)

type StacktraceError struct {
	err       error
	backtrace Backtrace
	opts      stackErrOpts
}

func New(err error, opts ...StackErrOption) *StacktraceError {
	return newStacktraceError(err, opts...)
}

func newStacktraceError(err error, opts ...StackErrOption) *StacktraceError {
	stErr := &StacktraceError{
		err: err,
	}

	for _, opt := range opts {
		opt.apply(&stErr.opts)
	}

	if stErr.opts.autoStacktrace {
		stErr.backtrace.Frames = callers(3, 0)
	}

	return stErr
}

func (self *StacktraceError) Error() string {
	if self == nil || self.err == nil {
		return NilErrorString
	}

	return self.err.Error()
}

func (self *StacktraceError) Throw() Error {
	if self == nil {
		return nil
	}

	if self.opts.autoStacktrace {
		return self
	}

	self.backtrace.Frames = append(self.backtrace.Frames, caller(1))

	return self
}

func (self *StacktraceError) Backtrace() Backtrace {
	if self == nil || self.err == nil {
		return Backtrace{}
	}

	return self.backtrace
}

func (self *StacktraceError) Unwrap() error {
	if self == nil {
		return nil
	}

	return self.err
}

func (self *StacktraceError) MarshalJSON() ([]byte, error) {
	if self == nil {
		return json.Marshal(nil)
	}

	return json.Marshal(map[string]any{
		"error": self.err.Error(),
		"stack": self.backtrace.Frames,
	})
}

func (self *StacktraceError) String() string {
	var sb strings.Builder

	formatter := defaultErrorFormatter
	if ErrorFormatter != nil {
		formatter = ErrorFormatter
	}

	sb.WriteString(formatter(self))

	btFormatter := defaultBacktraceFormatter
	if BacktraceFormatter != nil {
		btFormatter = BacktraceFormatter
	}

	sb.WriteString(btFormatter(self.Backtrace()))

	return sb.String()
}
