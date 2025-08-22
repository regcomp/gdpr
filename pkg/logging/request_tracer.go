package logging

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/google/uuid"
)

var RT IRequestTracer

func TraceRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cw := CreateCustomWriter(w)
		RT.NewRequestTrace(cw, r)
		next.ServeHTTP(cw, r)
		RT.DumpRequestTrace(r)
	})
}

type CustomWriter struct {
	http.ResponseWriter
	Code int
	Body bytes.Buffer
}

func CreateCustomWriter(w http.ResponseWriter) *CustomWriter {
	return &CustomWriter{ResponseWriter: w, Code: http.StatusOK}
}

func (cw *CustomWriter) WriteHeader(code int) {
	cw.Code = code
	cw.ResponseWriter.WriteHeader(code)
}

func (cw *CustomWriter) Write(data []byte) (int, error) {
	cw.Body.Write(data)
	return cw.ResponseWriter.Write(data)
}

type IRequestTracer interface {
	NewRequestTrace(*CustomWriter, *http.Request)
	UpdateRequestTrace(*http.Request, string) error
	DumpRequestTrace(*http.Request) error
}

type RequestTracer struct {
	mu               sync.RWMutex
	requestIDToTrace map[string]*RequestTrace
	responseBody     bool
}

func NewRequestTracer(enabled, displayResponses bool) {
	if enabled {
		RT = createRequestTracer(displayResponses)
	} else {
		RT = &NoOpRequestTracer{}
	}
}

func createRequestTracer(displayResponses bool) *RequestTracer {
	return &RequestTracer{
		mu:               sync.RWMutex{},
		requestIDToTrace: make(map[string]*RequestTrace),
		responseBody:     displayResponses,
	}
}

type contextKey string

const RequestIDKey contextKey = "tracer_id"

func (rts *RequestTracer) NewRequestTrace(cw *CustomWriter, r *http.Request) {
	id := generateRequestID()

	ctx := context.WithValue(r.Context(), RequestIDKey, id)
	*r = *r.WithContext(ctx)

	newTrace := newRequestTrace(cw, r)

	rts.mu.Lock()
	defer rts.mu.Unlock()
	rts.storeTrace(id, newTrace)
}

func generateRequestID() string {
	return uuid.New().String()
}

func (rts *RequestTracer) UpdateRequestTrace(r *http.Request, function string) error {
	requestID, ok := r.Context().Value(RequestIDKey).(string)
	if !ok {
		return fmt.Errorf("no request ID found in context")
	}
	rts.mu.Lock()
	defer rts.mu.Unlock()
	trace, err := rts.getTrace(requestID)
	if err != nil {
		return err
	}
	trace.logCurrentFunction(function)
	return nil
}

func (rts *RequestTracer) DumpRequestTrace(r *http.Request) error {
	id, ok := r.Context().Value(RequestIDKey).(string)
	if !ok {
		return fmt.Errorf("no request ID found in context")
	}

	rts.mu.Lock()
	defer rts.mu.Unlock()
	trace, err := rts.getTrace(id)
	if err != nil {
		return err
	}
	trace.printTrace(rts.responseBody)
	err = rts.deleteTrace(id)
	if err != nil {
		return err
	}
	return nil
}

func (rts *RequestTracer) storeTrace(id string, trace *RequestTrace) {
	rts.requestIDToTrace[id] = trace
}

func (rts *RequestTracer) getTrace(id string) (*RequestTrace, error) {
	if _, ok := rts.requestIDToTrace[id]; !ok {
		return nil, fmt.Errorf("could not find trace for request at=%s", id)
	}
	return rts.requestIDToTrace[id], nil
}

func (rts *RequestTracer) deleteTrace(id string) error {
	trace, ok := rts.requestIDToTrace[id]
	if !ok {
		return fmt.Errorf("could not find trace for request at=%s", id)
	}

	trace.zeroOutTrace()
	delete(rts.requestIDToTrace, id)
	return nil
}

type RequestTrace struct {
	path    string
	method  string
	qParams url.Values
	header  http.Header
	trace   []string
	cw      *CustomWriter
}

func newRequestTrace(cw *CustomWriter, r *http.Request) *RequestTrace {
	return &RequestTrace{
		path:    r.URL.Path,
		method:  r.Method,
		qParams: r.URL.Query(),
		header:  r.Header.Clone(),
		trace:   make([]string, 0, 32),
		cw:      cw,
	}
}

func (rt *RequestTrace) logCurrentFunction(t string) {
	rt.trace = append(rt.trace, t)
}

func (rt *RequestTrace) printTrace(responseBody bool) {
	b := &strings.Builder{}
	rt.constructPrintTraceOutput(b, responseBody)

	fmt.Print(b.String())
}

func (rt *RequestTrace) constructPrintTraceOutput(b *strings.Builder, responseBody bool) {
	fmt.Fprint(b, "[REQUEST]\n")
	fmt.Fprintf(b, "%s | %s\n", rt.path, rt.method)
	for k, vs := range rt.qParams {
		fmt.Fprintf(b, "%s: ", k)
		for i, v := range vs {
			fmt.Fprintf(b, "%s", v)
			if i != len(vs)-1 {
				fmt.Fprint(b, ",")
			}
		}
		fmt.Fprint(b, "\n")
	}
	fmt.Fprint(b, "[REQUEST HEADER]\n")
	for name, values := range rt.header {
		if name == "Cookie" {
			continue
		}
		fmt.Fprintf(b, "%s: %s\n", name, values)
	}
	fmt.Fprint(b, "[CALLS]\n")
	for _, call := range rt.trace {
		fmt.Fprintf(b, "%s\n", call)
	}
	fmt.Fprintf(b, "[RESPONSE HEADER %d]\n", rt.cw.Code)
	for name, values := range rt.cw.Header() {
		fmt.Fprintf(b, "%s: %s\n", name, values)
	}
	if responseBody {
		fmt.Fprint(b, "[BODY]\n")
		fmt.Fprintf(b, "%s\n", rt.cw.Body.String())
	}
	fmt.Fprint(b, "\n")
}

func (rt *RequestTrace) zeroOutTrace() {
	rt.trace = nil
	rt.cw = nil
}

type NoOpRequestTracer struct{}

func (not *NoOpRequestTracer) NewRequestTrace(*CustomWriter, *http.Request)   {}
func (not *NoOpRequestTracer) UpdateRequestTrace(*http.Request, string) error { return nil }
func (not *NoOpRequestTracer) DumpRequestTrace(*http.Request) error           { return nil }
