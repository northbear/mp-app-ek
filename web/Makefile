APP_VERSION    	:= 0.0.1

DOCKER		:= docker
APP_NAME  	?= mp-app-ek
IMAGE_NAME     	:= $(APP_NAME)/web

EXECUTABLE      := build/web

.PHONY: prepare clean all build-image

all: prepare $(EXECUTABLE)

prepare:
	mkdir -p build

$(EXECUTABLE): main.go
	go build -o $(EXECUTABLE) -ldflags="-X main.APP_VERSION=$(APP_VERSION)" $<

build-image:
	$(DOCKER) build --rm --no-cache -t $(IMAGE_NAME):$(APP_VERSION) .

clean:
	rm -rf build
	$(DOCKER) rmi --force $(IMAGE_NAME):$(APP_VERSION)
