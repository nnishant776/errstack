package errstack

import (
	"strings"
)

const (
	_MIN_STR_BYTES_PER_FRAME_STACKTRACE int = 256
)

type StackTrace struct {
	Frames []Frame `json:"stack"`
}

func (self StackTrace) String() string {
	sb := strings.Builder{}
	sb.Grow(_MIN_STR_BYTES_PER_FRAME_STACKTRACE * len(self.Frames))
	DefaultStackTraceFormatter.FormatBuffer(&sb, self)
	return sb.String()
}
