package errstack

import (
	"fmt"
	"strconv"
	"strings"
)

type Backtrace struct {
	Frames []Frame `json:"stack"`
}

var BacktraceFormatter = defaultBacktraceFormatter

var defaultBacktraceFormatter = func(bt Backtrace) string {
	var btStr strings.Builder

	for cnt, i := len(bt.Frames), 0; i < cnt; i++ {
		f := bt.Frames[i]
		btStr.WriteRune('\t')
		btStr.WriteRune('#')
		btStr.WriteString(strconv.FormatInt(int64(cnt-1-i), 10))
		btStr.WriteRune(':')
		btStr.WriteRune(' ')
		btStr.WriteString(f.Function)
		btStr.WriteRune(' ')
		btStr.WriteRune('[')
		btStr.WriteString(f.File)
		btStr.WriteRune(':')
		btStr.WriteString(strconv.FormatInt(int64(f.Line), 10))
		btStr.WriteRune(']')
		btStr.WriteRune('\n')
	}

	return btStr.String()
}

func (self Backtrace) String() string {
	return fmt.Sprintf("%+v", self.Frames)
}
