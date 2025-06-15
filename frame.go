package errstack

import (
	"strings"
)

const (
	_MIN_STR_BYTES_PER_FRAME_STANDALONE int = 128
)

type Frame struct {
	File     string `json:"file"`
	Function string `json:"function"`
	Line     string `json:"line"`
}

func (self Frame) String() string {
	sb := strings.Builder{}
	DefaultStackFrameFormatter.FormatBuffer(&sb, self)
	return sb.String()
}
