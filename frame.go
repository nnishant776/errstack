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
type FrameFormatter func(f Frame, w ...io.Writer)

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

func (self Frame) Print(w ...io.Writer) {
	self.formatter()(self, w...)
}

func (self Frame) String() string {
	sb := strings.Builder{}
	self.Print(&sb)
	return sb.String()
}

var callFrameFormatter = defaultCallFrameFormatter

var defaultCallFrameFormatter = func(f Frame, w ...io.Writer) {
	for _, outBuf := range w {
		outStr, strOk := outBuf.(io.StringWriter)

		// Stack frame format: <function> [file:line]
		if strOk {
			outStr.WriteString(f.Function)
			outStr.WriteString(" [")
			outStr.WriteString(f.File)
		} else {
			outBuf.Write([]byte(f.Function))
			outBuf.Write([]byte(" ["))
			outBuf.Write([]byte(f.File))
		}
		outBuf.Write([]byte{':'})
		if strOk {
			outStr.WriteString(strconv.FormatInt(int64(f.Line), 10))
		} else {
			outBuf.Write([]byte(strconv.FormatInt(int64(f.Line), 10)))
		}
		outBuf.Write([]byte{']'})
	}
}
