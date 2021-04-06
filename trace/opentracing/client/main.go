package main

import (
	"fmt"
	"net/http"

	"github.com/Augustu/go-draft/trace/opentracing/common"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

func main() {
	common.InitTrace()

	tracer := opentracing.GlobalTracer()

	clientSpan := tracer.StartSpan("client")
	defer clientSpan.Finish()

	url := "http://localhost:8082/publish"
	req, _ := http.NewRequest("GET", url, nil)

	ext.SpanKindRPCClient.Set(clientSpan)
	ext.HTTPUrl.Set(clientSpan, url)
	ext.HTTPMethod.Set(clientSpan, "GET")

	tracer.Inject(clientSpan.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
	resp, _ := http.DefaultClient.Do(req)
	fmt.Println(resp)
}
