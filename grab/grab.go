package grab

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

func CurlAndReturn(target string) []byte {
	curlTarget, curlErr := http.NewRequest("GET", target, nil)
	if curlErr != nil {
		fmt.Println(curlErr) // really I don't care, but you're welcome to uncomment
	}
	curlRes, _ := http.DefaultClient.Do(curlTarget)
	rawCurlBody, _ := ioutil.ReadAll(curlRes.Body)

	defer curlRes.Body.Close()

	return rawCurlBody
}

// Checkport returns an int value to get tossed into os.Exit
func Checkport(address string, port int, timeout int) (result int) {
	result = 0 // forcing this to "good"
	target := fmt.Sprintf("%v:%v", address, port)
	timeOutSeconds := time.Duration(timeout) * time.Second
	conn, err := net.DialTimeout("tcp", target, timeOutSeconds)
	if err != nil {
		fmt.Println("Crit -", err)
		result = 2
	}
	conn.Close()
	return result
}
