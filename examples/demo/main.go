package main

import "os"

var (
	serviceToRun = os.Getenv("SERVICE_TO_RUN")
)

func main() {
	if serviceToRun == "publisher" {
		Publisher()
	} else if serviceToRun == "subscriber" {
		Subscriber()
	} else if serviceToRun == "publisherlocal" {
		PublisherLocal()
	} else if serviceToRun == "subscriberlocal" {
		SubscriberLocal()
	}
}
