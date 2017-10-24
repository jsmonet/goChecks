package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jsmonet/goChecks/validify"
)

var (
	checkType      = flag.String("type", "", "What kind of check?")
	hostAddress    = flag.String("host", "", "Enter a host address")
	portNumber     = flag.Int("port", 22, "Enter a TCP port number")
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
		validify.Port(*portNumber)
	case "neo4j":
		roleIsValid, roleErr := validify.Neorole(*neoRole)
		authIsValid, authErr := validify.Authb64(*authString)
		portIsValid, portErr := validify.Port(*portNumber)
		if roleIsValid && authIsValid && portIsValid {
			fmt.Println("congrats on hitting submit")
		} else {
			fmt.Println(roleErr)
			fmt.Println(authErr)
			fmt.Println(portErr)
			os.Exit(2)
		}
	case "elasticsearch":
		portIsValid, portErr := validify.Port(*portNumber)
		if portIsValid {
			fmt.Println("hooray, it works")
		} else {
			fmt.Println(portErr)
		}
	}
}
