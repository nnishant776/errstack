package errstack

var DefaultStackErrorFormatter ErrorFormatter = &errorFormatter{
	stFmt: DefaultStackTraceFormatter,
	opts: ErrorFormatterOptions{
		// ErrorSeparator:      "",
		StackTraceSeparator: "=>",
		// ErrorPrefix: "Error: ",
	},
}

var DefaultChainErrorFormatter ErrorFormatter = &chainErrorFormatter{
	sfmt: DefaultStackTraceFormatter,
	opts: ErrorFormatterOptions{
		ErrorSeparator:      ", ",
		StackTraceSeparator: "=>",
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
