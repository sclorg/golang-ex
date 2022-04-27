package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
)

// this is plain dummy example code only
// not intended to be "good" go code :-)

var tracerProvider *sdktrace.TracerProvider

func initTracer() {

	ctx := context.Background()

	client := otlptracehttp.NewClient()

	otlpTraceExporter, err := otlptrace.New(ctx, client)
	if err != nil {
		log.Fatal(err)
	}

	batchSpanProcessor := sdktrace.NewBatchSpanProcessor(otlpTraceExporter)
	// The service.name attribute is required.
	resource :=
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("Service-2"),
		)

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(batchSpanProcessor),
		//trace.WithSampler(sdktrace.AlwaysSample()), - please check TracerProvider.WithSampler() implementation for details.
		sdktrace.WithResource(resource),
	)

	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)
}

func readURL(client http.Client, ctx context.Context, url string) string {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	strBody := string(body)
	log.Println(strBody)
	return strBody

}

func randomOutput(ctx context.Context, url string) string {
	const maxIterations int = 10

	rand.Seed(time.Now().UnixNano())

	iterations := rand.Intn(maxIterations)

	var sb strings.Builder

	log.Printf("Creating output %d iterations.", iterations)
	log.Printf("Service URL to call: '%s' ", url)

	client := http.Client{
		Timeout: 60 * time.Second,
	}

	for i := 1; i < iterations; i++ {
		minSleep := 100
		maxSleep := 500
		randSleep := rand.Intn(maxSleep-minSleep+1) + minSleep
		log.Printf("Iteration %d: Sleeping %d seconds, then adding next string fragment to output\n", i, randSleep)
		time.Sleep(time.Duration(randSleep) * time.Millisecond)

		/// service 3
		if len(url) == 0 {
			log.Println("No service-3 url in env configured")
			sb.WriteString("No service-3 url in env configured\n")
		} else {
			log.Print("Calling service on URL: ")
			log.Println(url)
			sb.WriteString("Result from Service-3, iteration #")
			sb.WriteString(strconv.Itoa(i))
			sb.WriteString(":\n")
			sb.WriteString(readURL(client, ctx, url))
			sb.WriteString("\n\n\n")
		}

	}

	result := sb.String()
	log.Println(result)
	return result
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tracer := otel.GetTracerProvider().Tracer("goex/svc2")

	var span trace.Span
	ctx, span = tracer.Start(ctx, "svc2")
	// span := trace.SpanFromContext(ctx)
	defer span.End()

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(r.Header))
	// otel instrumentation

	svc3url := os.Getenv("SERVICE3_URL")
	response := ""
	if len(response) == 0 {
		response = randomOutput(ctx, svc3url)
	}

	fmt.Fprintln(w, response)
	log.Println("Servicing request.")
}

func listenAndServe(port string) {
	log.Printf("serving on %s\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

func main() {
	initTracer()
	http.HandleFunc("/", helloHandler)
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}
	go listenAndServe(port)

	select {}
}
