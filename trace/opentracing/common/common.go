package common

import (
	"io"

	log "github.com/Augustu/go-micro/v2/logger"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
)

func InitTrace() io.Closer {
	cfg := jaegercfg.Configuration{
		// ServiceName: "demo client",
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LocalAgentHostPort: "127.0.0.1:6831",
			LogSpans:           true,
		},
	}

	jLogger := jaegerlog.StdLogger
	jMetricsFactory := metrics.NullFactory

	closer, err := cfg.InitGlobalTracer(
		"demo client",
		jaegercfg.Logger(jLogger),
		jaegercfg.Metrics(jMetricsFactory),
	)
	if err != nil {
		log.Warnf("init global tracer failed: %s", err.Error())
	}

	return closer

	// defer closer.Close()

	// tracer, _, _ := cfg.NewTracer(
	// 	jaegercfg.Logger(jLogger),
	// 	jaegercfg.Metrics(jMetricsFactory),
	// )

	// opentracing.SetGlobalTracer(tracer)
	// defer closer.Close()
}
