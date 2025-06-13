package errstack

import (
	"io"
	"strconv"
	"strings"
	"unsafe"
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

var defaultStackTraceFormatter = func(bt StackTrace, opts StackFrameFormatOptions, w ...io.Writer) {
	if len(bt.Frames) <= 0 {
		return
	}

	cnt := len(bt.Frames)

	for _, outBuf := range w {
		for i, f := range bt.Frames {
			if opts.SkipStackIndex {
				outBuf.Write([]byte{'\t'})
			} else {
				prefix := "\t#"
				index := strconv.FormatInt(int64(cnt-1-i), 10)
				suffix := ": "
				outBuf.Write(unsafe.Slice(unsafe.StringData(prefix), len(prefix)))
				outBuf.Write(unsafe.Slice(unsafe.StringData(index), len(index)))
				outBuf.Write(unsafe.Slice(unsafe.StringData(suffix), len(suffix)))
			}
			f.Print(opts, outBuf)
			outBuf.Write([]byte{'\n'})
		}
	}
}

func (self StackTrace) formatter() StackFormatter {
	return GetStackFormatter()
}

func (self StackTrace) Print(opts StackFrameFormatOptions, w ...io.Writer) {
	self.formatter()(self, opts, w...)
}

func (self StackTrace) String() string {
	sb := strings.Builder{}
	self.Print(StackFrameFormatOptions{}, &sb)
	return sb.String()
}
