package errstack

import (
	"io"
	"strings"
)

var _ ErrorFormatter = (*chainErrorFormatter)(nil)

type chainErrorFormatter struct {
	opts ErrorFormatterOptions
	sfmt StackTraceFormatter
}

func (self *chainErrorFormatter) format(w io.Writer, err error) {
	if err == nil {
		w.Write(string2Slice(NilErrorString))
		return
	}

	for err != nil {
		prefix, errStr := self.opts.ErrorPrefix, ""

		if chErr, ok := err.(ChainedError); ok {
			errStr = chErr.Inner().Error()
		} else {
			errStr = err.Error()
		}

		switch o := w.(type) {
		case io.StringWriter:
			o.WriteString(prefix)
			o.WriteString(errStr)
		default:
			w.Write(string2Slice(prefix))
			w.Write(string2Slice(errStr))
		}

		switch {
		case self.sfmt == nil, self.opts.StackTraceSeparator == "":
		default:
			stackTrace := StackTrace{}
			switch stErr := err.(type) {
			case ChainedError:
				stackTrace = stErr.Inner().StackTrace()
			case StackTracer:
				stackTrace = stErr.StackTrace()
			}

			if len(stackTrace.Frames) > 0 {
				switch o := w.(type) {
				case io.StringWriter:
					o.WriteString(self.opts.StackTraceSeparator)
				default:
					w.Write(string2Slice(self.opts.StackTraceSeparator))
				}

				self.sfmt.FormatBuffer(w, stackTrace)
			}
		}

		if chErr, ok := err.(ChainedError); !ok || chErr.Next() == nil {
			err = nil
		} else {
			switch o := w.(type) {
			case io.StringWriter:
				o.WriteString(self.opts.ErrorSeparator)
			default:
				w.Write(string2Slice(self.opts.ErrorSeparator))
			}

			err = chErr.Next()
		}
	}
}

func (self *chainErrorFormatter) WithOptions(opts ErrorFormatterOptions) ErrorFormatter {
	ef := &chainErrorFormatter{
		opts: opts,
		sfmt: self.sfmt,
	}

	return ef
}

func (self *chainErrorFormatter) WithStackTraceFormatter(stFmt StackTraceFormatter) ErrorFormatter {
	return &errorFormatter{
		opts:  self.opts,
		stFmt: stFmt,
	}
}

func (self *chainErrorFormatter) Options() ErrorFormatterOptions {
	return self.opts
}

func (self *chainErrorFormatter) StackTraceFormatter() StackTraceFormatter {
	return self.sfmt
}

func (self *chainErrorFormatter) Format(e error) string {
	sb := strings.Builder{}
	self.format(&sb, e)
	return sb.String()
}

func (self *chainErrorFormatter) FormatBuffer(w io.Writer, e error) {
	self.format(w, e)
}

func (self *chainErrorFormatter) Clone() ErrorFormatter {
	return &chainErrorFormatter{
		opts: self.opts,
		sfmt: self.sfmt.Clone(),
	}
}

func (self *chainErrorFormatter) Copy() ErrorFormatter {
	return &chainErrorFormatter{
		opts: self.opts,
		sfmt: self.sfmt,
	}
}

func (self *chainErrorFormatter) SetOptions(opts ErrorFormatterOptions) ErrorFormatter {
	self.opts = opts
	return self
}

func (self *chainErrorFormatter) SetStackTraceFormatter(stFmt StackTraceFormatter) ErrorFormatter {
	self.sfmt = stFmt
	return self
}
