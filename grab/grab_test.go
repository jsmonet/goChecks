package grab

import (
	"encoding/json"
	"fmt"
	"net"
	"testing"

	"gopkg.in/jarcoal/httpmock.v1"
)

func TestCurlAndReturnJson(t *testing.T) {

	cases := []struct {
		targetUrl      string
		jsonResponse   string
		expectedStatus string
	}{
		{
			targetUrl:      "http://localhost/health.json",
			jsonResponse:   `{"status": "yellow", "number_of_nodes": 14, "unassigned_shards": 2}`,
			expectedStatus: "yellow",
		},
		{
			targetUrl: "http://localhost/health.json",
			jsonResponse: `{"status": "green",	"number_of_nodes": 3,	"unassigned_shards": 0}`,
			expectedStatus: "green",
		},
	}

	for _, c := range cases {
		// create the elasticsearch struct. Ok, just steal it from func Elasjson
		type Elastest struct {
			Stat     string `json:"status"`
			Numnodes int    `json:"number_of_nodes"`
			Unshards int    `json:"unassigned_shards"`
		}
		var elastest Elastest

		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		// Let's make an elasticsearch endpoint
		httpmock.RegisterResponder("GET", "http://localhost/health.json", httpmock.NewStringResponder(200, c.jsonResponse))

		jsonBodyToBeUnmarshalled := CurlAndReturnJson(c.targetUrl)

		marshallerr := json.Unmarshal(jsonBodyToBeUnmarshalled, &elastest)
		if marshallerr != nil {
			fmt.Println(marshallerr)
		} // not sure I care about this error yet

		if elastest.Stat != c.expectedStatus {
			t.Error("Expected %v but got %v", c.expectedStatus, elastest.Stat)
		}
	}
}

func TestCheckport(t *testing.T) {
	cases := []struct {
		port           int
		doesItWork     bool
		expectedResult int
	}{
		{
			port:           45932,
			doesItWork:     true,
			expectedResult: 0,
		},
		{
			port:           40025,
			doesItWork:     true, // to-do: handle the panic in grab.Checkport when you feed it an incorrect port
			expectedResult: 0,
		},
	}

	for _, c := range cases {
		var testPort int
		testPort = 45000 // just understand you aren't going to use this port in a test case
		var testURL string
		if c.doesItWork {
			testURL = fmt.Sprintf("localhost:%v", c.port)
		} else {
			testURL = fmt.Sprintf("localhost:%v", testPort)
		}

		// testRequestURL := fmt.Sprintf("http://localhost:%v", c.port)
		fmt.Println(testURL)
		l, err := net.Listen("tcp", testURL)
		if err != nil {
			t.Fatal(err)
		}
		resultCheckPort := Checkport("localhost", c.port, 3)
		if resultCheckPort != c.expectedResult {
			t.Error("You messed up. resultCheckPort is", resultCheckPort, "and expected is", c.expectedResult, "for port", c.port, "and testport is ", testPort, "and doesitwork is", c.doesItWork)
		}
		l.Close()
	}
}

// func TestAuthcurlJson(t *testing.T) {
// 	// i'll hit this in a bit
// }

// func TestCheckneo(t *testing.T) {
// 	cases := struct {
// 		role string
// 		addressRole string
// 		boolStringResponse string
// 		expectedResult int
// 	}{
// 		{
// 			role: "master",
// 			addressRole: "master",
// 			boolStringResponse: "true",
// 			expectedResult: 0,
// 		},
// 		{
// 			role: "slave",
// 			addressRole: "master",
// 			boolStringResponse: "false",
// 			expectedResult: 2,
// 		},
// 		{
// 			role: "master",
// 			addressRole: "slave",
// 			boolStringResponse: "false",
// 			expectedResult: 2,
// 		},
// 		{
// 			role: "slave",
// 			addressRole: "slave",
// 			boolStringResponse: "true",
// 			expectedResult: 0,
// 		},
// 	}
// 	// add in address and auth as static values of some kind
// 	for _, c in range cases {
// 		auth := "bmVvNGo6bmVvNGo=" // neo4j:neo4j
// 		target := fmt.Sprintf("http://localhost:7474/db/manage/server/ha/%v", c.addressRole)
// 		httpmock.Activate()
// 		defer httpmock.DeactivateAndReset()
// 		// make a simple string endpoint that says true at 'master'
// 		// Neo4j HA clustering shows "true" when you go to http://localhost:7474/db/manage/server/ha/master and localhost, in this case, is the HA master
// 		httpmock.RegisterResponder("GET", target, httpmock.NewStringResponder(200, c.boolStringResponse))

// 	}

// }
