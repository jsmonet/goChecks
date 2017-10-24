package grab

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
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
	result = 0 //explicitly zeroing
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

// Checkneo returns a result int
func Checkneo(address string, role string, auth string) (result int) {
	result = 0 //explicitly zeroing
	target := fmt.Sprintf("http://%v:7474/db/manage/server/ha/%v", address, auth)
	req, _ := http.NewRequest("GET", target, nil) // I probably should catch an error here, but I don't really care about it right now
	req.Header.Set("Authorization", auth)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		result = 2
	}
	defer res.Body.Close()

	rawBody, _ := ioutil.ReadAll(res.Body)
	body := strings.ToLower(string(rawBody))
	if body != "true" {
		result = 2
	}
	return result
}
