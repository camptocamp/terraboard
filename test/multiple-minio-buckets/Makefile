
UID=$(shell id -u)
GID=$(shell id -g)

build:
	UID="${UID}" GID="${GID}" docker-compose build 

test:
	UID="${UID}" GID="${GID}" docker-compose up 

all: build test
