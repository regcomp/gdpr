BUILD_TARGETS=./.bin
SOURCE=./cmd/server/*.go
TARGET=server

build-api:
	@ mkdir $(BUILD_TARGETS) -p
	@ go build -o $(BUILD_TARGETS)/$(TARGET) $(SOURCE)

run-api: build-api
	@ $(BUILD_TARGETS)/$(TARGET)

clean:
	@ rm -rf $(BUILD_TARGETS)
