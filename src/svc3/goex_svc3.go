package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

// this is plain dummy example code only
// not intended to be "good" go code :-)

func randomOutput() string {
	const outputPart string = "this-is-a-dummy-output-line-which-will-be-concatenated"
	const maxIterations int = 20

	rand.Seed(time.Now().UnixNano())

	iterations := rand.Intn(maxIterations) + 1

	var sb strings.Builder

	log.Printf("Creating output %d iterations.", iterations)

	for i := 1; i < iterations; i++ {
		minSleep := 10
		maxSleep := 100
		randSleep := rand.Intn(maxSleep-minSleep+1) + minSleep
		log.Printf("Iteration %d: Sleeping %d seconds, then adding next string fragment to output\n", i, randSleep)
		time.Sleep(time.Duration(randSleep) * time.Millisecond)

		sb.WriteString("- ")
		sb.WriteString(strconv.Itoa(i))
		sb.WriteString(" - ")
		sb.WriteString(outputPart)
		sb.WriteString(" | ")
	}

	result := sb.String()
	log.Println(result)
	return result
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tracer := otel.GetTracerProvider().Tracer("goex/svc3")

	var span trace.Span
	_, span = tracer.Start(ctx, "svc2")

	// span := trace.SpanFromContext(ctx)
	defer span.End()
	// defer span.End()

	response := randomOutput()

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
