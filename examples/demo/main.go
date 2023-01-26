package main

import "os"

var (
	serviceToRun = os.Getenv("SERVICE_TO_RUN")
)

func main() {
	if serviceToRun == "publisher" {
		Publisher()
	} else {
		Subscriber()
	}
}
