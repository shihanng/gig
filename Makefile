.PHONY: clean test
build: gi

test:
	go test -v ./... -count=1

integ-test:
	go test -v ./... --tags=integration -count=1
	
gi: 
	go build .

clean:
	rm gi
