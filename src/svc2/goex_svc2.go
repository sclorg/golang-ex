package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

// this is plain dummy example code only
// not intended to be "good" go code :-)

func readURL(url string) string {
	client := http.Client{
		Timeout: 60 * time.Second,
	}
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

func randomOutput() string {
	const outputPart string = "this-is-a-dummy-output-line-which-will-be-concatenated"
	const maxIterations int = 10

	svc3url := os.Getenv("SERVICE31_URL")

	rand.Seed(time.Now().UnixNano())

	iterations := rand.Intn(maxIterations)

	var sb strings.Builder

	log.Printf("Creating output %d iterations.", iterations)

	for i := 1; i < iterations; i++ {
		minSleep := 1
		maxSleep := 3
		randSleep := rand.Intn(maxSleep-minSleep+1) + minSleep
		log.Printf("Iteration %d: Sleeping %d seconds, then adding next string fragment to output\n", i, randSleep)
		time.Sleep(time.Duration(randSleep) * time.Second)

		/// service 3
		if len(svc3url) == 0 {
			log.Println("No service-3 url in env configured")
			sb.WriteString("No service-3 url in env configured\n")
		} else {
			log.Print("Calling service on URL: ")
			log.Println(svc3url)
			sb.WriteString("Result from Service-3:\n")
			sb.WriteString(readURL(svc3url))
			sb.WriteString("\n\n\n")
		}

	}

	result := sb.String()
	log.Println(result)
	return result
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	response := ""
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

	port = os.Getenv("SECOND_PORT")
	if len(port) == 0 {
		port = "8888"
	}
	go listenAndServe(port)

	select {}
}
