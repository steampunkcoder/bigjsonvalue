test:
	go test -v -cover

bench:
	go test -v -bench=.

lint:
	go vet
	golint -set_exit_status
	gofmt -s -d *.go

