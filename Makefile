.PHONY: clean
build: gi
	
gi: 
	go build .

.PHONY: clean
clean:
	rm gi
