package trace

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/opentracing/opentracing-go"
	opentracinglog "github.com/opentracing/opentracing-go/log"
	"gitlab.meiyou.com/biz-lib/jaeger-client-go"
	"gitlab.meiyou.com/biz-lib/jaeger-client-go/config"
)

func Init() {
	//NewTrace()
}

// InitJaeger ...
func InitJaeger(service string) (opentracing.Tracer, io.Closer) {

	c := config.Configuration{
		Sampler: &config.SamplerConfig{Type: jaeger.SamplerTypeRemote}, // SamplingServerURL: "http://localhost:5778/sampling"

		Reporter: &config.ReporterConfig{
			LogSpans:            false,
			BufferFlushInterval: 60 * time.Second,
			LocalAgentHostPort:  "127.0.0.1:6831",
		}}

	tracer, closer, err := c.New(service, config.Logger(jaeger.StdLogger))
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}

	return tracer, closer
}

var Trace opentracing.Tracer
var Closer io.Closer

func NewTrace(srv string) {
	Trace, Closer = InitJaeger(srv)
	if Closer != nil {

	}
	opentracing.InitGlobalTracer(Trace)
}

//记录推送 的trace info
type PushTrace struct {
	Ctx      context.Context
	IsRecord bool //是否记录  true 记录
}

func NewPushContext(isRecord bool, span opentracing.Span) (trace PushTrace) {
	if !isRecord {
		trace = PushTrace{IsRecord: isRecord}
	} else {
		trace = PushTrace{
			Ctx:      opentracing.ContextWithSpan(context.Background(), span),
			IsRecord: isRecord,
		}
	}
	return trace
}

//简单追加span
func AppendSpan(trace PushTrace, operationName string, event, val string) {
	if trace.IsRecord {
		span, ctx := opentracing.StartSpanFromContext(trace.Ctx, operationName)
		opentracing.ContextWithSpan(ctx, span)
		span.LogFields(
			opentracinglog.String("event", event),
			opentracinglog.String("value", val),
		)
		defer span.Finish()
	}
}
