package errstack

import (
	"encoding/json"
	"fmt"
	"strings"
)

var _ ChainedError = (*ChainedStacktraceError)(nil)

func NewChain(err error, opts ...StackErrOption) *ChainedStacktraceError {
	return newChainedStacktraceError(err, opts...)
}

func NewChainString(errStr string, opts ...StackErrOption) *ChainedStacktraceError {
	return newChainedStacktraceErrorString(errStr, opts...)
}

func newChainedStacktraceError(err error, opts ...StackErrOption) *ChainedStacktraceError {
	chainErr := &ChainedStacktraceError{}

	if stErr, ok := err.(Error); ok {
		chainErr.currErr = stErr
	} else {
		chainErr.currErr = New(err, opts...)
	}

	return chainErr
}

func newChainedStacktraceErrorString(errStr string, opts ...StackErrOption) *ChainedStacktraceError {
	return &ChainedStacktraceError{
		currErr: NewString(errStr, opts...),
	}
}

func Chain(err1, err2 error) ChainedError {
	chErr1, ok1 := err1.(ChainedError)
	if ok1 {
		return chErr1.Chain(err2)
	}

	return newChainedStacktraceError(err1).Chain(err2)
}

type ChainedStacktraceError struct {
	nextErr ChainedError
	currErr Error
}

func (self *ChainedStacktraceError) Chain(err error) ChainedError {
	if self == nil {
		return nil
	}

	chErr, ok := err.(ChainedError)
	if ok {
		self.nextErr = chErr
	} else {
		self.nextErr = newChainedStacktraceError(err)
	}

	return self
}

func (self *ChainedStacktraceError) Unwrap() []error {
	if self == nil {
		return nil
	}

	errList := ([]error)(nil)
	chainElem := (ChainedError)(self)

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

//go:noinline
func (self *ChainedStacktraceError) Throw() ChainedError {
	if self == nil {
		return nil
	}

	self.currErr.Throw(1)

	return self

}

func (self *ChainedStacktraceError) Error() string {
	return fmt.Sprintf("%s", self)
}

func (self *ChainedStacktraceError) Inner() Error {
	if self == nil {
		return nil
	}

	return self.currErr
}

func (self *ChainedStacktraceError) MarshalJSON() ([]byte, error) {
	if self == nil {
		return json.Marshal(nil)
	}

	errList := ([]Error)(nil)

	for elem := (ChainedError)(self); elem != nil; elem = elem.Next() {
		if elem.Inner() != nil {
			errList = append(errList, elem.Inner())
		}
	}

	return json.Marshal(errList)
}

func (self *ChainedStacktraceError) String() string {
	if self == nil {
		return NilErrorString
	}

	sb := strings.Builder{}
	chainElem := (ChainedError)(self)

	for {
		sb.WriteString(chainElem.Inner().String())
		chainElem = chainElem.Next()
		if chainElem == nil {
			break
		}
	}

	return sb.String()
}

// Format formats the frame according to the fmt.Formatter interface. Format also accepts
// flags that alter the printing of some verbs. The allowed combinations are as follows:
//
//	%s	Plain error string, no stack trace, separated by the separator configured in the
//		default formatter
//
//	%+s	Same as %s, except the specified error formatter will be used for printing the string.
//		The errors are separated by ": "
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
//		are kept intact. '+' can be followed by an arbitrary number which will represent the count
//		of spaces used to indent the stack trace
//
//	%#v	Same as %+(n)v, except it will print stack indices as well
//
//	%j	Same as %-v, except it will be printed as a json string
//
//	%+j	Same as %j, except it will be pretty printed. '+' can be followed by an arbitrary number
//		to indicate the indentation in the json output
//
// NOTE: Every verb defined above will always use the error and stack formatters defined in the package.
// It will only override the options mentioned as part of the flags and the rest will be used as is. The user
// is free to define other options of their choosing or provide entirely different implmentations as long as
// the interfaces are satisfied.
func (self *ChainedStacktraceError) Format(s fmt.State, verb rune) {
	erFmt := DefaultChainErrorFormatter
	stFmt := erFmt.StackTraceFormatter()
	ffFmt := stFmt.FrameFormatter()

	eOpts := erFmt.Options()
	fOpts := ffFmt.Options()
	sOpts := stFmt.Options()

	switch verb {
	case 's':
		if s.Flag('+') {
			eOpts.ErrorSeparator = ": "
			erFmt = erFmt.Copy().SetOptions(eOpts)
		}
		erFmt.FormatBuffer(s, self)

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

		eOpts.StackTraceSeparator = "=>"
		fOpts.SkipLocation = flags <= 1
		sOpts.SkipStackIndex = flags&(1<<3) == 0

		erFmt = erFmt.Copy()
		stFmt = stFmt.Copy()
		ffFmt = ffFmt.Copy()

		if flags&0x0d > 0 {
			eOpts.ErrorSeparator = "\n"
			eOpts.StackTraceSeparator = "\n"
			sOpts.FrameSeparator = "\n"
			if w, ok := s.Width(); ok {
				w = max(2, w)
				sOpts.FrameIndent = strings.Repeat(" ", w)
			}
		}

		ffFmt.SetOptions(fOpts)
		stFmt.SetOptions(sOpts).SetFrameFormatter(ffFmt)
		erFmt.SetOptions(eOpts).SetStackTraceFormatter(stFmt)
		erFmt.FormatBuffer(s, self)

	case 'j':
		enc := json.NewEncoder(s)
		if s.Flag('+') {
			w, _ := s.Width()
			enc.SetIndent("", strings.Repeat(" ", max(2, w)))
		}
		enc.Encode(self)
	}

}
