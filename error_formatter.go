package errstack

import (
	"io"
	"strings"
)

type ErrorFormatter interface {
	Options() ErrorFormatterOptions
	StackTraceFormatter() StackTraceFormatter
	Format(e error) string
	FormatBuffer(w io.Writer, e error)
	Clone() ErrorFormatter
	Copy() ErrorFormatter
	SetOptions(opts ErrorFormatterOptions) ErrorFormatter
	SetStackTraceFormatter(stFmt StackTraceFormatter) ErrorFormatter
}

type ErrorFormatterOptions struct {
	ErrorPrefix         string
	ErrorSeparator      string
	StackTraceSeparator string
}

var _ ErrorFormatter = (*errorFormatter)(nil)

type errorFormatter struct {
	opts  ErrorFormatterOptions
	stFmt StackTraceFormatter
}

func (self *errorFormatter) format(w io.Writer, err error) {
	if self == nil {
		w.Write([]byte(NilErrorString))
		return
	}

	prefix, errStr := self.opts.ErrorPrefix, err.Error()

	switch o := w.(type) {
	case io.StringWriter:
		o.WriteString(prefix)
		o.WriteString(errStr)
	default:
		w.Write(string2Slice(prefix))
		w.Write(string2Slice(errStr))
	}

	switch {
	case self.stFmt == nil, self.opts.StackTraceSeparator == "":
	default:
		switch o := w.(type) {
		case io.StringWriter:
			o.WriteString(self.opts.StackTraceSeparator)
		default:
			w.Write(string2Slice(self.opts.StackTraceSeparator))
		}

		if stErr, ok := err.(StackTracer); ok {
			stackTrace := stErr.StackTrace()
			if len(stackTrace.Frames) > 0 {
				self.stFmt.FormatBuffer(w, stackTrace)
			}
		}
	}
}

func (self *errorFormatter) Options() ErrorFormatterOptions {
	return self.opts
}

func (self *errorFormatter) StackTraceFormatter() StackTraceFormatter {
	return self.stFmt
}

func (self *errorFormatter) Format(e error) string {
	sb := strings.Builder{}
	self.format(&sb, e)
	return sb.String()
}

func (self *errorFormatter) FormatBuffer(w io.Writer, e error) {
	self.format(w, e)
}

func (self *errorFormatter) Clone() ErrorFormatter {
	return &errorFormatter{
		opts:  self.opts,
		stFmt: self.stFmt.Clone(),
	}
}

func (self *errorFormatter) Copy() ErrorFormatter {
	return &errorFormatter{
		opts:  self.opts,
		stFmt: self.stFmt,
	}
}

func (self *errorFormatter) SetOptions(opts ErrorFormatterOptions) ErrorFormatter {
	self.opts = opts
	return self
}

func (self *errorFormatter) SetStackTraceFormatter(stFmt StackTraceFormatter) ErrorFormatter {
	self.stFmt = stFmt
	return self
}
