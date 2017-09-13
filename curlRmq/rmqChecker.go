package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// Example curl:
// curl --header 'Authorization: Basic Z3Vlc3Q6Z3Vlc3Q=' http://localhost:15672/api/healthchecks/node/rabbit@ip-172-31-27-58

func main() {
	rawHostAddress := flag.String("host", "localhost", "enter a host address or IP") // using fqdn depends on healthy DNS resolution
	rawCurlUserEncrypted := flag.String("auth", "Z3Vlc3Q6Z3Vlc3Q=", "Enter base64-encoded user:password. Default Z3Vlc3Q6Z3Vlc3Q= is guest:guest")
	// rawCheckType := flag.String("check", "", "which check do you want to run?")

	// parse those flags
	flag.Parse()
	// Let's set some variables
	flagAuthContent := *rawCurlUserEncrypted
	hostAddress := *rawHostAddress
	hostName, _ := os.Hostname()
	defaultRmqNodeName := fmt.Sprintf("rabbit@%v", hostName)
	// debug
	fmt.Println("Debug: printing string content of defaultRmqNodeName:", defaultRmqNodeName)

	// create the curl
	// var with the target. param this out later for the check types
	curlTarget := fmt.Sprintf("http://%v:15672/api/healthchecks/node/%v", hostAddress, defaultRmqNodeName)

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
	// defer closing the session
	defer res.Body.Close()

	// get the content and lowercase it as a string
	rawBody, _ := ioutil.ReadAll(res.Body)
	body := strings.ToLower(string(rawBody))
	fmt.Println(body)
}
