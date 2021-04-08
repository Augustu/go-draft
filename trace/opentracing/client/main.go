package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Augustu/go-draft/trace/opentracing/common"
	log "github.com/Augustu/go-micro/v2/logger"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-lib/metrics"

	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
)

// func init() {
// 	common.InitTrace()
// }

func main() {
	// queue cached, need explicit call closer.Close
	closer := common.InitTrace()
	defer closer.Close()

	// parent := opentracing.GlobalTracer().StartSpan("hello")
	// defer parent.Finish()

	// child := opentracing.GlobalTracer().StartSpan(
	// 	"world", opentracing.ChildOf(parent.Context()),
	// )
	// defer child.Finish()

	tracer := opentracing.GlobalTracer()

	clientSpan := opentracing.GlobalTracer().StartSpan("client")
	time.Sleep(1 * time.Second)

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

func main1() {
	cfg := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LocalAgentHostPort: "127.0.0.1:6831",
			LogSpans:           true,
		},
	}

	// Example logger and metrics factory. Use github.com/uber/jaeger-client-go/log
	// and github.com/uber/jaeger-lib/metrics respectively to bind to real logging and metrics
	// frameworks.
	jLogger := jaegerlog.StdLogger
	jMetricsFactory := metrics.NullFactory

	// Initialize tracer with a logger and a metrics factory
	closer, err := cfg.InitGlobalTracer(
		"serviceName",
		jaegercfg.Logger(jLogger),
		jaegercfg.Metrics(jMetricsFactory),
	)
	if err != nil {
		log.Warnf("Could not initialize jaeger tracer: %s", err.Error())
		return
	}
	defer closer.Close()

	parent := opentracing.GlobalTracer().StartSpan("hello")
	defer parent.Finish()

	child := opentracing.GlobalTracer().StartSpan(
		"world", opentracing.ChildOf(parent.Context()),
	)
	defer child.Finish()
}
