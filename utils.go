package errstack

import (
	"runtime"
)

func caller(skip int) Frame {
	pc, _, _, ok := runtime.Caller(skip + 1)
	if !ok {
		return Frame{}
	}

	frames := runtime.CallersFrames([]uintptr{pc})
	f, _ := frames.Next()

	return Frame{
		File:     f.File,
		Function: f.Function,
		Line:     int32(f.Line),
	}
}

func callers(skip int, cnt int) []Frame {
	if cnt < 0 {
		return nil
	}

	pcs := make([]uintptr, 20)
	if cnt > 0 && cnt < 20 {
		pcs = pcs[:cnt]
	}

	count := runtime.Callers(skip+1, pcs[:])
	if count == 0 {
		return nil
	}

	frames := make([]Frame, 0, count)
	callFrames := runtime.CallersFrames(pcs[:count])

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
