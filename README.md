### What is this?

A small experiment with Go.

### Impatient coder

```
go get && make
```

### Usage

```
# Build
go get
go build

# n can take values from 1 to 50. Default = 1
./p-primes -n 50

# prof is for cpu profiling
./p-primes -prof=./cpuprofile

# To clear cache
make clean
```

### What can you hack with it ?

This has access to the first 50 MILLION prime numbers

### from where ?

Here - https://primes.utm.edu/lists/small/millions/

### Why download ? Generating would be faster

I mean, why not ? this is just some weird experiment
