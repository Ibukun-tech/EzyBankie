build:
	go build -o bin/EzyBankie
run: build
	./bin/EzyBankie
test:
	go test -v ./..go