package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	rawPortNumber := flag.Int("port", 22, "Enter a TCP port to check")
	rawHostAddress := flag.String("host", "localhost", "Enter a host address")
	rawTimeOutSeconds := flag.Int("timeout", 1, "Integer for the timeout of the check in seconds")
	flag.Parse()

	if *rawPortNumber < 1 || *rawPortNumber > 65535 {
		fmt.Println(*rawPortNumber, "is out of range. Please try with a port number between 1 and 65535")
		os.Exit(2) // let's really piss off sensu if you choose a bad port
	}

	hostAddress := fmt.Sprintf("%v:%v", *rawHostAddress, *rawPortNumber)
	timeOutSeconds := *rawTimeOutSeconds
	timeOut := time.Duration(timeOutSeconds) * time.Second // this parses to '1s', but you can't put that in for the timeout duration. You can actually just use var*time.Seconds to get the proper value here, so I may simplify

	conn, err := net.DialTimeout("tcp", hostAddress, timeOut)
	if err != nil {
		fmt.Println("Critical - connection likely refused. See error:", err)
		os.Exit(2)
	} else {
		fmt.Println("OK - Successfully connected to", hostAddress)
		os.Exit(0)
	}

	conn.Close()

}
