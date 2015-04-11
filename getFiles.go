package main

import (
	_ "os/exec"
	"net/http"
	"io"
	"strconv"
	"crypto/tls"
	"os"
)

func getPrimes(n int) (string, error) {

	out, err := os.Create("primes" + strconv.Itoa(n) + ".zip")
	if err != nil {
		return "", err
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
		return "", err
	}
	io.Copy(out, resp.Body)
	return "", nil
}