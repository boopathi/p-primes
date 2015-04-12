package main

import (
	"bufio"
	"crypto/tls"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"runtime/pprof"
	"strconv"
	"strings"
	"time"
)

var (
	cpuprofile string
	N          int
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.StringVar(&cpuprofile, "prof", "", "Write CPU Profile to file")
	flag.IntVar(&N, "n", 1, "Filenumber to download")
}

func main() {
	start := time.Now()
	defer fmt.Println(time.Since(start))

	flag.Parse()

	if cpuprofile != "" {
		f, err := os.Create(cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if err := downloadFiles(N); err != nil {
		log.Fatal(err)
	}

	primes, err := getPrimes(N)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(len(primes))

}

func downloadFiles(n int) error {
	log.Println("download <start>")
	defer log.Println("download <end>")

	out, err := os.Create("primes" + strconv.Itoa(n) + ".zip")
	if err != nil {
		return err
	}
	defer out.Close()

	url := "https://primes.utm.edu/lists/small/millions/primes" + strconv.Itoa(n) + ".zip"
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(url)
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	io.Copy(out, resp.Body)

	unzip := exec.Command("unzip", "primes"+strconv.Itoa(n)+".zip")
	if err := unzip.Run(); err != nil {
		return err
	}

	return nil
}

func getPrimes(n int) ([]string, error) {
	primesFile, err := os.Open("primes" + strconv.Itoa(n) + ".txt")
	if err != nil {
		return nil, err
	}
	defer primesFile.Close()

	s := bufio.NewScanner(primesFile)
	// skip first line - contains header
	if !s.Scan() {
		return nil, errors.New("Nothing to read")
	}

	primes := make([]string, 0)

	for s.Scan() {
		fields := strings.Fields(s.Text())
		primes = append(primes, fields...)
	}

	return primes, nil
}
