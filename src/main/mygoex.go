package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

// this is plain dummy example code only
// not intended to be "good" go code :-)

func initTracer() {

	ctx := context.Background()

	client := otlptracehttp.NewClient()

	otlpTraceExporter, err := otlptrace.New(ctx, client)
	if err != nil {
		log.Fatal(err)
	}

	batchSpanProcessor := sdktrace.NewBatchSpanProcessor(otlpTraceExporter)

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(batchSpanProcessor),
		//trace.WithSampler(sdktrace.AlwaysSample()), - please check TracerProvider.WithSampler() implementation for details.
	)

	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)
}

func readURL(client http.Client, url string) string {
	resp, err := client.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	// read the response body on the line below
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	// convert the body to string and log
	strBody := string(body)
	log.Println(strBody)
	return strBody
}

func MainServiceHandler(w http.ResponseWriter, r *http.Request) {
	// otel instrumentation
	client := http.Client{
		Timeout:   60 * time.Second,
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}

	//newCtx, span := otel.Tracer("MainServiceHandler").Start(context.Background(), "MainServiceHandler")
	ctx := r.Context()
	span := trace.SpanFromContext(ctx)
	// otel instrumentation

	log.Println("Servicing request.")

	responseEnv := os.Getenv("RESPONSE")
	svc1url := os.Getenv("SERVICE1_URL")
	svc2url := os.Getenv("SERVICE2_URL")
	response := "Hello Caller!"
	response += "\n\n\n"
	if len(responseEnv) == 0 {
		log.Println("No response value in env configured")
		response += "No response value in env configured\n"
	} else {
		log.Print("Response value in env: ")
		log.Println(responseEnv)
		response += responseEnv
		response += "'"
		response += "\n\n\n"
	}

	/// service 1
	if len(svc1url) == 0 {
		log.Println("No service-1 url in env configured")
		response += "No service-1 url in env configured\n"
	} else {
		log.Print("Calling service on URL: ")
		log.Println(svc1url)
		response += "Result from Service-1:\n"
		response += readURL(client, svc1url)
		response += "\n\n\n"
	}

	/// service 2
	if len(svc2url) == 0 {
		log.Println("No service-2 url in env configured")
		response += "No service-2 url in env configured\n"
	} else {
		log.Print("Calling service on URL: ")
		log.Println(svc2url)
		response += "Result from Service-2:\n"
		response += readURL(client, svc2url)
	}

	fmt.Fprintln(w, response)

	// otel instrumentation
	span.End()
	// otel instrumentation
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
	http.HandleFunc("/", MainServiceHandler)
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}
	go listenAndServe(port)

	port = os.Getenv("SECOND_PORT")
	if len(port) == 0 {
		port = "8888"
	}
	go listenAndServe(port)

	select {}
}
