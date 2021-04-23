// Trace a transaction across more than one microservice
// Pass the context between processes using Inject and Extract
// Apply OpenTracing-recommended tags

package main

import (
	"context"
	"net/http"
	"net/url"
	"os"

	xhttp "github.com/goProjects/tracing/lib/http"
	"github.com/goProjects/tracing/lib/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
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
	defer span.Finish()

	ctx := context.Background()
	ctx = opentracing.ContextWithSpan(ctx, span)

	helloStr := formatString(ctx, helloTo)
	printHello(ctx, helloStr)
}

func formatString(ctx context.Context, helloTo string) string {
	span, _ := opentracing.StartSpanFromContext(ctx, "formatString")
	// not using ctx value. If new functions called from where, then passed this ctx
	defer span.Finish()

	v := url.Values{
		"helloTo": []string{helloTo},
	}
	url := "http://localhost:8081/format?" + v.Encode()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err.Error())
	}

	// to pass child span context over the HTTP request, we need ext
	// setting tags as metadata about the HTTP request
	ext.SpanKindRPCClient.Set(span)
	ext.HTTPUrl.Set(span, url)
	ext.HTTPMethod.Set(span, "GET")

	// To propagate the span context over the RPCs and process boundaries,
	// tracing instrumentation uses Inject and Extract
	span.Tracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header),
	)

	resp, err := xhttp.Do(req)
	if err != nil {
		ext.LogError(span, err)
		panic(err.Error())
	}

	helloStr := string(resp)

	span.LogFields(
		log.String("event", "string-format"),
		log.String("value", helloStr),
	)
	return helloStr
}

func printHello(ctx context.Context, helloStr string) {
	span, _ := opentracing.StartSpanFromContext(ctx, "printHello")
	defer span.Finish()

	v := url.Values{
		"helloStr": []string{helloStr},
	}
	url := "http://localhost:8082/publish?" + v.Encode()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err.Error())
	}

	// to pass child span context over the HTTP request, we need ext
	// setting tags as metadata about the HTTP request
	ext.SpanKindRPCClient.Set(span)
	ext.HTTPUrl.Set(span, url)
	ext.HTTPMethod.Set(span, "GET")

	// To propagate the span context over the RPCs and process boundaries,
	// tracing instrumentation uses Inject and Extract
	span.Tracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header),
	)

	if _, err := xhttp.Do(req); err != nil {
		ext.LogError(span, err)
		panic(err.Error())
	}
}
