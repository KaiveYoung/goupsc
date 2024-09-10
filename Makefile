TARGET="./release/goupsc"
all: build

build: app

app:
	echo "Building ...";\
	mkdir -p release
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags -s -o ${TARGET}  ./cmd/goupsc
	echo "Build done";\

