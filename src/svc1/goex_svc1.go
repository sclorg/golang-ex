package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tracer := otel.GetTracerProvider().Tracer("goex/svc1")

	var span trace.Span
	_, span = tracer.Start(ctx, "svc1")

	// span := trace.SpanFromContext(ctx)
	defer span.End()
	// defer span.End()

	response := os.Getenv("RESPONSE")
	if len(response) == 0 {
		response = "Hello OpenShift!"
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
	http.HandleFunc("/", helloHandler)
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}
	go listenAndServe(port)

	select {}
}
