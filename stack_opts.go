package errstack

type StackErrOption interface {
	apply(*stackErrOpts)
}

type stackErrOpts struct {
	autoStacktrace bool
}

type stackErrOptFunc func(*stackErrOpts)

func (self stackErrOptFunc) apply(o *stackErrOpts) {
	self(o)
}

func WithTraceback() StackErrOption {
	return stackErrOptFunc(func(o *stackErrOpts) {
		o.autoStacktrace = true
	})
}
