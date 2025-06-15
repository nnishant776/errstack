package errstack

var DefaultErrorFormatter ErrorFormatter = &errorFormatter{
	opts: ErrorFormatterOptions{
		ErrorSeparator: "=>",
		// ErrorPrefix: "Error: ",
	},
}

var DefaultStackFrameFormatter FrameFormatter = &frameFormatter{
	opts: FrameFormatterOptions{
		LocationPrefix: "@",
		// LocationSuffix:    "]",
		FileLineSeparator: ":",
	},
}

var DefaultStackTraceFormatter StackTraceFormatter = &stackTraceFormatter{
	ffmt: DefaultStackFrameFormatter,
	opts: StackTraceFormatOptions{
		// FrameIndent: "\t",
		FrameSeparator: ";",
		IndexPrefix:    "#",
		IndexSuffix:    ": ",
	},
}
