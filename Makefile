build:
	go build
	go build ./opt/hydra-head

install:
	go install .
	go install ./opt/hydra-head
