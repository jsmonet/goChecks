package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jsmonet/goChecks/grab"

	"github.com/jsmonet/goChecks/validify"
)

var (
	checkType      = flag.String("type", "", "What kind of check?")
	hostAddress    = flag.String("host", "", "Enter a host address")
	portNumber     = flag.Int("port", 0, "Enter a TCP port number. Leaving this out will throw an exception on checks requiring ports")
	timeOutSeconds = flag.Int("timeout", 5, "Enter a timeout in seconds")
	volumeLocation = flag.String("vol", "/", "Enter a volume to check (disk space)")
	volSizeWarn    = flag.Float64("volwarn", 75, "Percentage full that triggers a warning")
	volSizeCrit    = flag.Float64("volcrit", 90, "Percentage full that triggers a critical")
	authString     = flag.String("auth", "Z3Vlc3Q6Z3Vlc3Q=", "Curl Auth string. Default parses to guest:guest")
	neoRole        = flag.String("role", "", "Neo4j Role: master or slave")
)

func main() {
	// Parse for great justice
	flag.Parse()
	switch *checkType {
	case "port":
		portIsValid, portErr := validify.Port(*portNumber)
		if !portIsValid {
			panic(portErr)
		}
		portResult := grab.Checkport(*hostAddress, *portNumber, *timeOutSeconds)
		if portResult == 0 {
			fmt.Println("OK - port is open")
		} else {
			fmt.Println("Crit - port is closed or operation timed out")
		}
		os.Exit(portResult)
	case "neo4j":
		roleIsValid, roleErr := validify.Neorole(*neoRole)
		authIsValid, authErr := validify.Authb64(*authString)
		var neoIsUp int
		if roleIsValid && authIsValid {
			neoIsUp = grab.Checkneo(*hostAddress, *neoRole, *authString)
		} else {
			errorContent := fmt.Sprintf("\n%v\n%v", roleErr, authErr)
			panic(errorContent)
		}
		if neoIsUp == 0 {
			fmt.Println("OK - Neo4j role is", *neoRole, "as expected")
		} else {
			fmt.Println("Critical - Neo4j role is wrong")
		}
		os.Exit(neoIsUp)
	case "elasticsearch":
		portIsValid, portErr := validify.Port(*portNumber)
		var esIsUp int
		if !portIsValid {
			panic(portErr)
		}
		esTarget := fmt.Sprintf("http://%v:%v/_cluster/health", *hostAddress, *portNumber)
		esBody := grab.CurlAndReturn(esTarget)
		esIsUp, esStatus, esNodes, esUnshards := grab.Elasjson(esBody)
		esStatOutput := fmt.Sprintf("status: %v, nodes: %v, unassigned shards: %v", esStatus, esNodes, esUnshards)
		if esIsUp == 2 {
			fmt.Println("Crit - warnings are", esStatOutput)
		}
		fmt.Println("OK -", esStatOutput)
		os.Exit(esIsUp)
	}
}
