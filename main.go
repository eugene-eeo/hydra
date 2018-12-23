package main

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	events := make(chan string, 5)
	must(pactlEvents(events))
	must(nmcliEvents(events))
	server(events)
}
