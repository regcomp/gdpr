package logging

import (
	"bytes"
	"fmt"
	"net/http"
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

type CustomWriter struct {
	http.ResponseWriter
	Code int
	Body bytes.Buffer
}

func (cw *CustomWriter) WriteHeader(code int) {
	cw.Code = code
	cw.ResponseWriter.WriteHeader(code)
}

func (cw *CustomWriter) Write(data []byte) (int, error) {
	cw.Body.Write(data)
	return cw.ResponseWriter.Write(data)
}
