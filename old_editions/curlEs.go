package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func main() {
	rawCheckType := flag.String("type", "health", "Check type. Possible values are health, and nothing else right now")
	rawHostAddress := flag.String("host", "127.0.0.1", "IP or url to the host you wish to check.")

	// for great justice, parse off all flags
	flag.Parse()
	// debug: print them things
	// fmt.Println(*rawCheckType, *rawHostAddress)

	// just in case this isn't an IP, let's lowercase it
	hostAddress := strings.ToLower(*rawHostAddress)
	// this is pure laziness on my part in not wanting to case every possible way of writing these
	checkType := strings.ToLower(*rawCheckType)

	// Let's instantiate a variable or two
	var curlUri string

	if checkType == "health" {
		// _cat/health returns a simple one-liner with color-based
		// health representation and some metrics
		// way easier to mess with than json right now
		curlUri = fmt.Sprintf("/_cat/health")
	}

	curlTarget := fmt.Sprintf("http://%v:9200%v", hostAddress, curlUri)
	// debug: let's see what we've got so far
	// fmt.Println(curlTarget)

	req, err := http.NewRequest("GET", curlTarget, nil)
	if err != nil {
		fmt.Println("something went wrong", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("error:", err)
	}

	defer res.Body.Close() // just so I don't forget to close it

	rawBody, _ := ioutil.ReadAll(res.Body)
	body := strings.ToLower(string(rawBody))
	bodyWords := strings.Fields(body)

	// debug: pring out body
	// fmt.Println(body)
	// fmt.Println(bodyWords[3])
	if bodyWords[3] == "green" {
		fmt.Println("OK, the cluster health is", bodyWords[3])
		os.Exit(0)
	} else if bodyWords[3] == "yellow" {
		fmt.Println("Warning, cluster health is", bodyWords[3])
	} else if bodyWords[3] == "red" {
		fmt.Println("Critical, cluster health is,", bodyWords[3])
	}

}
