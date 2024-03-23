package errstack

import (
	"runtime"
)

func caller(skip int) Frame {
	pc, _, _, ok := runtime.Caller(skip + 1)
	if !ok {
		return Frame{}
	}

	var frames = runtime.CallersFrames([]uintptr{pc})
	f, _ := frames.Next()

	return Frame{
		File:     f.File,
		Function: f.Function,
		Line:     int32(f.Line),
	}
}
