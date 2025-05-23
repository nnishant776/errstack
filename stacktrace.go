package errstack

import (
	"io"
	"strconv"
	"strings"
)

const (
	_MIN_STR_BYTES_PER_FRAME_STACKTRACE int = 256
)

// This package assumes that the backing implementation of io.Writer(s)
// do not fail on call to `Write`
type StackFormatter func(s StackTrace, w ...io.Writer)

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

var defaultStackTraceFormatter = func(bt StackTrace, w ...io.Writer) {
	if len(bt.Frames) <= 0 {
		return
	}

	cnt := len(bt.Frames)

	for _, outBuf := range w {
		outStr, strOk := outBuf.(io.StringWriter)
		for i, f := range bt.Frames {
			if strOk {
				outStr.WriteString("\t#")
				outStr.WriteString(strconv.FormatInt(int64(cnt-1-i), 10))
				outStr.WriteString(": ")
			} else {
				outBuf.Write([]byte("\t#"))
				outBuf.Write([]byte(strconv.FormatInt(int64(cnt-1-i), 10)))
				outBuf.Write([]byte(": "))
			}
			f.Print(outBuf)
			outBuf.Write([]byte{'\n'})
		}
	}
}

func (self StackTrace) formatter() StackFormatter {
	return GetStackFormatter()
}

func (self StackTrace) Print(w ...io.Writer) {
	self.formatter()(self, w...)
}

func (self StackTrace) String() string {
	sb := strings.Builder{}
	self.Print(&sb)
	return sb.String()
}
