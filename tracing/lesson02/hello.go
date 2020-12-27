// Trace individual functions
// Combine multiple spans into a single trace
// Propagate the in-process context

// A trace is a DAG
// Nodes are span and edges are causal relationships between them
// ChildOf() helps to create edge btw span and rootSpan
// ChildOf indicates that rootSpan has logical dependency on child span
// Root span is reported last as it is the last one to finish

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/goProjects/tracing/lib/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

func main() {
	if len(os.Args) != 2 {
		panic("ERROR: expecting 1 argument!")
	}

	// initialize tracer
	tracer, closer := tracing.Init("hello-world")
	defer closer.Close()
	// StartSpanFromContext function uses opentracing.GlobalTracer() to start the new spans
	// Thus, initialize global variable to instance of Jaeger tracer
	opentracing.SetGlobalTracer(tracer)

	helloTo := os.Args[1]

	span := tracer.StartSpan("say-hello")
	// metadata about the span
	// use tags when we want to describe attributes of the span
	// that apply to the whole duration of the span
	span.SetTag("hello-to", helloTo)

	ctx := context.Background()
	ctx = opentracing.ContextWithSpan(ctx, span)

	helloStr := formatString(ctx, helloTo)
	printHello(ctx, helloStr)

	span.Finish()
}

func formatString(ctx context.Context, helloTo string) string {
	span, _ := opentracing.StartSpanFromContext(ctx, "formatString")
	// not using ctx value. If new functions called from where, then passed this ctx
	defer span.Finish()

	helloStr := fmt.Sprintf("Hello, %s!", helloTo)
	span.LogFields(
		log.String("event", "string-format"),
		log.String("value", helloStr),
	)
	return helloStr
}

func printHello(ctx context.Context, helloStr string) {
	span, _ := opentracing.StartSpanFromContext(ctx, "printHello")
	defer span.Finish()

	println(helloStr)
	span.LogKV("event", "println")
}
