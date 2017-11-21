package grab

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
)

func TestCurlAndReturn(t *testing.T) {
	httpmock.RegisterResponder("GET", "http://localhost/", func(req *http.Request) (*http.Response, error) {
		resp, err := httpmock.NewJsonResponse(200, resultJson)
		if err != nil {
			return httpmock.NewStringResponse(500, ""), nil
		}
		return resp, nil
	},
	)
	CurlAndReturnJson()
}
