BUILD_TARGETS=./.bin
SOURCE=./cmd/server/*.go
TARGET=server

gen-certs:
	@ ${MAKE} -C ./auth/local_certs/ gen_certs

gen-shared:
	@ go generate ./...

build-local: gen-shared
	@ templ generate
	@ mkdir $(BUILD_TARGETS) -p
	@ go build -o $(BUILD_TARGETS)/$(TARGET) $(SOURCE)

run-local: build-local
	@ clear
	@ $(BUILD_TARGETS)/$(TARGET)

clean:
	@ rm -rf $(BUILD_TARGETS)
