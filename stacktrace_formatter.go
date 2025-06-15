package errstack

import (
	"io"
	"strconv"
	"strings"
)

type StackTraceFormatter interface {
	WithOptions(opts StackTraceFormatOptions) StackTraceFormatter
	WithFrameFormatter(ff FrameFormatter) StackTraceFormatter
	Options() StackTraceFormatOptions
	Format(s StackTrace) string
	FormatBuffer(w io.Writer, s StackTrace)
}

var _ StackTraceFormatter = (*stackTraceFormatter)(nil)

type StackTraceFormatOptions struct {
	FrameIndent    string
	FrameSeparator string
	IndexPrefix    string
	IndexSuffix    string
	SkipStackIndex bool
}

type stackTraceFormatter struct {
	ffmt FrameFormatter
	opts StackTraceFormatOptions
}

func (self *stackTraceFormatter) format(w io.Writer, s StackTrace) {
	if len(s.Frames) <= 0 {
		return
	}

	cnt := len(s.Frames)

	for i, f := range s.Frames {
		if self.opts.SkipStackIndex {
			switch o := w.(type) {
			case io.StringWriter:
				o.WriteString(self.opts.FrameIndent)
			default:
				w.Write(string2Slice(self.opts.FrameIndent))
			}
		} else {
			switch o := w.(type) {
			case io.StringWriter:
				o.WriteString(self.opts.FrameIndent)
				o.WriteString(self.opts.IndexPrefix)
				o.WriteString(strconv.FormatInt(int64(cnt-i-1), 10))
				o.WriteString(self.opts.IndexSuffix)
			default:
				w.Write(string2Slice(self.opts.FrameIndent))
				w.Write(string2Slice(self.opts.IndexPrefix))
				w.Write(string2Slice(strconv.FormatInt(int64(cnt-i-1), 10)))
				w.Write(string2Slice(self.opts.IndexSuffix))
			}
		}

		self.ffmt.FormatBuffer(w, f)

		if i == cnt-1 {
			continue
		}

		switch o := w.(type) {
		case io.StringWriter:
			o.WriteString(self.opts.FrameSeparator)
		default:
			w.Write(string2Slice(self.opts.FrameSeparator))
		}
	}
}

func (self *stackTraceFormatter) WithOptions(opts StackTraceFormatOptions) StackTraceFormatter {
	stf := &stackTraceFormatter{
		opts: opts,
		ffmt: self.ffmt,
	}

	return stf
}

func (self *stackTraceFormatter) WithFrameFormatter(ff FrameFormatter) StackTraceFormatter {
	stf := &stackTraceFormatter{
		opts: self.opts,
		ffmt: ff,
	}

	return stf
}

func (self *stackTraceFormatter) Options() StackTraceFormatOptions {
	return self.opts
}

func (self *stackTraceFormatter) Format(s StackTrace) string {
	sb := strings.Builder{}
	self.format(&sb, s)
	return sb.String()
}

func (self *stackTraceFormatter) FormatBuffer(w io.Writer, s StackTrace) {
	self.format(w, s)
}
