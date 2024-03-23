package errstack

type StackErrOption interface {
	apply(*stackErrOpts)
}

type stackErrOpts struct {
}

type stackErrOptFunc func(*stackErrOpts)

func (self stackErrOptFunc) apply(o *stackErrOpts) {
	self(o)
}
