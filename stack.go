package errstack

import (
	"encoding/json"
	"strings"
)

var _ Error = (*StacktraceError)(nil)

type StacktraceError struct {
	err        error
	str        string
	stackTrace StackTrace
	opts       stackErrOpts
}

func New(err error, opts ...StackErrOption) *StacktraceError {
	return newStacktraceError(err, opts...)
}

func NewString(errStr string, opts ...StackErrOption) *StacktraceError {
	return newStacktraceErrorString(errStr, opts...)
}

func newStacktraceError(err error, opts ...StackErrOption) *StacktraceError {
	stErr := &StacktraceError{
		err: err,
	}

	for _, f := range opts {
		stErr.opts = f(stErr.opts)
	}

	if stErr.opts.autoStacktrace {
		stErr.stackTrace.Frames = callers(3, _MAX_CALL_DEPTH)
	}

	return stErr
}

func newStacktraceErrorString(errStr string, opts ...StackErrOption) *StacktraceError {
	stErr := &StacktraceError{
		str: errStr,
	}

	for _, f := range opts {
		stErr.opts = f(stErr.opts)
	}

	if stErr.opts.autoStacktrace {
		stErr.stackTrace.Frames = callers(3, _MAX_CALL_DEPTH)
	}

	return stErr
}

func (self *StacktraceError) Error() string {
	if self == nil {
		return NilErrorString
	}

	if self.err != nil {
		return self.err.Error()
	}

	return self.str
}

func (self *StacktraceError) Throw() Error {
	if self == nil {
		return nil
	}

	if self.opts.autoStacktrace {
		return self
	}

	self.stackTrace.Frames = append(self.stackTrace.Frames, caller(1))

	return self
}

func (self *StacktraceError) StackTrace() StackTrace {
	if self == nil {
		return StackTrace{}
	}

	return self.stackTrace
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
		"error": self.Error(),
		"stack": self.stackTrace.Frames,
	})
}

func (self *StacktraceError) String() string {
	var sb strings.Builder

	formatter := defaultErrorFormatter
	if ErrorFormatter != nil {
		formatter = ErrorFormatter
	}

	sb.WriteString(formatter(self))

	stFormatter := defaultStackTraceFormatter
	if StackTraceFormatter != nil {
		stFormatter = StackTraceFormatter
	}

	sb.WriteString(stFormatter(self.StackTrace()))

	return sb.String()
}
