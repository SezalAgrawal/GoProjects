package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/goProjects/tracing/lib/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	logx "github.com/opentracing/opentracing-go/log"
)

func main() {
	// initialize tracer
	tracer, closer := tracing.Init("formatter")
	defer closer.Close()

	http.HandleFunc("/format", func(w http.ResponseWriter, r *http.Request) {
		spanCtx, _ := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
		span := tracer.StartSpan("format", ext.RPCServerOption(spanCtx))
		defer span.Finish()

		helloTo := r.FormValue("helloTo")
		helloStr := fmt.Sprintf("Hello, %s!", helloTo)
		span.LogFields(
			logx.String("event", "string-format"),
			logx.String("value", helloStr),
		)
		w.Write([]byte(helloStr))
	})
	log.Fatal(http.ListenAndServe(":8081", nil))
}
