build:
	go build
	go build ./opt/hydra-head
	go build ./opt/hydra-watch-battery

install:
	go install .
	go install ./opt/hydra-head
	go install ./opt/hydra-watch-battery
