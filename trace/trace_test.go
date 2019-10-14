package trace

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"strings"
	"testing"
	"time"
)

func TestInitJaeger(t *testing.T) {

	NewTrace("test")
	fmt.Println(Trace)
	for i := 0; i <= 1; i++ {
		opentracing.InitGlobalTracer(Trace)
		span := Trace.StartSpan("sunwei-debug")
		span.SetTag("value", "1_2_3")
		ctx := context.Background()
		ctx = opentracing.ContextWithSpan(ctx, span)

		fmt.Println(GetRequestID(ctx))
		span.Finish()

	}
	defer Closer.Close()

	time.Sleep(time.Duration(10) * time.Second)
}

func TestStartSpan(t *testing.T) {

	span := Trace.StartSpan("handle_basic_push")
	span.SetTag("push_id", "111")
	span.Finish()
	trace := NewPushContext(true, span)

	AppendSpan(trace, "wtesttst", "android", "nopush")

}
func otherSpan(trace PushTrace) {

}

func t(ctx context.Context) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "test3")

	opentracing.ContextWithSpan(ctx, span)
	fmt.Println(GetRequestID(ctx))
	defer span.Finish()
}

func GetRequestID(ctx context.Context) interface{} {
	span := opentracing.SpanFromContext(ctx)
	if span == nil {
		return nil
	}
	return strings.SplitN(fmt.Sprintf("%s", span.Context()), ":", 2)[0]
}
