package errstack

var DefaultErrorFormatter ErrorFormatter = &errorFormatter{
	opts: ErrorFormatterOptions{
		// ErrorSeparator: "\n",
		// ErrorPrefix: "Error: ",
	},
}

var DefaultStackFrameFormatter FrameFormatter = &frameFormatter{}

var DefaultStackTraceFormatter StackTraceFormatter = &stackTraceFormatter{
	ffmt: DefaultStackFrameFormatter,
	opts: StackTraceFormatOptions{
		// FrameIndent: "\t",
		// FrameSeparator: "\n",
		IndexPrefix: "#",
		IndexSuffix: ": ",
	},
}
