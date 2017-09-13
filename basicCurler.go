package main

import (
	"net/http"
	"fmt"
  //
	// "os"
	// "io"
)

func main() {

	resp, err := http.Get("https://ifconfig.co")
	if err != nil {
		fmt.Println("lolz you screwed up")
	}

	fmt.Println(resp.StatusCode)

	//theBody, _ := ioutil.ReadAll(os.Stdout)
	// theBody, _ := io.Copy(os.Stdout, resp.Body)
  theBody := resp.Body
	fmt.Println(theBody)
	defer resp.Body.Close()
}
