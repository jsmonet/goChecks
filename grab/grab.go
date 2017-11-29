package grab

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"syscall"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
)

// CurlAndReturn returns a byte slice from a curl target input
// deprecating in favor of CurlAndReturnJson below
// func CurlAndReturn(target string) []byte {
// 	curlTarget, curlErr := http.NewRequest("GET", target, nil)
// 	if curlErr != nil {
// 		fmt.Println(curlErr) // really I don't care, but you're welcome to uncomment
// 	}
// 	curlRes, _ := http.DefaultClient.Do(curlTarget)
// 	rawCurlBody, _ := ioutil.ReadAll(curlRes.Body)

// 	defer curlRes.Body.Close()

// 	return rawCurlBody
// }

// CurlAndReturnJson returns a byte slice from a curl target input. This must only be used to curl JSON content.
func CurlAndReturnJson(target string) []byte {
	curlTarget, curlErr := http.NewRequest("GET", target, nil)
	if curlErr != nil {
		fmt.Println(curlErr) // really I don't care, but you're welcome to uncomment
	}
	curlTarget.Header.Set("Content-Type", "application/json")
	curlRes, _ := http.DefaultClient.Do(curlTarget)
	rawCurlBody, _ := ioutil.ReadAll(curlRes.Body)

	defer curlRes.Body.Close()

	return rawCurlBody
}

// Authcurl returns a byte slice from a curl with authentication
// deprecating in favor of AuthcurlJson below
// func Authcurl(target string, auth string) []byte {
// 	req, err := http.NewRequest("GET", target, nil)
// 	if err != nil {
// 		fmt.Println("req error:", err)
// 	}
// 	authString := fmt.Sprintf("Basic: %v", auth) // this doesn't jive with how I do neo curl
// 	req.Header.Set("Authorization", authString)
// 	res, _ := http.DefaultClient.Do(req)
// 	rawCurlBody, _ := ioutil.ReadAll(res.Body)

// 	defer res.Body.Close()

// 	return rawCurlBody
// }

// AuthcurlJson explicitly sets Content-Type application/json and returns a byte slice from a curl with authentication. This must only be used to curl JSON content.
func AuthcurlJson(target string, auth string) []byte {
	req, err := http.NewRequest("GET", target, nil)
	if err != nil {
		fmt.Println("req error:", err)
	}
	authString := fmt.Sprintf("Basic: %v", auth)
	req.Header.Set("Authorization", authString)
	req.Header.Set("Content-Type", "application/json")
	res, _ := http.DefaultClient.Do(req)
	rawCurlBody, _ := ioutil.ReadAll(res.Body)

	defer res.Body.Close()

	return rawCurlBody
}

// Checkport returns an int value to get tossed into os.Exit.
// address is the hostname or IP you wish to poll.
// port is the TCP port to be polled.
// timeout is how long, in seconds, to wait before giving up on the poll.
// result is an int to feed os.Exit(result) in the main scheck program.
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
func Rmqjstat(jbody []byte) (result int) {
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

// RmqQueueStat uses authcurl byte slice to check failed queues. I still need to figure out the key to check
func RmqQueueStat(jBody []byte) (messageCount int, result int) {
	type QueueFailedCount struct {
		Count int `json:"messages"` // figure out which of these we need
	}
	result = 0 // explicitly zeroing
	var queueFailedCount QueueFailedCount
	marshalerr := json.Unmarshal(jBody, &queueFailedCount)
	if marshalerr != nil {
		fmt.Println(marshalerr)
	}
	messageCount = queueFailedCount.Count
	if queueFailedCount.Count > 0 && queueFailedCount.Count < 20 {
		result = 1
	} else if queueFailedCount.Count > 19 {
		result = 2
	}
	return messageCount, result
}

// Diskuse returns usage stats based on the volume given
func Diskuse(path string) (cap uint64, used uint64) {

	fs := syscall.Statfs_t{}
	err := syscall.Statfs(path, &fs)
	if err != nil {
		fmt.Println(err)
	}

	cap = fs.Blocks * uint64(fs.Bsize)
	free := fs.Bfree * uint64(fs.Bsize) // yup, I just did that
	used = cap - free
	return cap, used
}

// Procload returns the 1 minute average, 5 minute average as float64's and an int (0, 1, 2) intended for use as an exit code
// seriously, I'm just taking shirou's excellent lib and adding an exit int var to lazy boat the exit code and
// keep the main.go file clean...er
func Procload() (load1 float64, load5 float64, load15 float64, loadPercentage float64, cores int, exit int) {
	exit = 0 // explicit zeroing so much
	loads, _ := load.Avg()
	cores, _ = cpu.Counts(true)
	fCores := float64(cores)
	// basing alerts off 1 minute avg
	load1 = loads.Load1
	// these are here to return for a cooler looking output
	load5 = loads.Load5
	load15 = loads.Load15
	loadPercentage = load1 / fCores
	warnload := 1.1
	critload := 2.5
	if loadPercentage > warnload && loadPercentage < critload {
		exit = 1
	} else if loadPercentage > critload {
		exit = 2
	}
	return load1, load5, load15, loadPercentage, cores, exit
}

// Memload just like Procload just rips off shirou's excellent lib (see imports) for the
// purpose of tacking on an exit code (int variable) and wrapping the %-used val
func Memload() (percentUsed float64, exit int) {
	exit = 0
	warnMem := 80.00
	critMem := 98.00
	memory, _ := mem.VirtualMemory()
	percentUsed = memory.UsedPercent
	if percentUsed > warnMem && percentUsed < critMem {
		exit = 1
	} else if percentUsed > critMem {
		exit = 2
	}
	return percentUsed, exit
}
