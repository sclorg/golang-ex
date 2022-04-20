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
)

// this is plain dummy example code only
// not intended to be "good" go code :-)

func randomOutput() string {
	const outputPart string = "this-is-a-dummy-output-line-which-will-be-concatenated"
	const maxIterations int = 10

	rand.Seed(time.Now().UnixNano())

	iterations := rand.Intn(maxIterations) + 1

	var sb strings.Builder

	log.Printf("Creating output %d iterations.", iterations)

	for i := 1; i < iterations; i++ {
		time.Sleep(time.Second)

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
	response := os.Getenv("RESPONSE")
	if len(response) == 0 {
		response = randomOutput()
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
