package errstack

import (
	"encoding/json"
	"strings"
)

var _ Error = (*StacktraceError)(nil)

type StacktraceError struct {
	err       error
	backtrace Backtrace
}

func New(err error) *StacktraceError {
	return newStacktraceError(err)
}

func newStacktraceError(err error) *StacktraceError {
	return &StacktraceError{
		err: err,
	}
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
