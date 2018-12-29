build:
	go build
	go build ./opt/hydra-head
	go build ./opt/hydra-get-vol

install:
	go install .
	go install ./opt/hydra-head
	go install ./opt/hydra-get-vol
