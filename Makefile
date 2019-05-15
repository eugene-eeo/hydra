INSTALL_FLAGS = -ldflags='-s -w'

build:
	go build
	go build ./opt/hydra-head
	go build ./opt/hydra-watch-battery
	go build ./opt/hydra-timer
	go build ./opt/hydra-emit

install:
	go install ${INSTALL_FLAGS} .
	go install ${INSTALL_FLAGS} ./opt/hydra-head
	go install ${INSTALL_FLAGS} ./opt/hydra-watch-battery
	go install ${INSTALL_FLAGS} ./opt/hydra-timer
	go install ${INSTALL_FLAGS} ./opt/hydra-emit
