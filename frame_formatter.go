package errstack

import (
	"io"
	"strings"
)

type FrameFormatterOptions struct {
	SkipFunctionName bool
	SkipLocation     bool
}

type FrameFormatter interface {
	WithOptions(opts FrameFormatterOptions) FrameFormatter
	Options() FrameFormatterOptions
	Format(f Frame) string
	FormatBuffer(w io.Writer, f Frame)
}

var _ FrameFormatter = (*frameFormatter)(nil)

type frameFormatter struct {
	opts FrameFormatterOptions
}

func (self *frameFormatter) format(w io.Writer, f Frame) {
	if self.opts.SkipFunctionName && self.opts.SkipLocation {
		return
	}

	if !self.opts.SkipFunctionName {
		switch o := w.(type) {
		case io.StringWriter:
			o.WriteString(f.Function)
		default:
			w.Write(string2Slice(f.Function))
		}
	}

	if !self.opts.SkipLocation {
		switch o := w.(type) {
		case io.StringWriter:
			if !self.opts.SkipFunctionName {
				o.WriteString(" [")
			}
			o.WriteString(f.File)
			o.WriteString(":")
			o.WriteString(f.Line)
			if !self.opts.SkipFunctionName {
				o.WriteString("]")
			}
		default:
			if !self.opts.SkipFunctionName {
				w.Write(string2Slice(" ["))
			}
			w.Write(string2Slice(f.File))
			w.Write([]byte{':'})
			w.Write(string2Slice(f.Line))
			if !self.opts.SkipFunctionName {
				w.Write([]byte{']'})
			}
		}
	}
}

func (self *frameFormatter) WithOptions(opts FrameFormatterOptions) FrameFormatter {
	ff := &frameFormatter{
		opts: opts,
	}

	return ff
}

func (self *frameFormatter) Options() FrameFormatterOptions {
	return self.opts
}

func (self *frameFormatter) Format(f Frame) string {
	sb := strings.Builder{}
	self.format(&sb, f)
	return sb.String()
}

func (self *frameFormatter) FormatBuffer(w io.Writer, f Frame) {
	self.format(w, f)
}
