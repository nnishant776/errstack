package errstack

import "io"

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
	for _, outBuf := range w {
		outStr, strOk := outBuf.(io.StringWriter)
		if strOk {
			outStr.WriteString("Error: ")
			outStr.WriteString(err.Error())
			outStr.WriteString("\n")
		} else {
			outBuf.Write([]byte("Error: "))
			outBuf.Write([]byte(err.Error()))
			outBuf.Write([]byte("\n"))
		}
	}
}
