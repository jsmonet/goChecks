package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	// "os"
	"strings"
)

// https://github.com/jmoiron/jsonq maybe consider using that library now that
// I've accomplished the same thing after shuffling it all around slightly

type healthValue struct {
	Status []byte `json:"status"`
}

func getJson(curlTarget string) (hv healthValue) {
	req, err := http.NewRequest("GET", curlTarget, nil)
	if err != nil {
		fmt.Println("Error:", err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error at the defaultclient stage:", err)
	}

	defer res.Body.Close()
	hv.Status, err = ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}
	return
}

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
		curlUri = fmt.Sprintf("/_cluster/health")
	}

	curlTargetLower := fmt.Sprintf("http://%v:9200%v", hostAddress, curlUri)
	// debug: let's see what we've got so far
	// fmt.Println(curlTargetLower)
	healthOutput := getJson(curlTargetLower)

	fmt.Println(string(healthOutput.Status[3:15])) // haahahahaa oh my I can just leave this and be the worst person in history
}
