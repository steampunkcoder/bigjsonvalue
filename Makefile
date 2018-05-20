test:
	go test -v -cover

lint:
	go vet
	golint -set_exit_status
	gofmt -s -d *.go

