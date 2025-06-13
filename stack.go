package errstack

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"strings"
)

var _ Error = (*StacktraceError)(nil)

type StacktraceError struct {
	err        error
	str        string
	stackTrace StackTrace
	pcList     [_MAX_CALL_DEPTH]uintptr
	opts       stackErrOpts
	frameCount int
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

	for i := range _MAX_CALL_DEPTH {
		stErr.pcList[i] = math.MaxUint64
	}

	for _, f := range opts {
		stErr.opts = f(stErr.opts)
	}

	if stErr.opts.autoStacktrace {
		stErr.frameCount = len(callersPCsBuf(3, _MAX_CALL_DEPTH, stErr.pcList[:]))
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
		stErr.frameCount = len(callersPCsBuf(3, _MAX_CALL_DEPTH, stErr.pcList[:]))
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

	if pc := callerPC(1); pc != math.MaxUint64 {
		if self.frameCount < _MAX_CALL_DEPTH {
			self.pcList[self.frameCount] = pc
			self.frameCount++
		}
		self.stackTrace = StackTrace{}
	}

	return self
}

func (self *StacktraceError) StackTrace() StackTrace {
	if self == nil {
		return StackTrace{}
	}

	if len(self.stackTrace.Frames) > 0 {
		return self.stackTrace
	}

	self.stackTrace.Frames = genStackTraceFromPCs(self.pcList[:self.frameCount])

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
		"stack": self.StackTrace(),
	})
}

func (self *StacktraceError) String() string {
	sb := strings.Builder{}

	stackTrace := self.StackTrace()
	sb.Grow(_MIN_STR_BYTES_PER_FRAME_STACKTRACE * len(stackTrace.Frames))

	GetErrorFormatter()(self, &sb)
	stackTrace.Print(StackFrameFormatOptions{}, &sb)

	return sb.String()
}

func (self *StacktraceError) Print(opts StackFrameFormatOptions, w io.Writer) {
	stackTrace := self.StackTrace()
	GetErrorFormatter()(self, w)
	stackTrace.Print(opts, w)
}

// Format formats the frame according to the fmt.Formatter interface.
//
//	%s	Plain error string, no stack trace
//
//	%n	Formatted error string, with just the function name in which error was generated
//
//	%v	Formatted error string, with the function name and the source location at which error was generated
//
//	%j	Plain error string, with the function name and the source location at which error was generated, except
//		it is printed as a json string
//
// Format accepts flags that alter the printing of some verbs, as follows:
//
//	%+s	Same as %s, except the specified error formatter will be used for printing the string
//
//	%+n	Same as %n, except it will print the stack trace with the name of each function in the call stack
//
//	%-n	Same as %+n, except it will not print the stack index
//
//	%+v	Same as %v, except it will print the stack trace with the name of each function and source location of
//		the function call in the call stack
//
//	%-v	Same as %+v, except it will not print the stack index
//
//	%+j	Same as %j, except it will print the stack trace with the name of each function and source location of
//		the function call in the call stack as a json string
//
// NOTE: Every verb except 's' and 'j' will always use the error and stack formatters defined in the package
func (self *StacktraceError) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		switch {
		case s.Flag('+'):
			io.WriteString(s, self.Error())
		default:
			GetErrorFormatter()(self, s)
		}

	case 'n':
		opts := StackFrameFormatOptions{
			SkipStackIndex:   false,
			SkipFunctionName: false,
			SkipLocation:     true,
		}
		stackTrace := self.StackTrace()
		switch {
		case s.Flag('+'):
		case s.Flag('-'):
			opts.SkipStackIndex = true
		default:
			stackTrace.Frames = stackTrace.Frames[:1]
		}
		GetErrorFormatter()(self, s)
		stackTrace.Print(opts, s)

	case 'v':
		opts := StackFrameFormatOptions{
			SkipStackIndex:   false,
			SkipFunctionName: false,
			SkipLocation:     false,
		}
		stackTrace := self.StackTrace()
		switch {
		case s.Flag('+'):
		case s.Flag('-'):
			opts.SkipStackIndex = true
		default:
			stackTrace.Frames = stackTrace.Frames[:1]
		}
		GetErrorFormatter()(self, s)
		stackTrace.Print(opts, s)

	case 'j':
		stackTrace := self.StackTrace()
		switch {
		case s.Flag('+'):
		case s.Flag('-'):
		default:
			stackTrace.Frames = stackTrace.Frames[:1]
		}
		data := map[string]any{
			"error": self.Error(),
			"stack": stackTrace,
		}
		json.NewEncoder(s).Encode(data)
	}
}
