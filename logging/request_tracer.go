package logging

import (
	"bytes"
	"fmt"
	"net/http"
	"os"

	"github.com/regcomp/gdpr/config"
)

var RT IRequestTracer

type IRequestTracer interface {
	NewRequestTrace(*CustomWriter, *http.Request)
	UpdateRequestTrace(*http.Request, string) error
	DumpRequestTrace(*http.Request) error
}

type RequestTracer struct {
	requestToTrace   map[*http.Request]*RequestTrace
	displayResponses bool
}

func NewRequestTracer(config *config.RequestTracerConfig) {
	if config.TracerOn == true {
		RT = createRequestTracer(config)
	} else {
		RT = &NoOpRequestTracer{}
	}
}

func (rts *RequestTracer) NewRequestTrace(cw *CustomWriter, r *http.Request) {
	newTrace := newRequestTrace()
	err := newTrace.initRequestTrace(cw, r)
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
	trace.printTrace(rts.displayResponses)
	trace.zeroOutTrace()
	rts.deleteTrace(r)
	return nil
}

func createRequestTracer(config *config.RequestTracerConfig) *RequestTracer {
	var displayResponses bool
	if config.DisplayResponses == true {
		displayResponses = true
	} else {
		displayResponses = false
	}
	return &RequestTracer{
		requestToTrace:   make(map[*http.Request]*RequestTrace),
		displayResponses: displayResponses,
	}
}

func (rt *RequestTracer) AddResponseToOutput(cw *CustomWriter) {
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

type RequestTrace struct {
	req   *http.Request
	trace []string
	cw    *CustomWriter
}

func newRequestTrace() *RequestTrace {
	return &RequestTrace{
		req:   nil,
		trace: nil,
		cw:    nil,
	}
}

func (rt *RequestTrace) initRequestTrace(cw *CustomWriter, r *http.Request) error {
	rt.req = r.Clone(r.Context())
	rt.trace = make([]string, 0, 32)
	rt.cw = cw

	return nil
}

func (rt *RequestTrace) logCurrentFunction(t string) {
	rt.trace = append(rt.trace, t)
}

func (rt *RequestTrace) printTrace(displayResponses bool) {
	var output []byte
	buf := bytes.NewBuffer(output)

	rt.constructPrintTraceOutput(buf, displayResponses)

	fmt.Print(buf.String())
}

func (rt *RequestTrace) constructPrintTraceOutput(buf *bytes.Buffer, displayResponses bool) {
	fmt.Fprintf(buf, "[REQUEST]\n")
	fmt.Fprintf(buf, "%s | %s\n", rt.req.URL.Path, rt.req.Method)
	fmt.Fprintf(buf, "[REQUEST HEADER]\n")
	for name, values := range rt.req.Header {
		if name == "Cookie" {
			continue
		}
		fmt.Fprintf(buf, "%s: %s\n", name, values)
	}
	fmt.Fprintf(buf, "[CALLS]\n")
	for _, call := range rt.trace {
		fmt.Fprintf(buf, "%s\n", call)
	}
	if displayResponses {
		fmt.Fprintf(buf, "[RESPONSE HEADER %d]\n", rt.cw.Code)
		rt.cw.Header().Write(os.Stdout)
		for name, values := range rt.cw.Header() {
			fmt.Fprintf(buf, "%s: %s\n", name, values)
		}
		fmt.Fprintf(buf, "[BODY]\n")
		fmt.Fprintf(buf, "%s\n", rt.cw.Body.String())
		fmt.Println("")
	}
}

func (rt *RequestTrace) zeroOutTrace() {
	rt.req = nil
	rt.trace = nil
}

type NoOpRequestTracer struct{}

func (not *NoOpRequestTracer) NewRequestTrace(*CustomWriter, *http.Request)   {}
func (not *NoOpRequestTracer) UpdateRequestTrace(*http.Request, string) error { return nil }
func (not *NoOpRequestTracer) DumpRequestTrace(*http.Request) error           { return nil }
