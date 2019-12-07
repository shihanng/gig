.PHONY: clean test
build: gi

test:
	go test -v ./... -count=1

integ-test:
	go test -v ./... --tags=integration -count=1

lint:
	docker run --rm -v $$(pwd):/app -w /app golangci/golangci-lint:latest golangci-lint run -v
	
gi: 
	go build .

clean:
	rm gi
