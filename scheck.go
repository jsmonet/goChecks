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
	volSizeWarn    = flag.Int("volwarn", 75, "Percentage full that triggers a warning")
	volSizeCrit    = flag.Int("volcrit", 90, "Percentage full that triggers a critical")
	authString     = flag.String("auth", "Z3Vlc3Q6Z3Vlc3Q=", "Curl Auth string. Default parses to guest:guest")
	neoRole        = flag.String("role", "", "Neo4j Role: master or slave")
	rmqNodeName    = flag.String("rmqname", "", "node name for RMQ curls")
	rmqQueueName   = flag.String("queue", "", "Which rabbit queue do you want to check?")
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
		if !roleIsValid || !authIsValid {
			errorContent := fmt.Sprintf("\n%v\n%v", roleErr, authErr)
			panic(errorContent)
		}
		neoIsUp = grab.Checkneo(*hostAddress, *neoRole, *authString)
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
		esBody := grab.CurlAndReturnJson(esTarget)
		esIsUp, esStatus, esNodes, esUnshards := grab.Elasjson(esBody)
		esStatOutput := fmt.Sprintf("status: %v, nodes: %v, unassigned shards: %v", esStatus, esNodes, esUnshards)
		if esIsUp == 2 {
			fmt.Println("Crit - warnings are", esStatOutput)
		} else {
			fmt.Println("OK -", esStatOutput)
		}
		os.Exit(esIsUp)
	case "rmq":
		portIsValid, portErr := validify.Port(*portNumber)
		authIsValid, authErr := validify.Authb64(*authString)
		var rmqIsUp int
		if !portIsValid || !authIsValid {
			errorContent := fmt.Sprintf("\n%v\n%v", portErr, authErr)
			panic(errorContent)
		}
		hostName, _ := os.Hostname()
		var defaultRmqNodeName string
		if len(*rmqNodeName) == 0 {
			defaultRmqNodeName = fmt.Sprintf("rabbit@%v", hostName)
		} else {
			defaultRmqNodeName = fmt.Sprintf("rabbit@%v", *rmqNodeName)
		}
		rmqTarget := fmt.Sprintf("http://%v:%v/api/healthchecks/node/%v", *hostAddress, *portNumber, defaultRmqNodeName)
		rBody := grab.AuthcurlJson(rmqTarget, *authString)
		rmqIsUp = grab.Rmqjstat(rBody)
		if rmqIsUp != 0 {
			fmt.Println("Crit - rmq status not ok")
		} else {
			fmt.Println("OK - rmq status ok")
		}
		os.Exit(rmqIsUp)
	case "rmqconsumers":
		// commenting the validation out for now because technically I don't need it
		// and technically it doesn't work. not even a little. I'll fix it next.
		// queueExists, queueExistsErr := validify.RmqQueueExists(*rmqQueueName)
		// if !queueExists {
		// 	queueIsOk = 2
		// 	errorStatement := fmt.Sprintf("Queue doesn't exist. Error output:\n%v", queueExistsErr)
		// 	fmt.Println(errorStatement)
		// 	os.Exit(queueIsOk)
		// }
		var queueCount, queueIsOk int
		queueTarget := fmt.Sprintf("http://%v:%v/api/queues/%%2F/%v", *hostAddress, *portNumber, *rmqQueueName)
		rqBody := grab.AuthcurlJson(queueTarget, *authString)
		queueCount, queueIsOk = grab.RmqQueueStat(rqBody)
		if queueIsOk == 1 {
			fmt.Println("Warn -", *rmqQueueName, "failed count is", queueCount)
			os.Exit(queueIsOk)
		} else if queueIsOk == 2 {
			fmt.Println("Crit -", *rmqQueueName, "failed count is", queueCount)
		}
		fmt.Println("OK -", *rmqQueueName, "failed count is", queueCount)
		os.Exit(queueIsOk)
	case "disk":
		percentIsValid, percentErr := validify.Percentages(*volSizeWarn, *volSizeCrit)
		if !percentIsValid {
			fmt.Println(percentErr)
			os.Exit(2)
		}
		diskCapacity, diskUsed := grab.Diskuse(*volumeLocation)
		percentUsed := float64(diskUsed) / float64(diskCapacity) * 100
		perUsedString := fmt.Sprintf("%.2f%% used", percentUsed)
		if percentUsed < float64(*volSizeWarn) {
			fmt.Println("Ok -", perUsedString)
			os.Exit(0)
		} else if percentUsed > float64(*volSizeWarn) && percentUsed < float64(*volSizeCrit) {
			fmt.Println("Warn -", perUsedString)
			os.Exit(1)
		} else {
			fmt.Println("Crit -", perUsedString)
			os.Exit(2)
		}
	case "load":
		cpuLoad1, cpuLoad5, cpuLoad15, smoothedLoad, numCores, cpuIsOk := grab.Procload()
		loadString := fmt.Sprintf("Load: %.3v, %.3v, %.3v, with metric based on avg load %.3v across %v cores", cpuLoad1, cpuLoad5, cpuLoad15, smoothedLoad, numCores)
		var statString string
		if cpuIsOk == 0 {
			statString = fmt.Sprintf("OK -")
		} else if cpuIsOk == 1 {
			statString = fmt.Sprintf("Warn -")
		} else if cpuIsOk == 2 {
			statString = fmt.Sprintf("Crit -")
		}
		outString := fmt.Sprintf("%v %v", statString, loadString)
		fmt.Println(outString)
		os.Exit(cpuIsOk)
	case "mem":
		memPercentUsed, memIsOk := grab.Memload()
		var statString string
		memPercentUsedString := fmt.Sprintf("%.4v%% memory used", memPercentUsed)
		if memIsOk == 0 {
			statString = fmt.Sprintf("OK -")
		} else if memIsOk == 1 {
			statString = fmt.Sprintf("Warn -")
		} else if memIsOk == 2 {
			statString = fmt.Sprintf("Crit -")
		}
		fmt.Println(statString, memPercentUsedString)
		os.Exit(memIsOk)
	}
}
