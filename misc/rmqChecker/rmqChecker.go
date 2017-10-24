package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/elgs/gojq"
)

// Example curl:
// curl --header 'Authorization: Basic Z3Vlc3Q6Z3Vlc3Q=' http://localhost:15672/api/healthchecks/node/rabbit@ip-169-254-27-13

func main() {
	rawHostAddress := flag.String("host", "localhost", "enter a host address or IP") // using fqdn depends on healthy DNS resolution
	rawPortNumber := flag.Int("port", 15672, "enter a port number for RMQ management")
	rawCurlUserEncrypted := flag.String("auth", "Z3Vlc3Q6Z3Vlc3Q=", "Enter base64-encoded user:password. Default Z3Vlc3Q6Z3Vlc3Q= is guest:guest")
	// rawCheckType := flag.String("check", "", "which check do you want to run?")

	// parse those flags
	flag.Parse()
	// Let's set some variables

	if *rawPortNumber < 1 || *rawPortNumber > 65535 {
		fmt.Println(*rawPortNumber, "is out of range. Pick a port between 1 and 65535 and try again")
		os.Exit(2) // smack. This should panic if you feed it a dumb value
	}
	flagAuthContent := fmt.Sprintf("basic, %v", *rawCurlUserEncrypted)
	hostAddress := *rawHostAddress
	hostName, _ := os.Hostname()
	defaultRmqNodeName := fmt.Sprintf("rabbit@%v", hostName)

	// create the curl
	// var with the target. param this out later for the check types
	curlTarget := fmt.Sprintf("http://%v:%v/api/healthchecks/node/%v", hostAddress, *rawPortNumber, defaultRmqNodeName)

	req, err := http.NewRequest("GET", curlTarget, nil)
	if err != nil {
		// nothing to put here right now
	}
	req.Header.Set("Authorization", flagAuthContent)

	// give it a method to ingest the body
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("something went wrong at the http.Defaultclient.Do(req) stage")

	}

	rawBody, err := ioutil.ReadAll(res.Body)
	body := string(rawBody) // this really should be parsing the json output of this curl instead of being so lazy

	parsed, err := gojq.NewStringQuery(body)
	if err != nil {
		fmt.Println(err)
		return
	}

	rawHealthCheckResult, _ := parsed.Query("status")
	if rawHealthCheckResult == "ok" {
		fmt.Println("OK - all is well")
		os.Exit(0)
	} else {
		fmt.Println("Critical - ", rawHealthCheckResult)
	}
	fmt.Println(rawHealthCheckResult)

	defer res.Body.Close()
}