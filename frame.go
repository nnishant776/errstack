package errstack

import (
	"io"
	"strings"
)

const (
	_MIN_STR_BYTES_PER_FRAME_STANDALONE int = 128
)

// This package assumes that the backing implementation of io.Writer(s)
// do not fail on call to `Write`
type FrameFormatter func(f Frame, opts StackFrameFormatOptions, w ...io.Writer)

func SetFrameFormatter(f FrameFormatter) {
	callFrameFormatter = f
}

func GetFrameFormatter() FrameFormatter {
	formatter := defaultCallFrameFormatter
	if callFrameFormatter != nil {
		formatter = callFrameFormatter
	}

	return formatter
}

type Frame struct {
	File     string `json:"file"`
	Function string `json:"function"`
	Line     string `json:"line"`
}

func (self Frame) Print(opts StackFrameFormatOptions, w ...io.Writer) {
	GetFrameFormatter()(self, opts, w...)
}

func (self Frame) String() string {
	sb := strings.Builder{}
	self.Print(StackFrameFormatOptions{}, &sb)
	return sb.String()
}

var callFrameFormatter = defaultCallFrameFormatter

func defaultCallFrameFormatter(f Frame, opts StackFrameFormatOptions, w ...io.Writer) {
	if opts.SkipFunctionName && opts.SkipLocation {
		return
	}

	for _, outBuf := range w {
		// Stack frame format: <function> [file:line]
		if !opts.SkipFunctionName {
			switch o := outBuf.(type) {
			case io.StringWriter:
				o.WriteString(f.Function)
				o.WriteString(" [")
			default:
				outBuf.Write(string2Slice(f.Function))
				outBuf.Write(string2Slice(" ["))
			}
		}
		if !opts.SkipLocation {
			switch o := outBuf.(type) {
			case io.StringWriter:
				o.WriteString(f.File)
				o.WriteString(":")
				o.WriteString(f.Line)
			default:
				outBuf.Write(string2Slice(f.File))
				outBuf.Write([]byte{':'})
				outBuf.Write(string2Slice(f.Line))
			}
		}
		if !opts.SkipFunctionName {
			switch o := outBuf.(type) {
			case io.StringWriter:
				o.WriteString("]")
			default:
				outBuf.Write([]byte{']'})
			}
		}
	}
}
