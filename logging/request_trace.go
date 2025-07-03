package logging

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
)

const shouldTraceRequestsConfigString = "SHOULD_TRACE_REQUESTS"

type Tracer interface {
	NewRequestTrace(*http.Request)
	UpdateActiveTrace(string) error
	DumpActiveTrace() error
}

type RequestTraces struct {
	traces []*RequestTrace
}

func NewTracer(getenv func(string) string) Tracer {
	if getenv(shouldTraceRequestsConfigString) == "TRUE" {
		return &RequestTraces{
			traces: make([]*RequestTrace, 0, 32),
		}
	}
	return &NoOpTracer{}
}

func (rts *RequestTraces) NewRequestTrace(r *http.Request) {
	newTrace := newRequestTrace()
	err := newTrace.initRequestTrace(r)
	if err != nil {
		// TODO:
	}
	rts.traces = append(rts.traces, newTrace)
}

func (rts *RequestTraces) UpdateActiveTrace(function string) error {
	activeIdx := len(rts.traces) - 1
	if activeIdx < 0 {
		return fmt.Errorf("could not log function. no active trace")
	}
	rts.traces[activeIdx].logCurrentFunction(function)
	return nil
}

func (rts *RequestTraces) DumpActiveTrace() error {
	activeIdx := len(rts.traces) - 1
	if activeIdx < 0 {
		return fmt.Errorf("could not print trace. no active trace")
	}
	rts.traces[activeIdx].printTrace(activeIdx)

	rts.traces[activeIdx].zeroOutTrace()
	rts.traces[activeIdx] = nil
	rts.traces = rts.traces[:activeIdx]

	return nil
}

type RequestTrace struct {
	req   *http.Request
	trace []string
}

func newRequestTrace() *RequestTrace {
	return &RequestTrace{
		req:   nil,
		trace: nil,
	}
}

func (rt *RequestTrace) initRequestTrace(r *http.Request) error {
	rt.req = r.Clone(r.Context())
	rt.trace = make([]string, 0, 32)

	return nil
}

func (rt *RequestTrace) logCurrentFunction(t string) {
	rt.trace = append(rt.trace, t)
}

func (rt *RequestTrace) printTrace(depth int) {
	var output []byte
	buf := bytes.NewBuffer(output)
	padding := strings.Repeat("\t", depth)

	rt.constructPrintTraceOutput(buf, padding)

	fmt.Print(buf.String())
}

func (rt *RequestTrace) constructPrintTraceOutput(buf *bytes.Buffer, padding string) {
	fmt.Fprintf(buf, "%s[REQUEST]\n", padding)
	fmt.Fprintf(buf, "%s%s | %s\n", padding, rt.req.URL.Path, rt.req.Method)
	fmt.Fprintf(buf, "%s[CALLS]\n", padding)
	for _, call := range rt.trace {
		fmt.Fprintf(buf, "%s%s\n", padding, call)
	}
}

func (rt *RequestTrace) zeroOutTrace() {
	rt.req = nil
	rt.trace = nil
}

type NoOpTracer struct{}

func (not *NoOpTracer) NewRequestTrace(*http.Request)  {}
func (not *NoOpTracer) UpdateActiveTrace(string) error { return nil }
func (not *NoOpTracer) DumpActiveTrace() error         { return nil }
