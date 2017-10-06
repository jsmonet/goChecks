package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// this curler is currently purpose-built to hit neo4j hosts and figure out
// if the host is master/slave, or up at all.
// I'll expand it's utility just a bit later.
// Also, this is ONLY functional with HA neo4j, not causal clustering or
// singletons at the moment. I'll fix that later. Right now I'm all POC all the time
// like, you know, POC to prod without any changes ftw

func main() {
	rawHostAddress := flag.String("host", "localhost", "enter a host address or IP") // using fqdn relies on faithful DNS resolution
	rawNeo4jRole := flag.String("role", "master", "enter master or slave")           // Neo4j HA cluster has 3 roles: master, slave, arbiter. We won't test for the latter
	rawFlagAuthContent := flag.String("auth", "", "enter JUST the base64-encoded auth content. We will add the rest for you. Example: aaAaaaAAaaaAbbbBb678bV==")

	flag.Parse() // parse those flags!

	hostAddress := strings.ToLower(*rawHostAddress)
	neo4jRole := strings.ToLower(*rawNeo4jRole)
	flagAuthContent := *rawFlagAuthContent
	// disabling debug printing
	// fmt.Println("debug: printing hostAddress", hostAddress)
	// fmt.Println("debug: printing neo4jRole", neo4jRole)
	// fmt.Println("debug: printing flagAuthContent", flagAuthContent)
	// make sure neo4jRole is master or slave
	if neo4jRole != "master" && neo4jRole != "slave" {
		fmt.Println("Please ONLY use master or slave with the -role flag")
		// since this is a Sensu check, I really shouldn't appropriate exit codes it likes for this.
		// os.Exit(1)
	}
	// rawHostAddress/hostAddress cannot be empty
	if hostAddress == "" {
		fmt.Println("Enter an address for your curl target")
		// since this is a Sensu check, I really shouldn't appropriate exit codes it likes for this.
		// os.Exit(1)
	}
	// auth content cannot be empty in this specific case
	if flagAuthContent == "" {
		fmt.Println("You must enter a base64-encoded auth string")
		// since this is a Sensu check, I really shouldn't appropriate exit codes it likes for this.
		// os.Exit(1)
	}

	curlTarget := fmt.Sprintf("http://%v:7474/db/manage/server/ha/%v", hostAddress, neo4jRole)
	// disabling debug printing
	// fmt.Println(curlTarget)

	// create the curl
	req, err := http.NewRequest("GET", curlTarget, nil)
	if err != nil {
		os.Exit(1)
	}
	// add in authorization
	req.Header.Set("Authorization", flagAuthContent)
	// give it a method to ingest the response
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		// I seriously don't see a point to a handler here right now
	}
	// don't forget to eventually close the session
	defer res.Body.Close()
	// read the body into 'body'
	rawBody, _ := ioutil.ReadAll(res.Body)
	body := strings.ToLower(string(rawBody))
	// disabling debug printing
	// fmt.Println("debug: printing contents of 'body' variable", body)

	if body != "true" {
		fmt.Println("Warning: Neo4j Server currently set to wrong role. Role should be", neo4jRole)
		os.Exit(1)
	} else {
		fmt.Println("OK - host matches expected role of", neo4jRole)
		os.Exit(0)
	}
}
