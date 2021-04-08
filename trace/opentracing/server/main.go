package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Augustu/go-draft/trace/opentracing/common"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

func main() {
	common.InitTrace()

	tracer := opentracing.GlobalTracer()

	http.HandleFunc("/publish", func(w http.ResponseWriter, r *http.Request) {
		spanCtx, _ := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
		serverSpan := tracer.StartSpan("server", ext.RPCServerOption(spanCtx))
		time.Sleep(1 * time.Second)
		defer serverSpan.Finish()

	})

	log.Fatal(http.ListenAndServe("127.0.0.1:8082", nil))
}
