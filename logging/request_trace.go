package logging

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
)

const shouldTraceRequestsConfigString = "SHOULD_TRACE_REQUESTS"

type IRequestTracer interface {
	NewRequestTrace(*http.Request)
	UpdateRequestTrace(*http.Request, string) error
	DumpRequestTrace(*http.Request) error
}

type RequestTracer struct {
	requestToTrace map[*http.Request]*RequestTrace
}

func createRequestTraces() *RequestTracer {
	return &RequestTracer{
		requestToTrace: make(map[*http.Request]*RequestTrace),
	}
}

func (rts *RequestTracer) addRequestTrace(r *http.Request, rt *RequestTrace) {
	rts.requestToTrace[r] = rt
}

func (rts *RequestTracer) getTrace(r *http.Request) (*RequestTrace, error) {
	if _, ok := rts.requestToTrace[r]; !ok {
		return nil, fmt.Errorf("could not find trace in getTrace")
	}
	return rts.requestToTrace[r], nil
}

func (rts *RequestTracer) deleteTrace(r *http.Request) error {
	if _, ok := rts.requestToTrace[r]; !ok {
		return fmt.Errorf("could not find trace in deleteTrace")
	}

	delete(rts.requestToTrace, r)
	return nil
}

func NewTracer(getenv func(string) string) IRequestTracer {
	if getenv(shouldTraceRequestsConfigString) == "TRUE" {
		return createRequestTraces()
	}
	return &NoOpRequestTracer{}
}

func (rts *RequestTracer) NewRequestTrace(r *http.Request) {
	newTrace := newRequestTrace()
	err := newTrace.initRequestTrace(r)
	if err != nil {
		// TODO:
	}

	rts.addRequestTrace(r, newTrace)
}

func (rts *RequestTracer) UpdateRequestTrace(r *http.Request, function string) error {
	trace, err := rts.getTrace(r)
	if err != nil {
		// TODO: error, that trace doesnt exist
		return err
	}
	trace.logCurrentFunction(function)
	return nil
}

func (rts *RequestTracer) DumpRequestTrace(r *http.Request) error {
	trace, err := rts.getTrace(r)
	if err != nil {
		// TODO: error, that trace doesnt exist
		return err
	}
	trace.printTrace(0)
	trace.zeroOutTrace()
	rts.deleteTrace(r)
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

func (rt *RequestTrace) printTrace(depth int) { // NOTE: depth is not currently used in a real way
	var output []byte
	buf := bytes.NewBuffer(output)
	padding := strings.Repeat("\t", depth)

	rt.constructPrintTraceOutput(buf, padding)

	fmt.Print(buf.String())
}

func (rt *RequestTrace) constructPrintTraceOutput(buf *bytes.Buffer, padding string) {
	fmt.Fprintf(buf, "%s[REQUEST]\n", padding)
	fmt.Fprintf(buf, "%s%s | %s\n", padding, rt.req.URL.Path, rt.req.Method)
	fmt.Fprintf(buf, "%s[REQUEST HEADER]\n", padding)
	for name, values := range rt.req.Header {
		if name == "Cookie" {
			continue
		}
		fmt.Fprintf(buf, "%s%s: %s\n", padding, name, values)
	}
	fmt.Fprintf(buf, "%s[CALLS]\n", padding)
	for _, call := range rt.trace {
		fmt.Fprintf(buf, "%s%s\n", padding, call)
	}
}

func (rt *RequestTrace) zeroOutTrace() {
	rt.req = nil
	rt.trace = nil
}

type NoOpRequestTracer struct{}

func (not *NoOpRequestTracer) NewRequestTrace(*http.Request)                  {}
func (not *NoOpRequestTracer) UpdateRequestTrace(*http.Request, string) error { return nil }
func (not *NoOpRequestTracer) DumpRequestTrace(*http.Request) error           { return nil }
