package errstack

import (
	"math"
	"runtime"
)

const (
	_MAX_CALL_DEPTH int = 32
)

func caller(skip int) Frame {
	pc := callerPC(skip)
	if pc == math.MaxUint64 {
		return Frame{}
	}

	frames := genStackTraceFromPCs([]uintptr{pc})
	return frames[0]
}

func callerPC(skip int) uintptr {
	pc, _, _, ok := runtime.Caller(skip + 1)
	if !ok {
		return math.MaxUint64
	}

	return pc
}

func callers(skip int, cnt int) []Frame {
	if cnt < 0 {
		return nil
	}

	pcs := callersPCs(skip, cnt)
	if len(pcs) <= 0 {
		return nil
	}

	return genStackTraceFromPCs(pcs)
}

func callersPCs(skip int, cnt int) []uintptr {
	if cnt <= 0 {
		return nil
	}

	pcs := make([]uintptr, cnt)
	count := runtime.Callers(skip+1, pcs[:])
	if count == 0 {
		return nil
	}

	return pcs[:count]
}

func genStackTraceFromPCs(pcs []uintptr) []Frame {
	if len(pcs) <= 0 {
		return nil
	}

	frames := make([]Frame, 0, len(pcs))
	callFrames := runtime.CallersFrames(pcs)

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
