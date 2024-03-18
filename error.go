package errstack

var NilErrorString = "<nil>"
var ErrorChainSeparator = ": "

var ErrorFormatter = defaultErrorFormatter

var defaultErrorFormatter = func(err error) string {
	return "Error: " + err.Error() + "\n"
}

type Backtracer interface {
	Backtrace() Backtrace
}

type Unwrapper interface {
	Unwrap() error
}

type Thrower interface {
	Throw() Error
}

type Chainer interface {
	Chain(err error) ChainedError
}
