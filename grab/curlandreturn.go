package grab

import (
	"fmt"
	"io/ioutil"
	"net/http"
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
