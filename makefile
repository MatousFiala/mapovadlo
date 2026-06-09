build:
	go build -o "bin/mapovadlo"

run: build 
	bin/mapovadlo 

test:
	go test -v ./... -count=1
