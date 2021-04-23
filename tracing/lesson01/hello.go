// Ref: https://github.com/yurishkuro/opentracing-tutorial/tree/master/go/lesson01

// Instantiate a Tracer
// Create a simple trace
// Annotate the trace

package main

import (
	"fmt"
	"os"

	"github.com/goProjects/tracing/lib/tracing"
	"github.com/opentracing/opentracing-go/log"
)

func main() {
	if len(os.Args) != 2 {
		panic("ERROR: expecting 1 argument!")
	}

	// initialize tracer
	tracer, closer := tracing.Init("hello-world")
	defer closer.Close()

	helloTo := os.Args[1]

	span := tracer.StartSpan("say-hello")
	// metadata about the span
	// use tags when we want to describe attributes of the span
	// that apply to the whole duration of the span
	span.SetTag("hello-to", helloTo) 
	
	helloStr := fmt.Sprintf("Hello, %s!", helloTo)
	span.LogFields(
		log.String("event", "string-format"),
		log.String("value", helloStr),
	)

	println(helloStr)
	span.LogKV("event", "println")

	span.Finish()
}


// mixed logs