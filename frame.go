package errstack

import (
	"strings"
)

const (
	_MIN_STR_BYTES_PER_FRAME_STANDALONE int = 128
)

type Frame struct {
	Function string `json:"function"`
	File     string `json:"file"`
	Line     string `json:"line"`
}

func (self Frame) String() string {
	sb := strings.Builder{}
	DefaultStackFrameFormatter.FormatBuffer(&sb, self)
	return sb.String()
}
