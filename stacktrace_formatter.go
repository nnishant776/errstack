package errstack

import (
	"io"
	"strconv"
	"strings"
)

type StackTraceFormatter interface {
	Options() StackTraceFormatOptions
	FrameFormatter() FrameFormatter
	Format(s StackTrace) string
	FormatBuffer(w io.Writer, s StackTrace)
	Clone() StackTraceFormatter
	SetOptions(opts StackTraceFormatOptions) StackTraceFormatter
	SetFrameFormatter(ffFmt FrameFormatter) StackTraceFormatter
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
	return &stackTraceFormatter{
		opts: self.opts,
		ffmt: ff,
	}
}

func (self *stackTraceFormatter) Options() StackTraceFormatOptions {
	return self.opts
}

func (self *stackTraceFormatter) FrameFormatter() FrameFormatter {
	return self.ffmt
}

func (self *stackTraceFormatter) Format(s StackTrace) string {
	sb := strings.Builder{}
	self.format(&sb, s)
	return sb.String()
}

func (self *stackTraceFormatter) FormatBuffer(w io.Writer, s StackTrace) {
	self.format(w, s)
}

func (self *stackTraceFormatter) Clone() StackTraceFormatter {
	return &stackTraceFormatter{
		ffmt: self.ffmt.Clone(),
		opts: self.opts,
	}
}

func (self *stackTraceFormatter) SetOptions(opts StackTraceFormatOptions) StackTraceFormatter {
	self.opts = opts
	return self
}

func (self *stackTraceFormatter) SetFrameFormatter(ffFmt FrameFormatter) StackTraceFormatter {
	self.ffmt = ffFmt
	return self
}
