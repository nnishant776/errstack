package errstack

import (
	"encoding/json"
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
