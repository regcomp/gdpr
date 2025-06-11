BUILD_TARGETS=./.bin
SOURCE=./cmd/server/*.go
TARGET=server

build-local:
	@ templ generate
	@ mkdir $(BUILD_TARGETS) -p
	@ go build -o $(BUILD_TARGETS)/$(TARGET) $(SOURCE)

run-local: build-local
	@ $(BUILD_TARGETS)/$(TARGET)

clean:
	@ rm -rf $(BUILD_TARGETS)
