package logging

import (
	"fmt"
	"runtime"
)

// WARN: This makes cooky names and doesn't work well with closures
func GetCallingNameAtDepth(depth int) string {
	pc, _, _, ok := runtime.Caller(depth)
	if !ok {
		return fmt.Sprintf("exceeded stack depth with %d", depth)
	}
	return runtime.FuncForPC(pc).Name()
}
