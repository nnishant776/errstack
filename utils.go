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

func callers(skip int, cnt int) []Frame {
	var frames []Frame

	if cnt < 0 {
		return frames
	}

	var pcs = make([]uintptr, 20)
	if cnt > 0 && cnt < 20 {
		pcs = pcs[:cnt]
	}

	var count = runtime.Callers(skip+1, pcs[:])
	frames = make([]Frame, 0, count)

	if count == 0 {
		return frames
	}

	var callFrames = runtime.CallersFrames(pcs[:count])

	for {
		f, ok := callFrames.Next()
		frames = append(frames, Frame{
			File:     f.File,
			Function: f.Function,
			Line:     int32(f.Line),
		})

		if !ok {
			break
		}
	}

	return frames
}
