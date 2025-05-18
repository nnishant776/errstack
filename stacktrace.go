package errstack

import (
	"strconv"
	"strings"
)

type StackTrace struct {
	Frames []Frame `json:"stack"`
}

var StackTraceFormatter = defaultStackTraceFormatter

var defaultStackTraceFormatter = func(bt StackTrace) string {
	const strBytesPerFrame = 256
	btStr := strings.Builder{}
	btStr.Grow(256 * len(bt.Frames))
	cnt := len(bt.Frames)

	for i, f := range bt.Frames {
		btStr.WriteString("\t#")
		btStr.WriteString(strconv.FormatInt(int64(cnt-1-i), 10))
		btStr.WriteString(": ")
		btStr.WriteString(f.Function)
		btStr.WriteString(" [")
		btStr.WriteString(f.File)
		btStr.WriteByte(':')
		btStr.WriteString(strconv.FormatInt(int64(f.Line), 10))
		btStr.WriteString("]\n")
	}

	return btStr.String()
}

func (self StackTrace) String() string {
	formatter := defaultStackTraceFormatter
	if StackTraceFormatter != nil {
		formatter = StackTraceFormatter
	}

	return formatter(self)
}
