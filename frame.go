package errstack

import (
	"io"
	"strconv"
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
	Line     int32  `json:"line"`
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

	funcSlice := string2Slice(f.Function)
	prefixSlice := string2Slice(" [")
	fileSlice := string2Slice(f.File)
	lineSlice := string2Slice(strconv.FormatInt(int64(f.Line), 10))

	suffixSlice := []byte{']'}
	sepSlice := []byte{':'}

	for _, outBuf := range w {
		// Stack frame format: <function> [file:line]
		if !opts.SkipFunctionName {
			outBuf.Write(funcSlice)
			outBuf.Write(prefixSlice)
		}
		if !opts.SkipLocation {
			outBuf.Write(fileSlice)
			outBuf.Write(sepSlice)
			outBuf.Write(lineSlice)
		}
		if !opts.SkipFunctionName {
			outBuf.Write(suffixSlice)
		}
	}
}
