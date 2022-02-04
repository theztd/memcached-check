export BINARY_NAME=memcached-check

help:
	@echo "There are multiple build options"
	@echo ""
	@echo " linux        ... Build binary for limux x86_64"
	@echo " linux-arm    ... Build binary for limux arm64"
	@echo " darwin       ... Build binary for OSX x86_64"
	@echo " darwin-arm   ... Build binary for OSX apple m1"
	@echo ""



linux:
	GOOS=linux GOARCH=amd64 go build -ldflags "-extldflags '-static'" -o build/${BINARY_NAME}-linux-amd64

linux-arm:
	GOOS=linux GOARCH=arm64 go build -ldflags "-extldflags '-static'" -o build/${BINARY_NAME}-linux-arm64

darwin:
	GOOS=darwin GOARCH=amd64 go build -ldflags "-extldflags '-static'" -o build/${BINARY_NAME}-darwin-amd64

darwin-arm:
	GOOS=darwin GOARCH=arm64 go build -ldflags "-extldflags '-static'" -o build/${BINARY_NAME}-darwin-arm64

all: linux linux-arm darwin darwin-arm
	@echo "Done"
	
