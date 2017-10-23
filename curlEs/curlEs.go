package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func main() {
	// flags
	rawHostAddress := flag.String("host", "127.0.0.1", "Enter an address to hit for this check")
	rawPortNumber := flag.Int("port", 9200, "What port is ES running on?")

	// say, what's in those flags? LETTUCE PARSING
	flag.Parse()

	// if it.IsLegit? { return congratulations }
	validPort, err := validatePort(*rawPortNumber)
	if !validPort {
		fmt.Println("Try using a valid port number between 1 and 65535. Thrown error:", err)
		os.Exit(2)
	}

	// we're going to start off simple, just curling address:port/_cluster/health
	// which gives us a little info that's fairly indicative of general node/cluster health

	curlAddress := fmt.Sprintf("http://%v:%v/_cluster/health", *rawHostAddress, *rawPortNumber)

	jsonBody := curlAndReturn(curlAddress)

	// set up the struct to parse the JSON
	// changing the left value arbitrarily just to demonstrate to myself
	// how unbound that is to the `json:"fieldname"` value. NBD, this should have been obv to me
	type Elascheck struct {
		Stat     string `json:"status"`
		Numnodes int    `json:"number_of_nodes"`
		Unshards int    `json:"unassigned_shards"`
	}
	// and now we have a way of accessing it once we marshal the data
	var elascheck Elascheck

	marshalerr := json.Unmarshal(jsonBody, &elascheck)
	if marshalerr != nil {
		fmt.Println("You messed up bad:", marshalerr)
	}
	fmt.Println("status output:", elascheck.Stat)
	fmt.Println("number of nodes:", elascheck.Numnodes)
	fmt.Println("and finally unassigned shards:", elascheck.Unshards)

	if strings.ToLower(elascheck.Stat) == "green" && elascheck.Unshards == 0 {
		fmt.Println("OK - all is clear")
		os.Exit(0)
	} else {
		fmt.Println("Critical - something is messed up") // some overcomplicated if/else if or case statement here
		os.Exit(2)
	}
}

func curlAndReturn(target string) []byte {
	esCurlReq, esCurlErr := http.NewRequest("GET", target, nil)
	if esCurlErr != nil {
		fmt.Println(esCurlErr) // really I don't care, but you're welcome to uncomment
	}
	esCurlRes, _ := http.DefaultClient.Do(esCurlReq)
	rawEsCurlBody, _ := ioutil.ReadAll(esCurlRes.Body)

	defer esCurlRes.Body.Close()

	return rawEsCurlBody
}

func validatePort(port int) (bool, error) {
	if port < 1 || port > 65535 {
		return false, errors.New("Port out of range")

	}

	return true, errors.New("no errors")

}
