package errstack

import (
	"io"
	"unsafe"
)

var NilErrorString = "<nil>"
var ErrorChainSeparator = ": "

type ErrorFormatter func(err error, w ...io.Writer)

func SetErrorFormatter(f ErrorFormatter) {
	errorValueFormatter = f
}

func GetErrorFormatter() ErrorFormatter {
	errFormatter := defaultErrorValueFormatter
	if errorValueFormatter != nil {
		errFormatter = errorValueFormatter
	}

	return errFormatter
}

var errorValueFormatter = defaultErrorValueFormatter

var defaultErrorValueFormatter = func(err error, w ...io.Writer) {
	prefix, errStr := "Error: ", err.Error()
	prefixSlice := unsafe.Slice(unsafe.StringData(prefix), len(prefix))
	errStrSlice := unsafe.Slice(unsafe.StringData(errStr), len(errStr))
	sepSlice := []byte{'\n'}

	for _, outBuf := range w {
		outBuf.Write(prefixSlice)
		outBuf.Write(errStrSlice)
		outBuf.Write(sepSlice)
	}
}
