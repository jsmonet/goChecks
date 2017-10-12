package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	// flags
	rawHostAddress := flag.String("host", "127.0.0.1", "Enter an address to hit for this check")
	rawPortNumber := flag.Int("port", 9200, "What port is ES running on?")

	// say, what's in those flags? LETTUCE PARSING
	flag.Parse()

	// if it.IsLegit? { return congratulations }
	if *rawPortNumber < 1 || *rawPortNumber > 65535 {
		fmt.Println("Enter a valid port number. You entered:", *rawPortNumber)
		os.Exit(2)
	}

	// we're going to start off simple, just curling address:port/_cluster/health
	// which gives us a little info that's fairly indicative of general node/cluster health

	curlAddress := fmt.Sprintf("http://%v:%v/_cluster/health", *rawHostAddress, *rawPortNumber)

	// debug
	fmt.Println(curlAddress)

	jsonBody := curlAndReturn(curlAddress)
	// debug output. coding equiv of chimping your display on your dSLR
	fmt.Println(string(jsonBody))

	// set up the struct to parse the JSON
	type Elascheck struct {
		Status           string `json:"status"`
		Numberofnodes    int    `json:"number_of_nodes"`
		Unassignedshards int    `json:"unassigned_shards"`
	}
	// and now we have a way of accessing it once we marshal the data
	var elascheck Elascheck

	marshalerr := json.Unmarshal(jsonBody, &elascheck)
	if marshalerr != nil {
		fmt.Println("You messed up bad:", marshalerr)
	}
	fmt.Println("status output:", elascheck.Status, "number of nodes:", elascheck.Numberofnodes, "and finally unassigned shards:", elascheck.Unassignedshards)
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
