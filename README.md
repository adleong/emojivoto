# Distributed Tracing with Emoji.voto

This is a fork of the [Emoji.voto](https://github.com/BuoyantIO/emojivoto)
project with distributed tracing instrumentation built in.  This is intended to
serve as a reference architecture for how to add distributed tracing to an
application.

For more information about Emoji.voto and for instructions on how to build and
run the project, see [the upstream
README](https://github.com/BuoyantIO/emojivoto/blob/master/README.md).

To see the changes that were necessary to instrument Emoji.voto, see [the
diff](https://github.com/BuoyantIO/emojivoto/compare/master...adleong:master).

## Quick Start

```
kubectl apply -f tracing.yml
kubectl apply -f emojivoto.yml
kubectl apply -f ingress.yml
kubectl -n tracing port-forward deploy/jaeger 16686 &
open http://localhost:16686
```

## Architecture

We use the [OpenCensus Service
collector](https://opencensus.io/service/components/collector/) to collect
traces and [Jaeger](https://www.jaegertracing.io/) to store and display them. 
Each Emoji.voto service is instrumented with the [OpenCensus Go
client](https://github.com/census-instrumentation/opencensus-go) to emit trace
data to the collector.  [Nginx](https://www.nginx.com/) acts as an ingress and
makes all sampling decisions about when to initiate a trace (in this example we
use a 50% sample rate.)  HTTP communication uses the [Zipkin trace propagation
headers](https://github.com/openzipkin/b3-propagation) and gRPC communication
uses the [gRPC trace
metadata](https://github.com/census-instrumentation/opencensus-specs/blob/master/trace/gRPC.md).

### Nginx

Nginx is deployed as an ingress controller.  For each request it receives, it
has a 50% change of sampling that trace.  For any traces that it samples, it
sends span data to the OpenCensus collector using the Zipkin reporting protocol
and sets the `X-b3-*` headers on the request to mark that downstream services
should sample it as well in order to produce a full trace.

### Emoji.voto

All Emoji.voto services use the OpenCensus Go client to propagate trace context
from the incoming requests to the outgoing requests.  They also honor the
sampling decision made by the ingress.  If a trace should be sampled, they
report span data to the OpenCensus collector using the OpenCensus agent
protocol.

### OpenCensus Collector

The collector is an aggregation and translation layer which receives span data
from Nginx and the Emoji.voto services and forwards that data to Jaeger.

### Jaeger

Jaeger receives traces from the OpenCensus collector, stores them, and displays
them in a web UI.
