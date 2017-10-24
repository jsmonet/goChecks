package grab

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

// CurlAndReturn returns a byte slice from a curl target input
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

// Authcurl returns a byte slice from a curl with authentication
func Authcurl(target string, auth string) []byte {
	req, err := http.NewRequest("GET", target, nil)
	if err != nil {
		fmt.Println("req error:", err)
	}
	authString := fmt.Sprintf("Basic: %v", auth) // this doesn't jive with how I do neo curl
	req.Header.Set("Authorization", authString)
	res, _ := http.DefaultClient.Do(req)
	rawCurlBody, _ := ioutil.ReadAll(res.Body)

	defer res.Body.Close()

	return rawCurlBody
}

// Checkport returns an int value to get tossed into os.Exit
func Checkport(address string, port int, timeout int) (result int) {
	result = 0 // explicitly zeroing
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
	result = 0 // explicitly zeroing
	target := fmt.Sprintf("http://%v:7474/db/manage/server/ha/%v", address, role)
	req, _ := http.NewRequest("GET", target, nil) // I probably should catch an error here, but I don't really care about it right now
	req.Header.Set("Authorization", auth)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()

	rawBody, _ := ioutil.ReadAll(res.Body)
	body := strings.ToLower(string(rawBody))
	if body != "true" {
		result = 2
	}
	return result
}

// Elasjson returns a result int and some strings
func Elasjson(jbody []byte) (result int, status string, nodes int, unshards int) {
	type Elascheck struct {
		Stat     string `json:"status"`
		Numnodes int    `json:"number_of_nodes"`
		Unshards int    `json:"unassigned_shards"`
	}
	result = 0 // explicitly zeroing
	var elascheck Elascheck
	marshalerr := json.Unmarshal(jbody, &elascheck)
	if marshalerr != nil {
		result = 2
	}
	status = strings.ToLower(elascheck.Stat)
	nodes = elascheck.Numnodes
	unshards = elascheck.Unshards
	if strings.ToLower(elascheck.Stat) == "yellow" {
		result = 1
	} else if strings.ToLower(elascheck.Stat) == "red" {
		result = 2
	}
	return result, status, nodes, unshards
}

// Rmqjson uses authenticated curl to return a result int based on the value of key "status"
func Rmqjson(jbody []byte) (result int) {
	type Rmqcheck struct {
		Stat string `json:"status"`
	}
	result = 0 // explicitly zeroing
	var rmqcheck Rmqcheck
	marshalerr := json.Unmarshal(jbody, &rmqcheck)
	if marshalerr != nil {
		fmt.Println(marshalerr)
	}
	if strings.ToLower(rmqcheck.Stat) != "ok" {
		result = 2
	}
	return result
}
