package errstack

type stackErrOpts struct {
	extraFrameSkip int
	autoStacktrace bool
}

type StackErrOption func(stackErrOpts) stackErrOpts

func WithStack() StackErrOption {
	return func(o stackErrOpts) stackErrOpts {
		o.autoStacktrace = true
		return o
	}
}

func withExtraFrameSkip(n int) StackErrOption {
	return func(o stackErrOpts) stackErrOpts {
		o.extraFrameSkip = max(n, 0)
		return o
	}
}
