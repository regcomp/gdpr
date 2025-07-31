BUILD_TARGETS=./.bin
SOURCE=./cmd/server/*.go
TARGET=server

gen-certs:
	@ ${MAKE} -C ./auth/local_certs/ gen_certs

gen-config:
	@ go generate ./config/codegen/

build-local: gen-variables
	@ templ generate
	@ mkdir $(BUILD_TARGETS) -p
	@ go build -o $(BUILD_TARGETS)/$(TARGET) $(SOURCE)

run-local: build-local
	@ $(BUILD_TARGETS)/$(TARGET)

clean:
	@ rm -rf $(BUILD_TARGETS)
