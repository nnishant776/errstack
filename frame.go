package errstack

import (
	"strconv"
	"strings"
)

type Frame struct {
	File     string `json:"file"`
	Function string `json:"function"`
	Line     int32  `json:"line"`
}

func (self Frame) String() string {
	formatter := defaultCallFrameFormatter
	if CallFrameFormatter != nil {
		formatter = CallFrameFormatter
	}

	return formatter(self)
}

var CallFrameFormatter = defaultCallFrameFormatter

var defaultCallFrameFormatter = func(f Frame) string {
	var fStr strings.Builder
	fStr.Grow(128)

	// Stack frame format: <function> [file:line]
	fStr.WriteString(f.Function)
	fStr.WriteString(" [")
	fStr.WriteString(f.File)
	fStr.WriteRune(':')
	fStr.WriteString(strconv.FormatInt(int64(f.Line), 10))
	fStr.WriteRune(']')

	return fStr.String()
}
