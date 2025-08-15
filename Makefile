BUILD_TARGETS=./.bin
SOURCE=./cmd/server/*.go
TARGET=server

gen-tlscerts:
	@ ${MAKE} -C ./auth/local_certs/ gen_certs

gen-constants:
	@ go generate ./internal/config/codegen/

gen-templ:
	@ templ generate 

build-local: gen-constants gen-templ
	@ mkdir $(BUILD_TARGETS) -p
	@ go build -o $(BUILD_TARGETS)/$(TARGET) $(SOURCE)

run-local: build-local
	@ $(BUILD_TARGETS)/$(TARGET)

clean:
	@ rm -rf $(BUILD_TARGETS)
