package errstack

import (
	"io"
	"strconv"
	"strings"
	"unsafe"
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
	Line     int32  `json:"line"`
}

func (self Frame) formatter() FrameFormatter {
	return GetFrameFormatter()
}

func (self Frame) Print(opts StackFrameFormatOptions, w ...io.Writer) {
	self.formatter()(self, opts, w...)
}

func (self Frame) String() string {
	sb := strings.Builder{}
	self.Print(StackFrameFormatOptions{}, &sb)
	return sb.String()
}

var callFrameFormatter = defaultCallFrameFormatter

var defaultCallFrameFormatter = func(f Frame, opts StackFrameFormatOptions, w ...io.Writer) {
	if opts.SkipFunctionName && opts.SkipLocation {
		return
	}

	for _, outBuf := range w {
		// Stack frame format: <function> [file:line]
		if !opts.SkipFunctionName {
			outBuf.Write(unsafe.Slice(unsafe.StringData(f.Function), len(f.Function)))
			outBuf.Write(unsafe.Slice(unsafe.StringData(" ["), len(" [")))
		}
		if !opts.SkipLocation {
			outBuf.Write(unsafe.Slice(unsafe.StringData(f.File), len(f.File)))
			outBuf.Write([]byte{':'})
			line := strconv.FormatInt(int64(f.Line), 10)
			outBuf.Write(unsafe.Slice(unsafe.StringData(line), len(line)))
		}
		if !opts.SkipFunctionName {
			outBuf.Write([]byte{']'})
		}
	}
}
