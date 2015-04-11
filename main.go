package main

import (
	"fmt"
	"flag"
	"log"
	"runtime"
	"runtime/pprof"
	"time"
	"os"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	start := time.Now()

	var cpuprofile string
	flag.StringVar(&cpuprofile, "prof", "", "Write CPU Profile to file")
	flag.Parse()

	if cpuprofile != "" {
		f, err := os.Create(cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	log.Println(getPrimes(1))

	fmt.Println(time.Since(start))
}
