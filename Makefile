all: clean test

build:
	go build

test:
	go build && ./p-primes

clean:
	rm -f *.txt *.zip
