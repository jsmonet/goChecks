package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	rawPortNumber := flag.String("port", "", "Enter a TCP port to check")
	rawHostAddress := flag.String("host", "localhost", "Enter a host address")
	rawTimeOutSeconds := flag.Int("timeout", 1, "Integer for the timeout of the check in seconds")
	flag.Parse()
	hostAddress := fmt.Sprintf("%v:%v", *rawHostAddress, *rawPortNumber)
	timeOutSeconds := *rawTimeOutSeconds
	timeOut := time.Duration(timeOutSeconds) * time.Second

	fmt.Println(timeOut)
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
