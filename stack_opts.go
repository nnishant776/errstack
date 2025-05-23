package errstack

type stackErrOpts struct {
	autoStacktrace bool
}

type StackErrOption func(stackErrOpts) stackErrOpts

func WithStack() StackErrOption {
	return func(o stackErrOpts) stackErrOpts {
		o.autoStacktrace = true
		return o
	}
}
