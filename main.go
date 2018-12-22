package main

func main() {
	pactl_events, err := pactlEvents()
	if err != nil {
		panic(err)
	}
	nmcli_events, err := nmcliEvents()
	if err != nil {
		panic(err)
	}
	events := make(chan string)
	go func() {
		for {
			select {
			case <-pactl_events:
				events <- "pactl"
			case <-nmcli_events:
				events <- "nmcli"
			}
		}
	}()
	server(events)
}
