package errstack

import (
	"io"
	"strconv"
	"strings"
)

const (
	_MIN_STR_BYTES_PER_FRAME_STACKTRACE int = 256
)

type StackFrameFormatOptions struct {
	StackLineSeparator rune // Default: '\n'
	SkipStackIndex     bool
	SkipFunctionName   bool
	SkipLocation       bool
}

// This package assumes that the backing implementation of io.Writer(s)
// do not fail on call to `Write`
type StackFormatter func(s StackTrace, opts StackFrameFormatOptions, w ...io.Writer)

type StackTrace struct {
	Frames []Frame `json:"stack"`
}

func SetStackFormatter(f StackFormatter) {
	stackTraceFormatter = f
}

func GetStackFormatter() StackFormatter {
	formatter := defaultStackTraceFormatter
	if stackTraceFormatter != nil {
		formatter = stackTraceFormatter
	}

	return formatter
}

var stackTraceFormatter = defaultStackTraceFormatter

func defaultStackTraceFormatter(bt StackTrace, opts StackFrameFormatOptions, w ...io.Writer) {
	if len(bt.Frames) <= 0 {
		return
	}

	cnt := len(bt.Frames)

	sepSlice := []byte{'\t', '#', '\n', ':', ' '}
	for _, outBuf := range w {
		for i, f := range bt.Frames {
			if opts.SkipStackIndex {
				switch o := outBuf.(type) {
				case io.StringWriter:
					o.WriteString("\t")
				default:
					outBuf.Write(sepSlice[:1])
				}
			} else {
				switch o := outBuf.(type) {
				case io.StringWriter:
					o.WriteString("\t#")
					o.WriteString(strconv.FormatInt(int64(cnt-i-1), 10))
					o.WriteString(": ")
				default:
					outBuf.Write(sepSlice[:2])
					outBuf.Write(string2Slice(strconv.FormatInt(int64(cnt-i-1), 10)))
					outBuf.Write(sepSlice[3:])
				}
			}
			f.Print(opts, outBuf)
			switch o := outBuf.(type) {
			case io.StringWriter:
				o.WriteString("\n")
			default:
				outBuf.Write(sepSlice[2:3])
			}
		}
	}
}

func (self StackTrace) Print(opts StackFrameFormatOptions, w ...io.Writer) {
	GetStackFormatter()(self, opts, w...)
}

func (self StackTrace) String() string {
	sb := strings.Builder{}
	sb.Grow(_MIN_STR_BYTES_PER_FRAME_STACKTRACE * len(self.Frames))
	self.Print(StackFrameFormatOptions{}, &sb)
	return sb.String()
}
