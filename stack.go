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

func (self *StacktraceError) StackTraceN(n int) StackTrace {
	if self == nil {
		return StackTrace{}
	}

	n = min(n, self.frameCount)

	self.stackTrace.Frames = genStackTraceFromPCs(self.pcList[:n])

	return self.stackTrace
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
	DefaultErrorFormatter.FormatBuffer(&sb, self)
	DefaultStackTraceFormatter.FormatBuffer(&sb, stackTrace)
	return sb.String()
}

// Format formats the frame according to the fmt.Formatter interface. Format also accepts
// flags that alter the printing of some verbs. The allowed combinations are as follows:
//
//	%s	Plain error string, no stack trace
//
//	%+s	Same as %s, except the specified error formatter will be used for printing the string
//
//	%v	Formatted error string, with just the function name in the stack trace.
//
//	% v	Same as %v, except the stack trace is printed on separate lines. This overrides the
//		stack frame separator and the error separator to '\n'. Rest of the options are kept intact
//
//	%-v	Same as %v, except it will print the source location as well
//
//	%+v	Same as %-v, except it will print the stack trace on a separate line than the error. This
//		overrides the stack frame separator and the error separator to '\n'. Rest of the options
//		are kept intact
//
//	%#v	Same as %+v, except it will print stack indices as well
//
//	%j	Same as %-v, except it will be printed as a json string
//
//	%+j	Same as %-j, except it will be pretty printed. '+' can be followed by an arbitrary number
//		to indicate the indentation in the json output
//
// NOTE: Every verb defined above will always use the error and stack formatters defined in the package.
// It will only override the options mentioned as part of the flags and the rest will be used as is. The user
// is free to define other options of their choosing or provide entirely different implmentations as long as
// the interfaces are satisfied.
func (self *StacktraceError) Format(s fmt.State, verb rune) {
	erFmt := DefaultErrorFormatter
	ffFmt := DefaultStackFrameFormatter
	stFmt := DefaultStackTraceFormatter

	eOpts := erFmt.Options()
	fOpts := ffFmt.Options()
	sOpts := stFmt.Options()

	stackTrace := StackTrace{}
	switch verb {
	case 's':
		switch {
		case s.Flag('+'):
			eOpts.ErrorSeparator = ""
			erFmt.WithOptions(eOpts).FormatBuffer(s, self)
		default:
			errStr := self.Error()
			switch o := s.(type) {
			case io.StringWriter:
				o.WriteString(errStr)
			default:
				s.Write(string2Slice(errStr))
			}
		}

	case 'v':
		flags := byte(0)
		switch {
		case s.Flag(' '):
			flags = flags | 1<<0
		case s.Flag('-'):
			flags = flags | 1<<1
		case s.Flag('+'):
			flags = flags | 1<<2
		case s.Flag('#'):
			flags = flags | 1<<3
		default:
		}

		fOpts.SkipLocation = flags <= 1
		sOpts.SkipStackIndex = flags&(1<<3) == 0
		if flags&0x0d > 0 {
			eOpts.ErrorSeparator = "\n"
			sOpts.FrameSeparator = "\n"
		}

		stackTrace = self.StackTrace()

		erFmt.WithOptions(eOpts).FormatBuffer(s, self)
		stFmt.WithOptions(sOpts).
			WithFrameFormatter(ffFmt.WithOptions(fOpts)).
			FormatBuffer(s, stackTrace)

	case 'j':
		indentEnabled := s.Flag('+')
		stackTrace = self.StackTrace()
		data := map[string]any{
			"error": self.Error(),
			"trace": stackTrace,
		}

		if indentEnabled {
			w, _ := s.Width()
			w = max(2, w)
			enc := json.NewEncoder(s)
			enc.SetIndent("", strings.Repeat(" ", w))
			enc.Encode(data)
		} else {
			json.NewEncoder(s).Encode(data)
		}
	}
}
