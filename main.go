package main

import (
	"bufio"
	"crypto/tls"
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

const (
	CACHEDIR = ".cache"
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

	flag.Parse()

	if cpuprofile != "" {
		f, err := os.Create(cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	ns := make(chan int)
	downloader(N, ns)

	for i := 1; i <= N; i++ {
		n := <-ns
		primes := atoiPipe(primeGenerator(n))
		x := 0
		for _ = range primes {
			x++
		}
		log.Println(n, "=>", x)
	}

	fmt.Println(time.Since(start))
}

func downloader(n int, out chan int) {
	for i := 1; i <= n; i++ {
		go func(i int) {
			err := downloadFiles(i)
			if err != nil {
				log.Fatal("Error downloading - "+strconv.Itoa(i), err)
			}
			out <- i
		}(i)
	}
}

func primeGenerator(n int) <-chan string {
	out := make(chan string)
	go func() {
		filepath := CACHEDIR + "/primes" + strconv.Itoa(n) + ".txt"
		primesFile, err := os.Open(filepath)
		if err != nil {
			log.Fatal(err)
		}
		defer primesFile.Close()

		s := bufio.NewScanner(primesFile)
		// skip first line - contains header
		if !s.Scan() {
			log.Fatal("No content found in file - ", filepath)
		}
		for s.Scan() {
			fields := strings.Fields(s.Text())
			for _, f := range fields {
				out <- f
			}
		}
		close(out)
	}()
	return out
}

func atoiPipe(in <-chan string) <-chan int {
	out := make(chan int)
	go func() {
		for str := range in {
			val, err := strconv.Atoi(str)
			if err != nil {
				log.Fatal(err)
			}
			out <- val
		}
		close(out)
	}()
	return out
}

func downloadFiles(n int) error {
	nPrime := strconv.Itoa(n)

	log.Println("download " + nPrime + " <start>")
	defer log.Println("download " + nPrime + " <end>")

	if _, err := os.Stat(CACHEDIR); err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(CACHEDIR, 0755)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	if _, err := os.Stat(CACHEDIR + "/primes" + nPrime + ".zip"); err != nil {
		if os.IsNotExist(err) {
			out, err := os.Create(CACHEDIR + "/primes" + strconv.Itoa(n) + ".zip")
			if err != nil {
				return err
			}
			url := "https://primes.utm.edu/lists/small/millions/primes" + strconv.Itoa(n) + ".zip"
			tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
			client := &http.Client{Transport: tr}
			resp, err := client.Get(url)
			if err != nil {
				return err
			}
			io.Copy(out, resp.Body)
			out.Close()
			resp.Body.Close()
		} else {
			return err
		}
	} else {
		log.Println("download " + nPrime + " <msg> primes" + nPrime + ".zip already exists, skipping download")
	}

	if _, err := os.Stat(CACHEDIR + "/primes" + nPrime + ".txt"); err != nil {
		if os.IsNotExist(err) {
			unzip := exec.Command("unzip", "primes"+strconv.Itoa(n)+".zip")
			unzip.Dir = CACHEDIR
			if err := unzip.Run(); err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		log.Println("download " + nPrime + " <msg> primes" + nPrime + ".txt already exists, skipping unzip")
	}

	return nil
}
