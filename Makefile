all: test

build:
	go build

test:
	go build && ./p-primes

clean:
	rm -f .cache/*.txt .cache/*.zip
