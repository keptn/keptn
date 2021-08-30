package events

// TODO: Move this into the go-utils package
import (
	"context"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/extensions"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type CloudEventTraceContext struct {
	traceContext propagation.TraceContext
}

func newCloudEventTraceContext() CloudEventTraceContext {
	return CloudEventTraceContext{traceContext: propagation.TraceContext{}}
}

func (etc CloudEventTraceContext) extract(ctx context.Context, carrier CloudEventCarrier) context.Context {
	// TODO: Is there a better way to check if ctx already has a current span on it?
	span := trace.SpanFromContext(ctx)
	if span.IsRecording() {
		// if the context already has an active span so just return that
		return ctx
	}

	// extract the traceparent from the cloud event
	// this is in case the context is not propagated (e.g reading from the queue or something)
	return etc.traceContext.Extract(ctx, carrier)
}

func (etc CloudEventTraceContext) inject(ctx context.Context, carrier CloudEventCarrier) {
	etc.traceContext.Inject(ctx, carrier)
}

func (tc CloudEventTraceContext) Fields() []string {
	return []string{extensions.TraceParentExtension, extensions.TraceStateExtension}
}

type CloudEventCarrier struct {
	Extension *extensions.DistributedTracingExtension
}

func NewCloudEventCarrier() CloudEventCarrier {
	return CloudEventCarrier{Extension: &extensions.DistributedTracingExtension{}}
}

func NewCloudEventCarrierWithEvent(event cloudevents.Event) CloudEventCarrier {
	var te, ok = extensions.GetDistributedTracingExtension(event)
	if !ok {
		return CloudEventCarrier{Extension: &extensions.DistributedTracingExtension{}}
	}
	return CloudEventCarrier{Extension: &te}
}

// Get returns the value associated with the passed key.
func (cec CloudEventCarrier) Get(key string) string {
	switch key {
	case extensions.TraceParentExtension:
		return cec.Extension.TraceParent
	case extensions.TraceStateExtension:
		return cec.Extension.TraceParent
	default:
		return ""
	}
}

// Set stores the key-value pair.
func (cec CloudEventCarrier) Set(key string, value string) {
	switch key {
	case extensions.TraceParentExtension:
		cec.Extension.TraceParent = value
	case extensions.TraceStateExtension:
		cec.Extension.TraceState = value
	}
}

// Keys lists the keys stored in this carrier.
func (cec CloudEventCarrier) Keys() []string {
	return []string{extensions.TraceParentExtension, extensions.TraceStateExtension}
}

// InjectDistributedTracingExtension injects the tracecontext into the event as a DistributedTracingExtension
func InjectDistributedTracingExtension(ctx context.Context, event cloudevents.Event) {

	// TODO: Should we validate if there's already a tracecontext in the event?
	// Calling it will override any existing value..
	tc := newCloudEventTraceContext()
	carrier := NewCloudEventCarrier()
	tc.inject(ctx, carrier)
	carrier.Extension.AddTracingAttributes(&event)
}

// ExtractDistributedTracingExtension reads tracecontext from the cloudevent DistributedTracingExtension into a returned Context.
//
// The returned Context will be a copy of ctx and contain the extracted
// tracecontext as the remote SpanContext. If the extracted tracecontext is
// invalid, the passed ctx will be returned directly instead.
func ExtractDistributedTracingExtension(ctx context.Context, event cloudevents.Event) context.Context {
	tc := newCloudEventTraceContext()
	carrier := NewCloudEventCarrierWithEvent(event)

	ctx = tc.extract(ctx, carrier)

	return ctx
}
