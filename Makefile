APP_NAME?=matcher
BUILD_DIR=build

clean:
	rm -f ${BUILD_DIR}/${APP_NAME}

build: clean
	go build -o ${BUILD_DIR}/${APP_NAME} ./...

run: build
	${BUILD_DIR}/${APP_NAME}