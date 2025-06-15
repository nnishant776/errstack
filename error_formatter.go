package errstack

import (
	"io"
	"strings"
)

type ErrorFormatter interface {
	WithOptions(opts ErrorFormatterOptions) ErrorFormatter
	Options() ErrorFormatterOptions
	Format(e error) string
	FormatBuffer(w io.Writer, e error)
}

type ErrorFormatterOptions struct {
	ErrorPrefix    string
	ErrorSeparator string
}

var _ ErrorFormatter = (*errorFormatter)(nil)

type errorFormatter struct {
	opts ErrorFormatterOptions
}

func (self *errorFormatter) format(w io.Writer, err error) {
	prefix, errStr := self.opts.ErrorPrefix, err.Error()

	switch o := w.(type) {
	case io.StringWriter:
		o.WriteString(prefix)
		o.WriteString(errStr)
		o.WriteString(self.opts.ErrorSeparator)
	default:
		w.Write(string2Slice(prefix))
		w.Write(string2Slice(errStr))
		w.Write(string2Slice(self.opts.ErrorSeparator))
	}
}

func (self *errorFormatter) WithOptions(opts ErrorFormatterOptions) ErrorFormatter {
	ef := &errorFormatter{
		opts: opts,
	}

	return ef
}

func (self *errorFormatter) Options() ErrorFormatterOptions {
	return self.opts
}

func (self *errorFormatter) Format(e error) string {
	sb := strings.Builder{}
	self.format(&sb, e)
	return sb.String()
}

func (self *errorFormatter) FormatBuffer(w io.Writer, e error) {
	self.format(w, e)
}
