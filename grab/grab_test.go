package grab

import (
	"encoding/json"
	"fmt"
	"testing"

	"gopkg.in/jarcoal/httpmock.v1"
)

// deprecating curlandreturn 
// func TestCurlAndReturn(t *testing.T) {

// 	cases := []struct {
// 		targetUrl      string
// 		jsonResponse   string
// 		expectedStatus string
// 	}{
// 		{
// 			targetUrl:      "http://localhost/health.json",
// 			jsonResponse:   `{"status": "yellow", "number_of_nodes": 14, "unassigned_shards": 2}`,
// 			expectedStatus: "yellow",
// 		},
// 		{
// 			targetUrl: "http://localhost/health.json",
// 			jsonResponse: `{"status": "green",	"number_of_nodes": 3,	"unassigned_shards": 0}`,
// 			expectedStatus: "green",
// 		},
// 	}

// 	for _, c := range cases {
// 		// create the elasticsearch struct. Ok, just steal it from func Elasjson
// 		type Elastest struct {
// 			Stat     string `json:"status"`
// 			Numnodes int    `json:"number_of_nodes"`
// 			Unshards int    `json:"unassigned_shards"`
// 		}
// 		var elastest Elastest

// 		httpmock.Activate()
// 		defer httpmock.DeactivateAndReset()
// 		// Let's make an elasticsearch endpoint
// 		httpmock.RegisterResponder("GET", "http://localhost/health.json", httpmock.NewStringResponder(200, c.jsonResponse))

// 		jsonBodyToBeUnmarshalled := CurlAndReturn(c.targetUrl)

// 		marshallerr := json.Unmarshal(jsonBodyToBeUnmarshalled, &elastest)
// 		if marshallerr != nil {
// 			fmt.Println(marshallerr)
// 		} // not sure I care about this error yet

// 		if elastest.Stat != c.expectedStatus {
// 			t.Error("Expected %v but got %v", c.expectedStatus, elastest.Stat)
// 		}
// 	}
// }

// I want to fully deprecate curlandreturn without json
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

// deprecating Authcurl in favor of AuthcurlJson, so skipping writing a test against it
func TestAuthcurlJson(t *testing.T) {
	// i'll hit this in a bit
}

func TestCheckneo(t *testing.T) {
	cases := struct {
		role string
		addressRole string
		boolStringResponse string
		expectedResult int
	}{
		{
			role: "master",
			addressRole: "master",
			boolStringResponse: "true",
			expectedResult: 0,
		},
		{
			role: "slave",
			addressRole: "master",
			boolStringResponse: "false",
			expectedResult: 2,
		},
		{
			role: "master",
			addressRole: "slave",
			boolStringResponse: "false",
			expectedResult: 2,
		},
		{
			role: "slave",
			addressRole: "slave",
			boolStringResponse: "true",
			expectedResult: 0,
		},
	}
	// add in address and auth as static values of some kind
	for _, c in range cases {
		auth := "bmVvNGo6bmVvNGo=" // neo4j:neo4j
		target := fmt.Sprintf("http://localhost:7474/db/manage/server/ha/%v", c.addressRole)
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		// make a simple string endpoint that says true at 'master'
		// Neo4j HA clustering shows "true" when you go to http://localhost:7474/db/manage/server/ha/master and localhost, in this case, is the HA master
		httpmock.RegisterResponder("GET", target, httpmock.NewStringResponder(200, c.boolStringResponse))
		
	}

}